build:
	GOOS=linux GOARCH=amd64 go build
	docker build -t greeter-cli .

run:
	docker run -e MICRO_REGISTRY=mdns greeter-cli