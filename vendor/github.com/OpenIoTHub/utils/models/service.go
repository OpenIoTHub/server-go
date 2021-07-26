package models

import "net"

type FindmDNS struct {
	Service string
	Domain  string
	Second  int
}

type MDNSResult struct {
	Instance string   `json:"name"`
	Service  string   `json:"type"`
	Domain   string   `json:"domain"`
	HostName string   `json:"hostname"`
	Port     int      `json:"port"`
	Text     []string `json:"text"`
	TTL      uint32   `json:"ttl"`
	AddrIPv4 []net.IP `json:"addripv4"`
	AddrIPv6 []net.IP `json:"addripv6"`
}

type ScanPort struct {
	Host      string
	StartPort int
	EndPort   int
}

type ScanPortResult []int
