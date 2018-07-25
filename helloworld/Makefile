build:
	protoc -I. --micro_out=. --go_out=. ./proto/greeter.proto

run-server:
	go run main.go --registry=mdns
run-cli:
	go run cli/client.go --registry=mdns