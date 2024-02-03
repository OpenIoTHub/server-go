package manager

import (
	"context"
	pb "github.com/OpenIoTHub/openiothub_grpc_api/pb-go/proto/manager"
	"github.com/OpenIoTHub/server-go/config"
	"github.com/OpenIoTHub/utils/models"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
)

const IoTManagerAddr = "iot-manager.iotserv.com:8881"

//const IoTManagerAddr = "127.0.0.1:8881"

func LoadConfigFromIoTManager() (err error) {
	conn, err := grpc.Dial(IoTManagerAddr, grpc.WithInsecure())
	if err != nil {
		log.Println("grpc.Dial:", err)
		return
	}
	defer conn.Close()
	c := pb.NewPortManagerClient(conn)
	jwt, err := models.GetUuidToken(config.ConfigMode.Security.LoginKey, config.ConfigMode.ServerUuid, "server-go", []string{}, map[string]string{}, 100)
	if err != nil {
		log.Println(err)
		return
	}
	//metadata传递jwt
	md := metadata.Pairs("jwt", jwt)
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	rst, err := c.GetAllHttpInfoListByServerUuid(ctx, &emptypb.Empty{})
	if err != nil {
		log.Println(err)
		return
	}
	for _, info := range rst.HttpInfoList {
		log.Println("add domain:", info)
		SessionsCtl.AddOrUpdateHttpProxy(&HttpProxy{
			Domain:           info.Domain,
			UserName:         info.Username,
			Password:         info.Password,
			RunId:            info.GatewayUuid,
			RemoteIP:         info.RemoteAddr,
			RemotePort:       int(info.RemotePort),
			IfHttps:          info.ApplicationProtocol == "https",
			Description:      info.Description,
			RemotePortStatus: false,
		})
	}
	log.Println("LoadConfigFromIoTManager:OK!")
	return
}
