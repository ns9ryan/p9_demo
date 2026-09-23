# node-agent

P9 节点代理服务。

负责连接 node-dispatch、接收调度消息、调用本机业务服务，并将执行结果返回调度中心。

node-agent 不处理具体业务，也不保存具体业务数据。

## 服务职责

- 维护与 node-dispatch 的 WebSocket 长连接。
- 接收异步调度任务和同步节点请求。
- 根据目标服务将请求路由到本机 RPC 服务。
- 将任务确认、执行结果和同步响应返回 node-dispatch。

## 目录结构

```text
node-agent
├── etc                         # 服务运行配置
│   └── node_agent.yaml
│
├── internal
│   ├── config                  # 服务配置结构
│   ├── svc                     # 服务上下文及依赖初始化
│   ├── protocol                # WebSocket 通信协议
│   ├── websocket               # WebSocket 连接与消息收发
│   ├── router                  # 本机服务路由
│   └── client                  # 本机 RPC 服务客户端
│
├── node_agent.go               # node-agent 启动入口
├── go.mod                      # Go Module 定义
└── README.md
```

> 尚未实际使用的目录不提前创建，按对应功能开发进度补充。

## 主要目录说明

- `etc`：node-agent 运行配置。
- `internal/config`：节点编码、认证密钥、调度中心连接等配置结构。
- `internal/svc`：初始化并保存 WebSocket、RPC Client 等公共依赖。
- `internal/protocol`：定义 node-agent 与 node-dispatch 之间的 WebSocket 消息结构。
- `internal/websocket`：负责连接、重连、消息接收和消息发送。
- `internal/router`：根据目标服务将调度请求路由到对应本机服务。
- `internal/client`：封装 operator-base、core 等本机 RPC 客户端。

## 主要文件

- `etc/node_agent.yaml`：node-agent 运行配置。
- `internal/config/config.go`：服务配置结构。
- `internal/svc/service_context.go`：初始化并保存服务依赖。
- `internal/websocket/client.go`：维护与 node-dispatch 的 WebSocket 连接。
- `internal/websocket/message_handler.go`：接收并分发调度中心消息。
- `internal/router/router.go`：根据目标服务执行本机路由。
- `node_agent.go`：node-agent 服务启动入口。

## 服务边界

node-agent 负责：

```text
节点连接
消息收发
本机服务路由
任务结果回传
```

具体 operator、会员、代理、游戏等业务逻辑由对应本机业务服务处理，node-agent 不承载具体业务逻辑。