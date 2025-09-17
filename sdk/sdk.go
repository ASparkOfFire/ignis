package sdk

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"os"

	types "github.com/ASparkOfFire/ignis/internal/proto"
	_ "github.com/breml/rootcerts"
	_ "github.com/ignis-runtime/net/http"
	"google.golang.org/protobuf/proto"
)

type Response struct {
	Headers    http.Header
	Body       []byte
	StatusCode int
	Length     int
}

func NewFDResponse() *Response {
	return &Response{
		StatusCode: http.StatusOK,
		Headers:    make(http.Header),
	}
}

func Handle(h http.Handler, stdin io.Reader) {
	// If no stdin provided, use os.Stdin as fallback
	if stdin == nil {
		stdin = os.Stdin
	}

	HandleWithIO(h, stdin, os.Stdout)
}

func HandleWithIO(h http.Handler, stdin io.Reader, stdout io.Writer) {
	b, err := io.ReadAll(stdin)
	if err != nil {
		log.Fatal(err)
	}
	w := NewFDResponse()

	var req types.FDRequest
	if err := proto.Unmarshal(b, &req); err != nil {
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

	// Write to provided stdout
	n, err := stdout.Write(b)
	if err != nil || n != len(b) {
		log.Printf("Error writing response: %s, bytes written: %d", err, n)
	}
}

func (w *Response) Header() http.Header {
	return w.Headers
}

func (w *Response) Write(b []byte) (n int, err error) {
	w.Body = append(w.Body, b...) // Store as []byte
	return len(b), nil
}

func (w *Response) WriteHeader(status int) {
	w.StatusCode = status
}
