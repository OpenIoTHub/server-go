package session

import (
	"fmt"
	"net"
)

func RunTCP(port int) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		fmt.Println(err)
		return
	}
	listenerHdl(listener)
}
