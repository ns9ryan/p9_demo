# P9 服务端口

P9 默认开发端口按部署域划分。

同一服务的 API 和 RPC 原则上使用相同尾号，其他通信端口按服务实际需要单独分配。

## 端口段

| 部署域 | HTTP / WebSocket 端口段 | RPC 端口段 |
|---|---:|---:|
| Platform | `18000-18999` | `19000-19999` |
| Node Dispatch | `28000-28999` | `29000-29999` |
| Operator | `38000-38999` | `39000-39999` |

## Platform

| 服务 | API | RPC |
|---|---:|---:|
| core | `18000` | `19000` |
| platform-base | `18001` | `19001` |
| platform-operator | `18002` | `19002` |
| platform-game | `18003` | `19003` |
| platform-message | `18004` | `19004` |
| integration | `18005` | `19005` |
| oss | `18006` | `19006` |

## Node Dispatch

| 服务 | API | WebSocket | RPC |
|---|---:|---:|---:|
| node-dispatch | `28001` | `28002` | `29001` |

说明：

- `28001` 用于总网后台访问 node-dispatch HTTP API
- `28002` 用于 node-agent 与 node-dispatch 建立 WebSocket 长连接
- `29001` 用于 platform 等内部服务调用 node-dispatch RPC

## Operator

厅侧当前预留 `38xxx / 39xxx` 端口段。

`operator-base` 等厅侧服务的具体端口，在对应服务运行方式和接口边界确认后再补充。

`node-agent` 主动连接 node-dispatch WebSocket，当前不分配固定对外监听端口。