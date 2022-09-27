gen:
		protoc --proto_path=proto proto/*.proto --go_out=. --plugin=grpc --go-grpc_out=.

clean:
		rm -rf pb/*

server:
		go run cmd/server/main.go -port 50051

client:
		go run cmd/client/main.go -address localhost:50051

test:
		go test -cover -race ./...