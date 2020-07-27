package config

import (
	"github.com/OpenIoTHub/utils/models"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

var ConfigMode models.ServerConfig

func LoadConfig() (err error) {
	//是否是snapcraft应用，如果是则从snapcraft指定的工作目录保存配置文件
	appDataPath, havaAppDataPath := os.LookupEnv("SNAP_USER_DATA")
	if havaAppDataPath {
		DefaultConfigFilePath = filepath.Join(appDataPath, DefaultConfigFileName)
	}
	_, err = os.Stat(DefaultConfigFilePath)
	if err != nil {
		InitConfigFile()
	}
	log.Println("使用配置文件：", DefaultConfigFilePath)
	ConfigMode, err = GetConfig(DefaultConfigFilePath)
	if err != nil {
		return
	}
	return
}

func InitConfigFile() {
	var err error
	log.Println("没有找到配置文件：", DefaultConfigFilePath)
	log.Println("开始生成默认的空白配置文件")
	ConfigMode.Common.BindAddr = DefaultBindAddr
	ConfigMode.Common.KcpPort = DefaultKcpPort
	ConfigMode.Common.TcpPort = DefaultTcpPort
	ConfigMode.Common.TlsPort = DefaultTlsPort
	ConfigMode.Common.GrpcPort = DefaultGrpcPort
	ConfigMode.Common.HttpPort = DefaultHttpPort
	ConfigMode.Common.HttpsPort = DefaultHttpsPort
	ConfigMode.Common.UdpApiPort = DefaultUdpApiPort
	ConfigMode.Common.KcpApiPort = DefaultKcpApiPort
	ConfigMode.Security.LoginKey = DefaultLoginKey
	//	生成配置文件模板
	err = writeConfigFile(ConfigMode, DefaultConfigFilePath)
	if err != nil {
		log.Printf("写入默认的配置文件失败：%s\n", err.Error())
		return
	}
	log.Println("配置文件写入成功,路径为：", DefaultConfigFilePath)
	log.Println("你也可以修改上述配置文件后在运行")
}

//从配置文件路径解析配置文件的内容
func GetConfig(configFilePath string) (configMode models.ServerConfig, err error) {
	content, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		log.Println(err.Error())
		return
	}
	err = yaml.Unmarshal(content, &configMode)
	if err != nil {
		log.Println(err.Error())
		return
	}
	return
}

//将配置内容写入指定的配置文件
func writeConfigFile(configMode models.ServerConfig, path string) (err error) {
	configByte, err := yaml.Marshal(configMode)
	if err != nil {
		log.Println(err.Error())
		return
	}
	err = os.MkdirAll(filepath.Dir(path), 0644)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(path, configByte, 0644)
	return
}
