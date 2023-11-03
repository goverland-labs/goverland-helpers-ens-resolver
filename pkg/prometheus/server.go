package prometheus

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"helpers-ens-resolver/pkg/middleware"
)

func NewPrometheusServer(listen, path string) *http.Server {
	handler := mux.NewRouter()
	handler.Use(middleware.Panic)
	handler.Handle(path, promhttp.Handler())

	server := &http.Server{
		Addr:    listen,
		Handler: handler,
	}

	return server
}
