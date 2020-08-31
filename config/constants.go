package config

import "fmt"

var DefaultConfigFileName = "server.yaml"
var DefaultConfigFilePath = fmt.Sprintf("./%s", DefaultConfigFileName)

var DefaultBindAddr = "0.0.0.0"

var DefaultKcpPort = 34320
var DefaultTcpPort = 34320
var DefaultTlsPort = 34321
var DefaultUdpApiPort = 34321
var DefaultKcpApiPort = 34322
var DefaultGrpcPort = 34322

var DefaultLoginKey = "HLLdsa544&*S"

var DefaultHttpPort = 80
var DefaultHttpsPort = 443
