package middleware

import (
	"log/slog"
	"math/rand"
	"net/http"
	"strconv"
)

func DevMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		next.ServeHTTP(w, r)
	})
}

type LoggerMiddleware struct {
	log *slog.Logger
}

func NewLoggerMiddleware(l *slog.Logger) *LoggerMiddleware {
	l = l.WithGroup("logger_middleware")
	return &LoggerMiddleware{l}
}

func (l *LoggerMiddleware) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := randHash(5)

		l.log.Info("NEW REQUEST",
			slog.String("id", id),
			slog.String("status", "xxx"),
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
		)

		next.ServeHTTP(w, r)

		s := w.Header().Get("Status")
		if s == "" {
			s = strconv.Itoa(http.StatusOK)
		}

		sc, err := strconv.Atoi(s)
		if err != nil {
			l.log.Error("INVALID REQUEST STATUS",
				slog.String("id", id),
				slog.String("status", s),
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
			)
			return
		}
		if sc >= 400 {
			l.log.Warn("ERR REQUEST",
				slog.String("id", id),
				slog.String("status", s),
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
			)
			return
		}

		l.log.Info("END REQUEST",
			slog.String("id", id),
			slog.String("status", s),
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
		)
	})
}

const HASH_CHARS = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// This is not the most performant function, as a TODO we could
// improve based on this Stackoberflow thread:
// https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
func randHash(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = HASH_CHARS[rand.Int63()%int64(len(HASH_CHARS))]
	}
	return string(b)
}
