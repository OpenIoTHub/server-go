package msg

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/OpenIoTHub/utils/models"
	"github.com/libp2p/go-msgio"
	"io"
	"reflect"
	"time"
)

func readMsg(c io.Reader) (typeString string, buffer []byte, err error) {
	if c == nil {
		return "", nil, fmt.Errorf("conn is nil")
	}
	mrdr := msgio.NewReader(c)
	typeByte, err := mrdr.ReadMsg()
	if err != nil || typeByte == nil {
		return "", nil, err
	}
	typeString = string(typeByte)
	if _, ok := models.TypeMap[typeString]; !ok {
		err = fmt.Errorf("Message type error")
		return
	}
	buffer, err = mrdr.ReadMsg()
	return
}

func ReadMsg(c io.Reader) (msg models.Message, err error) {
	typeString, buffer, err := readMsg(c)
	if err != nil {
		return
	}
	return UnPack(typeString, buffer)
}

//读取Msg超时错误返回
func ReadMsgWithTimeOut(c io.Reader, t time.Duration) (msg models.Message, err error) {
	var typeString string
	var buffer []byte
	var ch = make(chan struct{}, 1)
	go func() {
		typeString, buffer, err = readMsg(c)
		if err != nil {
			return
		}
		ch <- struct{}{}
	}()
	select {
	case <-ch:
		return UnPack(typeString, buffer)
	case <-time.After(time.Second * 3):
		return nil, errors.New("Read Msg TimeOut")
	}
}

func ReadMsgInto(c io.Reader, msg models.Message) (err error) {
	_, buffer, err := readMsg(c)
	if err != nil {
		return
	}
	return UnPackInto(buffer, msg)
}

func WriteMsg(c io.Writer, msg interface{}) (err error) {
	if c == nil {
		return fmt.Errorf("写入消息的连接为nil")
	}
	typeString, ok := models.TypeStringMap[reflect.TypeOf(msg).Elem()]
	if !ok {
		return errors.New("message type not found")
	}
	mwtr := msgio.NewWriter(c)
	err = mwtr.WriteMsg([]byte(typeString))
	if err != nil {
		return err
	}

	content, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	err = mwtr.WriteMsg(content)
	return
}
