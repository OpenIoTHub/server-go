package session

import (
	"context"
	"errors"
	"fmt"
	"github.com/OpenIoTHub/server-go/config"
	"github.com/OpenIoTHub/utils/models"
	"github.com/OpenIoTHub/utils/msg"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log"
)

//检查远端内网端口的可用性，可用true
func (sm *SessionsManager) CheckRemoteStatus(targetType, runId, remoteIp string, remotePort int) (bool, error) {
	stream, err := sm.GetStream(runId)
	defer func() {
		if stream != nil {
			err := stream.Close()
			if err != nil {
				log.Println(err.Error())
			}
		}
	}()
	if err != nil {
		fmt.Printf("get stream err: " + err.Error())
		return false, err
	}
	msgsd := &models.CheckStatusRequest{
		Type: targetType,
		Addr: fmt.Sprintf("%s:%d", remoteIp, remotePort),
	}
	err = msg.WriteMsg(stream, msgsd)
	if err != nil {
		fmt.Printf(err.Error())
		return false, err
	}
	//:TODO 可能不会及时返还
	rawMsg, err := msg.ReadMsg(stream)
	if err != nil {
		return false, err
	}
	switch m := rawMsg.(type) {
	case *models.CheckStatusResponse:
		{
			if m.Code == 0 {
				return true, nil
			}
			return false, errors.New(m.Message)
		}
	default:
		break
	}
	return false, nil
}

func checkOpenIoTHubToken(key, tokenStr, id string) (token *models.TokenClaims, err error) {
	token, err = models.DecodeToken(key, tokenStr)
	if err != nil {
		log.Println(err.Error())
		return
	}
	if token.Permission&1 != 1 {
		log.Println("token type err ,not n")
		return nil, errors.New("not gateway token")
	}
	if token.Id != id {
		log.Println("token type err ,not n")
		return nil, errors.New("id check error")
	}
	return
}

func authOpenIoTHubGrpc(ctx context.Context, id string) (err error) {
	var jwt string
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Errorf(codes.Unauthenticated, "metadata.FromIncomingContext err")
	}
	if jwts, ok := md["jwt"]; ok {
		jwt = jwts[0]
	} else {
		return status.Errorf(codes.Unauthenticated, "jwt is empty")
	}

	_, err = checkOpenIoTHubToken(config.ConfigMode.Security.LoginKey, jwt, id)
	if err != nil {
		return status.Errorf(codes.Unauthenticated, err.Error())
	}
	return nil
}
