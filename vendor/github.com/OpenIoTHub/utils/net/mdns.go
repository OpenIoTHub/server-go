package nettool

import (
	"context"
	"fmt"
	"github.com/grandcat/zeroconf"
	"github.com/satori/go.uuid"
	"log"
	"strings"
	"time"
)

var MDNSServiceBaseInfo = map[string]string{
	"name":                 "OpenIoTHub服务",
	"model":                "com.iotserv.services.web",
	"author":               "Farry",
	"email":                "newfarry@126.com",
	"home-page":            "https://github.com/OpenIoTHub",
	"firmware-respository": "https://github.com/iotdevice",
	"firmware-version":     "1.0",
}

func CheckComponentExist(model string) (bool, error) {
	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		log.Fatalln("Failed to initialize resolver:", err.Error())
		return false, err
	}

	entries := make(chan *zeroconf.ServiceEntry)
	//TODO 是否需要手动关闭channel？
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(2))
	defer cancel()
	err = resolver.Browse(ctx, "_iotdevice._tcp", "local", entries)
	if err != nil {
		log.Fatalln("Failed to browse:", err.Error())
		return false, err
	}

	<-ctx.Done()
	for entry := range entries {
		for _, text := range entry.Text {
			keyValue := strings.Split(text, "=")
			if len(keyValue) == 2 && keyValue[0] == "model" && keyValue[1] == model {
				return true, nil
			}
		}
	}
	return false, nil
}

//检查mdns发布的服务的信息是否存在错误
func CheckmDNSServiceInfo(info map[string]string) error {
	var keyRequire = []string{
		"name",
		"model",
		"mac",
		"id",
		"author",
		"email",
		"home-page",
		"firmware-respository",
		"firmware-version",
	}
	for _, name := range keyRequire {
		if _, ok := info[name]; !ok {
			return fmt.Errorf("mDNSServiceInfo: %s not exist error", name) //存在
		}
	}
	return nil
}

func RegistermDNSService(info map[string]string, port int) (*zeroconf.Server, error) {
	if _, ok := info["mac"]; !ok {
		macs, err := GetMacs()
		if err == nil && len(macs) > 0 {
			for _, vMac := range macs {
				if vMac != "" {
					info["mac"] = vMac
					break
				}
			}
		}
	}
	if _, ok := info["mac"]; !ok {
		info["mac"] = uuid.Must(uuid.NewV4()).String()
	}
	if _, ok := info["id"]; !ok {
		info["id"] = uuid.Must(uuid.NewV4()).String()
	}
	err := CheckmDNSServiceInfo(info)
	if err != nil {
		return nil, err
	}
	var txt = []string{}
	for key, value := range info {
		txt = append(txt, fmt.Sprintf("%s=%s", key, value))
	}
	return zeroconf.Register(fmt.Sprintf("%s-%s", info["model"], info["mac"]), "_iotdevice._tcp", "local.", port, txt, nil)
}
