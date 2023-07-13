package session

import (
	"errors"
	"github.com/OpenIoTHub/utils/models"
	"github.com/OpenIoTHub/utils/msg"
	"github.com/libp2p/go-yamux"
	"log"
	"net"
	"time"
)

type Session struct {
	Id             string
	OS             string
	ARCH           string
	Version        string
	DisableMuxer   bool
	Conn           *net.Conn
	GatewaySession *yamux.Session
	WorkConn       chan net.Conn
}

func (sess *Session) GetStream() (net.Conn, error) {
	if sess.GatewaySession == nil || sess.DisableMuxer {
		return sess.GetNewWorkConn()
	}
	return sess.GatewaySession.OpenStream()
}

func (sess *Session) RequestNewWorkConn() error {
	if sess.GatewaySession != nil && !sess.DisableMuxer {
		stream, err := sess.GetStream()
		if err != nil {
			return err
		}
		return msg.WriteMsg(stream, &models.RequestNewWorkConn{
			Type:   "tcp",
			Config: "",
		})
	} else if sess.DisableMuxer {
		return msg.WriteMsg(*sess.Conn, &models.RequestNewWorkConn{
			Type:   "tcp",
			Config: "",
		})
	}
	return errors.New("RequestNewWorkConn err")
}

func (sess *Session) GetNewWorkConn() (net.Conn, error) {
	//TODO 考虑提前缓存连接以提高性能，但是得做好保活
	var workConn net.Conn
	err := sess.RequestNewWorkConn()
	if err != nil {
		log.Println(err.Error())
		return workConn, err
	}
	//超时返回错误
	select {
	case workConn = <-sess.WorkConn:
		log.Println("获取工作连接成功！")
		return workConn, err
	case <-time.After(time.Second * 3):
		return workConn, errors.New("获取WorkConn超时")
	}
}

// :TODO 存活检测
func (sess *Session) Task() {
	//defer DelSession(sess.Id)
	//Loop:
	//for {
	//	select {
	//		case <-sess.heartbeat.C:
	//			stream,err:=sess.GatewaySession.OpenStream()
	//			if err != nil{
	//				log.Printf(err.Error())
	//				break Loop
	//			}
	//			err=msg.WriteMsg(stream, &models.Ping{})
	//			if err != nil{
	//				log.Printf(err.Error())
	//				break Loop
	//			}
	//			stream.Close()
	//		}
	//}
	log.Printf("end session Task")
}
