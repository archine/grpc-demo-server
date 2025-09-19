package hello

import (
	"context"

	"github.com/archine/gin-plus/v4/component/ioc"
	"github.com/archine/grpc-demo-proto/hello"
	"google.golang.org/grpc"
)

func init() {
	ioc.RegisterBeanDef(&Server{})
}

type Server struct {
	hello.UnimplementedHelloServiceServer
	server *grpc.Server `autowire:""`
}

func (s *Server) BeanPostConstruct() {
	hello.RegisterHelloServiceServer(s.server, s)
}

func (s *Server) SayHello(ctx context.Context, req *hello.HelloRequest) (*hello.HelloResponse, error) {
	if req.Name == "宋卢生" {
		return &hello.HelloResponse{Message: "SB，宋卢生！"}, nil
	}

	return &hello.HelloResponse{Message: "Hello " + req.Name}, nil
}
