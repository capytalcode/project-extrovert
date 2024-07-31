package router

import (
	"log"
	"net/http"
)

type Middleware interface {
	Serve(r http.HandlerFunc) http.HandlerFunc
}

type MiddlewaredResponse struct {
	w          http.ResponseWriter
	status     int
	bodyWrites [][]byte
}

func NewMiddlewaredResponse(w http.ResponseWriter) *MiddlewaredResponse {
	return &MiddlewaredResponse{w, 200, [][]byte{[]byte("")}}
}

func (m *MiddlewaredResponse) WriteHeader(s int) {
	log.Printf("Status changed %v", s)
	m.status = s
}

func (m *MiddlewaredResponse) Header() http.Header {
	return m.w.Header()
}

func (m *MiddlewaredResponse) Write(b []byte) (int, error) {
	m.bodyWrites = append(m.bodyWrites, b)
	return len(b), nil
}

func (m *MiddlewaredResponse) ReallyWriteHeader() (int, error) {
	m.w.WriteHeader(m.status)
	bytes := 0
	for _, b := range m.bodyWrites {
		by, err := m.w.Write(b)
		bytes += by
		if err != nil {
			return bytes, err
		}
	}
	return bytes, nil
}
