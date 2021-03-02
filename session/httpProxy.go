package session

import (
	"context"
	"fmt"
	"github.com/OpenIoTHub/server-grpc-api/pb-go"
	"github.com/OpenIoTHub/utils/io"
	"github.com/OpenIoTHub/utils/net/httpUtil"
	"github.com/libp2p/go-yamux"
	"golang.org/x/net/websocket"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"
)

//服务结构
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

//grpc
func (sm *SessionsManager) CreateOneHTTP(ctx context.Context, in *pb.HTTPConfig) (*pb.HTTPConfig, error) {
	err := authOpenIoTHubGrpc(ctx, in.RunId)
	if err != nil {
		return in, status.Errorf(codes.Unauthenticated, err.Error())
	}

	return in, sm.AddHttpProxy(&HttpProxy{
		Domain:      in.Domain,
		RunId:       in.RunId,
		RemoteIP:    in.RemoteIP,
		RemotePort:  int(in.RemotePort),
		UserName:    in.UserName,
		Password:    in.Password,
		IfHttps:     in.IfHttps,
		Description: in.Description,
	})
}

func (sm *SessionsManager) DeleteOneHTTP(ctx context.Context, in *pb.HTTPConfig) (*pb.Empty, error) {
	err := authOpenIoTHubGrpc(ctx, in.RunId)
	if err != nil {
		return &pb.Empty{}, status.Errorf(codes.Unauthenticated, err.Error())
	}
	sm.DelHttpProxy(in.Domain)
	return &pb.Empty{}, nil

}

func (sm *SessionsManager) GetOneHTTP(ctx context.Context, in *pb.HTTPConfig) (*pb.HTTPConfig, error) {
	err := authOpenIoTHubGrpc(ctx, in.RunId)
	if err != nil {
		return in, status.Errorf(codes.Unauthenticated, err.Error())
	}
	config, err := sm.GetOneHttpProxy(in.Domain)
	if err != nil {
		return &pb.HTTPConfig{}, err
	}
	return &pb.HTTPConfig{
		Domain:           config.Domain,
		RunId:            config.RunId,
		RemoteIP:         config.RemoteIP,
		RemotePort:       int32(config.RemotePort),
		UserName:         config.UserName,
		Password:         config.Password,
		IfHttps:          config.IfHttps,
		Description:      config.Description,
		RemotePortStatus: config.RemotePortStatus,
	}, err
}

func (sm *SessionsManager) GetAllHTTP(ctx context.Context, in *pb.Device) (*pb.HTTPList, error) {
	var cfgs []*pb.HTTPConfig
	err := authOpenIoTHubGrpc(ctx, in.RunId)
	if err != nil {
		return &pb.HTTPList{HTTPConfigs: cfgs}, status.Errorf(codes.Unauthenticated, err.Error())
	}
	for _, config := range sm.GetAllHttpProxy() {
		if config.RunId == in.RunId && config.RemoteIP == in.Addr {
			cfgs = append(cfgs, &pb.HTTPConfig{
				Domain:           config.Domain,
				RunId:            config.RunId,
				RemoteIP:         config.RemoteIP,
				RemotePort:       int32(config.RemotePort),
				UserName:         config.UserName,
				Password:         config.Password,
				IfHttps:          config.IfHttps,
				Description:      config.Description,
				RemotePortStatus: config.RemotePortStatus,
			})
		}
	}
	return &pb.HTTPList{HTTPConfigs: cfgs}, nil
}

//grpc end

//监听服务
//type Handle struct{}

func (sm *SessionsManager) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//:TODO 当前不支持websocket的代理（websocket支持basic Auth），或者从beego分发？
	//：TODO 非80端口r.Host的支持情况
	log.Printf("Serve host:" + r.Host) //"http://127.0.0.1:8080/"
	var hostInfo *HttpProxy
	var err error
	hostPort := strings.Replace(strings.Replace(r.Host, "http://", "", -1), "https://", "", -1)
	hostInfo, err = sm.GetOneHttpProxy(strings.Split(hostPort, ":")[0])
	if err != nil {
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	if hostInfo.UserName != "" || hostInfo.Password != "" {
		ok := httpUtil.Auth(w, r, hostInfo.UserName, hostInfo.Password)
		if !ok {
			return
		}
	}
	//是websocket
	if v, ok := r.Header["Upgrade"]; ok {
		if v[0] == "websocket" {
			fn := func(ws *websocket.Conn) {
				var pro = "ws"
				var orgpro = "http"
				if hostInfo.IfHttps {
					pro = "wss"
					orgpro = "https"
				}
				conn, err := sm.ConnectToWs(hostInfo.RunId, fmt.Sprintf("%s://%s:%d%s", pro, hostInfo.RemoteIP, hostInfo.RemotePort, r.URL.String()),
					"", fmt.Sprintf("%s://%s:%d%s", orgpro, hostInfo.RemoteIP, hostInfo.RemotePort, r.URL.String()))
				if err != nil {
					log.Printf(err.Error())
					ws.Write([]byte(err.Error()))
					return
				}
				io.Join(ws, conn)
				w.Write([]byte("over"))
				return
			}
			websocket.Handler(fn).ServeHTTP(w, r)
		}
	}
	//是普通http的情况
	proxy := httputil.ReverseProxy{Director: func(request *http.Request) {
		//	修改代理的request
		if ip, _, err := net.SplitHostPort(strings.TrimSpace(request.RemoteAddr)); err == nil {
			request.Header.Add("REMOTE_ADDR", ip)
			//request.Header.Add("X-Forwarded-For", ip)
			request.Header.Add("X-Real-Ip", ip)
		}
	}}
	var pTransport http.RoundTripper = &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		DialContext:           sm.dial,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	proxy.Transport = pTransport

	proxy.ServeHTTP(w, r)
}

func (sm *SessionsManager) dial(c context.Context, network, address string) (net.Conn, error) {
	log.Printf("请求的地址addr：%s", address)
	hostPort := strings.Replace(strings.Replace(address, "http://", "", -1), "https://", "", -1)
	hostInfo, err := sm.GetOneHttpProxy(strings.Split(hostPort, ":")[0]) //id
	if err != nil {
		return nil, err
	}
	var stream *yamux.Stream
	if hostInfo.IfHttps {
		stream, err = sm.ConnectToTls(hostInfo.RunId, hostInfo.RemoteIP, hostInfo.RemotePort)
	} else {
		stream, err = sm.ConnectToTcp(hostInfo.RunId, hostInfo.RemoteIP, hostInfo.RemotePort)
	}
	return stream, err
}
