package session

import (
	"crypto/tls"
	"fmt"
	"github.com/OpenIoTHub/server-go/config"
	"github.com/xtaci/kcp-go"
	"log"
	"net"
	"os"
)

func (sess SessionsManager) RunKCP(port int) {
	listener, err := kcp.ListenWithOptions(fmt.Sprintf(":%d", port), nil, 10, 3)
	if err != nil {
		log.Println(err)
		return
	}
	sess.kcpListenerHdl(listener)
}

func (sess SessionsManager) RunTCP(port int) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Println(err)
		return
	}
	sess.listenerHdl(listener)
}

func (sess SessionsManager) RunTLS(port int) {
	_, err := os.Stat(config.ConfigMode.Security.TlsCertFilePath)
	if err != nil {
		log.Println("warning:File Path:%s Not Exist! So tls server NOT Available!", config.ConfigMode.Security.TlsCertFilePath)
		return
	}
	_, err = os.Stat(config.ConfigMode.Security.TlsKeyFilePath)
	if err != nil {
		log.Println("warning:File Path:%s Not Exist!  So tls server NOT Available!", config.ConfigMode.Security.TlsKeyFilePath)
		return
	}
	cer, err := tls.LoadX509KeyPair(config.ConfigMode.Security.TlsCertFilePath, config.ConfigMode.Security.TlsKeyFilePath)
	//cer, err := tls.LoadX509KeyPair("./cert.pem", "./key.pem")
	if err != nil {
		log.Println(err)
		return
	}
	tlsConfig := &tls.Config{Certificates: []tls.Certificate{cer}}
	listener, err := tls.Listen("tcp", fmt.Sprintf(":%d", port), tlsConfig)
	if err != nil {
		log.Println(err)
		return
	}
	sess.listenerHdl(listener)
}

///////////////////////////
//////
//////  Listenner处理
//////
///////////////////////////
func (sess SessionsManager) listenerHdl(listener net.Listener) {
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err.Error())
			continue
		}
		go sess.connHdl(conn)
	}
}

func (sess SessionsManager) kcpListenerHdl(listener *kcp.Listener) {
	defer listener.Close()
	for {
		conn, err := listener.AcceptKCP()
		if err != nil {
			log.Println(err.Error())
			continue
		}
		conn.SetStreamMode(true)
		conn.SetWriteDelay(false)
		conn.SetNoDelay(0, 40, 2, 1)
		conn.SetWindowSize(1024, 1024)
		conn.SetMtu(1472)
		conn.SetACKNoDelay(true)
		go sess.connHdl(conn)
	}
}
