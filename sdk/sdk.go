package sdk

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
)

type Request struct {
	Body   []byte      `json:"body"`
	Method string      `json:"method"`
	URL    string      `json:"url"`
	Header http.Header `json:"header"`
}

type FDResponse struct {
	buffer     bytes.Buffer
	statusCode int
	header     http.Header
}

func NewFDResponse() *FDResponse {
	return &FDResponse{
		statusCode: http.StatusOK,
		header:     make(http.Header),
	}
}

func Handle(h http.Handler) {
	b, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	w := NewFDResponse()

	var req Request
	if err := json.Unmarshal(b, &req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	r, err := http.NewRequest(req.Method, req.URL, bytes.NewReader(req.Body))
	if err != nil {
		log.Fatal(err)
	}

	h.ServeHTTP(w, r) // execute the user's handler

	// Write response body
	if _, err := os.Stdout.Write(w.buffer.Bytes()); err != nil {
		log.Printf("Error Writing Response Body: %s\n", err)
	}

	// Write status code and body length
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint32(buf[0:4], uint32(w.statusCode))
	binary.LittleEndian.PutUint32(buf[4:8], uint32(w.buffer.Len()))

	if _, err := os.Stdout.Write(buf); err != nil {
		log.Printf("Error Writing Status Code and Content-Length: %s\n", err)
	}
}

func (w *FDResponse) Header() http.Header {
	return w.header
}

func (w *FDResponse) Write(b []byte) (n int, err error) {
	return w.buffer.Write(b)
}

func (w *FDResponse) WriteHeader(status int) {
	w.statusCode = status
}
