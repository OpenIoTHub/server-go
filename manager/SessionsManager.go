package manager

import (
	"errors"
	pb "github.com/OpenIoTHub/openiothub_grpc_api/pb-go/proto/server"
	"github.com/OpenIoTHub/server-go/config"
	"github.com/OpenIoTHub/server-go/iface/runtimeStorage"
	"github.com/OpenIoTHub/server-go/imp/runtimeStorage/memImp"
	"github.com/OpenIoTHub/server-go/imp/runtimeStorage/redisImp"
	"github.com/OpenIoTHub/server-go/session"
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
	Session map[string]*session.Session
	//TODO 设置和使用，还有用于黑白名单使用的存储
	HttpProxyRuntimeStorage runtimeStorage.RuntimeStorageIfce
	pb.UnimplementedHttpManagerServer
}

var SessionsCtl SessionsManager

func InitSessionsCtl() {
	var httpProxyRuntimeStorage runtimeStorage.RuntimeStorageIfce
	if config.ConfigMode.RedisConfig.Enabled {
		httpProxyRuntimeStorage = redisImp.NewRuntimeStorageRedisImp(
			&redis.Pool{
				MaxIdle:     256,
				MaxActive:   0,
				IdleTimeout: time.Duration(120),
				Dial: func() (redis.Conn, error) {
					redisConfigs := []redis.DialOption{redis.DialReadTimeout(time.Duration(1000) * time.Millisecond),
						redis.DialWriteTimeout(time.Duration(1000) * time.Millisecond),
						redis.DialConnectTimeout(time.Duration(1000) * time.Millisecond),
						redis.DialDatabase(config.ConfigMode.RedisConfig.Database),
					}
					if config.ConfigMode.RedisConfig.NeedAuth {
						redisConfigs = append(redisConfigs, redis.DialPassword(config.ConfigMode.RedisConfig.Password))
					}
					return redis.Dial(
						config.ConfigMode.RedisConfig.Network,
						config.ConfigMode.RedisConfig.Address,
						redisConfigs...,
					)
				},
			})
	} else {
		httpProxyRuntimeStorage = memImp.NewRuntimeStorageMemImp()
	}

	SessionsCtl = SessionsManager{
		Session:                 make(map[string]*session.Session),
		HttpProxyRuntimeStorage: httpProxyRuntimeStorage,
	}
}

func (sess *SessionsManager) GetSessionByID(id string) (*session.Session, error) {
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

func (sess *SessionsManager) GetStreamByID(id string) (net.Conn, error) {
	mysession, err := sess.GetSessionByID(id)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return mysession.GetStream()
}

func (sess *SessionsManager) GetNewWorkConnByID(id string) (net.Conn, error) {
	session, err := sess.GetSessionByID(id)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return session.GetNewWorkConn()
}

func (sess *SessionsManager) SetSession(id string, session *session.Session) {
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
		if sess.Session[id].WorkConn != nil {
			close(sess.Session[id].WorkConn)
		}
	}
	delete(sess.Session, id)
}

// connHdl
func (sess *SessionsManager) connHdl(conn net.Conn) {
	var yamuxSession *yamux.Session
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
			//token, err = models.DecodeToken(config.ConfigMode.Security.LoginKey, m.Token)
			//if err != nil {
			//	log.Println(err.Error())
			//	conn.Close()
			//	return
			//}
			//if !token.IfContainPermission(models.PermissionGatewayLogin) {
			//	log.Println("token type err ,not n")
			//	conn.Close()
			//	return
			//}
			//log.Printf("新Gateway登录： runId：%s 系统(OS)：%s 芯片架构(CPU Arch)：%s Version：%s 禁止muxer：%t",
			//	token.RunId, m.Os, m.Arch, m.Version, m.DisableMuxer)
			if m.DisableMuxer {
				//禁止muxer的网关
				//sess[token.RunId]=session
				gatewaySession := &session.Session{
					Id:             token.RunId,
					OS:             m.Os,
					ARCH:           m.Arch,
					Version:        m.Version,
					DisableMuxer:   m.DisableMuxer,
					Conn:           &conn,
					GatewaySession: nil,
					WorkConn:       make(chan net.Conn, 5)}
				//:TODO 新的登录存储之前先清除旧的同id登录
				sess.SetSession(token.RunId, gatewaySession)
				return
			}
			config := yamux.DefaultConfig()
			//config.EnableKeepAlive = false
			yamuxSession, err = yamux.Client(conn, config)
			if err != nil {
				log.Println(err.Error())
				conn.Close()
				return
			}
			//TODO 添加上线、下线日志储存以供用户查询
			//sess[token.RunId]=session
			gatewaySession := &session.Session{
				Id:             token.RunId,
				OS:             m.Os,
				ARCH:           m.Arch,
				Version:        m.Version,
				Conn:           &conn,
				GatewaySession: yamuxSession,
				WorkConn:       make(chan net.Conn, 5)}
			//:TODO 新的登录存储之前先清除旧的同id登录
			sess.SetSession(token.RunId, gatewaySession)
		}

	case *models.GatewayWorkConn:
		{
			//内网主动新创建的用来接收数据传输业务的连接
			//TODO 添加上线、下线日志储存以供用户查询
			log.Println("获取到一个Gateway主动发起的工作连接")
			log.Println("GatewayWorkConn:", m.RunId, "@", m.Version)
			//TODO 验证Secret
			token, _ = models.DecodeUnverifiedToken(config.ConfigMode.Security.LoginKey)
			//if err != nil {
			//	log.Println(err.Error())
			//	conn.Close()
			//	return
			//}
			//if !token.IfContainPermission(models.PermissionGatewayLogin) {
			//	log.Println("token type err ,not n")
			//	conn.Close()
			//	return
			//}
			//验证Secret Over
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
			if !token.IfContainPermission(models.PermissionOpenIoTHubLogin) {
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

// 访问器的登录处理 conn : 访问器 stream ： 网关
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
	//TODO 模拟嵌入式设备不支持mux的情况
	gatewaySess, _ := sess.GetSessionByID(id)
	if gatewaySess != nil && gatewaySess.DisableMuxer {
		resp(nil)
		//mux conn并分别与workConn桥接
		//config := yamux.DefaultConfig()
		//config.EnableKeepAlive = false
		session, err := yamux.Server(conn, yamux.DefaultConfig())
		if err != nil {
			log.Println(err.Error())
			conn.Close()
			return
		}
		for {
			mobileConn, err := session.Accept()
			if err != nil {
				log.Println(err)
				return
			}
			workConn, err := sess.GetStreamByID(id)
			if err != nil {
				log.Println(err.Error())
				resp(err)
				return
			}
			go io.Join(workConn, mobileConn)
		}
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
