# node-dispatch

P9 调度中心服务。

负责节点管理、operator 与节点部署关系、调度任务管理，以及总网与 node-agent 之间的调度通信。

node-dispatch 不处理具体业务逻辑，具体 operator、会员、代理、游戏等业务由对应业务服务负责。

## 主要职责

- 管理业务节点及节点认证信息。
- 维护 operator 与节点的部署关系。
- 创建和管理调度任务。
- 记录任务执行状态和执行结果。
- 维护与 node-agent 的 WebSocket 长连接。
- 向指定节点下发任务并接收执行结果。
- 为总网提供节点和调度相关 RPC 能力。

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
node-dispatch
├── api
│   └── desc                    # 总网后台API定义
│
├── pkg
│   └── database                # 数据库基础配置
│
├── rpc
│   ├── client                  # goctl生成的RPC客户端
│   ├── ent                     # Ent数据模型及生成代码
│   ├── etc                     # RPC服务配置
│   │
│   ├── internal
│   │   ├── config              # 服务配置
│   │   ├── connection          # 节点在线连接管理
│   │   ├── logic               # RPC接口逻辑
│   │   ├── protocol            # WebSocket通信协议
│   │   ├── server              # gRPC Server
│   │   ├── svc                 # 服务依赖初始化
│   │   ├── task                # 调度任务核心逻辑
│   │   └── websocket           # WebSocket服务与消息处理
│   │
│   ├── pb                      # Proto生成代码
│   ├── proto                   # RPC Proto定义
│   └── node_dispatch.go        # RPC和WebSocket启动入口
│
├── Makefile
├── go.mod
└── README.md
```

## 目录说明

- `api`：总网后台节点和调度管理接口。
- `pkg/database`：数据库连接基础配置。
- `rpc/ent`：节点、operator 节点关系、调度任务和执行记录的数据模型。
- `rpc/internal/connection`：维护当前在线的 node-agent 连接。
- `rpc/internal/logic`：RPC 接口逻辑。
- `rpc/internal/protocol`：node-dispatch 与 node-agent 的 WebSocket 通信协议。
- `rpc/internal/task`：任务创建、下发、状态流转和执行结果处理。
- `rpc/internal/websocket`：WebSocket 连接、认证、心跳和消息处理。
- `rpc/internal/svc`：初始化数据库、连接管理器和任务服务等依赖。
- `rpc/proto`：node-dispatch RPC Proto 定义。

## 服务边界

node-dispatch 负责：

```text
节点管理
operator部署关系
任务调度
任务状态管理
节点通信
执行结果记录
```

具体业务数据和业务规则由对应业务服务处理。