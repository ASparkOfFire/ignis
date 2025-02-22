package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	types "github.com/ASparkOfFire/ignis/proto"
	"github.com/ASparkOfFire/ignis/runtime"
	"github.com/ASparkOfFire/ignis/sdk"
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

func WASIWrapper(wasmFile string) gin.HandlerFunc {
	return func(c *gin.Context) {
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			if !errors.Is(err, io.EOF) || len(body) != 0 { // it is perfectly okay to have empty body
				log.Printf("Error reading body: %s\n", err)
				c.JSON(500, gin.H{"error": "Error reading body"})
				return
			}
		}

		msg := sdk.Request{
			Method:           c.Request.Method,
			Header:           c.Request.Header,
			Body:             body,
			ContentLength:    c.Request.ContentLength,
			TransferEncoding: c.Request.TransferEncoding,
			Host:             c.Request.Host,
			RemoteAddr:       c.Request.RemoteAddr,
			RequestURI:       c.Request.RequestURI,
			Pattern:          c.Request.Pattern,
		}
		req, err := json.Marshal(msg)
		if err != nil {
			log.Fatal(err)
		}

		b, err := os.ReadFile(wasmFile)
		if err != nil {
			log.Printf("Error reading WASM file: %s\n", err)
			c.JSON(500, gin.H{"error": "Failed to read WASM file"})
			return
		}

		// Create an in-memory file descriptor
		fd, err := memfd.Create()
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to create in-memory file descriptor"})
			return
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
			c.JSON(500, gin.H{"error": "Failed to initialize runtime"})
			return
		}

		if err := rt.Invoke(bytes.NewReader(req), map[string]string{}); err != nil {
			c.JSON(500, gin.H{"error": "Failed to invoke runtime"})
			return
		}

		// SEEK to the beginning of fd before reading
		if _, err := fd.Seek(0, 0); err != nil {
			c.JSON(500, gin.H{"error": "Failed to seek file descriptor"})
			return
		}

		// Read the response from fd
		buf := make([]byte, 4096) // Adjust buffer size if needed
		n, err := fd.Read(buf)
		if err != nil && err != io.EOF {
			c.JSON(500, gin.H{"error": "Failed to read response"})
			return
		}
		buf = buf[:n]
		var protoResp types.FDResponse

		err = proto.Unmarshal(buf, &protoResp)
		if err != nil {
			log.Printf("Failed to decode response: %v\n", err)
			c.JSON(500, gin.H{"error": "Failed to decode response"})
			return
		}
		header := make(http.Header)
		for k, v := range protoResp.Header {
			header[k] = v.Fields
		}
		resp := sdk.FDResponse{
			Headers:    header,
			Body:       protoResp.Body,
			StatusCode: int(protoResp.StatusCode),
			Length:     int(protoResp.Length),
		}
		// Set response headers
		c.Status(resp.StatusCode)
		c.Header("Content-Length", strconv.Itoa(resp.Length))
		for k, v := range resp.Headers {
			for _, val := range v {
				c.Header(k, val)
			}
		}

		// Send response body
		c.Data(resp.StatusCode, "application/octet-stream", resp.Body)
		return
	}
}
