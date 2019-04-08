package main

import (
	"fmt"
	"git.iotserv.com/iotserv/server/config"
	"git.iotserv.com/iotserv/server/session"
	"git.iotserv.com/iotserv/server/udpapi"
	"git.iotserv.com/iotserv/utils/models"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"time"
)

func main() {
	var configFilePath = "./server.yaml"
	var configMode = models.ServerConfig{}
	_, err := os.Stat(configFilePath)
	if err != nil {
		fmt.Println("没有找到配置文件：", configFilePath)
		fmt.Println("开始生成默认的空白配置文件，请填写配置文件后重复运行本程序")
		//	生成配置文件模板
		err = writeConfigFile(configMode, configFilePath)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		return
	}
	configMode, err = config.GetConfig(configFilePath)
	if err != nil {
		fmt.Printf(err.Error())
		return
	}
	go session.RunTLS(configMode.Common.TlsPort)
	go session.RunTCP(configMode.Common.TcpPort)
	go session.RunKCP(configMode.Common.KcpPort)
	go udpapi.RunApiServer(configMode.Common.UdpApiPort)
	for {
		time.Sleep(time.Hour * 99999)
	}
}

func writeConfigFile(configMode models.ServerConfig, path string) (err error) {
	configByte, err := yaml.Marshal(configMode)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if ioutil.WriteFile(path, configByte, 0644) == nil {
		fmt.Println("写入配置文件文件成功！\n")
		return
	}
	return
}
