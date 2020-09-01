package models

type GatewayConfig struct {
	ConnectionType string `yaml:"connection_type"`
	LastId         string `yaml:"last_id"`
	Server         *Srever
}

type Srever struct {
	ServerHost string `yaml:"server_host"`
	TcpPort    int    `yaml:"tcp_port"`
	KcpPort    int    `yaml:"kcp_port"`
	UdpApiPort int    `yaml:"udp_p2p_port"`
	KcpApiPort int    `yaml:"kcp_p2p_port"`
	TlsPort    int    `yaml:"tls_port"`
	GrpcPort   int    `yaml:"grpc_port"`
	LoginKey   string `yaml:"login_key"`
}

type ServerConfig struct {
	Common struct {
		BindAddr   string `yaml:"bind_addr"`
		TcpPort    int    `yaml:"tcp_port"`
		KcpPort    int    `yaml:"kcp_port"`
		UdpApiPort int    `yaml:"udp_p2p_port"`
		KcpApiPort int    `yaml:"kcp_p2p_port"`
		TlsPort    int    `yaml:"tls_port"`
		GrpcPort   int    `yaml:"grpc_port"`
		HttpPort   int    `yaml:"http_port"`
		HttpsPort  int    `yaml:"https_port"`
	}
	Security struct {
		LoginKey          string `yaml:"login_key"`
		TlsCertFilePath   string `yaml:"tls_Cert_file_path"`
		TlsKeyFilePath    string `yaml:"tls_key_file_path"`
		HttpsCertFilePath string `yaml:"https_cert_file_path"`
		HttpsKeyFilePath  string `yaml:"https_key_file_path"`
	}
	RedisConfig struct {
		Network  string `yaml:"network"`
		Address  string `yaml:"address"`
		Database int    `yaml:"database"`
		NeedAuth bool   `yaml:"needAuth"`
		Password string `yaml:"password"`
	}
}
