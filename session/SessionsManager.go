package session

import (
	"errors"
	"github.com/OpenIoTHub/server-go/config"
	"github.com/OpenIoTHub/utils/io"
	"github.com/OpenIoTHub/utils/models"
	"github.com/OpenIoTHub/utils/msg"
	"github.com/libp2p/go-yamux"
	"log"
	"net"
	"time"
)

type SessionsManager struct {
	Session      map[string]*Session
	HttpProxyMap map[string]*HttpProxy
}

var SessionsCtl = SessionsManager{
	Session:      make(map[string]*Session),
	HttpProxyMap: make(map[string]*HttpProxy),
}

func (sess *SessionsManager) GetSession(id string) (*Session, error) {
	if _, ok := sess.Session[id]; ok {
		if sess.Session[id].GatewaySession == nil || sess.Session[id].GatewaySession.IsClosed() {
			sess.DelSession(id)
			return nil, errors.New("Session 处于断线状态")
		}
		return sess.Session[id], nil //存在
	} else {
		return nil, errors.New("Session id未注册")
	}
}

func (sess *SessionsManager) GetStream(id string) (*yamux.Stream, error) {
	mysession, err := sess.GetSession(id)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	log.Println("get session ok")
	stream, err := mysession.GatewaySession.OpenStream()
	log.Println("open stream")
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return stream, err
}

func (sess *SessionsManager) SetSession(id string, session *Session) {
	sess.DelSession(id)
	sess.Session[id] = session
}

func (sess *SessionsManager) DelSession(id string) {
	if _, ok := sess.Session[id]; ok {
		if sess.Session[id].GatewaySession != nil && !sess.Session[id].GatewaySession.IsClosed() {
			sess.Session[id].GatewaySession.Close()
		}
		if sess.Session[id].Conn != nil {
			myconn := *sess.Session[id].Conn
			myconn.Close()
		}
	}
	delete(sess.Session, id)
}

//connHdl
func (sess *SessionsManager) connHdl(conn net.Conn) {
	var session *yamux.Session
	var token *models.TokenClaims
	var err error
	rawMsg, err := msg.ReadMsg(conn)
	if err != nil {
		log.Println(err)
		return
	}
	switch m := rawMsg.(type) {
	case *models.GatewayLogin:
		{
			//log.Println(m.Token)
			token, err = models.DecodeToken(config.ConfigMode.Security.LoginKey, m.Token)
			if err != nil {
				log.Println(err.Error())
				conn.Close()
				return
			}
			if token.Permission&1 != 1 {
				log.Println("token type err ,not n")
				conn.Close()
				return
			}
			config := yamux.DefaultConfig()
			//config.EnableKeepAlive = false
			session, err = yamux.Client(conn, config)
			if err != nil {
				log.Println(err.Error())
				conn.Close()
				return
			}
			log.Printf("新Gateway登录： runId：%s 系统：%s 芯片架构：%s", token.RunId, m.Os, m.Arch)
			//sess[token.RunId]=session
			session := &Session{Id: token.RunId, Conn: &conn, GatewaySession: session, WorkConn: make(chan net.Conn, 5)}
			//:TODO 新的登录存储之前先清除旧的同id登录
			sess.SetSession(token.RunId, session)
		}

	case *models.GatewayWorkConn:
		{
			//:TODO	内网主动新创建的用来接收数据传输业务的连接
			log.Println("获取到一个Gateway主动发起的工作连接")
			log.Println(m.RunId)
			log.Println(m.Secret)
			session, err := sess.GetSession(m.RunId)
			if err != nil {
				log.Println(err)
				conn.Close()
				return
			}
			session.WorkConn <- conn
		}

	case *models.OpenIoTHubLogin:
		{
			token, err = models.DecodeToken(config.ConfigMode.Security.LoginKey, m.Token)
			if err != nil {
				log.Println(err.Error())
				conn.Close()
				return
			}
			if token.Permission != 2 {
				log.Println("token type err ,not 2")
				conn.Close()
				return
			}
			log.Printf("新OpenIoTHub登录： runId：%s 系统：%s 芯片架构：%s", token.RunId, m.Os, m.Arch)
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
func (sess *SessionsManager) openIoTHubLoginHdl(id string, conn net.Conn) {
	resp := func(err error) {
		code := 0
		msgStr := ""
		if err != nil {
			code = 1
			msgStr = err.Error()
			defer conn.Close()
		}
		msg.WriteMsg(conn, &models.CheckStatusResponse{
			Code:    code,
			Message: msgStr,
		})
		time.Sleep(time.Millisecond * 100)

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
