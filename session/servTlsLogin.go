package session

import (
	"crypto/tls"
	"fmt"
	"git.iotserv.com/iotserv/server/config"
	"os"
)

func RunTLS(port int) {
	_, err := os.Stat(config.ConfigMode.Security.TlsCertFilePath)
	if err != nil {
		fmt.Printf("File Path:%s Not Exist! So tls server NOT Available!\n", config.ConfigMode.Security.TlsCertFilePath)
		return
	}
	_, err = os.Stat(config.ConfigMode.Security.TlsKeyFilePath)
	if err != nil {
		fmt.Printf("File Path:%s Not Exist!  So tls server NOT Available!\n", config.ConfigMode.Security.TlsKeyFilePath)
		return
	}
	cer, err := tls.LoadX509KeyPair(config.ConfigMode.Security.TlsCertFilePath, config.ConfigMode.Security.TlsKeyFilePath)
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
