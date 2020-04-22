package session

import (
	"fmt"
	"github.com/OpenIoTHub/server-go/config"
	"github.com/OpenIoTHub/utils/models"
	"github.com/OpenIoTHub/utils/msg"
	"github.com/OpenIoTHub/utils/mux"
	"github.com/xtaci/kcp-go"
	"log"
	"net"
)

//listenerHdl
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

//connHdl
func connHdl(conn net.Conn) {
	var session *mux.Session
	var token *models.TokenClaims
	var err error
	rawMsg, err := msg.ReadMsg(conn)
	if err != nil {
		return
	}
	switch m := rawMsg.(type) {
	case *models.GatewayLogin:
		{
			//log.Println(m.Token)
			token, err = models.DecodeToken(config.ConfigMode.Security.LoginKey, m.Token)
			if err != nil {
				fmt.Printf(err.Error())
				conn.Close()
				return
			}
			if token.Permission&1 != 1 {
				fmt.Printf("token type err ,not n")
				conn.Close()
				return
			}
			config := mux.DefaultConfig()
			//config.EnableKeepAlive = false
			session, err = mux.Client(conn, config)
			if err != nil {
				fmt.Printf(err.Error())
				conn.Close()
				return
			}
			log.Println("新内网客户端登录： runId：" + token.RunId + " 系统：" + m.Os + "芯片架构：" + m.Arch)
			//sessions[token.RunId]=session
			sess := &Session{Id: token.RunId, Conn: &conn, Ssession: session, WorkConn: make(chan net.Conn, 5)}
			//:TODO 新的登录存储之前先清除旧的同id登录
			sessions.SetSession(token.RunId, sess)
		}

	case *models.GatewayWorkConn:
		{
			//:TODO	内网主动新创建的用来接收数据传输业务的连接
			log.Println("获取到一个内网主动发起的工作连接")
			sess, err := sessions.GetSession(m.RunId)
			if err != nil {
				conn.Close()
				return
			}
			sess.WorkConn <- conn
		}

	case *models.OpenIoTHubLogin:
		{
			token, err = models.DecodeToken(config.ConfigMode.Security.LoginKey, m.Token)
			if err != nil {
				fmt.Printf(err.Error())
				conn.Close()
				return
			}
			if token.Permission != 2 {
				fmt.Printf("token type err ,not n")
				conn.Close()
				return
			}
			log.Println("新访问器登录上线： runId：" + token.RunId + " 系统：" + m.Os + "芯片架构：" + m.Arch)
			//sessions[token.RunId]=session
			//sess := &Session{Id: token.RunId, Conn: &conn, Ssession: session}
			//SetSession(token.RunId, sess)
			go openIoTHubLoginHdl(token.RunId, conn)
		}
	default:
		{
			//:TODO 为什么重连会跑到
			log.Println("从端口获取两种登录类别之一错误")
			conn.Close()
		}
	}
}
