package main

import (
	"fmt"
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
	fmt.Println("服务器正在运行，内网端配置请根据本服务器配置填写！")
	for {
		time.Sleep(time.Hour * 99999)
	}
}
