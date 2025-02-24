run: example-go example-js
	@go run main.go

example-go: proto
	@GOOS=wasip1 GOARCH=wasm go build -o example/go/example.wasm example/go/example.go

example-js: proto
	@npm i --prefix=example/js
	@esbuild example/js/example.js --bundle --platform=neutral --outfile=example/js/dist/example.js

proto:
	@protoc --go_out=. --go_opt=paths=source_relative --proto_path=. internal/proto/types.proto

.PHONY: example-go run proto