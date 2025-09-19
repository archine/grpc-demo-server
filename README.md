## gRPC DEMO 服务端

此仓库主要是作为[grpc-demo-proto](https://github.com/archine/grpc-demo-proto)仓库的服务端实现

### 必须

- Go 1.24 or later

### 项目结构

```
grpc-demo-server
├── internal/
│   ├── server/             # gRPC 服务端实现
├── listener/               # 项目监听器
├── main.go                 # 项目入口
```

### 注意：
每当你在internal/server/目录下新增一个服务实现文件时，``请确保在main.go中正确导入该文件，以便服务能够被注册和使用``。