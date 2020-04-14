package nettool

import (
	"net"
)

func GetMacs() ([]string, error) {
	// 获取本机的MAC地址
	var macs []string
	interfaces, err := net.Interfaces()
	if err != nil {
		return macs, err
	}
	for _, inter := range interfaces {
		mac := inter.HardwareAddr //获取本机MAC地址
		macs = append(macs, mac.String())
	}
	return macs, nil
}
