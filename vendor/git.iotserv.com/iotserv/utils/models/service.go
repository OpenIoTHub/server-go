package models

import "github.com/iotdevice/zeroconf"

type FindmDNS struct {
	Service string
	Domain  string
	Second  int
}

type MDNSResult []*zeroconf.ServiceEntry

type ScanPort struct {
	Host      string
	StartPort int
	EndPort   int
}

type ScanPortResult []int
