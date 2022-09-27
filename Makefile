gen:
		protoc --proto_path=proto proto/*.proto --go_out=. --plugin=grpc --go-grpc_out=.

clean:
		rm -rf pb/*

run:
		go run main.go

test:
		go test -cover -race ./...