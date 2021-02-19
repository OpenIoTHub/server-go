package ip

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/OpenIoTHub/getip/iputils"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
)

//获取所有内网ip
func GetIntranetIp() string {
	intranetIps := ""
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		log.Println(err)
		return intranetIps
	}

	for _, address := range addrs {

		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				//log.Println("ip:", ipnet.IP.String())
				if intranetIps == "" {
					intranetIps = ipnet.IP.String()
				} else {
					intranetIps = intranetIps + "," + ipnet.IP.String()
				}
			}
		}
	}
	//fmt.Printf("所有内网ip：" + intranetIps)
	return intranetIps
}

//淘宝接口：获取ip信息
type IPInfo struct {
	Code int `json:"code"`
	Data IP  `json:"data"`
}

type IP struct {
	Ip        string `json:"ip"`
	Country   string `json:"country"`
	Area      string `json:"area"`
	Region    string `json:"region"`
	City      string `json:"city"`
	County    string `json:"county"`
	Isp       string `json:"isp"`
	CountryId string `json:"country_id"`
	AreaId    string `json:"area_id"`
	RegionId  string `json:"region_id"`
	CityId    string `json:"city_id"`
	CountyId  string `json:"county_id"`
	IspId     string `json:"isp_id"`
}

//获取自己的公网ip
func GetMyPublicIpInfo() (string, error) {
	return GetMyPublicIpv4()
}

func GetMyPublicIpv4() (string, error) {
	ip := iputils.GetMyPublicIpv4()
	if ip != "" {
		return ip, nil
	}
	return "", fmt.Errorf("获取ipv4地址失败")
}

func GetMyPublicIpv6() (string, error) {
	ip := iputils.GetMyPublicIpv6()
	if ip != "" {
		return ip, nil
	}
	return "", fmt.Errorf("获取ipv6地址失败")
}

func GetIpInfo(ip string) (*IPInfo, error) {
	if ip == "" {
		ip = "myip"
	}
	url := "http://ip.taobao.com/service/getIpInfo.php?ip="
	url += ip

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	out, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var result IPInfo
	if err := json.Unmarshal(out, &result); err != nil {
		return nil, err
	}

	return &result, err
}

//10.0.0.0/8：10.0.0.0～10.255.255.255
//172.16.0.0/12：172.16.0.0～172.31.255.255
//192.168.0.0/16：192.168.0.0～192.168.255.255
func IsPublicIP(IP net.IP) bool {
	if IP.IsLoopback() || IP.IsLinkLocalMulticast() || IP.IsLinkLocalUnicast() {
		return false
	}
	if ip4 := IP.To4(); ip4 != nil {
		switch true {
		case ip4[0] == 10:
			return false
		case ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31:
			return false
		case ip4[0] == 192 && ip4[1] == 168:
			return false
		default:
			return true
		}
	}
	return false
}

func IsChinaIP(IP net.IP) bool {
	data, err := Asset("chn_ip.txt")
	if err != nil {
		log.Println(err.Error())
		return true
	}
	r := bytes.NewReader(data)
	buf := bufio.NewReader(r)
	for {
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)
		if err != nil {
			return false
		}
		n := strings.Split(line, " ")
		ip1 := net.ParseIP(n[0])
		ip2 := net.ParseIP(n[1])
		if bytes.Compare(IP, ip1) >= 0 && bytes.Compare(IP, ip2) <= 0 {
			return true
		}
	}
}
