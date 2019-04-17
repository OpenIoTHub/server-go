package config

import (
	"fmt"
	"git.iotserv.com/iotserv/utils/models"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

var ConfigMode models.ServerConfig

func init() {
	var err error
	var configFilePath = "./server.yaml"
	_, err = os.Stat(configFilePath)
	if err != nil {
		fmt.Println("没有找到配置文件：", configFilePath)
		fmt.Println("开始生成默认的空白配置文件，请填写配置文件后重复运行本程序")
		ConfigMode.Common.BindAddr = "0.0.0.0"
		ConfigMode.Common.KcpPort = 7200
		ConfigMode.Common.TcpPort = 7200
		ConfigMode.Common.TlsPort = 7300
		ConfigMode.Common.UdpApiPort = 7300
		ConfigMode.Security.LoginKey = "SETByYourSelf."
		//	生成配置文件模板
		err = writeConfigFile(ConfigMode, configFilePath)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		os.Exit(1)
	}
	ConfigMode, err = GetConfig(configFilePath)
	if err != nil {
		fmt.Printf(err.Error())
		os.Exit(1)
	}
}

func GetConfig(configFilePath string) (configMode models.ServerConfig, err error) {
	content, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	err = yaml.Unmarshal(content, &configMode)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	return
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
