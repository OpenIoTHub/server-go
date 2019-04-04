package session

import (
	"errors"
	"fmt"
	//"github.com/xtaci/smux"
	"git.iotserv.com/iotserv/utils/mux"
	"net"
	//"time"
	//"git.iotserv.com/iotserv/utils/msg"
	//"git.iotserv.com/iotserv/utils/models"
)

type Session struct {
	Id       string
	Conn     *net.Conn
	Ssession *mux.Session
	WorkConn chan net.Conn
}

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

var sessions = make(map[string]*Session)

func GetSession(id string) (*Session, error) {
	if _, ok := sessions[id]; ok {
		if sessions[id].Ssession.IsClosed() {
			DelSession(id)
			return nil, errors.New("sessions 处于断线状态")
		}
		return sessions[id], nil //存在
	} else {
		return nil, errors.New("sessions id未注册")
	}
}

func GetStream(id string) (*mux.Stream, error) {
	mysession, err := GetSession(id)
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

func SetSession(id string, session *Session) {
	DelSession(id)
	sessions[id] = session
}

func DelSession(id string) {
	if _, ok := sessions[id]; ok {
		if sessions[id].Ssession != nil && !sessions[id].Ssession.IsClosed() {
			sessions[id].Ssession.Close()
		}
		if sessions[id].Conn != nil {
			myconn := *sessions[id].Conn
			myconn.Close()
		}
	}
	delete(sessions, id)
}
