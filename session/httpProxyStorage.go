package session

import (
	"errors"
	"fmt"
	"log"
)

//TODO 根据配置文件确定从内存获取映射表还是redis
func (sm *SessionsManager) GetOneHttpProxy(domain string) (*HttpProxy, error) {
	if _, ok := sm.HttpProxyMap[domain]; ok {
		go sm.HttpProxyMap[domain].UpdateRemotePortStatus()
		return sm.HttpProxyMap[domain], nil //存在
	}
	log.Printf("httpProxy id未注册")
	return nil, errors.New("httpProxy id未注册")
}

func (sm *SessionsManager) GetAllHttpProxy() map[string]*HttpProxy {
	for _, hp := range sm.HttpProxyMap {
		go hp.UpdateRemotePortStatus()
	}
	return sm.HttpProxyMap
}

func (sm *SessionsManager) AddHttpProxy(httpProxy *HttpProxy) error {
	if _, ok := sm.HttpProxyMap[httpProxy.Domain]; ok {
		return fmt.Errorf("域名%s已经被占用！", httpProxy.Domain) //存在
	}
	go httpProxy.UpdateRemotePortStatus()
	sm.HttpProxyMap[httpProxy.Domain] = httpProxy
	return nil
}

func (sm *SessionsManager) DelHttpProxy(domain string) {
	delete(sm.HttpProxyMap, domain)
}
