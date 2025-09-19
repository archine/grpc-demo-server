package main

import (
	_ "grpc-demo-server/internal/server/hello"
	"grpc-demo-server/listener"

	ginplus "github.com/archine/gin-plus/v4"
)

func main() {
	ginplus.New().
		With(
			ginplus.WithEvent(
				listener.NewGrpcServerListener(),
			),
		).
		Run(ginplus.ContainerMode)
}
