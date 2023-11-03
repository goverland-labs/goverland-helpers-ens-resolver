package config

type GRPC struct {
	Listen string `env:"GRPC_LISTEN" envDefault:":20200"`
}
