package errors

import (
	"errors"
	"fmt"
)

var (
	ErrMsgType = errors.New("message type error")
)

func PanicToError(fn func()) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("Panic error: %v", r)
		}
	}()

	fn()
	return
}
