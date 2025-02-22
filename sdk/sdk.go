package sdk

import (
	"bytes"
	"encoding/json"
	types "github.com/ASparkOfFire/ignis/proto"
	"google.golang.org/protobuf/proto"
	"io"
	"log"
	"net/http"
	"os"
)

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
	Body       []byte
	StatusCode int
	Length     int
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
	protoResp := types.FDResponse{
		Body:       w.Body,
		StatusCode: int32(w.StatusCode),
		Length:     int32(w.Length),
		Header:     make(map[string]*types.HeaderFields),
	}
	for k, v := range w.Headers {
		protoResp.Header[k] = &types.HeaderFields{Fields: v}
	}

	b, err = proto.Marshal(&protoResp)
	if err != nil {
		log.Printf("Error encoding response: %s", err)
	}
	n, err := os.Stdout.Write(b)
	if err != nil || n != len(b) {
		log.Printf("Error writing response: %s, bytes written: %d", err, n)
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
