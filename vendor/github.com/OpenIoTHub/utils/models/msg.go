package models

import (
	"github.com/jacobsa/go-serial/serial"
	"reflect"
)

const (
	//////1
	TypeLogin                 = 'l'
	TypeConnectTCP            = 't'
	TypeConnectSTCP           = 's'
	TypeConnectUDP            = 'u'
	TypeConnectWs             = '1'
	TypeConnectWss            = '2'
	TypeConnectSerialPort     = 'p'
	TypeConnectToLogin        = 'c'
	TypePing                  = '3'
	TypePong                  = '4'
	TypeNewSubSession         = 'n'
	TypeReqNewP2PCtrl         = 'a'
	TypeRemoteNetInfo         = 'b'
	TypeReqNewP2PCtrlAsClient = 'd'
	TypeOK                    = 'e'
	TypeCheckStatusRequest    = 'f'
	TypeCheckStatusResponse   = 'g'
	TypeNewService            = 'h'
	TypeConnectSSH            = 'i'
	TypeRequestNewWorkConn    = 'j'
	TypeNewWorkConn           = 'k'
	TypeJsonResponse          = 'm'
)

var (
	TypeMap       map[byte]reflect.Type
	TypeStringMap map[reflect.Type]byte
)

func init() {
	TypeMap = make(map[byte]reflect.Type)
	TypeStringMap = make(map[reflect.Type]byte)
	////////2
	TypeMap[TypeLogin] = reflect.TypeOf(Login{})

	TypeMap[TypeConnectTCP] = reflect.TypeOf(ConnectTCP{})
	TypeMap[TypeConnectSTCP] = reflect.TypeOf(ConnectSTCP{})
	TypeMap[TypeConnectUDP] = reflect.TypeOf(ConnectUDP{})
	TypeMap[TypeConnectWs] = reflect.TypeOf(ConnectWs{})
	TypeMap[TypeConnectWss] = reflect.TypeOf(ConnectWss{})
	TypeMap[TypeConnectSerialPort] = reflect.TypeOf(ConnectSerialPort{})
	TypeMap[TypeConnectToLogin] = reflect.TypeOf(ConnectToLogin{})
	TypeMap[TypePing] = reflect.TypeOf(Ping{})
	TypeMap[TypePong] = reflect.TypeOf(Pong{})
	TypeMap[TypeNewSubSession] = reflect.TypeOf(NewSubSession{})
	TypeMap[TypeReqNewP2PCtrl] = reflect.TypeOf(ReqNewP2PCtrl{})
	TypeMap[TypeRemoteNetInfo] = reflect.TypeOf(RemoteNetInfo{})
	TypeMap[TypeReqNewP2PCtrlAsClient] = reflect.TypeOf(ReqNewP2PCtrlAsClient{})
	TypeMap[TypeOK] = reflect.TypeOf(OK{})
	TypeMap[TypeCheckStatusRequest] = reflect.TypeOf(CheckStatusRequest{})
	TypeMap[TypeCheckStatusResponse] = reflect.TypeOf(CheckStatusResponse{})
	//新的服务
	TypeMap[TypeNewService] = reflect.TypeOf(NewService{})
	TypeMap[TypeConnectSSH] = reflect.TypeOf(ConnectSSH{})
	TypeMap[TypeRequestNewWorkConn] = reflect.TypeOf(RequestNewWorkConn{})
	TypeMap[TypeNewWorkConn] = reflect.TypeOf(NewWorkConn{})
	TypeMap[TypeJsonResponse] = reflect.TypeOf(JsonResponse{})
	for k, v := range TypeMap {
		TypeStringMap[v] = k
	}
}

// Message wraps socket packages for communicating between frpc and frps.
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
