# node-agent

P9 节点代理服务。

部署在业务节点上，负责连接 node-dispatch、接收调度任务、调用当前节点业务 RPC 服务，并将执行结果返回调度中心。

node-agent 不保存业务数据，也不实现具体业务规则。

## 主要职责

- 维护与 node-dispatch 的 WebSocket 长连接。
- 接收并处理调度中心下发的任务。
- 调用总网或当前节点业务 RPC 服务。
- 返回任务接收确认和执行结果。
- 负责节点侧任务执行编排，不承载具体业务逻辑。

## 调用关系

```text
Platform
    ↓ RPC
node-dispatch
    ↓ WebSocket
node-agent
    ↓ RPC
当前节点业务服务
```

## 目录结构

```text
node-agent
├── etc
│   └── node_agent.yaml            # 服务配置
│
├── internal
│   ├── config                     # 配置结构
│   ├── protocol                   # WebSocket通信协议
│   ├── svc                        # 服务依赖初始化
│   ├── task                       # 调度任务执行
│   └── websocket                  # WebSocket连接与消息处理
│
├── node_agent.go                  # 服务启动入口
├── go.mod
└── README.md
```

## 目录说明

- `etc`：服务运行配置。
- `internal/config`：节点、WebSocket、RPC 等配置结构。
- `internal/protocol`：node-agent 与 node-dispatch 的通信协议。
- `internal/svc`：初始化 WebSocket、RPC Client、任务服务等依赖。
- `internal/task`：任务识别、参数解析和执行编排。
- `internal/websocket`：WebSocket 连接、重连、消息收发和消息处理。

## 服务边界

node-agent 负责：

```text
节点通信
任务接收
任务执行编排
RPC调用
结果回传
```

具体业务数据和业务规则由对应业务服务负责，例如 operator 相关业务由 `operator-base` 处理。