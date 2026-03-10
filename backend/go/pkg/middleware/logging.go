package middleware

import(
	"net/http"
	"time"
	"github.com/rs/zerolog/log"
)

func LoggingMiddleware(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()

		next.ServeHTTP(w,r)

		duration := time.Since(start).Milliseconds()

		log.Info().
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Int64("duration_ms", duration).
			Msg("request completed")
	})
}