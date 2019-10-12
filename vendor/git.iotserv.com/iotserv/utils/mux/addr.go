package mux

import (
	"fmt"
	"net"
)

// hasAddr is used to get the address from the underlying connection
type hasAddr interface {
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
}

// muxAddr is used when we cannot get the underlying address
type muxAddr struct {
	Addr string
}

func (*muxAddr) Network() string {
	return "mux"
}

func (y *muxAddr) String() string {
	return fmt.Sprintf("mux:%s", y.Addr)
}

// Addr is used to get the address of the listener.
func (s *Session) Addr() net.Addr {
	return s.LocalAddr()
}

// LocalAddr is used to get the local address of the
// underlying connection.
func (s *Session) LocalAddr() net.Addr {
	addr, ok := s.conn.(hasAddr)
	if !ok {
		return &muxAddr{"local"}
	}
	return addr.LocalAddr()
}

// RemoteAddr is used to get the address of remote end
// of the underlying connection
func (s *Session) RemoteAddr() net.Addr {
	addr, ok := s.conn.(hasAddr)
	if !ok {
		return &muxAddr{"remote"}
	}
	return addr.RemoteAddr()
}

// LocalAddr returns the local address
func (s *Stream) LocalAddr() net.Addr {
	return s.session.LocalAddr()
}

// LocalAddr returns the remote address
func (s *Stream) RemoteAddr() net.Addr {
	return s.session.RemoteAddr()
}
