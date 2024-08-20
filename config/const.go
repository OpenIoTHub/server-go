package config

import "fmt"

const DefaultConfigFileName = "server-go.yaml"

var DefaultConfigFilePath = fmt.Sprintf("./%s", DefaultConfigFileName)

const DefaultBindAddr = "0.0.0.0"

const DefaultKcpPort = 34320
const DefaultTcpPort = 34320
const DefaultTlsPort = 34321
const DefaultUdpApiPort = 34321
const DefaultKcpApiPort = 34322
const DefaultGrpcPort = 34322

const DefaultLoginKey = "HLLdsa544&*S"

const DefaultHttpPort = 80
const DefaultHttpsPort = 443

const DefaultRedisNetwork = "tcp"
const DefaultRedisAddress = "127.0.0.1:6379"

const IoTManagerAddr = "api.iot-manager.iothub.cloud:50051"
