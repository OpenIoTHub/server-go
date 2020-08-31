package session

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/OpenIoTHub/server-go/config"
	"github.com/OpenIoTHub/utils/file"
	"github.com/xtaci/kcp-go"
	"golang.org/x/crypto/acme/autocert"
	"log"
	"net"
	"net/http"
	"os"
)

func (sess SessionsManager) RunKCP() {
	listener, err := kcp.ListenWithOptions(fmt.Sprintf(":%d", config.ConfigMode.Common.KcpPort), nil, 10, 3)
	if err != nil {
		log.Println(err)
		return
	}
	sess.kcpListenerHdl(listener)
}

func (sess SessionsManager) RunTCP() {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", config.ConfigMode.Common.TcpPort))
	if err != nil {
		log.Println(err)
		return
	}
	sess.listenerHdl(listener)
}

func (sess SessionsManager) RunTLS() {
	_, err := os.Stat(config.ConfigMode.Security.TlsCertFilePath)
	if err != nil {
		log.Printf("warning:File Path:%s Not Exist! So tls server NOT Available!", config.ConfigMode.Security.TlsCertFilePath)
		return
	}
	_, err = os.Stat(config.ConfigMode.Security.TlsKeyFilePath)
	if err != nil {
		log.Printf("warning:File Path:%s Not Exist!  So tls server NOT Available!", config.ConfigMode.Security.TlsKeyFilePath)
		return
	}
	cer, err := tls.LoadX509KeyPair(config.ConfigMode.Security.TlsCertFilePath, config.ConfigMode.Security.TlsKeyFilePath)
	//cer, err := tls.LoadX509KeyPair("./cert.pem", "./key.pem")
	if err != nil {
		log.Println(err)
		return
	}
	tlsConfig := &tls.Config{Certificates: []tls.Certificate{cer}}
	listener, err := tls.Listen("tcp", fmt.Sprintf(":%d", config.ConfigMode.Common.TlsPort), tlsConfig)
	if err != nil {
		log.Println(err)
		return
	}
	sess.listenerHdl(listener)
}

//http(s)代理端口监听
func (sess SessionsManager) StartHttpListenAndServ() {
	var err error
	m := &autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: func(_ context.Context, host string) error { return nil },
	}
	dir := file.CacheDir()
	if err := os.MkdirAll(dir, 0700); err != nil {
		log.Printf("没有使用缓存目录来存放https证书: %v", err)
	} else {
		m.Cache = autocert.DirCache(dir)
	}

	go func() {
		serverHttp := http.Server{
			Addr:    fmt.Sprintf(":%d", config.ConfigMode.Common.HttpPort),
			Handler: &sess,
		}
		log.Printf("请访问浏览器访问http://127.0.0.1:%d/查看管理界面\n", config.ConfigMode.Common.HttpPort)
		err = serverHttp.ListenAndServe()
		if err != nil {
			log.Println(err.Error())
			serverHttp := http.Server{
				Addr:    fmt.Sprintf(":%s", "1083"),
				Handler: &sess,
			}
			log.Printf("%d端口被占用，请访问http://127.0.0.1:1083/\n", config.DefaultHttpPort)
			err = serverHttp.ListenAndServe()
			if err != nil {
				log.Println(err.Error())
			}
		}
	}()

	go func() {
		serverHttps := http.Server{
			Addr:      fmt.Sprintf(":%d", config.ConfigMode.Common.HttpsPort),
			Handler:   &sess,
			TLSConfig: m.TLSConfig(),
		}
		err = serverHttps.ListenAndServeTLS("", "")
		if err != nil {
			log.Println(err.Error())
			serverHttps := http.Server{
				Addr:      fmt.Sprintf(":%s", "1443"),
				Handler:   &sess,
				TLSConfig: m.TLSConfig(),
			}
			log.Println("1443端口被占用，请访问https://127.0.0.1:1443/")
			err = serverHttps.ListenAndServeTLS("", "")
			if err != nil {
				log.Println(err.Error())
			}
		}
	}()
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
