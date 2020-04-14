package msg

import (
	"encoding/json"
	"fmt"
	"github.com/OpenIoTHub/utils/models"
	"reflect"
)

func unpack(typeByte string, buffer []byte, msgIn models.Message) (msg models.Message, err error) {
	if msgIn == nil {
		t, ok := models.TypeMap[typeByte]
		if !ok {
			err = fmt.Errorf("Unsupported message type %b", typeByte)
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

func UnPack(typeByte string, buffer []byte) (msg models.Message, err error) {
	return unpack(typeByte, buffer, nil)
}
