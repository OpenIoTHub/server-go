package session

import (
	"fmt"
	"log"
	"net"
)

func RunTCP(port int) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Println(err)
		return
	}
	listenerHdl(listener)
}
