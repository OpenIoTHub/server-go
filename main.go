package main

import (
	"fmt"
	"github.com/OpenIoTHub/server-go/config"
	"github.com/OpenIoTHub/server-go/session"
	web_ui "github.com/OpenIoTHub/server-go/web-ui"
	"github.com/OpenIoTHub/utils/models"
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
	myApp.Version = buildVersion(version, commit, date, builtBy)
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
	myApp.Commands = []*cli.Command{
		{
			Name:    "generate",
			Aliases: []string{"g"},
			Usage:   "generate one token for gateway and one token for OpenIoTHub",
			Action: func(c *cli.Context) error {
				err := config.LoadConfig()
				if err != nil {
					log.Println(err)
					return err
				}
				GateWayToken, OpenIoTHubToken, err := models.GetTokenByServerConfig(&config.ConfigMode, 1, 200000000000)
				if err != nil {
					return err
				}
				fmt.Println("Generated one pair of token:")
				fmt.Println("注意不要复制了换行符:")
				fmt.Println("====================Gateway Token:->====================")
				fmt.Println(GateWayToken)
				fmt.Println("====================OpenIoTHub Token:->=================")
				fmt.Println(OpenIoTHubToken)
				fmt.Println("========================================")
				return nil
			},
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
	go web_ui.RunWebStatic()
	log.Println("服务器正在运行，内网端配置请根据本服务器配置填写！")
	return
}

func buildVersion(version, commit, date, builtBy string) string {
	var result = version
	if commit != "" {
		result = fmt.Sprintf("%s\ncommit: %s", result, commit)
	}
	if date != "" {
		result = fmt.Sprintf("%s\nbuilt at: %s", result, date)
	}
	if builtBy != "" {
		result = fmt.Sprintf("%s\nbuilt by: %s", result, builtBy)
	}
	return result
}
