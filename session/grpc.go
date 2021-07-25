package session

import (
	"context"
	"github.com/OpenIoTHub/server-grpc-api/pb-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//grpc
func (sm *SessionsManager) CreateOneHTTP(ctx context.Context, in *pb.HTTPConfig) (*pb.HTTPConfig, error) {
	err := authOpenIoTHubGrpc(ctx, in.RunId)
	if err != nil {
		return in, status.Errorf(codes.Unauthenticated, err.Error())
	}

	return in, sm.AddHttpProxy(&HttpProxy{
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
	return in, sm.AddHttpProxy(h)
}

func (sm *SessionsManager) DeleteOneHTTP(ctx context.Context, in *pb.HTTPConfig) (*pb.Empty, error) {
	err := authOpenIoTHubGrpc(ctx, in.RunId)
	if err != nil {
		return &pb.Empty{}, status.Errorf(codes.Unauthenticated, err.Error())
	}
	//TODO 验证要删除的域名的所属id是否和token的id一致
	sm.DelHttpProxy(in.Domain)
	return &pb.Empty{}, nil

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
