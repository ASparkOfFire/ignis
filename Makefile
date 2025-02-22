run: example-go
	@go run main.go

example-go:
	@GOOS=wasip1 GOARCH=wasm go build -o example/example.wasm example/example.go

.PHONY: example-go run