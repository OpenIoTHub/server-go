package session

import (
	"fmt"
	"github.com/OpenIoTHub/grpc-api/pb-go"
	"github.com/OpenIoTHub/server-go/config"
	"google.golang.org/grpc"
	"log"
	"net"
)

func init() {
	SessionsCtl.StartgRpcListenAndServ()
}

func (sm *SessionsManager) StartgRpcListenAndServ() {
	go func() {
		s := grpc.NewServer()
		pb.RegisterHttpManagerServer(s, sm)
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", config.DefaultGrpcPort))
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
			return
		}
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
}
