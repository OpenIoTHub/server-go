package msg

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/OpenIoTHub/utils/errors"
	"github.com/OpenIoTHub/utils/models"
)

func unpack(typeByte byte, buffer []byte, msgIn models.Message) (msg models.Message, err error) {
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
	_, err = unpack(' ', buffer, msg)
	return
}

func UnPack(typeByte byte, buffer []byte) (msg models.Message, err error) {
	return unpack(typeByte, buffer, nil)
}

func Pack(msg models.Message) ([]byte, error) {
	typeByte, ok := models.TypeStringMap[reflect.TypeOf(msg).Elem()]
	if !ok {
		return nil, errors.ErrMsgType
	}

	content, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	buffer := bytes.NewBuffer(nil)
	buffer.WriteByte(typeByte)
	binary.Write(buffer, binary.BigEndian, int64(len(content)))
	buffer.Write(content)
	return buffer.Bytes(), nil
}
