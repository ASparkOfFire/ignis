package main

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"github.com/ASparkOfFire/ignis/runtime"
	"github.com/google/uuid"
	"github.com/justincormack/go-memfd"
	"io"
	"net/http"
	"os"
	"strconv"
)

func main() {
	http.HandleFunc("/", handlewasi)
	fmt.Println("Listening on 6969")
	http.ListenAndServe(":6969", nil)
}

func handlewasi(w http.ResponseWriter, r *http.Request) {
	msg := `{
    "body": "SGVsbG8gV29ybGQ=", 
    "method": "GET",
    "url": "/",
    "headers": {
        "Content-Type": "application/json",
        "Authorization": "Bearer abc123",
        "User-Agent": "CustomClient/1.0"
    }
}`

	b, err := os.ReadFile("example/example.wasm")
	if err != nil {
		panic(err)
	}

	// Create an in-memory file descriptor
	fd, err := memfd.Create()
	if err != nil {
		panic(err)
	}
	defer fd.Close()

	rt, err := runtime.New(context.Background(), runtime.Args{
		Stdout:       fd,
		DeploymentID: uuid.UUID{},
		Engine:       "go",
		Blob:         b,
		Cache:        nil,
	})
	if err != nil {
		panic(err)
	}

	if err := rt.Invoke(bytes.NewReader([]byte(msg)), map[string]string{}); err != nil {
		panic(err)
	}

	// SEEK to the beginning of fd before reading
	if _, err := fd.Seek(0, 0); err != nil {
		panic(err)
	}

	// Read the response from fd
	buf := make([]byte, 4096) // Adjust buffer size if needed
	n, err := fd.Read(buf)
	if err != nil && err != io.EOF {
		panic(err)
	}
	resp := buf[:n]

	// Extract status code and length
	var statusCode, length uint32
	bStatusCode := resp[n-8 : n-4]
	bLength := resp[n-4 : n]

	statusCode = binary.LittleEndian.Uint32(bStatusCode)
	length = binary.LittleEndian.Uint32(bLength)

	body := resp[0 : n-8]

	w.WriteHeader(int(statusCode))
	w.Header().Set("Content-Length", strconv.Itoa(int(length)))
	w.Write(body)
}
