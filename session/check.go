package session

import (
	"errors"
	"fmt"
	"github.com/OpenIoTHub/utils/models"
	"github.com/OpenIoTHub/utils/msg"
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
