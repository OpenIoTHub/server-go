package main

import (
	"git.iotserv.com/iotserv/server/config"
	"git.iotserv.com/iotserv/server/session"
	"git.iotserv.com/iotserv/server/udpapi"
	"time"
)

func main() {
	go session.RunTLS(config.ConfigMode.Common.TlsPort)
	go session.RunTCP(config.ConfigMode.Common.TcpPort)
	go session.RunKCP(config.ConfigMode.Common.KcpPort)
	go udpapi.RunApiServer(config.ConfigMode.Common.UdpApiPort)
	for {
		time.Sleep(time.Hour * 99999)
	}
}
