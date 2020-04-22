package session

import (
	"errors"
	"fmt"
	"github.com/OpenIoTHub/server-go/config"
	"github.com/OpenIoTHub/utils/io"
	"github.com/OpenIoTHub/utils/models"
	"github.com/OpenIoTHub/utils/msg"
	"github.com/OpenIoTHub/utils/mux"
	"log"
	"net"
	"time"
)

type SessionsManager map[string]*Session

var SessionsCtl = make(SessionsManager)

func (sess SessionsManager) GetSession(id string) (*Session, error) {
	if _, ok := sess[id]; ok {
		if sess[id].GatewaySession.IsClosed() {
			sess.DelSession(id)
			return nil, errors.New("Session 处于断线状态")
		}
		return sess[id], nil //存在
	} else {
		return nil, errors.New("Session id未注册")
	}
}

func (sess SessionsManager) GetStream(id string) (*mux.Stream, error) {
	mysession, err := sess.GetSession(id)
	if err != nil {
		fmt.Printf(err.Error())
		return nil, err
	}
	fmt.Printf("get session ok")
	stream, err := mysession.GatewaySession.OpenStream()
	fmt.Printf("open stream")
	if err != nil {
		fmt.Printf(err.Error())
		return nil, err
	}
	return stream, err
}

func (sess SessionsManager) SetSession(id string, session *Session) {
	sess.DelSession(id)
	sess[id] = session
}

func (sess SessionsManager) DelSession(id string) {
	if _, ok := sess[id]; ok {
		if sess[id].GatewaySession != nil && !sess[id].GatewaySession.IsClosed() {
			sess[id].GatewaySession.Close()
		}
		if sess[id].Conn != nil {
			myconn := *sess[id].Conn
			myconn.Close()
		}
	}
	delete(sess, id)
}

//connHdl
func (sess SessionsManager) connHdl(conn net.Conn) {
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
			//sess[token.RunId]=session
			session := &Session{Id: token.RunId, Conn: &conn, GatewaySession: session, WorkConn: make(chan net.Conn, 5)}
			//:TODO 新的登录存储之前先清除旧的同id登录
			sess.SetSession(token.RunId, session)
		}

	case *models.GatewayWorkConn:
		{
			//:TODO	内网主动新创建的用来接收数据传输业务的连接
			log.Println("获取到一个内网主动发起的工作连接")
			session, err := sess.GetSession(m.RunId)
			if err != nil {
				conn.Close()
				return
			}
			session.WorkConn <- conn
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
			//sess[token.RunId]=session
			//sess := &Session{Id: token.RunId, Conn: &conn, GatewaySession: session}
			//SetSession(token.RunId, sess)
			go sess.openIoTHubLoginHdl(token.RunId, conn)
		}
	default:
		{
			//:TODO 为什么重连会跑到
			log.Println("从端口获取两种登录类别之一错误")
			conn.Close()
		}
	}
}

//访问器的登录处理 conn : 访问器 stream ： 内网端
func (sess SessionsManager) openIoTHubLoginHdl(id string, conn net.Conn) {
	resp := func(err error) {
		code := 0
		msgStr := ""
		if err != nil {
			code = 1
			msgStr = err.Error()
		}
		msg.WriteMsg(conn, &models.CheckStatusResponse{
			Code:    code,
			Message: msgStr,
		})
		time.Sleep(time.Millisecond * 100)
		conn.Close()
	}
	var workConn net.Conn
	stream, err := sess.GetStream(id)
	if err != nil {
		log.Println(err.Error())
		resp(err)
		return
	}
	err = msg.WriteMsg(stream, &models.RequestNewWorkConn{
		Type:   "kcp",
		Config: "",
	})
	if err != nil {
		log.Println(err.Error())
		resp(err)
		return
	}
	session, err := sess.GetSession(id)
	if err != nil {
		log.Println(err.Error())
		resp(err)
		return
	}
	//超时返回错误
	select {
	case workConn = <-session.WorkConn:
		resp(nil)
		go io.Join(workConn, conn)
		return
	case <-time.After(time.Second * 3):
		resp(errors.New("获取内网连接超时"))
		return
	}
}

//
type Session struct {
	Id             string
	Conn           *net.Conn
	GatewaySession *mux.Session
	WorkConn       chan net.Conn
}

//:TODO 存活检测
func (sess *Session) Task() {
	//defer DelSession(sess.Id)
	//Loop:
	//for {
	//	select {
	//		case <-sess.heartbeat.C:
	//			stream,err:=sess.GatewaySession.OpenStream()
	//			if err != nil{
	//				fmt.Printf(err.Error())
	//				break Loop
	//			}
	//			err=msg.WriteMsg(stream, &models.Ping{})
	//			if err != nil{
	//				fmt.Printf(err.Error())
	//				break Loop
	//			}
	//			stream.Close()
	//		}
	//}
	fmt.Printf("end session Task")
}
