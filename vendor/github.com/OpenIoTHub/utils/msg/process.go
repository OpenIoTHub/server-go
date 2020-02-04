package msg

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/OpenIoTHub/utils/models"
	"io"
	"time"
)

var (
	MaxMsgLength int64 = 10240
)

func readMsg(c io.Reader) (typeByte byte, buffer []byte, err error) {
	buffer = make([]byte, 1)
	if c == nil {
		return byte('a'), nil, fmt.Errorf("conn is nil")
	}
	_, err = c.Read(buffer)
	if err != nil {
		return
	}
	typeByte = buffer[0]
	if _, ok := models.TypeMap[typeByte]; !ok {
		err = fmt.Errorf("Message type error")
		return
	}

	var length int64
	err = binary.Read(c, binary.BigEndian, &length)
	if err != nil {
		return
	}
	if length > MaxMsgLength {
		err = fmt.Errorf("Message length exceed the limit")
		return
	}

	buffer = make([]byte, length)
	n, err := io.ReadFull(c, buffer)
	if err != nil {
		return
	}

	if int64(n) != length {
		err = fmt.Errorf("Message format error")
	}
	return
}

func ReadMsg(c io.Reader) (msg models.Message, err error) {
	typeByte, buffer, err := readMsg(c)
	if err != nil {
		return
	}
	return UnPack(typeByte, buffer)
}

//读取Msg超时错误返回
func ReadMsgWithTimeOut(c io.Reader, t time.Duration) (msg models.Message, err error) {
	var typeByte byte
	var buffer []byte
	var ch = make(chan struct{}, 1)
	go func() {
		typeByte, buffer, err = readMsg(c)
		if err != nil {
			return
		}
		ch <- struct{}{}
	}()
	select {
	case <-ch:
		return UnPack(typeByte, buffer)
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
	buffer, err := Pack(msg)
	if err != nil {
		return
	}
	if c == nil {
		return fmt.Errorf("写入消息的连接为nil")
	}
	if _, err = c.Write(buffer); err != nil {
		return
	}

	return nil
}
