package session

import (
	"fmt"
	"github.com/OpenIoTHub/utils/mux"
	"net"
)

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
