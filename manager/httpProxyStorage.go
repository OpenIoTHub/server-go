package manager

import (
	"encoding/json"
	"log"
)

// TODO 根据配置文件确定从内存获取映射表还是redis
func (sm *SessionsManager) GetOneHttpProxy(domain string) (*HttpProxy, error) {
	log.Println("query:", domain)
	hpBytes, err := sm.HttpProxyRuntimeStorage.GetValueByKeyToBytes(domain)
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

func (sm *SessionsManager) GetAllHttpProxy() map[string]*HttpProxy {
	var HttpProxyMap = make(map[string]*HttpProxy)
	keys, err := sm.HttpProxyRuntimeStorage.GetAllKeys()
	if err != nil {
		return HttpProxyMap
	}
	for _, key := range keys {
		hpBytes, err := sm.HttpProxyRuntimeStorage.GetValueByKeyToBytes(key)
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

func (sm *SessionsManager) AddOrUpdateHttpProxy(httpProxy *HttpProxy) error {
	httpProxyBytes, err := json.Marshal(httpProxy)
	if err != nil {
		return err
	}
	return sm.HttpProxyRuntimeStorage.SetValueByKey(httpProxy.Domain, httpProxyBytes)
}

func (sm *SessionsManager) DelHttpProxy(domain string) {
	sm.HttpProxyRuntimeStorage.DelValueByKey(domain)
}

func (sm *SessionsManager) UpdateHttpProxyByMap(HttpProxyMap map[string]*HttpProxy) {
	for _, hp := range HttpProxyMap {
		sm.AddOrUpdateHttpProxy(hp)
	}
}
