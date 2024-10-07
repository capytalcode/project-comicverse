package rerrors

import (
	"encoding/json"
	"fmt"
	"net/http"
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
		w.Write(j)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(j)
	w.WriteHeader(rerr.StatusCode)
}
