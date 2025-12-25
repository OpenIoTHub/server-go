package manager

import (
	"context"
	"net"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"
)

// 服务结构
type HttpProxy struct {
	Domain           string
	RunId            string
	RemoteIP         string
	RemotePort       int
	UserName         string
	Password         string
	IfHttps          bool
	Description      string
	RemotePortStatus bool
}

func (hp *HttpProxy) UpdateRemotePortStatus() {
	var online bool
	if hp.IfHttps {
		online, _ = SessionsCtl.CheckRemoteStatus("tls", hp.RunId, hp.RemoteIP, hp.RemotePort)
	} else {
		online, _ = SessionsCtl.CheckRemoteStatus("tcp", hp.RunId, hp.RemoteIP, hp.RemotePort)
	}
	hp.RemotePortStatus = online
}

//监听服务
//type Handle struct{}

func (sm *SessionsManager) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//:TODO 当前不支持websocket的代理（websocket支持basic Auth），或者从beego分发？
	//：TODO 非80端口r.Host的支持情况
	//log.Println("host:", r.Host)
	//log.Println("hostRequestURI:", r.RequestURI)
	//log.Println("hostRequestHEADER:", r.Header)
	//log.Println("r.URL.Scheme:", r.URL.Scheme)
	//if r.URL.Scheme == "http" {
	//	log.Println("是http请求")
	//} else if r.URL.Scheme == "https" {
	//	log.Println("是https请求")
	//}
	hostInfo, err := sm.GetOneHttpProxy(strings.Split(r.Host, ":")[0])
	if err != nil {
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	if hostInfo.UserName != "" || hostInfo.Password != "" {
		if u, p, ok := r.BasicAuth(); ok {
			if u != hostInfo.UserName || p != hostInfo.Password {
				w.Header().Set("WWW-Authenticate", `Basic realm="Dotcoo User Login"`)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
		} else {
			w.Header().Set("WWW-Authenticate", `Basic realm="Dotcoo User Login"`)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	}

	//是普通http的情况
	proxy := &httputil.ReverseProxy{}
	proxy.Transport = &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		DialContext:           sm.dialTcp,
		DialTLSContext:        sm.dialTls,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	proxy.ServeHTTP(w, r)
}

func (sm *SessionsManager) dialTcp(ctx context.Context, network, address string) (net.Conn, error) {
	//log.Printf("请求的地址addr：%s", address)
	end := strings.Index(address, ":")
	host := address[0:end]
	hostInfo, err := sm.GetOneHttpProxy(host) //id
	if err != nil {
		return nil, err
	}
	var stream net.Conn
	if hostInfo.IfHttps {
		stream, err = sm.ConnectToTls(hostInfo.RunId, hostInfo.RemoteIP, hostInfo.RemotePort)
	} else {
		stream, err = sm.ConnectToTcp(hostInfo.RunId, hostInfo.RemoteIP, hostInfo.RemotePort)
	}
	return stream, err
}

func (sm *SessionsManager) dialTls(ctx context.Context, network, address string) (net.Conn, error) {
	//log.Printf("请求的地址addr：%s", address)
	end := strings.Index(address, ":")
	host := address[0:end]
	hostInfo, err := sm.GetOneHttpProxy(host) //id
	if err != nil {
		return nil, err
	}
	var stream net.Conn
	if hostInfo.IfHttps {
		stream, err = sm.ConnectToTls(hostInfo.RunId, hostInfo.RemoteIP, hostInfo.RemotePort)
	} else {
		stream, err = sm.ConnectToTcp(hostInfo.RunId, hostInfo.RemoteIP, hostInfo.RemotePort)
	}
	return stream, err
}
