package main

import (
	"github.com/caarlos0/env/v8"
	"github.com/rs/zerolog"
	"github.com/s-larionov/process-manager"

	"helpers-ens-resolver/internal"
	"helpers-ens-resolver/internal/config"
)

var cfg config.App

func init() {
	err := env.Parse(&cfg)
	if err != nil {
		panic(err)
	}

	level, err := zerolog.ParseLevel(cfg.LogLevel)
	if err != nil {
		panic(err)
	}
	zerolog.SetGlobalLevel(level)
	process.SetLogger(&PMLogger{})
}

func main() {
	application, err := internal.NewApplication(cfg)
	if err != nil {
		panic(err)
	}

	application.Run()
}
