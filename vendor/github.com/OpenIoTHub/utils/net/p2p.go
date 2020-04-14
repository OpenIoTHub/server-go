package nettool

import (
	"github.com/OpenIoTHub/utils/crypto"
	"github.com/OpenIoTHub/utils/models"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

//获取一个随机UDP Dial的内部ip，端口，外部ip端口
func GetDialIpPort(token *crypto.TokenClaims) (localAddr net.Addr, externalIp string, externalPort int, err error) {
	udpaddr, err := net.ResolveUDPAddr("udp", token.Host+":"+strconv.Itoa(token.P2PApiPort))
	//udpaddr, err := net.ResolveUDPAddr("udp", "tencent-shanghai-v1.host.nat-cloud.com:34321")
	if err != nil {
		return nil, "", 0, err
	}
	udpconn, err := net.DialUDP("udp", nil, udpaddr)
	defer udpconn.Close()
	if err != nil {
		log.Println(err.Error())
		return nil, "", 0, err
	}
	err = udpconn.SetDeadline(time.Now().Add(time.Duration(3 * time.Second)))
	if err != nil {
		return nil, "", 0, err
	}
	_, err = udpconn.Write([]byte("getIpPort"))
	if err != nil {
		return nil, "", 0, err
	}
	data := make([]byte, 256)
	n, err := udpconn.Read(data)
	if err != nil {
		return nil, "", 0, err
	}
	ipPort := string(data[:n])
	ip := strings.Split(ipPort, ":")[0]
	port, err := strconv.Atoi(strings.Split(ipPort, ":")[1])
	if err != nil {
		return nil, "", 0, err
	}
	//return strings.Split(udpconn.LocalAddr().String(), ":")[0]
	localAddr = udpconn.LocalAddr()
	return localAddr, ip, port, nil
}

func GetP2PListener(token *crypto.TokenClaims) (localIps string, localPort int, externalIp string, externalPort int, listener *net.UDPConn, err error) {
	localIps = GetIntranetIp()
	//localPort = randint.GenerateRangeNum(10000, 60000)
	localPort = 0
	listener, err = net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP("0.0.0.0"), Port: localPort})
	if err != nil {
		return
	}
	//获取监听的端口的外部ip和端口
	externalIp, externalPort, err = GetExternalIpPort(listener, token)
	if err != nil {
		log.Println(err)
		return
	}
	localPort = listener.LocalAddr().(*net.UDPAddr).Port
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
