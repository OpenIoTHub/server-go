package nettool

import "net"

//获取一个随机空闲的tcp端口
func GetOneFreeTcpPort() (int, error) {
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}
