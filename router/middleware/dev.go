package middleware

import (
	"log/slog"
	"net/http"
)

func DevMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		w.Header().Set("Cache-Control", "no-store")
	})
}

type LoggerMiddleware struct {
	log *slog.Logger
}

func NewLoggerMiddleware(l *slog.Logger) *LoggerMiddleware {
	return &LoggerMiddleware{l}
}

func (l *LoggerMiddleware) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l.log.Debug("Handling request", slog.String("path", r.URL.Path))

		next.ServeHTTP(w, r)

		l.log.Debug("Request finished",
			slog.String("path", r.URL.Path),
			slog.String("status", w.Header().Get("Status")),
		)
	})
}
