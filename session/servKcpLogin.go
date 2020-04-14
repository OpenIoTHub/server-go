package session

import (
	"fmt"
	"github.com/xtaci/kcp-go"
	"log"
)

func RunKCP(port int) {
	listener, err := kcp.ListenWithOptions(fmt.Sprintf(":%d", port), nil, 10, 3)
	if err != nil {
		log.Println(err)
		return
	}
	kcpListenerHdl(listener)
}
