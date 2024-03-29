package middleware

import (
	"fmt"
	"io"
	"net/http"
	"runtime/debug"

	"github.com/rs/zerolog/log"
)

func Panic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func(w http.ResponseWriter, r *http.Request) {
			var message string

			rec := recover()
			if rec == nil {
				return
			}

			switch v := rec.(type) {
			case string:
				message = v
			case error:
				message = v.Error()
			default:
				message = "unknown error"
			}

			body, _ := io.ReadAll(r.Body)

			log.Error().Fields(map[string]interface{}{
				"message": message,
				"request": fmt.Sprintf("%s %s?%s", r.Method, r.URL.String(), r.URL.Query().Encode()),
				"body":    string(body),
				"stack":   string(debug.Stack()),
			}).Msg("panic raised")

			w.WriteHeader(http.StatusInternalServerError)
		}(w, r)

		next.ServeHTTP(w, r)
	})
}
