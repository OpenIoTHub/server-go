package nettool

import (
	"errors"
	"fmt"
	"github.com/OpenIoTHub/utils/crypto"
	"github.com/OpenIoTHub/utils/models"
	"github.com/OpenIoTHub/utils/msg"
	"github.com/xtaci/kcp-go/v5"
	"log"
	"net"
	"strconv"
	"time"
)

func RunApiServer(port int) {
	//listener, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP("0.0.0.0"), Port: port})
	listener, err := kcp.ListenWithOptions(fmt.Sprintf("0.0.0.0:%d", port), nil, 10, 3)
	if err != nil {
		log.Println(err)
		return
	}
	go kcpListenerHdl(listener)
}

func kcpListenerHdl(listener *kcp.Listener) {
	for {
		conn, err := listener.AcceptKCP()
		if err != nil {
			log.Println("UDP API 错误:")
			log.Println(err)
			continue
		}
		go connHdl(conn)
	}
}

func connHdl(conn *kcp.UDPSession) {
	//defer conn.Close()
	remoteAddr := conn.RemoteAddr().(*net.UDPAddr)
	rawMsg, err := msg.ReadMsg(conn)
	if err != nil {
		return
	}
	switch m := rawMsg.(type) {
	case *models.GetMyUDPPublicAddr:
		{
			_ = m
			_ = msg.WriteMsg(conn, remoteAddr)
		}

	default:
		{
			//:TODO 为什么重连会跑到
			log.Println("从端口获取两种登录类别之一错误")
			_ = msg.WriteMsg(conn, remoteAddr)
		}
	}
}

//获取一个listener的外部地址和端口
func GetExternalIpPort(listener *net.UDPConn, token *crypto.TokenClaims) (ip string, port int, err error) {
	//udpaddr, err := net.ResolveUDPAddr("udp", token.Host+":"+strconv.Itoa(token.P2PApiPort))
	conn, err := kcp.DialWithOptions(token.Host+":"+strconv.Itoa(token.P2PApiPort), nil, 10, 3)
	if err != nil {
		fmt.Printf("%s", err.Error())
		return "", 0, err
	}

	err = conn.SetDeadline(time.Now().Add(time.Duration(3 * time.Second)))
	if err != nil {
		fmt.Printf("%s", err.Error())
		return "", 0, err
	}

	err = msg.WriteMsg(conn, models.GetMyUDPPublicAddr{})
	if err != nil {
		fmt.Printf("%s", err.Error())
		return "", 0, err
	}

	log.Println("发送到服务器确定成功！等待确定外网ip和port")
	addr, err := msg.ReadMsg(conn)
	log.Println("获取api的UDP包成功，开始解析自己listener出口地址和端口")
	if err != nil {
		fmt.Printf("获取listener的出口出错: %s", err.Error())
		return "", 0, err
	}

	switch m := addr.(type) {
	case *net.UDPAddr:
		{
			return m.IP.String(), m.Port, err
		}

	case *models.Error:
		{
			return "", 0, errors.New(m.Message)
		}

	default:
		{
			//:TODO 为什么重连会跑到
			log.Println("从端口获取两种登录类别之一错误")
			return "", 0, errors.New("获取UDP的外网地址失败:错误的信息返回")
		}
	}
}
