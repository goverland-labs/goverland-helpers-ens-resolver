package health

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"github.com/s-larionov/process-manager"

	"github.com/goverland-labs/helpers-ens-resolver/pkg/middleware"
)

func NewHealthCheckServer(listen, path string, handler http.Handler) *http.Server {
	router := mux.NewRouter()
	router.Use(middleware.Panic, middleware.Prometheus)
	router.Handle(path, handler)

	server := &http.Server{
		Addr:    listen,
		Handler: router,
	}

	return server
}

func DefaultHandler(manager *process.Manager) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]interface{}{
			"process_manager": manager.IsRunning(),
		}

		body, err := json.Marshal(resp)
		if err != nil {
			log.Error().Err(err).Msg("unable to marshal health check")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(body)
	})
}
