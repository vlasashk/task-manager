package httpchi

import (
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
	"net/http"
)

func LoggerRequestID(logger zerolog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := logger.With().
				Str("request_id", middleware.GetReqID(r.Context())).
				Caller().
				Logger()
			r = r.WithContext(log.WithContext(r.Context()))
			next.ServeHTTP(w, r)
		})
	}
}
