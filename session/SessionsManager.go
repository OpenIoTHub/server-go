package session

import (
	"errors"
	"github.com/OpenIoTHub/server-go/config"
	"github.com/OpenIoTHub/utils/io"
	"github.com/OpenIoTHub/utils/models"
	"github.com/OpenIoTHub/utils/msg"
	"github.com/gomodule/redigo/redis"
	"github.com/libp2p/go-yamux"
	"log"
	"net"
	"time"
)

type SessionsManager struct {
	Session      map[string]*Session
	HttpProxyMap map[string]*HttpProxy
	RedisPool    *redis.Pool
}

var SessionsCtl SessionsManager

func InitSessionsCtl() {
	SessionsCtl = SessionsManager{
		Session:      make(map[string]*Session),
		HttpProxyMap: make(map[string]*HttpProxy),
		RedisPool: &redis.Pool{
			MaxIdle:     256,
			MaxActive:   0,
			IdleTimeout: time.Duration(120),
			Dial: func() (redis.Conn, error) {
				if config.ConfigMode.RedisConfig.NeedAuth {
					return redis.Dial(
						config.ConfigMode.RedisConfig.Network,
						config.ConfigMode.RedisConfig.Address,
						redis.DialReadTimeout(time.Duration(1000)*time.Millisecond),
						redis.DialWriteTimeout(time.Duration(1000)*time.Millisecond),
						redis.DialConnectTimeout(time.Duration(1000)*time.Millisecond),
						redis.DialDatabase(config.ConfigMode.RedisConfig.Database),
					)
				}
				return redis.Dial(
					config.ConfigMode.RedisConfig.Network,
					config.ConfigMode.RedisConfig.Address,
					redis.DialReadTimeout(time.Duration(1000)*time.Millisecond),
					redis.DialWriteTimeout(time.Duration(1000)*time.Millisecond),
					redis.DialConnectTimeout(time.Duration(1000)*time.Millisecond),
					redis.DialDatabase(config.ConfigMode.RedisConfig.Database),
					redis.DialPassword(config.ConfigMode.RedisConfig.Password),
				)
			},
		},
	}
}

func (sess *SessionsManager) GetRedisConn() (redis.Conn, error) {
	conn := sess.RedisPool.Get()
	if err := conn.Err(); err != nil {
		return conn, err
	}
	return conn, nil
}

func (sess *SessionsManager) GetSessionByID(id string) (*Session, error) {
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

func (sess *SessionsManager) GetStreamByID(id string) (*yamux.Stream, error) {
	mysession, err := sess.GetSessionByID(id)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	stream, err := mysession.GetStream()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return stream, err
}

func (sess *SessionsManager) GetNewWorkConnByID(id string) (net.Conn, error) {
	session, err := sess.GetSessionByID(id)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return session.GetNewWorkConn()
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
			//TODO 添加上线、下线日志储存以供用户查询
			log.Printf("新Gateway登录： runId：%s 系统(OS)：%s 芯片架构(CPU Arch)：%s Version：%s", token.RunId, m.Os, m.Arch, m.Version)
			//sess[token.RunId]=session
			session := &Session{
				Id:             token.RunId,
				OS:             m.Os,
				ARCH:           m.Arch,
				Version:        m.Version,
				Conn:           &conn,
				GatewaySession: session,
				WorkConn:       make(chan net.Conn, 5)}
			//:TODO 新的登录存储之前先清除旧的同id登录
			sess.SetSession(token.RunId, session)
		}

	case *models.GatewayWorkConn:
		{
			//内网主动新创建的用来接收数据传输业务的连接
			//TODO 添加上线、下线日志储存以供用户查询
			log.Println("获取到一个Gateway主动发起的工作连接")
			log.Println("GatewayWorkConn:", m.RunId, "@", m.Version)
			//TODO 验证Secret
			session, err := sess.GetSessionByID(m.RunId)
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
			//TODO 添加上线、下线日志储存以供用户查询
			log.Printf("新OpenIoTHub登录： runId：%s 系统(OS)：%s 芯片架构(Arch)：%s Version：%s", token.RunId, m.Os, m.Arch, m.Version)
			//sess[token.RunId]=session
			//sess := &Session{Id: token.RunId, Conn: &conn, GatewaySession: session}
			//SetSession(token.RunId, sess)
			go sess.openIoTHubLoginHdl(token.RunId, conn)
		}
	default:
		{
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
	workConn, err := sess.GetNewWorkConnByID(id)
	if err != nil {
		log.Println(err.Error())
		resp(err)
		return
	}
	resp(nil)
	go io.Join(workConn, conn)
	return
}
