run: example-go
	@go run main.go

example-go: proto
	@GOOS=wasip1 GOARCH=wasm go build -o example/example.wasm example/example.go

proto:
	@protoc --go_out=. --go_opt=paths=source_relative --proto_path=. proto/types.proto

.PHONY: example-go run proto