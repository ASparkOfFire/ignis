package utils

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/ASparkOfFire/ignis/internal/cache"
	types "github.com/ASparkOfFire/ignis/internal/proto"
	"github.com/ASparkOfFire/ignis/internal/runtime"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/justincormack/go-memfd"
	"google.golang.org/protobuf/proto"
)

// WASIWrapper initializes and executes a WASM runtime for processing HTTP requests.
func WASIWrapper(wasmFile string, cache cache.ModCache[uuid.UUID], engine runtime.RuntimeEngine, id uuid.UUID) gin.HandlerFunc {
	return func(c *gin.Context) {
		reqPayload, err := buildRequestPayload(c)
		if err != nil {
			logAndRespond(c, http.StatusInternalServerError, "Failed to process request", err)
			return
		}

		wasmBytes, err := os.ReadFile(wasmFile)
		if err != nil {
			logAndRespond(c, http.StatusInternalServerError, "Failed to read WASM file", err)
			return
		}

		respProto, err := executeWASM(reqPayload, wasmBytes, cache, engine, id)
		if err != nil {
			logAndRespond(c, http.StatusInternalServerError, "Failed to execute WASM", err)
			return
		}

		sendResponse(c, respProto)
	}
}

// buildRequestPayload constructs a protobuf request payload from the Gin context.
func buildRequestPayload(c *gin.Context) ([]byte, error) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil && (!errors.Is(err, io.EOF) || len(body) > 0) {
		return nil, fmt.Errorf("error reading request body: %w", err)
	}

	headers := make(map[string]*types.HeaderFields, len(c.Request.Header))
	for k, v := range c.Request.Header {
		headers[k] = &types.HeaderFields{Fields: v}
	}

	reqMsg := &types.FDRequest{
		Method:           c.Request.Method,
		Header:           headers,
		Body:             body,
		ContentLength:    c.Request.ContentLength,
		TransferEncoding: &types.StringSlice{Fields: c.Request.TransferEncoding},
		Host:             c.Request.Host,
		RemoteAddr:       c.Request.RemoteAddr,
		RequestURI:       c.Request.RequestURI,
		Pattern:          c.Request.Pattern,
	}

	return proto.Marshal(reqMsg)
}

// executeWASM runs the WASM binary and returns the parsed response.
func executeWASM(reqPayload, wasmBytes []byte, cache cache.ModCache[uuid.UUID], engine runtime.RuntimeEngine, id uuid.UUID) (*types.FDResponse, error) {
	fd, err := memfd.Create()
	if err != nil {
		return nil, fmt.Errorf("failed to create in-memory file descriptor: %w", err)
	}
	defer fd.Close()

	rt, err := runtime.New(context.Background(), runtime.Args{
		Stdout:       fd,
		DeploymentID: id,
		Engine:       engine,
		Blob:         wasmBytes,
		Cache:        cache,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize WASM runtime: %w", err)
	}

	var script []byte
	if engine == runtime.RuntimeEngineJS {
		script = wasmBytes
	}

	if err := rt.Invoke(bytes.NewReader(reqPayload), nil, script); err != nil {
		return nil, fmt.Errorf("failed to invoke WASM runtime: %w", err)
	}

	responseBytes, err := readMemfd(fd)
	if err != nil {
		return nil, err
	}

	return parseWASMResponse(responseBytes)
}

// readMemfd reads and returns data from the in-memory file descriptor.
func readMemfd(fd *memfd.Memfd) ([]byte, error) {
	if _, err := fd.Seek(0, 0); err != nil {
		return nil, fmt.Errorf("failed to seek file descriptor: %w", err)
	}

	buf := make([]byte, 4096) // Adjust buffer size as needed
	n, err := fd.Read(buf)
	if err != nil && err != io.EOF {
		return nil, fmt.Errorf("failed to read response from WASM: %w", err)
	}

	return buf[:n], nil
}

// parseWASMResponse unmarshals the WASM response into a struct.
func parseWASMResponse(data []byte) (*types.FDResponse, error) {
	var respProto types.FDResponse
	if err := proto.Unmarshal(data, &respProto); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return &respProto, nil
}

// sendResponse formats and sends the WASM response to the client.
func sendResponse(c *gin.Context, protoResp *types.FDResponse) {
	for k, v := range protoResp.Header {
		for _, val := range v.Fields {
			c.Header(k, val)
		}
	}

	c.Data(int(protoResp.StatusCode), "application/octet-stream", protoResp.Body)
}

// logAndRespond logs the error and sends an HTTP response with an error message.
func logAndRespond(c *gin.Context, status int, msg string, err error) {
	log.Printf("%s: %v\n", msg, err)
	c.JSON(status, gin.H{"error": msg})
}
