package nettool

import "github.com/xtaci/kcp-go/v5"

func SetYamuxConn(kcpconn *kcp.UDPSession) {
	//配置
	//kcpconn.SetDeadline(time.Now().Add(time.Second * 5))
	kcpconn.SetStreamMode(true)
	kcpconn.SetWriteDelay(false)
	kcpconn.SetNoDelay(0, 100, 1, 1)
	kcpconn.SetWindowSize(128, 256)
	kcpconn.SetMtu(1350)
	kcpconn.SetACKNoDelay(true)
}
