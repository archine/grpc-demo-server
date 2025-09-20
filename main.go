package main

import (
	"grpc-demo-server/base/listener"
	"grpc-demo-server/external/database/mysql"
	_ "grpc-demo-server/internal/server/hello"
	_ "grpc-demo-server/internal/server/user"

	ginplus "github.com/archine/gin-plus/v4"
)

func main() {
	ginplus.New().
		With(
			ginplus.WithEvent(
				mysql.NewStarter(),
				listener.NewGrpcEtcdServerListener(),
				//listener.NewGrpcServerListener(),
			),
		).
		Run(ginplus.ContainerMode)
}
