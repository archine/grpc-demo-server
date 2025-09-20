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
	cfg    *config
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
	var cfg config
	err := ctx.GetConfigProvider().Unmarshal("grpc.server", &cfg)
	if err != nil {
		gplog.Fatal(fmt.Sprintf("Starting gRPC Server failed, unable to parse configuration: %v", err))
	}

	cfg.verify()
	g.cfg = &cfg

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		gplog.Fatal(fmt.Sprintf("Starting gRPC Server on port %d failed: %v", cfg.Port, err))
	}
	gplog.Info(fmt.Sprintf("gRPC Server listening on %d", cfg.Port))

	errCh := make(chan error, 1)
	go func() {
		errCh <- g.server.Serve(lis)
	}()

	select {
	case err = <-errCh:
		if err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			gplog.Fatal(fmt.Sprintf("Starting gRPC Server failed: %v", err))
		}
	case <-time.After(10 * time.Millisecond):
		// 等待一点时间，确保服务启动成功
		gplog.Info("Starting gRPC Server succeeded")
	}

	stopSignalCh := make(chan os.Signal, 1)
	signal.Notify(stopSignalCh, syscall.SIGTERM, syscall.SIGINT)
	<-stopSignalCh

	gplog.Info("Received shutdown signal, starting graceful shutdown")
	g.server.GracefulStop()

	gplog.Info("gRPC Server exited properly")
}
