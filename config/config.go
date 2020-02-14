package config

import (
	"fmt"
	"github.com/OpenIoTHub/utils/models"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
)

var ConfigMode models.ServerConfig
var ConfigFileName = "server.yaml"
var ConfigFilePath = "./server.yaml"

func LoadConfig() (err error) {
	//是否是snapcraft应用，如果是则从snapcraft指定的工作目录保存配置文件
	appDataPath, havaAppDataPath := os.LookupEnv("SNAP_USER_DATA")
	if havaAppDataPath {
		ConfigFilePath = filepath.Join(appDataPath, ConfigFileName)
	}
	_, err = os.Stat(ConfigFilePath)
	if err != nil {
		fmt.Println("没有找到配置文件：", ConfigFilePath)
		fmt.Println("开始生成默认的空白配置文件，请填写配置文件后重复运行本程序")
		ConfigMode.Common.BindAddr = "0.0.0.0"
		ConfigMode.Common.KcpPort = 34320
		ConfigMode.Common.TcpPort = 34320
		ConfigMode.Common.TlsPort = 34321
		ConfigMode.Common.UdpApiPort = 34321
		ConfigMode.Security.LoginKey = "HLLdsa544&*S"
		//	生成配置文件模板
		err = writeConfigFile(ConfigMode, ConfigFilePath)
		if err != nil {
			return
		}
		fmt.Println("配置文件写入成功,路径为：", ConfigFilePath)
		fmt.Println("你也可以修改上述配置文件后在运行")
	}
	fmt.Println("使用配置文件：", ConfigFilePath)
	ConfigMode, err = GetConfig(ConfigFilePath)
	if err != nil {
		return
	}
	return
}

//从配置文件路径解析配置文件的内容
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

//将配置内容写入指定的配置文件
func writeConfigFile(configMode models.ServerConfig, path string) (err error) {
	configByte, err := yaml.Marshal(configMode)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if ioutil.WriteFile(path, configByte, 0644) == nil {
		fmt.Println("写入配置文件文件成功!")
		return
	}
	return
}
