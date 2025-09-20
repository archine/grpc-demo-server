## gRPC DEMO 服务端

此仓库主要是作为[grpc-demo-proto](https://github.com/archine/grpc-demo-proto)仓库的服务端实现

### 必须

- Go 1.24 or later

### 项目结构

```
grpc-demo-server
├── external/               # 外部依赖层 (如 etcd、mysql等)
├── internal/               # 应用层
│   ├── server/             # gRPC 服务端实现
│   ├── service/            # 逻辑层
│   ├── mapper/             # ORM层
│   ├── entity/             # 模型层
├── base/                   # 基础设施层
│   ├── listener/           # 服务监听器
├── main.go                 # 项目入口
```

### 注意：

每当你在internal/server/目录下新增一个服务实现文件时，``请确保在main.go中正确导入该文件，以便服务能够被注册和使用``。

### 安装

如果选择了带有 ``etcd``的grpc服务注册监听器，请务必启动 etcd 服务

* Docker方式

```shell
docker pull bitnami/etcd:3.6.4

docker run -d --name etcd \
  -p 2379:2379 -p 2380:2380 \
  -e ALLOW_NONE_AUTHENTICATION=yes \
  -e ETCD_ADVERTISE_CLIENT_URLS=http://0.0.0.0:2379 \
  bitnami/etcd:3.6.4
```

数据库

```shell
docker run -d --name mysql -p 3306:3306 -e MYSQL_ROOT_PASSWORD=root mysql:8
```

```sql
create database grpc_demo default charset = utf8mb4;

create table users
(
    id         int auto_increment primary key,
    name       varchar(32)      default ''  not null comment '用户名',
    age        tinyint unsigned default '1' null comment '年龄',
    email      varchar(64)      default ''  not null comment '邮箱',
    created_at int                          not null
)
    comment '用户表';
```