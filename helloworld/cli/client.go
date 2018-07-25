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
