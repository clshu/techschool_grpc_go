gen:
		protoc --proto_path=proto proto/*.proto --go_out=. --plugin=grpc

clean:
		rm -rf pb/*

run:
		go run main.go