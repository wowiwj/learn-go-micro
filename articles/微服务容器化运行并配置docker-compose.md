# 微服务容器化运行并配置docker-compose

上节代码参考: https://github.com/wowiwj/learn-go-micro/tree/v0.01


### Docker的使用及介绍

请参考： https://docker.rails365.net/chapters/1.html

### Docker化service服务

创建文件: helloworld/Dockerfile


```dockerfile
FROM debian:latest
RUN mkdir /app
WORKDIR /app
ADD helloworld /app/helloworld
CMD ["./helloworld"]
```

修改 `helloworld/Makefile` 文件

```makefile
build:
	protoc -I. --micro_out=. --go_out=. ./proto/greeter.proto
	GOOS=linux GOARCH=amd64 go build
	docker build -t greeter-service .


run:
	docker run -e MICRO_REGISTRY=mdns greeter-service
```

编译及构建容器:

```
make build
```

执行服务

```
make run
```

### docker化cli服务

进入cli目录

```
cd cli
```

创建文件: helloworld/cli/Dockerfile

```
FROM debian:latest
RUN mkdir /app
WORKDIR /app
ADD cli /app/cli
CMD ["./cli"]
```


修改 `helloworld/cli/Makefile` 文件

```
build:
	GOOS=linux GOARCH=amd64 go build
	docker build -t greeter-cli .

run:
	docker run -e MICRO_REGISTRY=mdns greeter-cli
```

编译及构建容器:

```
make build
```

执行服务

```
make run
```

可以看到，成功的输出

```
docker run -e MICRO_REGISTRY=mdns greeter-cli
hello 111
```

### 通过docker-compose构建

创建文件：helloworld/docker-compose.yaml

```yaml
version: "2"


services:
  greeter-server:
    build: .
    environment:
      MICRO_REGISTRY: "mdns"

  greeter-cli:
      build: ./cli
      environment:
        MICRO_REGISTRY: "mdns"
      depends_on:
        - greeter-server
```

执行：
```
docker-compose up
```

输出

```
docker-compose up
Starting helloworld_greeter-server_1 ... done
Recreating helloworld_greeter-cli_1 ... done
Attaching to helloworld_greeter-server_1, helloworld_greeter-cli_1
greeter-server_1  | 2018/07/25 07:32:15 Listening on [::]:35725
greeter-server_1  | 2018/07/25 07:32:15 Broker Listening on [::]:37259
greeter-server_1  | 2018/07/25 07:32:15 Registering node: greeter-da44f95e-8fdc-11e8-9ce6-0242ac140002
greeter-cli_1     | hello 111
helloworld_greeter-cli_1 exited with code 0

```

说明微服务通信成功

相关代码请查看: https://github.com/wowiwj/learn-go-micro/tree/v0.0.2






