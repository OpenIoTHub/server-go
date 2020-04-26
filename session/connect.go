package session

import (
	"fmt"
	"github.com/OpenIoTHub/utils/models"
	"github.com/OpenIoTHub/utils/msg"
	"github.com/libp2p/go-yamux"
)

//Connect to tcp
func (sm *SessionsManager) ConnectToTcp(runId, remoteIp string, remotePort int) (*yamux.Stream, error) {
	stream, err := sm.GetStream(runId)
	if err != nil {
		return nil, err
	}
	msgsd := &models.ConnectTCP{
		TargetIP:   remoteIp,
		TargetPort: remotePort,
	}
	err = msg.WriteMsg(stream, msgsd)
	if err != nil {
		fmt.Printf(err.Error())
		return nil, err
	}
	return stream, nil
}

func (sm *SessionsManager) ConnectToTls(runId, remoteIp string, remotePort int) (*yamux.Stream, error) {
	stream, err := sm.GetStream(runId)
	if err != nil {
		return nil, err
	}
	msgsd := &models.ConnectSTCP{
		TargetIP:   remoteIp,
		TargetPort: remotePort,
	}
	err = msg.WriteMsg(stream, msgsd)
	if err != nil {
		fmt.Printf(err.Error())
		return nil, err
	}
	return stream, nil
}

//Connect to udp
func (sm *SessionsManager) ConnectToUdp(runId, remoteIp string, remotePort int) (*yamux.Stream, error) {
	stream, err := sm.GetStream(runId)
	if err != nil {
		return nil, err
	}
	msgsd := &models.ConnectUDP{
		TargetIP:   remoteIp,
		TargetPort: remotePort,
	}
	err = msg.WriteMsg(stream, msgsd)
	if err != nil {
		fmt.Printf(err.Error())
		return nil, err
	}
	return stream, nil
}

//Connect to Serial Port
func (sm *SessionsManager) ConnectToSerialPort(runId string, msgsd *models.ConnectSerialPort) (*yamux.Stream, error) {
	stream, err := sm.GetStream(runId)
	if err != nil {
		return nil, err
	}
	//msgsd := &models.ConnectSerialPort{
	//	PortName: "COM4",
	//	BaudRate: 115200,
	//	DataBits: 8,
	//	StopBits: 1,
	//	MinimumReadSize: 4,
	//}
	err = msg.WriteMsg(stream, msgsd)
	if err != nil {
		fmt.Printf(err.Error())
		return nil, err
	}
	return stream, nil
}

func (sm *SessionsManager) ConnectToTapTun(runId string) (*yamux.Stream, error) {
	stream, err := sm.GetStream(runId)
	if err != nil {
		return nil, err
	}
	msgsd := &models.NewService{
		Type: "tap",
	}
	err = msg.WriteMsg(stream, msgsd)
	if err != nil {
		fmt.Printf(err.Error())
		return nil, err
	}
	return stream, nil
}

func (sm *SessionsManager) ConnectToSSH(runId, remoteIP string, remotePort int, userName, passWord string) (stream *yamux.Stream, err error) {
	stream, err = sm.GetStream(runId)
	if err != nil {
		return nil, err
	}
	msgsd := &models.ConnectSSH{
		TargetIP:   remoteIP,
		TargetPort: remotePort,
		UserName:   userName,
		PassWord:   passWord,
	}
	err = msg.WriteMsg(stream, msgsd)
	if err != nil {
		fmt.Printf(err.Error())
		return nil, err
	}
	return stream, nil
}

func (sm *SessionsManager) ConnectToWs(runId, targetUrl, protocol, origin string) (*yamux.Stream, error) {
	stream, err := sm.GetStream(runId)
	if err != nil {
		return nil, err
	}
	msgsd := &models.ConnectWs{
		TargetUrl: targetUrl,
		Protocol:  protocol,
		Origin:    origin,
	}
	err = msg.WriteMsg(stream, msgsd)
	if err != nil {
		fmt.Printf(err.Error())
		return nil, err
	}
	return stream, nil
}

func (sm *SessionsManager) ListenMulticastUDP(runId, ip string, port uint) (*yamux.Stream, error) {
	//"224.0.0.50:9898"
	stream, err := sm.GetStream(runId)
	if err != nil {
		return nil, err
	}
	msgsd := &models.NewService{
		Type:   "ListenMulticastUDP",
		Config: fmt.Sprintf("%s:%d", ip, port),
	}
	err = msg.WriteMsg(stream, msgsd)
	if err != nil {
		fmt.Printf(err.Error())
		return nil, err
	}
	return stream, nil
}
