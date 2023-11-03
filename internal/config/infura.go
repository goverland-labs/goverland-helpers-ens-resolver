package config

type Infura struct {
	Key      string `env:"INFURA_API_KEY,required"`
	Endpoint string `env:"INFURA_ENDPOINT" envDefault:"https://mainnet.infura.io/v3/"`
}
