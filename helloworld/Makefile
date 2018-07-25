build:
	protoc -I. --micro_out=. --go_out=. ./proto/greeter.proto
	GOOS=linux GOARCH=amd64 go build
	docker build -t greeter-service .


run:
	docker run -e MICRO_REGISTRY=mdns greeter-service

