package nettool

import (
	"github.com/OpenIoTHub/utils/models"
	"log"
	"net"
	"time"
)

//获取一个随机UDP Listen的内部ip，端口，外部ip端口
func GetP2PListener(token *models.TokenClaims) (externalUDPAddr *net.UDPAddr, listener *net.UDPConn, err error) {
	listener, err = net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP("0.0.0.0"), Port: 0})
	if err != nil {
		return
	}
	//获取监听的端口的外部ip和端口
	externalUDPAddr, err = GetExternalIpPort(listener, token)
	if err != nil {
		log.Println(err)
		return
	}
	return
}

//client通过指定listener发送数据到explorer指定的p2p请求地址
func SendPackToPeer(listener *net.UDPConn, ctrlmMsg *models.ReqNewP2PCtrl) {
	log.Println("发送包到远程：", ctrlmMsg.ExternalIp, ctrlmMsg.ExternalPort)
	//发送5次防止丢包，稳妥点
	for i := 1; i <= 5; i++ {
		listener.WriteToUDP([]byte("packFromPeer"), &net.UDPAddr{
			IP:   net.ParseIP(ctrlmMsg.ExternalIp),
			Port: ctrlmMsg.ExternalPort,
		})
		time.Sleep(time.Millisecond * 100)
	}
}
