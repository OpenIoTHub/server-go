package manager

import (
	"context"
	"net"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"

	"github.com/gorilla/websocket"
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
	// WebSocket upgrade for gateway login (gateway-chrome connects via wss://)
	// Only allow on *.iothub.cloud subdomains with the gateway path; others fall through to proxy.
	if websocket.IsWebSocketUpgrade(r) && isIoTHubDomain(r.Host) && r.URL.Path == "/api/v1/ws/gateway" {
		sm.handleWebSocketLogin(w, r)
		return
	}
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

const iotHubDomainSuffix = ".iothub.cloud"

func isIoTHubDomain(host string) bool {
	h := strings.Split(host, ":")[0]
	h = strings.ToLower(h)
	return h == "iothub.cloud" || strings.HasSuffix(h, iotHubDomainSuffix)
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
