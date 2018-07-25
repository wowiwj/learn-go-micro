# 通过go-micro搭建一个简单的问候服务器

### 安装Protobuf

> 安装过程请注意梯子的使用，不过，不要问我梯子是什么

```
go get github.com/micro/protoc-gen-micro
```

#### 关于Proto

由于每个微服务对应一个独立运行的代码库，一个很自然的问题就是如何在这些微服务之间通信。

我们可以使用传统的REST，用http传输JSON或者XML。但用这种方法的一个问题在于，当两个微服务A和B之间要通信时，A要先把数据编码成JSON/XML，然后发送一个大字符串给B，然后B在将数据从JSON/XML解码。这在大型应用中可能会造成大量的开销。尽管我们在和浏览器交互时必须使用这种方法，但微服务之间可以选择其他方式。

gRPC就是这另外一种方式。gRPC是谷歌出品的一个RPC通信工具，它很轻量，且其协议是基于二进制的。让我们来仔细研究下这个定义。gRPC将二进制当作其核心的编码格式。在我们使用JSON的RESTful例子中，我们的数据会以字符串的格式通过http传输。字符串包含了相对大量的元数据，用于描述其编码格式，长度，内容格式以及其他必要数据。之所以包含这些元数据，是因为要让传统的网页浏览器知道收到的数据会是怎样的。但是在两个微服务之间通信时，我们不一定需要这么多元数据。我们可以只需要更轻量的二进制数据。gRPC支持全新的HTTP 2协议，正好可以使用二进制数据。gRPC甚至可以建立双向的流数据。HTTP 2是gRPC的基础，如果你想了解更多HTTP 2的内容，可以看Google的这篇文章。

那么我们该怎么用二进制数据呢？gRPC使用protobuf来描述数据格式。使用Protobuf，你可以清晰的定义一个微服务的界面。关于gRPC，我建议你读一读这篇文章。


#### 编写一个基础的proto

文件目录: helloworld/proto/greeter.proto

```proto

syntax = "proto3";

service Greeter{
    rpc Hello(HelloRequest) returns (HelloResponse){}
}

message HelloRequest{
    string name = 1;
}

message HelloResponse{
    string greeting = 2;
}
```

首先，你得定义service。一个service定义了此服务暴露给外界的交互界面。然后，你得定义message。宽泛的讲，message就是你的数据结构

这个文件里，message由protobuf处理，而service则是由protobuf的grpc插件处理。这个grpc插件使我们定义的service能使用message。

#### 创建Makefile文件
文件目录: helloworld/Makefile

```makefile
build:
	protoc -I. --micro_out=. --go_out=. ./proto/greeter.proto
```

现在，如果你运行$ make build，然后前往文件夹proto，greeter.pb.go。这个文件是protoc自动生成的，它将proto文件中的service转化成了需要我们在Golang代码中需要编写的interface。


#### 创建服务

创建文件: helloworld/main.go

```go
package main

import (
	"github.com/micro/go-micro"
	pb "LearnMicro/helloworld/proto"
	"context"
	"fmt"
)

type Greeter struct {}

func (this *Greeter)Hello(ctx context.Context, req *pb.HelloRequest, res *pb.HelloResponse) error {
	res.Greeting = "hello " + req.Name
	return nil
}


func main() {

	// 创建一个服务，并添加相关的配置
	service := micro.NewService(
		micro.Name("greeter"),
	)

	//解析终端的命令进行初始化
	service.Init()

	// 注册句柄
	pb.RegisterGreeterHandler(service.Server(),new(Greeter))

	// 运行服务
	if err := service.Run();err != nil {
		fmt.Println(err)
	}

}
```

#### 创建测试客户端

创建文件: helloworld/cli/client.go

```go
package main

import (
	pb "LearnMicro/helloworld/proto"
	"context"
	"fmt"
	"github.com/micro/go-micro"
)

func main() {


	// 创建服务并添加响应的配置
	service := micro.NewService(
		micro.Name("greeter.client"),
	)

	service.Init()

	// 创建客户端
	greeter := pb.NewGreeterService("greeter",service.Client())

	rsp,err := greeter.Hello(context.TODO(),&pb.HelloRequest{Name:"111"})

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(rsp.Greeting)
}
```

#### 运行

由于是本地开发，所以先使用Multicast DNS 的方式做服务发现
后期线上项目可使用 Consul： https://www.consul.io/ 或者 ectd来做服务发现

通过 `--registry=mdns` 或者 环境变量注入 `MICRO_REGISTRY=mdns` 进行设置

更新Makefile文件

```makefile
build:
	protoc -I. --micro_out=. --go_out=. ./proto/greeter.proto

run-server:
	go run main.go --registry=mdns
run-cli:
	go run cli/client.go --registry=mdns
```

打开终端，执行 

```
make run-server
```

可以看到终端执行
```
make run-server
go run main.go --registry=mdns
2018/07/25 14:20:29 Listening on [::]:39869
2018/07/25 14:20:29 Broker Listening on [::]:45203
2018/07/25 14:20:29 Registering node: greeter-d3ff4052-8fd2-11e8-844d-6045cb9fda21
```

然后执行客户端

```
make run-cli
```

可以看到输出

```
make  run-cli   
go run cli/client.go --registry=mdns
hello 111
```

至次，我们通过微服务创建了一个客户端和服务段




