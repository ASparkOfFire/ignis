package sdk

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
)

var once sync.Once

func registerGobTypes() {
	once.Do(func() {
		gob.Register(FDResponse{})
	})
}

type Request struct {
	Method           string      `json:"method,omitempty"`
	Header           http.Header `json:"header,omitempty"`
	Body             []byte      `json:"body,omitempty"`
	ContentLength    int64       `json:"content_length,omitempty"`
	TransferEncoding []string    `json:"transfer_encoding,omitempty"`
	Host             string      `json:"host,omitempty"`
	RemoteAddr       string      `json:"remote_addr,omitempty"`
	RequestURI       string      `json:"request_uri,omitempty"`
	Pattern          string      `json:"pattern,omitempty"`
}

type FDResponse struct {
	Headers    http.Header
	Body       []byte // Change from bytes.Buffer to []byte
	StatusCode int
	Length     int
}

// Encode to binary
func (w *FDResponse) encodeToBinary() ([]byte, error) {
	registerGobTypes()
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	err := encoder.Encode(w)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Decode from binary
func DecodeFromBinary(data []byte) (FDResponse, error) {
	registerGobTypes()
	var resp FDResponse
	buf := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buf)
	err := decoder.Decode(&resp)
	return resp, err
}

func NewFDResponse() *FDResponse {
	return &FDResponse{
		StatusCode: http.StatusOK,
		Headers:    make(http.Header),
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
	r, err := http.NewRequest(req.Method, req.RequestURI, bytes.NewReader(req.Body))
	if err != nil {
		log.Fatal(err)
	}

	h.ServeHTTP(w, r) // execute the user's handler

	w.Length = len(w.Body)
	b, err = w.encodeToBinary()
	if err != nil {
		log.Printf("Error encoding response: %s", err)
	}
	if _, err := os.Stdout.Write(b); err != nil {
		log.Printf("Error writing response: %s", err)
	}
}

func (w *FDResponse) Header() http.Header {
	return w.Headers
}

func (w *FDResponse) Write(b []byte) (n int, err error) {
	w.Body = append(w.Body, b...) // Store as []byte
	return len(b), nil
}

func (w *FDResponse) WriteHeader(status int) {
	w.StatusCode = status
}
