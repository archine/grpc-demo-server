package listener

import (
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/archine/gin-plus/v4/app"
	"github.com/archine/gin-plus/v4/component/gplog"
	"google.golang.org/grpc"
)

// GrpcServerListener grpc服务启动监听器
// 负责在容器刷新前注册grpc server实例
// 以便grpc controller可以注入使用
// 例子: server/controller/hello/hello_grpc_api.go
type GrpcServerListener struct {
	server *grpc.Server
}

func NewGrpcServerListener() *GrpcServerListener {
	return &GrpcServerListener{}
}

func (g *GrpcServerListener) Order() int {
	return 0
}

func (g *GrpcServerListener) OnContainerRefreshBefore(ctx app.ApplicationContext) {
	g.server = grpc.NewServer()
	ctx.RegisterBean("gRPCServer", g.server)
}

func (g *GrpcServerListener) OnContainerRefreshAfter(ctx app.ApplicationContext) {
	port := ctx.GetConfigProvider().GetString("gin-plus.server.port")
	if port == "" {
		port = "4006"
	}
	appName := ctx.GetConfigProvider().GetString("gin-plus.server.name")
	if appName == "" {
		appName = "gRPC Service"
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		gplog.Fatal(fmt.Sprintf("Starting %s using gRPC on port %s failed: %v", appName, port, err))
	}

	errCh := make(chan error, 1)
	go func() {
		errCh <- g.server.Serve(lis)
	}()

	select {
	case err = <-errCh:
		if err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			gplog.Fatal(fmt.Sprintf("Starting %s using gRPC on port %s failed: %v", appName, port, err))
		}
	case <-time.After(10 * time.Millisecond):
		// wait some time for the server to start
		gplog.Info(fmt.Sprintf("Starting %s using gRPC on %s with PID %d", appName, lis.Addr().String(), os.Getpid()))
	}

	stopSignalCh := make(chan os.Signal, 1)
	signal.Notify(stopSignalCh, syscall.SIGTERM, syscall.SIGINT)
	<-stopSignalCh

	gplog.Info("Received shutdown signal, starting graceful shutdown")
	g.server.GracefulStop()

	gplog.Info("gRPC server shutdown completed successfully")
}
