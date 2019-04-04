package main

import (
	"git.iotserv.com/iotserv/server/config"
	"git.iotserv.com/iotserv/server/session"
	"git.iotserv.com/iotserv/server/udpapi"
	"time"
)

func main() {
	//crypto.Salt = "abc"
	go session.RunTLS(config.TlsPort)
	go session.RunTCP(config.TcpPort)
	go session.RunKCP(config.KcpPort)
	go udpapi.RunApiServer(config.UdpApiPort)
	for {
		time.Sleep(time.Hour * 99999)
	}
}
