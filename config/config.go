package config

import (
	"fmt"
	"git.iotserv.com/iotserv/utils/models"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

var ConfigMode models.ServerConfig

func init() {
	var err error
	ConfigMode, err = GetConfig("./server.yaml")
	if err != nil {
		fmt.Printf(err.Error())
		return
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
