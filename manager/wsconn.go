package manager

import (
	"io"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var wsUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// handleWebSocketLogin upgrades an HTTP request to WebSocket, wraps it as
// a net.Conn, and passes it to connHdl for normal gateway login processing.
// After Upgrade the connection is hijacked so the http.Server no longer owns
// it; the yamux session created by connHdl keeps it alive.
func (sm *SessionsManager) handleWebSocketLogin(w http.ResponseWriter, r *http.Request) {
	ws, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("[WebSocket] upgrade failed: %v", err)
		return
	}
	log.Printf("[WebSocket] gateway client connected from %s", r.RemoteAddr)
	sm.connHdl(newWSNetConn(ws))
}

// wsNetConn adapts a gorilla/websocket.Conn into a net.Conn interface.
// Read provides a continuous byte stream by chaining across WebSocket binary messages.
// Write sends each call as a single binary WebSocket message.
type wsNetConn struct {
	ws     *websocket.Conn
	reader io.Reader
	rMu    sync.Mutex
	wMu    sync.Mutex
}

func newWSNetConn(ws *websocket.Conn) net.Conn {
	return &wsNetConn{ws: ws}
}

func (c *wsNetConn) Read(b []byte) (int, error) {
	c.rMu.Lock()
	defer c.rMu.Unlock()
	for {
		if c.reader != nil {
			n, err := c.reader.Read(b)
			if n > 0 {
				return n, nil
			}
			if err != io.EOF {
				return 0, err
			}
			c.reader = nil
		}
		_, reader, err := c.ws.NextReader()
		if err != nil {
			return 0, err
		}
		c.reader = reader
	}
}

func (c *wsNetConn) Write(b []byte) (int, error) {
	c.wMu.Lock()
	defer c.wMu.Unlock()
	if err := c.ws.WriteMessage(websocket.BinaryMessage, b); err != nil {
		return 0, err
	}
	return len(b), nil
}

func (c *wsNetConn) Close() error {
	return c.ws.Close()
}

func (c *wsNetConn) LocalAddr() net.Addr {
	return c.ws.LocalAddr()
}

func (c *wsNetConn) RemoteAddr() net.Addr {
	return c.ws.RemoteAddr()
}

func (c *wsNetConn) SetDeadline(t time.Time) error {
	if err := c.ws.SetReadDeadline(t); err != nil {
		return err
	}
	return c.ws.SetWriteDeadline(t)
}

func (c *wsNetConn) SetReadDeadline(t time.Time) error {
	return c.ws.SetReadDeadline(t)
}

func (c *wsNetConn) SetWriteDeadline(t time.Time) error {
	return c.ws.SetWriteDeadline(t)
}
