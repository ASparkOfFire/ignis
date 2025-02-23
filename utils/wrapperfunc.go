package utils

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	types "github.com/ASparkOfFire/ignis/proto"
	"github.com/ASparkOfFire/ignis/runtime"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/justincormack/go-memfd"
	"google.golang.org/protobuf/proto"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

// WASIWrapper initializes and executes a WASM runtime for processing HTTP requests.
func WASIWrapper(wasmFile string) gin.HandlerFunc {
	return func(c *gin.Context) {
		reqPayload, err := buildRequestPayload(c)
		if err != nil {
			log.Printf("Request processing error: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process request"})
			return
		}

		wasmBytes, err := os.ReadFile(wasmFile)
		if err != nil {
			log.Printf("Error reading WASM file: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read WASM file"})
			return
		}

		respProto, err := executeWASM(reqPayload, wasmBytes)
		if err != nil {
			log.Printf("WASM execution error: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to execute WASM"})
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

	headers := make(map[string]*types.HeaderFields)
	for k, v := range c.Request.Header {
		headers[k] = &types.HeaderFields{
			Fields: v,
		}
	}

	msg := types.FDRequest{
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

	return proto.Marshal(&msg)
}

// executeWASM runs the WASM binary and returns the parsed response.
func executeWASM(reqPayload []byte, wasmBytes []byte) (*types.FDResponse, error) {
	fd, err := memfd.Create()
	if err != nil {
		return nil, fmt.Errorf("failed to create in-memory file descriptor: %w", err)
	}
	defer fd.Close()

	rt, err := runtime.New(context.Background(), runtime.Args{
		Stdout:       fd,
		DeploymentID: uuid.UUID{},
		Engine:       "go",
		Blob:         wasmBytes,
		Cache:        nil,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize WASM runtime: %w", err)
	}

	if err := rt.Invoke(bytes.NewReader(reqPayload), nil); err != nil {
		return nil, fmt.Errorf("failed to invoke WASM runtime: %w", err)
	}

	if _, err := fd.Seek(0, 0); err != nil {
		return nil, fmt.Errorf("failed to seek file descriptor: %w", err)
	}

	buf := make([]byte, 4096) // Adjust buffer size as needed
	n, err := fd.Read(buf)
	if err != nil && err != io.EOF {
		return nil, fmt.Errorf("failed to read response from WASM: %w", err)
	}

	var respProto types.FDResponse
	if err := proto.Unmarshal(buf[:n], &respProto); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &respProto, nil
}

// sendResponse formats and sends the WASM response to the client.
func sendResponse(c *gin.Context, protoResp *types.FDResponse) {
	headers := make(http.Header)
	for k, v := range protoResp.Header {
		headers[k] = v.Fields
	}

	c.Status(int(protoResp.StatusCode))
	c.Header("Content-Length", strconv.Itoa(int(protoResp.Length)))
	for k, v := range headers {
		for _, val := range v {
			c.Header(k, val)
		}
	}

	c.Data(int(protoResp.StatusCode), "application/octet-stream", protoResp.Body)
}
