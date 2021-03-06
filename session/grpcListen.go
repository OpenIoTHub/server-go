package session

import (
	"fmt"
	"github.com/OpenIoTHub/server-go/config"
	"github.com/OpenIoTHub/server-grpc-api/pb-go"
	"google.golang.org/grpc"
	"log"
	"net"
)

func (sm *SessionsManager) StartgRpcListenAndServ() {
	go func() {
		s := grpc.NewServer()
		pb.RegisterHttpManagerServer(s, sm)
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", config.ConfigMode.Common.GrpcPort))
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
			return
		}
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
}
