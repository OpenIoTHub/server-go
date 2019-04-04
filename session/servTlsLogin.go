package session

import (
	"crypto/tls"
	"fmt"
	"git.iotserv.com/iotserv/server/config"
	"os"
)

func RunTLS(port int) {
	_, err := os.Stat(config.TlsCertFilePath)
	if err != nil {
		fmt.Printf("File Path:%s Not Exist! So tls server NOT Available!\n", config.TlsCertFilePath)
		return
	}
	_, err = os.Stat(config.TlsKeyFilePath)
	if err != nil {
		fmt.Printf("File Path:%s Not Exist!  So tls server NOT Available!\n", config.TlsKeyFilePath)
		return
	}
	cer, err := tls.LoadX509KeyPair(config.TlsCertFilePath, config.TlsKeyFilePath)
	//cer, err := tls.LoadX509KeyPair("./cert.pem", "./key.pem")
	if err != nil {
		fmt.Println(err)
		return
	}
	tlsConfig := &tls.Config{Certificates: []tls.Certificate{cer}}
	listener, err := tls.Listen("tcp", fmt.Sprintf(":%d", port), tlsConfig)
	if err != nil {
		fmt.Println(err)
		return
	}
	listenerHdl(listener)
}
