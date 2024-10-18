package rerrors

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"forge.capytal.company/capytalcode/project-comicverse/router/middleware"
	"github.com/a-h/templ"
)

type RouteError struct {
	StatusCode int            `json:"status_code"`
	Error      string         `json:"error"`
	Info       map[string]any `json:"info"`
}

func NewRouteError(status int, error string, info ...map[string]any) RouteError {
	rerr := RouteError{StatusCode: status, Error: error}
	if len(info) > 0 {
		rerr.Info = info[0]
	} else {
		rerr.Info = map[string]any{}
	}
	return rerr
}

func (rerr RouteError) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if rerr.StatusCode == 0 {
		rerr.StatusCode = http.StatusNotImplemented
	}

	if rerr.Error == "" {
		rerr.Error = "MISSING ERROR DESCRIPTION"
	}

	if rerr.Info == nil {
		rerr.Info = map[string]any{}
	}

	w.Header().Set("Content-Type", "application/json")

	j, err := json.Marshal(rerr)
	if err != nil {
		j, _ := json.Marshal(RouteError{
			StatusCode: http.StatusInternalServerError,
			Error:      "Failed to marshal error message to JSON",
			Info: map[string]any{
				"source_value": fmt.Sprintf("%#v", rerr),
				"error":        err.Error(),
			},
		})
		w.WriteHeader(http.StatusInternalServerError)
		if _, err = w.Write(j); err != nil {
			_, _ = w.Write([]byte("Failed to write error JSON string to body"))
		}
		return
	}

	w.WriteHeader(rerr.StatusCode)
	if _, err = w.Write(j); err != nil {
		_, _ = w.Write([]byte("Failed to write error JSON string to body"))
	}
}

type ErrorMiddlewarePage func(err RouteError) templ.Component

type ErrorMiddleware struct {
	page     ErrorMiddlewarePage
	notfound ErrorMiddlewarePage
	log      *slog.Logger
}

func NewErrorMiddleware(
	p ErrorMiddlewarePage,
	l *slog.Logger,
	notfound ...ErrorMiddlewarePage,
) *ErrorMiddleware {
	var nf ErrorMiddlewarePage
	if len(notfound) > 0 {
		nf = notfound[0]
	} else {
		nf = p
	}

	l = l.WithGroup("error_middleware")

	return &ErrorMiddleware{p, nf, l}
}

func (m *ErrorMiddleware) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if uerr := r.URL.Query().Get("error"); uerr != "" {
			e, err := base64.URLEncoding.DecodeString(uerr)
			if err != nil {
				m.log.Error("Failed to decode \"error\" parameter from error redirect",
					slog.String("method", r.Method),
					slog.String("path", r.URL.Path),
					slog.Int("status", 0),
					slog.String("data", string(e)),
				)
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(
					fmt.Sprintf("Data %s\nError %s", string(e), err.Error()),
				))
				return
			}

			var rerr RouteError
			if err := json.Unmarshal(e, &rerr); err != nil {
				m.log.Error("Failed to decode \"error\" parameter from error redirect",
					slog.String("method", r.Method),
					slog.String("path", r.URL.Path),
					slog.Int("status", 0),
					slog.String("data", string(e)),
				)
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(
					fmt.Sprintf("Data %s\nError %s", string(e), err.Error()),
				))
				return
			}

			w.WriteHeader(rerr.StatusCode)
			if err := m.page(rerr).Render(r.Context(), w); err != nil {
				_, _ = w.Write(e)
			}

			return
		}

		var buf bytes.Buffer
		mw := middleware.MultiResponseWriter(w, &buf)

		next.ServeHTTP(mw, r)

		if mw.Header().Get("Status") == "" {
			m.log.Warn("Endpoint did not return a Status code",
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("status", mw.Header().Get("Status")),
				slog.Any("data", buf),
			)
			return
		}

		status, err := strconv.Atoi(mw.Header().Get("Status"))
		if err != nil {
			m.log.Warn("Failed to parse Status code to a integer",
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("status", mw.Header().Get("Status")),
				slog.String("data", err.Error()),
			)
			return
		}

		if status < 400 {
			return
		} else if status == 404 {
			rerr := NotFound()

			w.WriteHeader(rerr.StatusCode)

			b, err := json.Marshal(rerr)
			if err != nil {
				_, _ = w.Write([]byte(
					fmt.Sprintf("%#v", rerr),
				))
				return
			}

			if prefersHtml(r.Header) {
				u, _ := url.Parse(r.URL.String())

				q := u.Query()
				q.Add("error", base64.URLEncoding.EncodeToString(b))
				u.RawQuery = q.Encode()

				http.Redirect(w, r, u.String(), http.StatusTemporaryRedirect)
				return
			}

			_, _ = w.Write(b)

			return
		}

		body, err := io.ReadAll(&buf)
		if err != nil {
			m.log.Error("Failed to read from error body",
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("status", mw.Header().Get("Status")),
				slog.Any("data", buf),
			)
			return
		}

		if mw.Header().Get("Content-Type") != "application/json" {
			m.log.Warn("Endpoint didn't return a structured error",
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("status", mw.Header().Get("Status")),
				slog.String("data", string(body)),
			)
			return
		}

		if prefersHtml(r.Header) {
			u, _ := url.Parse(r.URL.String())

			q := u.Query()
			q.Add("error", base64.URLEncoding.EncodeToString(body))
			u.RawQuery = q.Encode()

			http.Redirect(w, r, u.String(), http.StatusTemporaryRedirect)
		}
	})
}

func prefersHtml(h http.Header) bool {
	if h.Get("Accept") == "" {
		return false
	}
	return (strings.Contains(h.Get("Accept"), "text/html") ||
		strings.Contains(h.Get("Accept"), "application/xhtml+xml") ||
		strings.Contains(h.Get("Accept"), "application/xml")) &&
		!strings.Contains(h.Get("Accept"), "application/json")
}
