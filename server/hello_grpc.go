package server

import (
	"golang.org/x/net/context"

	pb "gRPC-Gateway/proto"
)

type helloService struct{}

func NewHelloService() *helloService {
	return &helloService{}
}

// helloService.helloService 对应.proto 的 rpc SayHelloWorld
// ctx context.Context 接收上下文参数
// r *pb.HelloWorldRequest 用于接收protobuf的 message HelloWorldRequest 参数
func (h helloService) SayHelloWorld(
	ctx context.Context,
	r *pb.HelloWorldRequest) (*pb.HelloWorldResponse, error) {
	return &pb.HelloWorldResponse{
		Message: r.Referer,
	}, nil
}
