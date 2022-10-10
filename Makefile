submodule:
		git submodule init; git submodule update

gen:
		protoc --proto_path=proto proto/*.proto --go_out=. --plugin=grpc --go-grpc_out=.

clean:
		rm -rf pb/*

server:
		go run cmd/server/main.go -port 50051 -tls

client:
		go run cmd/client/main.go -address localhost:50051 -tls

test:
		go test -cover -race ./...

cert:
		cd cert; ./gen.sh; cd ..

.PHONY: gen clean server client test cert submodule