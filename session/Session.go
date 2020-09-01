package session

import (
	"github.com/OpenIoTHub/utils/models"
	"github.com/OpenIoTHub/utils/msg"
	"github.com/libp2p/go-yamux"
	"log"
	"net"
)

//
type Session struct {
	Id             string
	OS             string
	ARCH           string
	Version        string
	Conn           *net.Conn
	GatewaySession *yamux.Session
	WorkConn       chan net.Conn
}

func (sess *Session) GetStream() (*yamux.Stream, error) {
	return sess.GatewaySession.OpenStream()
}

func (sess *Session) RequestNewWorkConn() error {
	stream, err := sess.GetStream()
	if err != nil {
		return err
	}
	return msg.WriteMsg(stream, &models.RequestNewWorkConn{
		Type:   "tcp",
		Config: "",
	})
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
