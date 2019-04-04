package config

type ClientConfig struct {
	Common struct {
		Id           string `yaml:"id"`
		RegisterHost string `yaml:"register_host"`
	}
	LastToken struct {
		ClientToken   string `yaml:"client_token"`
		ExplorerToken string `yaml:"explorer_token"`
	}
}

type RegisterConfig struct {
}

type ExplorerConfig struct {
}
