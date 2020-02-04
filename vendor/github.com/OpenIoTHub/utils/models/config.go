package models

type ClientConfig struct {
	ExplorerTokenHttpPort int `yaml:"explorer_token_http_port"`
	Server                Srever
	LastId                string `yaml:"last_id"`
}

type Srever struct {
	ConnectionType string `yaml:"connection_type"`
	ServerHost     string `yaml:"server_host"`
	TcpPort        int    `yaml:"tcp_port"`
	KcpPort        int    `yaml:"kcp_port"`
	UdpApiPort     int    `yaml:"udp_p2p_port"`
	TlsPort        int    `yaml:"tls_port"`
	LoginKey       string `yaml:"login_key"`
}

type ServerConfig struct {
	Common struct {
		BindAddr   string `yaml:"bind_addr"`
		TcpPort    int    `yaml:"tcp_port"`
		KcpPort    int    `yaml:"kcp_port"`
		UdpApiPort int    `yaml:"udp_p2p_port"`
		TlsPort    int    `yaml:"tls_port"`
	}
	Security struct {
		LoginKey          string `yaml:"login_key"`
		TlsCertFilePath   string `yaml:"tls_Cert_file_path"`
		TlsKeyFilePath    string `yaml:"tls_key_file_path"`
		HttpsCertFilePath string `yaml:"https_cert_file_path"`
		HttpsKeyFilePath  string `yaml:"https_key_file_path"`
	}
}

type ClientFlat struct {
	ExplorerTokenHttpPort int    `yaml:"explorer_token_http_port"`
	LastId                string `yaml:"last_id"`
	ConnectionType        string `yaml:"connection_type"`
	ServerHost            string `yaml:"server_host"`
	TcpPort               string `yaml:"tcp_port"`
	KcpPort               string `yaml:"kcp_port"`
	UdpApiPort            string `yaml:"udp_p2p_port"`
	TlsPort               string `yaml:"tls_port"`
	LoginKey              string `yaml:"login_key"`
}
