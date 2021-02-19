package msg

import (
	"encoding/json"
	"fmt"
	"github.com/OpenIoTHub/utils/models"
	"reflect"
)

func unpack(typeStr string, buffer []byte, msgIn models.Message) (msg models.Message, err error) {
	if msgIn == nil {
		t, ok := models.TypeMap[typeStr]
		if !ok {
			err = fmt.Errorf("unsupported message type %s", typeStr)
			return
		}

		msg = reflect.New(t).Interface().(models.Message)
	} else {
		msg = msgIn
	}

	err = json.Unmarshal(buffer, &msg)
	return
}

func UnPackInto(buffer []byte, msg models.Message) (err error) {
	_, err = unpack("", buffer, msg)
	return
}

func UnPack(typeStr string, buffer []byte) (msg models.Message, err error) {
	return unpack(typeStr, buffer, nil)
}
