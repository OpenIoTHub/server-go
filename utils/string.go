package utils

import (
	"github.com/satori/go.uuid"
)

func GetOneUUID() string {
	uid, _ := uuid.NewV4()
	return uid.String()
}
