package models

import (
	"github.com/jacobsa/go-serial/serial"
	"reflect"
)

var (
	TypeMap       = make(map[string]reflect.Type)
	TypeStringMap = make(map[reflect.Type]string)

	Types = []reflect.Type{
		reflect.TypeOf(Login{}),
		reflect.TypeOf(ConnectTCP{}),
		reflect.TypeOf(ConnectSTCP{}),
		reflect.TypeOf(ConnectUDP{}),
		reflect.TypeOf(ConnectWs{}),
		reflect.TypeOf(ConnectWss{}),
		reflect.TypeOf(ConnectSerialPort{}),
		reflect.TypeOf(ConnectToLogin{}),
		reflect.TypeOf(Ping{}),
		reflect.TypeOf(Pong{}),
		reflect.TypeOf(NewSubSession{}),
		reflect.TypeOf(ReqNewP2PCtrl{}),
		reflect.TypeOf(RemoteNetInfo{}),
		reflect.TypeOf(ReqNewP2PCtrlAsClient{}),
		reflect.TypeOf(OK{}),
		reflect.TypeOf(CheckStatusRequest{}),
		reflect.TypeOf(CheckStatusResponse{}),
		//新的服务
		reflect.TypeOf(NewService{}),
		reflect.TypeOf(ConnectSSH{}),
		reflect.TypeOf(RequestNewWorkConn{}),
		reflect.TypeOf(NewWorkConn{}),
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
type Login struct {
	Token string `json:"token"`
	Os    string `json:"os"`
	Arch  string `json:"arch"`
}

// Connect TO
type ConnectToLogin struct {
	Token string `json:"token"`
	Os    string `json:"os"`
	Arch  string `json:"arch"`
}

type NewSubSession struct{}

// connect       //tcp,stcp,udp,serialport,ws,wss
type ConnectTCP struct {
	TargetIP   string `json:"target_ip"`
	TargetPort int    `json:"target_port"`
}

type ConnectSTCP struct {
	TargetIP   string `json:"target_ip"`
	TargetPort int    `json:"target_port"`
}

type ConnectUDP struct {
	TargetIP   string `json:"target_ip"`
	TargetPort int    `json:"target_port"`
}

type ConnectSerialPort serial.OpenOptions

type ConnectWs struct {
	TargetUrl string `json:"target_url"`
	Protocol  string `json:"protocol"`
	Origin    string `json:"origin"`
}

type ConnectWss struct {
	TargetUrl string `json:"target_url"`
	Protocol  string `json:"protocol"`
	Origin    string `json:"origin"`
}

type ConnectSSH struct {
	TargetIP   string `json:"target_ip"`
	TargetPort int    `json:"target_port"`
	UserName   string `json:"username"`
	PassWord   string `json:"password"`
}

///Ping
type Ping struct{}
type Pong struct{}

//P2P让远端以listener身份运行
type ReqNewP2PCtrl struct {
	IntranetIp   string `json:"intranet_ip"`
	IntranetPort int    `json:"intranet_port"`
	ExternalIp   string `json:"external_ip"`
	ExternalPort int    `json:"external_port"`
}

//让内网端以dial的身份连接我
type ReqNewP2PCtrlAsClient struct {
	IntranetIp   string `json:"intranet_ip"`
	IntranetPort int    `json:"intranet_port"`
	ExternalIp   string `json:"external_ip"`
	ExternalPort int    `json:"external_port"`
}

//TODO:NETINFO Model
type RemoteNetInfo struct {
	IntranetIp   string `json:"intranet_ip"`
	IntranetPort int    `json:"intranet_port"`
	ExternalIp   string `json:"external_ip"`
	ExternalPort int    `json:"external_port"`
}

type OK struct{}

type CheckStatusRequest struct {
	Type string `json:"type"`
	Addr string `json:"addr"`
}

type CheckStatusResponse struct {
	//Code:0:在线;1:离线
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type NewService struct {
	Type   string
	Config string
}

type RequestNewWorkConn struct {
	Type   string
	Config string
}

type NewWorkConn struct {
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
	TargetUrl string `json:"target_url"`
}

type UpgradePlugin struct {
	TargetUrl string `json:"target_url"`
}

type RemovePlugin struct {
	TargetUrl string `json:"target_url"`
}

type RunPlugin struct {
	TargetUrl string `json:"target_url"`
}

type StopPlugin struct {
	TargetUrl string `json:"target_url"`
}

///query installed plugin
type QueryInstalledPlugin struct{}
type RespInstalledPlugin struct{}

///rsponse Msg

type Msg struct {
	MsgType    string `json:"msg_type"`
	MsgContent string `json:"msg_content"`
}
