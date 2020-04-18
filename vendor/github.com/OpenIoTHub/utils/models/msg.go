package models

import (
	"github.com/jacobsa/go-serial/serial"
	"net"
	"reflect"
)

var (
	TypeMap       = make(map[string]reflect.Type)
	TypeStringMap = make(map[reflect.Type]string)

	Types = []reflect.Type{
		//服务器需要处理的消息
		reflect.TypeOf(GatewayLogin{}),
		reflect.TypeOf(GatewayWorkConn{}),
		reflect.TypeOf(OpenIoTHubLogin{}),
		//连接的消息
		reflect.TypeOf(ConnectTCP{}),
		reflect.TypeOf(ConnectSTCP{}),
		reflect.TypeOf(ConnectUDP{}),
		reflect.TypeOf(ConnectWs{}),
		reflect.TypeOf(ConnectWss{}),
		reflect.TypeOf(ConnectSerialPort{}),
		reflect.TypeOf(ConnectSSH{}),
		//P2P相关的消息
		reflect.TypeOf(NewSubSession{}),
		reflect.TypeOf(ReqNewP2PCtrl{}),
		reflect.TypeOf(RemoteNetInfo{}),
		reflect.TypeOf(ReqNewP2PCtrlAsClient{}),
		//状态验证消息
		reflect.TypeOf(CheckStatusRequest{}),
		reflect.TypeOf(CheckStatusResponse{}),
		//新的服务
		reflect.TypeOf(NewService{}),
		reflect.TypeOf(RequestNewWorkConn{}),

		// UDP API
		reflect.TypeOf(GetMyUDPPublicAddr{}),
		reflect.TypeOf(net.UDPAddr{}),

		reflect.TypeOf(Ping{}),
		reflect.TypeOf(Pong{}),

		reflect.TypeOf(OK{}),
		reflect.TypeOf(Error{}),

		reflect.TypeOf(JsonResponse{}),
	}
)

func init() {
	for _, v := range Types {
		TypeMap[v.String()] = v
		TypeStringMap[v] = v.String()
	}
}

type Message interface{}

// login
type GatewayLogin struct {
	Token string
	Os    string
	Arch  string
}

// Connect TO
type OpenIoTHubLogin struct {
	Token string
	Os    string
	Arch  string
}

type NewSubSession struct{}

// connect       //tcp,stcp,udp,serialport,ws,wss
type ConnectTCP struct {
	TargetIP   string
	TargetPort int
}

type ConnectSTCP struct {
	TargetIP   string
	TargetPort int
}

type ConnectUDP struct {
	TargetIP   string
	TargetPort int
}

type ConnectSerialPort serial.OpenOptions

type ConnectWs struct {
	TargetUrl string
	Protocol  string
	Origin    string
}

type ConnectWss struct {
	TargetUrl string
	Protocol  string
	Origin    string
}

type ConnectSSH struct {
	TargetIP   string
	TargetPort int
	UserName   string
	PassWord   string
}

///Ping
type Ping struct{}
type Pong struct{}

//P2P让远端以listener身份运行
type ReqNewP2PCtrl struct {
	IntranetIp   string
	IntranetPort int
	ExternalIp   string
	ExternalPort int
}

//让内网端以dial的身份连接我
type ReqNewP2PCtrlAsClient struct {
	IntranetIp   string
	IntranetPort int
	ExternalIp   string
	ExternalPort int
}

//TODO:NETINFO Model
type RemoteNetInfo struct {
	IntranetIp   string
	IntranetPort int
	ExternalIp   string
	ExternalPort int
}

type CheckStatusRequest struct {
	Type string
	Addr string
}

type CheckStatusResponse struct {
	//Code:0:在线;1:离线
	Code    int
	Message string
}

type NewService struct {
	Type   string
	Config string
}

type RequestNewWorkConn struct {
	Type   string
	Config string
}

type GatewayWorkConn struct {
	RunId  string
	Secret string
}

type JsonResponse struct {
	Code   int
	Msg    string
	Result string
}

///plugin
type InstallPlugin struct {
	TargetUrl string
}

type UpgradePlugin struct {
	TargetUrl string
}

type RemovePlugin struct {
	TargetUrl string
}

type RunPlugin struct {
	TargetUrl string
}

type StopPlugin struct {
	TargetUrl string
}

///query installed plugin
type QueryInstalledPlugin struct{}
type RespInstalledPlugin struct{}

///rsponse Msg

type Msg struct {
	MsgType    string
	MsgContent string
}

type GetMyUDPPublicAddr struct{}

type OK struct{}

type Error struct {
	Code    int
	Message string
}
