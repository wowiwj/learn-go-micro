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