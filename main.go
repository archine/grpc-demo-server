package main

import (
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
