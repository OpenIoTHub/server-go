package iputils

import (
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
)

var Ipv4APIUrls = []string{
	"http://members.3322.org/dyndns/getip",
	"http://ifconfig.me/ip", "http://ip.3322.net",
	"https://myexternalip.com/raw",
	"http://ipv4.ident.me",
	"http://ipv4.icanhazip.com",
	"http://nsupdate.info/myip",
	"http://whatismyip.akamai.com",
	"http://ipv4.myip.dk/api/info/IPv4Address",
	"http://checkip4.spdyn.de",
	"http://v4.ipv6-test.com/api/myip.php",
	"http://checkip.amazonaws.com",
	"http://ipinfo.io/ip",
	"http://bot.whatismyipaddress.com",
	"http://ipv4.ident.me",
	"http://ipv4.icanhazip.com",
	"http://nsupdate.info/myip",
	"http://whatismyip.akamai.com",
	"http://ipv4.myip.dk/api/info/IPv4Address",
	"http://checkip4.spdyn.de",
	"http://v4.ipv6-test.com/api/myip.php",
	"http://checkip.amazonaws.com",
	"http://ipinfo.io/ip http://bot.whatismyipaddress.com",
}
var Ipv6APIUrls = []string{
	"http://bbs6.ustc.edu.cn/cgi-bin/myip",
	"http://ipv6.ident.me",
	"http://ipv6.icanhazip.com",
	"http://ipv6.ident.me",
	"http://ipv6.icanhazip.com",
	"http://ipv6.yunohost.org",
	"http://v6.ipv6-test.com/api/myip.php",
	"http://ipv6.ident.me",
	"http://ipv6.icanhazip.com",
	"http://ipv6.yunohost.org",
	"http://v6.ipv6-test.com/api/myip.php",
	"http://ipv6.ident.me",
	"http://ipv6.icanhazip.com",
	"http://ipv6.ident.me",
	"http://ipv6.icanhazip.com",
	"http://ipv6.yunohost.org",
	"http://v6.ipv6-test.com/api/myip.php",
	"http://ipv6.ident.me",
	"http://ipv6.icanhazip.com",
	"http://ipv6.yunohost.org",
	"http://v6.ipv6-test.com/api/myip.php",
}

func GetMyPublicIpv4() string {
	for _, url := range Ipv4APIUrls {
		resp, err := http.Get(url)
		if err != nil {
			log.Printf("get public ipv4 err：%s", err)
			continue
		}
		bytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("get public ipv4 err：%s", err)
			_ = resp.Body.Close()
			continue
		}
		ipv4 := strings.Replace(string(bytes), "\n", "", -1)
		ip := net.ParseIP(ipv4)
		if ip != nil {
			_ = resp.Body.Close()
			return ip.String()
		}
	}
	return ""
}

func GetMyPublicIpv6() string {
	for _, url := range Ipv6APIUrls {
		resp, err := http.Get(url)
		if err != nil {
			log.Printf("get public ipv6 err：%s", err)
			continue
		}
		// 读取 IPv6
		bytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("get public ipv6 err：%s", err)
			_ = resp.Body.Close()
			continue
		}
		// 删除 document.write(xxx) (如有)
		tmp := strings.Replace(string(bytes), "document.write('", "", -1)
		tmp = strings.Replace(tmp, "');", "", -1)
		ipv6 := strings.Replace(tmp, "\n", "", -1)
		ip := net.ParseIP(ipv6)
		if ip != nil {
			_ = resp.Body.Close()
			return ip.String()
		}
	}
	return ""
}
