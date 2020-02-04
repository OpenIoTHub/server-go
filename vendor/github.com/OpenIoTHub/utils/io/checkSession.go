package io

import (
	"fmt"
	"github.com/OpenIoTHub/utils/models"
	"github.com/OpenIoTHub/utils/msg"
	"github.com/OpenIoTHub/utils/mux"
	"github.com/pkg/errors"
	"time"
)

func CheckSession(muxSession *mux.Session) error {
	//:TODO 当这个内网被删除时退出重连
	readPong := func(stream *mux.Stream) error {
		defer stream.Close()
		var ch = make(chan struct{}, 1)
		go func() {
			_, err := msg.ReadMsg(stream)
			if err != nil {
				fmt.Println(err)
				return
			}
			ch <- struct{}{}
		}()
		select {
		case <-ch:
			return nil
		case <-time.After(time.Second * 3):
			if muxSession != nil && !muxSession.IsClosed() {
				muxSession.Close()
			}
			return errors.New("Session Check TimeOut")
		}
	}
	if muxSession == nil {
		return errors.New("Session is nil")
	}
	if muxSession.IsClosed() {
		return errors.New("Session.IsClosed")
	}
	stream, err := muxSession.OpenStream()
	if err != nil {
		return err
	}
	err = msg.WriteMsg(stream, &models.Ping{})
	if err != nil {
		return err
	}
	return readPong(stream)
}
