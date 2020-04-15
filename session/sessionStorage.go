package session

import (
	"errors"
	"fmt"
	"github.com/OpenIoTHub/utils/mux"
	"net"
)

type Session struct {
	Id       string
	Conn     *net.Conn
	Ssession *mux.Session
	WorkConn chan net.Conn
}

type Sessions map[string]*Session

var sessions = make(Sessions)

//:TODO 存活检测
func (sess *Session) Task() {
	//defer DelSession(sess.Id)
	//Loop:
	//for {
	//	select {
	//		case <-sess.heartbeat.C:
	//			stream,err:=sess.Ssession.OpenStream()
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

func (sess Sessions) GetSession(id string) (*Session, error) {
	if _, ok := sess[id]; ok {
		if sess[id].Ssession.IsClosed() {
			sess.DelSession(id)
			return nil, errors.New("sessions 处于断线状态")
		}
		return sess[id], nil //存在
	} else {
		return nil, errors.New("sessions id未注册")
	}
}

func (sess Sessions) GetStream(id string) (*mux.Stream, error) {
	mysession, err := sess.GetSession(id)
	if err != nil {
		fmt.Printf(err.Error())
		return nil, err
	}
	fmt.Printf("get session ok")
	stream, err := mysession.Ssession.OpenStream()
	fmt.Printf("open stream")
	if err != nil {
		fmt.Printf(err.Error())
		return nil, err
	}
	return stream, err
}

func (sess Sessions) SetSession(id string, session *Session) {
	sess.DelSession(id)
	sess[id] = session
}

func (sess Sessions) DelSession(id string) {
	if _, ok := sess[id]; ok {
		if sess[id].Ssession != nil && !sess[id].Ssession.IsClosed() {
			sess[id].Ssession.Close()
		}
		if sess[id].Conn != nil {
			myconn := *sess[id].Conn
			myconn.Close()
		}
	}
	delete(sess, id)
}
