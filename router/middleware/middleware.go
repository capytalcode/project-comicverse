package middleware

import (
	"errors"
	"fmt"
	"net/http"
)

type Middleware interface {
	Serve(r http.HandlerFunc) http.HandlerFunc
}

type MiddlewaredReponse struct {
	w          http.ResponseWriter
	statuses   []int
	bodyWrites [][]byte
}

func NewMiddlewaredResponse(w http.ResponseWriter) *MiddlewaredReponse {
	return &MiddlewaredReponse{w, []int{500}, [][]byte{[]byte("")}}
}

func (m *MiddlewaredReponse) WriteHeader(s int) {
	m.statuses = append(m.statuses, s)
}

func (m *MiddlewaredReponse) Header() http.Header {
	return m.w.Header()
}

func (m *MiddlewaredReponse) Write(b []byte) (int, error) {
	m.bodyWrites = append(m.bodyWrites, b)
	return len(b), nil
}

func (m *MiddlewaredReponse) ReallyWriteHeader() (int, error) {
	status := m.statuses[len(m.statuses)-1]
	m.w.WriteHeader(status)
	bytes := 0
	for _, b := range m.bodyWrites {
		by, err := m.w.Write(b)
		if err != nil {
			return bytes, errors.Join(
				fmt.Errorf(
					"Failed to write to response in middleware."+
						"\nStatuses are %v"+
						"\nTried to write %v bytes"+
						"\nTried to write response:\n%s",
					m.statuses, bytes, string(b),
				),
				err,
			)
		}
		bytes += by
	}

	return bytes, nil
}
