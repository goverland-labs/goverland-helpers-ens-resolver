package internal

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/s-larionov/process-manager"

	"github.com/goverland-labs/goverland-helpers-ens-resolver/protocol/enspb"

	"github.com/goverland-labs/goverland-helpers-ens-resolver/internal/config"
	"github.com/goverland-labs/goverland-helpers-ens-resolver/internal/metrics"
	"github.com/goverland-labs/goverland-helpers-ens-resolver/internal/server"
	"github.com/goverland-labs/goverland-helpers-ens-resolver/pkg/grpcsrv"
	"github.com/goverland-labs/goverland-helpers-ens-resolver/pkg/health"
	"github.com/goverland-labs/goverland-helpers-ens-resolver/pkg/prometheus"
	"github.com/goverland-labs/goverland-helpers-ens-resolver/pkg/sdk/alchemy"
	"github.com/goverland-labs/goverland-helpers-ens-resolver/pkg/sdk/stamp"
)

type Application struct {
	sigChan <-chan os.Signal
	manager *process.Manager
	cfg     config.App

	stamp   *stamp.Client
	alchemy *alchemy.Client
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
	sdk := stamp.NewSDK(a.cfg.Stamp.Endpoint, &http.Client{
		Transport: metrics.NewRequestWatcher("stamp"),
	})
	sc, err := stamp.NewClient(sdk)
	if err != nil {
		return err
	}

	a.stamp = sc

	a.alchemy = alchemy.NewClient(a.cfg.Alchemy.APIKey, &http.Client{
		Transport: metrics.NewRequestWatcher("alchemy"),
	})

	return nil
}

func (a *Application) initGRPCWorker() error {
	srv := grpcsrv.NewGrpcServer(func(ctx context.Context) (context.Context, error) {
		return ctx, nil
	})

	enspb.RegisterEnsServer(srv, server.NewEnsHandler(a.stamp, a.alchemy))
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
