package internal

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/s-larionov/process-manager"

	"github.com/goverland-labs/helpers-ens-resolver/internal/config"
	"github.com/goverland-labs/helpers-ens-resolver/internal/infura"
	"github.com/goverland-labs/helpers-ens-resolver/internal/server"
	"github.com/goverland-labs/helpers-ens-resolver/pkg/grpcsrv"
	"github.com/goverland-labs/helpers-ens-resolver/pkg/health"
	"github.com/goverland-labs/helpers-ens-resolver/pkg/prometheus"
	"github.com/goverland-labs/helpers-ens-resolver/proto"
)

type Application struct {
	sigChan <-chan os.Signal
	manager *process.Manager
	cfg     config.App

	infura *infura.Client
}

func NewApplication(cfg config.App) (*Application, error) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	a := &Application{
		sigChan: sigChan,
		cfg:     cfg,
		manager: process.NewManager(),
	}

	err := a.bootstrap()
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *Application) Run() {
	a.manager.StartAll()
	a.registerShutdown()
}

func (a *Application) bootstrap() error {
	initializers := []func() error{
		// Init dependencies
		a.initServices,

		// Init Workers: Application
		a.initGRPCWorker,

		// Init Workers: System
		a.initHealthWorker,
		a.initPrometheusWorker,
	}

	for _, initializer := range initializers {
		if err := initializer(); err != nil {
			return err
		}
	}

	return nil
}

func (a *Application) initServices() error {
	inf, err := infura.NewClient(a.cfg.Infura)
	if err != nil {
		return err
	}

	a.infura = inf

	return nil
}

func (a *Application) initGRPCWorker() error {
	srv := grpcsrv.NewGrpcServer(func(ctx context.Context) (context.Context, error) {
		return ctx, nil
	})

	proto.RegisterEnsServer(srv, server.NewEnsHandler(a.infura))
	a.manager.AddWorker(grpcsrv.NewGrpcServerWorker("resolve api", srv, a.cfg.GRPC.Listen))

	return nil
}

func (a *Application) initPrometheusWorker() error {
	srv := prometheus.NewPrometheusServer(a.cfg.Prometheus.Listen, "/metrics")
	a.manager.AddWorker(process.NewServerWorker("prometheus", srv))

	return nil
}

func (a *Application) initHealthWorker() error {
	srv := health.NewHealthCheckServer(a.cfg.Health.Listen, "/status", health.DefaultHandler(a.manager))
	a.manager.AddWorker(process.NewServerWorker("health", srv))

	return nil
}

func (a *Application) registerShutdown() {
	go func(manager *process.Manager) {
		<-a.sigChan

		manager.StopAll()
	}(a.manager)

	a.manager.AwaitAll()
}
