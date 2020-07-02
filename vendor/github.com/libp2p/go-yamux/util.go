package yamux

import (
	"io"
	"sync"
	"sync/atomic"

	pool "github.com/libp2p/go-buffer-pool"
)

// asyncSendErr is used to try an async send of an error
func asyncSendErr(ch chan error, err error) {
	if ch == nil {
		return
	}
	select {
	case ch <- err:
	default:
	}
}

// asyncNotify is used to signal a waiting goroutine
func asyncNotify(ch chan struct{}) {
	select {
	case ch <- struct{}{}:
	default:
	}
}

// min computes the minimum of a set of values
func min(values ...uint32) uint32 {
	m := values[0]
	for _, v := range values[1:] {
		if v < m {
			m = v
		}
	}
	return m
}

type segmentedBuffer struct {
	cap     uint32
	pending uint32
	len     uint32
	bm      sync.Mutex
	b       [][]byte
}

// NewSegmentedBuffer allocates a ring buffer.
func NewSegmentedBuffer(initialCapacity uint32) segmentedBuffer {
	return segmentedBuffer{cap: initialCapacity, b: make([][]byte, 0)}
}

func (s *segmentedBuffer) Len() int {
	return int(atomic.LoadUint32(&s.len))
}

func (s *segmentedBuffer) Cap() uint32 {
	return atomic.LoadUint32(&s.cap)
}

// If the space to write into + current buffer size has grown to half of the window size,
// grow up to that max size, and indicate how much additional space was reserved.
func (s *segmentedBuffer) GrowTo(max uint32, force bool) (bool, uint32) {
	s.bm.Lock()
	defer s.bm.Unlock()

	currentWindow := atomic.LoadUint32(&s.len) + atomic.LoadUint32(&s.cap) + s.pending
	if currentWindow > max {
		// somewhat counter-intuitively not an error.
		// note that len+cap is the 'window' that shouldn't exceed max or a reservation
		// would fail, triggering an error.
		// We pre-count 'pending' data where we've read a header and are working on
		// reading it into available data here, so that we don't undercount the remaining
		// window size, but that can mean this sum ends up larger than max.
		return false, 0
	}
	delta := max - currentWindow

	if delta < (max/2) && !force {
		return false, 0
	}

	atomic.AddUint32(&s.cap, delta)
	return true, delta
}

func (s *segmentedBuffer) TryReserve(space uint32) bool {
	// It is noticable that the check-and-set of pending is not atomic,
	// Due to this, accesses to pending are protected by bm.
	s.bm.Lock()
	defer s.bm.Unlock()
	if atomic.LoadUint32(&s.cap) < s.pending+space {
		return false
	}
	s.pending += space
	return true
}

func (s *segmentedBuffer) Read(b []byte) (int, error) {
	s.bm.Lock()
	defer s.bm.Unlock()
	if len(s.b) == 0 {
		return 0, io.EOF
	}
	n := copy(b, s.b[0])
	if n == len(s.b[0]) {
		pool.Put(s.b[0])
		s.b[0] = nil
		s.b = s.b[1:]
	} else {
		s.b[0] = s.b[0][n:]
	}
	if n > 0 {
		atomic.AddUint32(&s.len, ^uint32(n-1))
	}
	return n, nil
}

func (s *segmentedBuffer) Append(input io.Reader, length int) error {
	dst := pool.Get(length)
	n := 0
	read := 0
	var err error
	for n < length && err == nil {
		read, err = input.Read(dst[n:])
		n += read
	}
	if err == io.EOF {
		if length == n {
			err = nil
		} else {
			err = ErrStreamReset
		}
	}

	s.bm.Lock()
	defer s.bm.Unlock()
	if n > 0 {
		atomic.AddUint32(&s.len, uint32(n))
		// cap -= n
		atomic.AddUint32(&s.cap, ^uint32(n-1))
		s.pending = s.pending - uint32(length)
		s.b = append(s.b, dst[0:n])
	}
	return err
}
