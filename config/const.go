package config

//const RegisterHost = "tencent-shanghai-v1.host.nat-cloud.com"

//const RegisterHost = "127.0.0.1"
//const RegisterHost = "167.179.92.55"
const RegisterHost = "s1.365hour.com"

//const RegisterHost = "iotserv.lu8.win"
//const RegisterHost = "nat-cloud.cn"
//const RegisterHost = "netipcam.com"

//转发注册服务器需要的端口
const TcpPort = 34320
const KcpPort = 34320
const UdpApiPort = 34321
const TlsPort = 34321

const ConfigFilePath = "./config.yaml"

const TlsCertFilePath = "./cert.pem"
const TlsKeyFilePath = "./key.pem"

const HttpsCertFilePath = ""
const HttpsKeyFilePath = ""

//加密使用的种子
const DefaultConnectToServerKey = "HLLdsa544&*S"
const DefaultExplorerKey = "Abc&&*DDhhA"
