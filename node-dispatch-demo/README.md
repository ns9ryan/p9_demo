# node-dispatch

P9 调度中心服务。

负责节点管理、operator 与节点部署关系、调度任务、任务执行记录，以及总网与节点之间的调度通信。

## 服务组成

- `API`：提供给总网前端使用，主要用于节点管理和调度任务管理。
- `RPC`：提供节点、部署关系、任务调度等内部服务能力。
- `WebSocket`：维护与 node-agent 的长连接，负责调度消息下发和节点消息接收。

## 目录结构

```text
node-dispatch
├── api                         # 总网后台 HTTP API
│   └── desc                    # API DSL 定义
│
├── pkg                         # 与具体服务入口解耦的公共基础代码
│   └── database                # 数据库连接配置
│
├── rpc                         # 调度中心核心 RPC 服务
│   ├── client                  # goctl 生成的 RPC 客户端
│   ├── ent                     # Ent 数据模型及生成代码
│   ├── etc                     # RPC 服务配置
│   │
│   ├── internal
│   │   ├── config              # 服务配置结构
│   │   ├── logic               # RPC 业务逻辑
│   │   ├── server              # gRPC Server 适配层
│   │   ├── svc                 # 服务上下文及依赖初始化
│   │   ├── connection          # 节点在线连接管理
│   │   ├── protocol            # WebSocket 通信协议
│   │   ├── task                # 调度任务核心逻辑
│   │   ├── request             # 同步节点请求管理
│   │   └── websocket           # WebSocket 服务
│   │
│   ├── pb                      # Proto 生成代码
│   ├── proto                   # RPC Proto 定义
│   └── node_dispatch.go        # RPC 与 WebSocket 启动入口
│
├── Makefile                    # 代码生成与构建命令
├── go.mod                      # Go Module 定义
└── README.md
```

> 尚未实际使用的目录不提前创建，按对应功能开发进度补充。

## 主要目录说明

- `api`：总网后台接口，节点管理、任务管理等前端功能从这里进入。
- `rpc`：调度中心核心服务。
- `rpc/ent`：节点、operator 节点关系、调度任务及执行记录的数据模型。
- `rpc/internal/logic`：RPC 接口逻辑，具体业务按 service 目录划分。
- `rpc/internal/connection`：维护当前在线的 node-agent 连接。
- `rpc/internal/protocol`：定义 node-dispatch 与 node-agent 之间的 WebSocket 消息结构。
- `rpc/internal/task`：负责任务创建、状态流转、执行结果等调度逻辑。
- `rpc/internal/request`：负责需要同步等待节点响应的请求。
- `rpc/internal/websocket`：负责 WebSocket 连接、认证、心跳和消息处理。
- `pkg`：保存与 API、RPC 服务入口解耦的可复用基础代码。

## 主要文件

- `api/desc/main.api`：API DSL 入口。
- `rpc/proto/node_dispatch.proto`：RPC Proto 入口。
- `rpc/internal/svc/service_context.go`：初始化并保存数据库、连接管理等服务依赖。
- `rpc/internal/connection/manager.go`：管理 node-agent 在线连接。
- `rpc/internal/websocket/server.go`：WebSocket 服务及连接生命周期。
- `rpc/internal/websocket/message_handler.go`：接收并分发节点消息。
- `rpc/internal/task/service.go`：调度任务核心服务。
- `rpc/node_dispatch.go`：node-dispatch 服务启动入口。

## 服务边界

node-dispatch 负责：

```text
节点管理
operator 与节点部署关系
任务调度
任务执行状态
节点通信
```

具体 operator、会员、代理、游戏等业务由对应业务服务处理，node-dispatch 不承载具体业务逻辑。