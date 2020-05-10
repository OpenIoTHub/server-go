package main

import (
	"fmt"
	"github.com/OpenIoTHub/server-go/config"
	"github.com/OpenIoTHub/server-go/session"
	"github.com/OpenIoTHub/utils/net"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"time"
)

var (
	version = "dev"
	commit  = ""
	date    = ""
	builtBy = ""
)

func main() {
	myApp := cli.NewApp()
	myApp.Name = "server-go"
	myApp.Usage = "-c [config File Path]"
	myApp.Version = fmt.Sprintf("%s(commit:%s,build on:%s,buildBy:%s)", version, commit, date, builtBy)
	myApp.Flags = []cli.Flag{
		//TODO 应该设置工作目录，各组件共享
		&cli.StringFlag{
			Name:        "config",
			Aliases:     []string{"c"},
			Value:       config.DefaultConfigFilePath,
			Usage:       "config file path",
			EnvVars:     []string{"ServerConfigFilePath"},
			Destination: &config.DefaultConfigFilePath,
		},
	}
	myApp.Action = func(c *cli.Context) error {
		err := run()
		if err != nil {
			os.Exit(1)
		}
		for {
			time.Sleep(time.Hour)
		}
	}
	err := myApp.Run(os.Args)
	if err != nil {
		log.Println(err.Error())
	}
}

func run() (err error) {
	err = config.LoadConfig()
	if err != nil {
		log.Println(err)
		return
	}
	go session.SessionsCtl.RunTLS()
	go session.SessionsCtl.RunTCP()
	go session.SessionsCtl.RunKCP()
	go session.SessionsCtl.StartgRpcListenAndServ()
	go session.SessionsCtl.StartHttpListenAndServ()
	go nettool.RunUDPApiServer(config.ConfigMode.Common.UdpApiPort)
	go nettool.RunKCPApiServer(config.ConfigMode.Common.KcpApiPort)
	log.Println("服务器正在运行，内网端配置请根据本服务器配置填写！")
	log.Println("config.ConfigMode.Common.KcpApiPort:", config.ConfigMode.Common.KcpApiPort)
	return
}
