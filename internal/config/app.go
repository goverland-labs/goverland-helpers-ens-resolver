package config

type App struct {
	LogLevel   string `env:"LOG_LEVEL" envDefault:"info"`
	Prometheus Prometheus
	Health     Health
	Infura     Infura
	Stamp      Stamp
	GRPC       GRPC
}
