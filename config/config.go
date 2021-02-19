package config

import (
	"fmt"
	"github.com/OpenIoTHub/utils/models"
	"gopkg.in/yaml.v2"
	"io"
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
	//解析日志配置
	writers := []io.Writer{}
	if ConfigMode.LogConfig != nil && ConfigMode.LogConfig.EnableStdout {
		writers = append(writers, os.Stdout)
	}
	if ConfigMode.LogConfig != nil && ConfigMode.LogConfig.LogFilePath != "" {
		f, err := os.OpenFile(ConfigMode.LogConfig.LogFilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			log.Fatal(err)
		}
		writers = append(writers, f)
	}
	fileAndStdoutWriter := io.MultiWriter(writers...)
	log.SetOutput(fileAndStdoutWriter)
	//
	if err != nil {
		return
	}
	return
}

func InitConfigFile() {
	var err error
	ConfigMode.LogConfig = &models.LogConfig{
		EnableStdout: true,
		LogFilePath:  "",
	}
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
	ConfigMode.RedisConfig.Network = DefaultRedisNetwork
	ConfigMode.RedisConfig.Address = DefaultRedisAddress
	//	生成配置文件模板
	err = writeConfigFile(ConfigMode, DefaultConfigFilePath)
	if err == nil {
		fmt.Println("config created")
		return
	} else {
		log.Println("配置文件路径为：", DefaultConfigFilePath)
		log.Println("写入配置文件失败：", err.Error())
	}
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
