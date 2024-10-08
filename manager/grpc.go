package manager

import (
	"context"
	"fmt"
	pb "github.com/OpenIoTHub/openiothub_grpc_api/pb-go/proto/server"
	"github.com/OpenIoTHub/server-go/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"net"
)

// grpc
func (sm *SessionsManager) StartgRpcListenAndServ() {
	go func() {
		//TODO 如果是IP则使用不安全的连接，如果是域名则使用安全的连接
		var s *grpc.Server
		if net.ParseIP(config.ConfigMode.PublicIp) == nil {
			s = grpc.NewServer(grpc.Creds(credentials.NewTLS(autocertManager.TLSConfig())))
		} else {
			s = grpc.NewServer()
		}
		pb.RegisterHttpManagerServer(s, sm)
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", config.ConfigMode.Common.GrpcPort))
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
			return
		}
		reflection.Register(s)
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
}

func (sm *SessionsManager) CreateOneHTTP(ctx context.Context, in *pb.HTTPConfig) (*pb.HTTPConfig, error) {
	err := authOpenIoTHubGrpc(ctx, in.RunId)
	if err != nil {
		return in, status.Errorf(codes.Unauthenticated, err.Error())
	}
	log.Println("CreateOneHTTP:", in.Domain)
	return in, sm.AddOrUpdateHttpProxy(&HttpProxy{
		Domain:      in.Domain,
		RunId:       in.RunId,
		RemoteIP:    in.RemoteIP,
		RemotePort:  int(in.RemotePort),
		UserName:    in.UserName,
		Password:    in.Password,
		IfHttps:     in.IfHttps,
		Description: in.Description,
	})
}

func (sm *SessionsManager) UpdateOneHTTP(ctx context.Context, in *pb.HTTPConfig) (*pb.HTTPConfig, error) {
	err := authOpenIoTHubGrpc(ctx, in.RunId)
	if err != nil {
		return in, status.Errorf(codes.Unauthenticated, err.Error())
	}
	log.Println("UpdateOneHTTP:", in.Domain)
	h := &HttpProxy{
		Domain:      in.Domain,
		RunId:       in.RunId,
		RemoteIP:    in.RemoteIP,
		RemotePort:  int(in.RemotePort),
		UserName:    in.UserName,
		Password:    in.Password,
		IfHttps:     in.IfHttps,
		Description: in.Description,
	}
	sm.DelHttpProxy(h.Domain)
	return in, sm.AddOrUpdateHttpProxy(h)
}

func (sm *SessionsManager) DeleteOneHTTP(ctx context.Context, in *pb.HTTPConfig) (*emptypb.Empty, error) {
	err := authOpenIoTHubGrpc(ctx, in.RunId)
	if err != nil {
		return &emptypb.Empty{}, status.Errorf(codes.Unauthenticated, err.Error())
	}
	log.Println("DeleteOneHTTP:", in.Domain)
	//TODO 验证要删除的域名的所属id是否和token的id一致
	sm.DelHttpProxy(in.Domain)
	return &emptypb.Empty{}, nil

}

func (sm *SessionsManager) GetOneHTTP(ctx context.Context, in *pb.HTTPConfig) (*pb.HTTPConfig, error) {
	err := authOpenIoTHubGrpc(ctx, in.RunId)
	if err != nil {
		return in, status.Errorf(codes.Unauthenticated, err.Error())
	}
	config, err := sm.GetOneHttpProxy(in.Domain)
	if err != nil {
		return &pb.HTTPConfig{}, err
	}
	return &pb.HTTPConfig{
		Domain:           config.Domain,
		RunId:            config.RunId,
		RemoteIP:         config.RemoteIP,
		RemotePort:       int32(config.RemotePort),
		UserName:         config.UserName,
		Password:         config.Password,
		IfHttps:          config.IfHttps,
		Description:      config.Description,
		RemotePortStatus: config.RemotePortStatus,
	}, err
}

func (sm *SessionsManager) GetAllHTTP(ctx context.Context, in *pb.Device) (*pb.HTTPList, error) {
	var cfgs []*pb.HTTPConfig
	err := authOpenIoTHubGrpc(ctx, in.RunId)
	if err != nil {
		return &pb.HTTPList{HTTPConfigs: cfgs}, status.Errorf(codes.Unauthenticated, err.Error())
	}
	for _, config := range sm.GetAllHttpProxy() {
		if config.RunId == in.RunId && config.RemoteIP == in.Addr {
			cfgs = append(cfgs, &pb.HTTPConfig{
				Domain:           config.Domain,
				RunId:            config.RunId,
				RemoteIP:         config.RemoteIP,
				RemotePort:       int32(config.RemotePort),
				UserName:         config.UserName,
				Password:         config.Password,
				IfHttps:          config.IfHttps,
				Description:      config.Description,
				RemotePortStatus: config.RemotePortStatus,
			})
		}
	}
	return &pb.HTTPList{HTTPConfigs: cfgs}, nil
}

//grpc end
