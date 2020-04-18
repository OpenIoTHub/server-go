package session

import (
	"github.com/xtaci/kcp-go"
	"log"
	"net"
)

func listenerHdl(listener net.Listener) {
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err.Error())
			continue
		}
		go connHdl(conn)
	}
}

func kcpListenerHdl(listener *kcp.Listener) {
	defer listener.Close()
	for {
		conn, err := listener.AcceptKCP()
		if err != nil {
			log.Println(err.Error())
			continue
		}
		conn.SetStreamMode(true)
		conn.SetWriteDelay(false)
		conn.SetNoDelay(0, 40, 2, 1)
		conn.SetWindowSize(1024, 1024)
		conn.SetMtu(1472)
		conn.SetACKNoDelay(true)
		go connHdl(conn)
	}
}
