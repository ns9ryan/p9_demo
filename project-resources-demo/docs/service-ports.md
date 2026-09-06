# P9 服务端口

P9 默认开发端口按部署域划分，同一业务域的 API 和 RPC 使用相同尾号。

| 部署域 | API 端口段 | RPC 端口段 |
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
| file | `18006` | `19006` |

## Node Dispatch

| 服务 | API | RPC |
|---|---:|---:|
| node-dispatch | `28001` | `29001` |

## Operator

厅侧当前只预留 `38xxx / 39xxx` 端口段，具体服务及编号在厅侧仓库划分确认后再补充。