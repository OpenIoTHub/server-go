package session

import (
	//"github.com/xtaci/smux"
	"errors"
	"github.com/OpenIoTHub/utils/io"
	"github.com/OpenIoTHub/utils/models"
	"github.com/OpenIoTHub/utils/msg"
	"log"
	"net"
	"time"
)

//访问器的登录处理 conn : 访问器 stream ： 内网端
func openIoTHubLoginHdl(id string, conn net.Conn) {
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
	stream, err := sessions.GetStream(id)
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
	sess, err := sessions.GetSession(id)
	if err != nil {
		log.Println(err.Error())
		resp(err)
		return
	}
	//超时返回错误
	select {
	case workConn = <-sess.WorkConn:
		resp(nil)
		go io.Join(workConn, conn)
		return
	case <-time.After(time.Second * 3):
		resp(errors.New("获取内网连接超时"))
		return
	}
}
