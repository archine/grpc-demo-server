package listener

import (
	"context"
	"errors"
	"fmt"
	"grpc-demo-server/external/etcd"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/archine/gin-plus/v4/app"
	"github.com/archine/gin-plus/v4/component/gplog"
	"google.golang.org/grpc"
)

type config struct {
	Port      int           `mapstructure:"port"`       // gRPC服务端口
	Name      string        `mapstructure:"name"`       // etcd服务名称
	Addr      string        `mapstructure:"addr"`       // etcd服务地址，必须是能被访问的地址，格式 ip:port
	ExitDelay time.Duration `mapstructure:"exit-delay"` // 服务退出延迟，单位秒，默认1s
}

func (c *config) verify() {
	if c.Port == 0 {
		c.Port = 50051
	}
	if c.Name == "" {
		c.Name = "gRPCService"
	}
	if c.Addr == "" {
		c.Addr = "localhost:50051"
	}
	if c.ExitDelay == 0 {
		c.ExitDelay = time.Second
	}
}

// GrpcEtcdServerListener grpc etcd服务启动监听器
// 负责在容器刷新后注册etcd服务实例
type GrpcEtcdServerListener struct {
	server      *grpc.Server
	cfg         *config
	etcdManager *etcd.Manager
}

func NewGrpcEtcdServerListener() *GrpcEtcdServerListener {
	return &GrpcEtcdServerListener{}
}

func (g *GrpcEtcdServerListener) Order() int {
	return 0
}

func (g *GrpcEtcdServerListener) OnContainerRefreshBefore(ctx app.ApplicationContext) {
	g.server = grpc.NewServer()
	ctx.RegisterBean("gRPCServer", g.server)
}

func (g *GrpcEtcdServerListener) OnContainerRefreshAfter(ctx app.ApplicationContext) {
	var cfg config
	err := ctx.GetConfigProvider().Unmarshal("grpc.server", &cfg)
	if err != nil {
		gplog.Fatal(fmt.Sprintf("Starting gRPC Server failed, unable to parse configuration: %v", err))
	}

	cfg.verify()
	g.cfg = &cfg

	// 1. 启动grpc服务
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

	// 2. 注册etcd服务
	manager, err := etcd.NewManager(ctx.GetConfigProvider())
	if err != nil {
		gplog.Fatal(fmt.Sprintf("gRPC Server register etcd service failed: %v", err))
	}

	err = manager.Register(context.Background(), cfg.Name, cfg.Addr)
	if err != nil {
		gplog.Fatal(fmt.Sprintf("gRPC Server register etcd service failed: %v", err))
	}

	// 3. 监听退出信号，优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	gplog.Info("Received shutdown signal, starting graceful shutdown")

	g.server.GracefulStop()
	manager.Unregister()

	time.Sleep(cfg.ExitDelay)
	gplog.Info("gRPC Server exited properly")
}
