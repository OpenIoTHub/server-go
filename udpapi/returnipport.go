package udpapi

import (
	"fmt"
	"log"
	"net"
)

func RunApiServer(port int) {
	listener, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP("0.0.0.0"), Port: port})
	if err != nil {
		log.Println(err)
		return
	}
	go udpListener(listener)
}

func udpListener(listener *net.UDPConn) {
	data := make([]byte, 256)
	for {
		_, remoteAddr, err := listener.ReadFromUDP(data)
		if err != nil {
			fmt.Printf("error during read: %s", err)
		}
		//fmt.Printf("<%s> %s\n", remoteAddr, data[:n])
		//:TODO 防阻塞
		go func() {
			_, err = listener.WriteToUDP([]byte(remoteAddr.String()), remoteAddr)
			if err != nil {
				fmt.Printf(err.Error())
			}
		}()
	}
}
