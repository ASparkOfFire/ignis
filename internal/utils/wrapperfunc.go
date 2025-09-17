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
	// create a new buffer for output
	fd := new(bytes.Buffer)
	
	// create a reader for the request payload
	stdin := bytes.NewReader(reqPayload)

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

	fmt.Printf("Invoking WASM with payload of size: %d\n", len(reqPayload))
	if err := rt.Invoke(stdin, nil, script); err != nil {
		return nil, fmt.Errorf("failed to invoke WASM runtime: %w", err)
	}
	fmt.Printf("WASM invocation completed\n")

	responseBytes, err := readStdout(fd)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Read %d bytes from WASM output\n", len(responseBytes))

	return parseWASMResponse(responseBytes)
}

// readStdout reads and returns data from the in-memory file descriptor.
func readStdout(fd *bytes.Buffer) ([]byte, error) {
	// Read all data from the buffer
	data := fd.Bytes()
	if len(data) == 0 {
		return nil, fmt.Errorf("no data read from WASM")
	}
	return data, nil
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
