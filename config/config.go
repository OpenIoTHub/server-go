package config

import "runtime"

var Setting = make(map[string]string)

func init() {
	//可以修改的配置参数

	//管理API的认证
	Setting["explorerWebUser"] = ""
	Setting["explorerWebPass"] = ""
	//用来访问管理API的域名
	Setting["apiHost"] = "mcunode.com"
	//是否提供流量转发能力
	Setting["canForward"] = "true"
	//不可修改的配置参数
	Setting["AppName"] = "nat-explorer"
	//设置api的监听端口和beego应用的运行模式
	Setting["apiPort"] = "1081"
	if runtime.GOOS == "android" {
		//用于代理http的http出口端口
		Setting["proxyHttpPort"] = "1080"
		//用于代理http的https出口端口
		Setting["proxyHttpsPort"] = "1443"
	} else {
		//用于代理http的http出口端口
		Setting["proxyHttpPort"] = "80"
		//用于代理http的https出口端口
		Setting["proxyHttpsPort"] = "43"
	}

}
