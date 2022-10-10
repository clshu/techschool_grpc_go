submodule:
		git submodule init; git submodule update

gen:
		protoc --proto_path=proto proto/*.proto --go_out=. --plugin=grpc --go-grpc_out=.

clean:
		rm -rf pb/*

server1-tls:
		go run cmd/server/main.go -port 50052 -tls

server1:
		go run cmd/server/main.go -port 50052

server2-tls:
		go run cmd/server/main.go -port 50053 -tls

server2:
		go run cmd/server/main.go -port 50053

server-tls:
		go run cmd/server/main.go -port 50051 -tls

server:
		go run cmd/server/main.go -port 50051

client-tls:
		go run cmd/client/main.go -address localhost:50051 -tls

client:
		go run cmd/client/main.go -address localhost:50051

client1-tls:
		go run cmd/client/main.go -address localhost:8080 -tls

client1:
		go run cmd/client/main.go -address localhost:8080

test:
		mkdir -p tmp; go test -cover -race ./...

cert:
		cd cert; ./gen.sh; cd ..

.PHONY: gen clean server server-tls client client-tls test cert submodule server1 server1-tls server2 server2-tls