package config

// Model definition for configuration

// Config the application's configuration
type Config struct {
	Config          string
	DB              DBConfig
	Port            string
	GRPCPort        string
	SomeWeirdConfig string `json:"some-weird-config" yaml:"SomeWeirdConfig"`
	SomeAddress     string `json:"some-address" yaml:"ProxyPass"`
	LogFile         string
	Env             string
	Verbose         bool
}

type DBConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Name     string `json:"name"`
}
