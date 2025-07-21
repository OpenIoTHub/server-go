package main

import (
	"fmt"
	"github.com/OpenIoTHub/server-go/config"
	"github.com/OpenIoTHub/server-go/manager"
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
	defer func() {
		if r := recover(); r != nil {
			log.Printf("main Recovered from panic: %v\n", r) // 记录日志
		}
	}()
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
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:        "config",
					Aliases:     []string{"c"},
					Value:       config.DefaultConfigFilePath,
					Usage:       "config file path",
					EnvVars:     []string{"GatewayConfigFilePath"},
					Destination: &config.DefaultConfigFilePath,
				},
			},
			Action: func(c *cli.Context) error {
				err := config.LoadConfig()
				if err != nil {
					log.Println(err)
					return err
				}
				GateWayToken, OpenIoTHubToken, err := models.GetTokenByServerConfig(&config.ConfigMode, 200000000000)
				if err != nil {
					return err
				}
				log.Println("Generated one pair of token:")
				log.Println("注意不要复制了换行符:")
				log.Println("====================Gateway Token:->====================")
				log.Println(GateWayToken)
				log.Println("====================OpenIoTHub Token:->=================")
				log.Println(OpenIoTHubToken)
				log.Println("========================================")
				return nil
			},
		},
		{
			Name:    "init",
			Aliases: []string{"i"},
			Usage:   "init config file",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:        "config",
					Aliases:     []string{"c"},
					Value:       config.DefaultConfigFilePath,
					Usage:       "config file path",
					EnvVars:     []string{"GatewayConfigFilePath"},
					Destination: &config.DefaultConfigFilePath,
				},
			},
			Action: func(c *cli.Context) error {
				config.InitConfigFile()
				return nil
			},
		},
		{
			Name:    "test",
			Aliases: []string{"t"},
			Usage:   "test this command",
			Action: func(c *cli.Context) error {
				fmt.Println("ok")
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
	manager.InitSessionsCtl()
	manager.LoadConfigFromIoTManager()
	go manager.SessionsCtl.RunTLS()
	go manager.SessionsCtl.RunTCP()
	go manager.SessionsCtl.RunKCP()
	go manager.SessionsCtl.StartgRpcListenAndServ()
	go manager.SessionsCtl.StartHttpListenAndServ()
	go nettool.RunUDPApiServer(config.ConfigMode.Common.UdpApiPort)
	go nettool.RunKCPApiServer(config.ConfigMode.Common.KcpApiPort)
	log.Println("服务器正在运行，网关配置请根据本服务器配置填写！")
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
