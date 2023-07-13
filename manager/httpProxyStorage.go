package manager

import (
	"encoding/json"
	"fmt"
	"github.com/OpenIoTHub/server-go/config"
	"log"
)

// TODO 根据配置文件确定从内存获取映射表还是redis
func (sm *SessionsManager) GetOneHttpProxy(domain string) (*HttpProxy, error) {
	log.Println("query:", domain)
	if config.ConfigMode.RedisConfig.Enabled {
		hpBytes, err := sm.GetRedisValueByKeyToBytes(domain)
		if err != nil {
			return nil, err
		}
		var httpProxyModel = &HttpProxy{}
		err = json.Unmarshal(hpBytes, httpProxyModel)
		if err != nil {
			return nil, err
		}
		//TODO Update Status
		return httpProxyModel, nil
	}
	if _, ok := sm.HttpProxyMap[domain]; ok {
		return sm.HttpProxyMap[domain], nil //存在
	}
	log.Printf("httpProxy id未注册%s", domain)
	return nil, fmt.Errorf("httpProxy id未注册:%s", domain)
}

func (sm *SessionsManager) GetAllHttpProxy() map[string]*HttpProxy {
	if config.ConfigMode.RedisConfig.Enabled {
		var HttpProxyMap = make(map[string]*HttpProxy)
		keys, err := sm.GetAllRedisKey()
		if err != nil {
			return HttpProxyMap
		}
		for _, key := range keys {
			hpBytes, err := sm.GetRedisValueByKeyToBytes(key)
			if err != nil {
				continue
			}
			var httpProxyModel = &HttpProxy{}
			err = json.Unmarshal(hpBytes, httpProxyModel)
			if err != nil {
				continue
			}
			HttpProxyMap[key] = httpProxyModel
		}
		return HttpProxyMap
	}
	return sm.HttpProxyMap
}

func (sm *SessionsManager) AddHttpProxy(httpProxy *HttpProxy) error {
	if config.ConfigMode.RedisConfig.Enabled {
		httpProxyBytes, err := json.Marshal(httpProxy)
		if err != nil {
			return err
		}
		return sm.SetRedisKeyValue(httpProxy.Domain, httpProxyBytes)
	}
	if _, ok := sm.HttpProxyMap[httpProxy.Domain]; ok {
		return fmt.Errorf("域名%s已经被占用！", httpProxy.Domain) //存在
	}
	sm.HttpProxyMap[httpProxy.Domain] = httpProxy
	return nil
}

func (sm *SessionsManager) DelHttpProxy(domain string) {
	if config.ConfigMode.RedisConfig.Enabled {
		sm.DelRedisByKey(domain)
	}
	delete(sm.HttpProxyMap, domain)
}

func (sm *SessionsManager) UpdateHttpProxyByMap(HttpProxyMap map[string]*HttpProxy) {
	if config.ConfigMode.RedisConfig.Enabled {
		for _, hp := range HttpProxyMap {
			sm.AddHttpProxy(hp)
		}
	}
}
