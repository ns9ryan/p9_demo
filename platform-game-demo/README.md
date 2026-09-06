# Platform-Game 游戏管理微服务

## 项目概述

**Platform-Game** 是P9总网的**游戏数据管理与同步微服务**，提供完整的游戏基础资料管理和上游游戏数据同步能力。

### 项目职责

✅ **负责**：
- 游戏分类、厂商、渠道等主数据管理
- 游戏基础信息CRUD操作
- 游戏币种(货币)管理
- 与上游游戏系统的数据同步（Preview + Run两阶段）
- 冲突检测与解决
- Kafka事件通知

❌ **不负责**：
- Operator本地游戏运营配置
- 第三方游戏运行时业务
- 游戏运营相关的统计和分析

### 核心特性

🎮 **完整的游戏管理**
- 游戏、分类、厂商、渠道、币种5个核心实体
- 标准的CRUD操作接口

🔄 **智能数据同步**
- 两阶段同步流程（Preview预检查 + Run执行）
- 冲突自动检测（编码重复、ID重复、数据不匹配等）
- 支持部分同步（按ID或编码列表选择性同步）

🌐 **多协议支持**
- REST API (HTTP) - 用于WEB管理平台
- gRPC服务 - 用于内部系统互调

📋 **Simple Admin标准**
- 严格遵循Simple Admin代码组织规范
- 统一的错误码和响应格式
- 完整的字段注释和Swagger支持

⚡ **高性能框架**
- 基于go-zero框架
- 支持自动代码生成
- 内置中间件和性能优化

## 项目结构

```
platform-game/
├── api/                          # REST API层 (HTTP接口)
│   ├── desc/                     # API定义
│   │   └── main.api              # go-zero API语法定义（5大模块 + 同步）
│   ├── platform_game.go          # API服务入口
│   └── internal/
│       ├── handler/              # HTTP请求处理器（auto-gen）
│       ├── logic/                # 业务逻辑层（auto-gen）
│       ├── svc/                  # 服务容器(依赖注入)
│       ├── config/               # 配置加载
│       └── types/                # 请求/响应类型（auto-gen）
│
├── rpc/                          # gRPC服务层
│   ├── proto/
│   │   ├── platform_game.proto   # 主Proto定义（所有服务集合）
│   │   ├── ping.proto            # Ping服务（旧）
│   │   └── desc/                 # 业务Proto参考（备份）
│   │       ├── game.proto        # 游戏服务定义
│   │       ├── category.proto    # 分类服务定义
│   │       ├── provider.proto    # 厂商服务定义
│   │       ├── channel.proto     # 渠道服务定义
│   │       ├── sync.proto        # 同步服务定义
│   │       └── vendor_service.proto # 第三方集成定义
│   ├── platform_game.go          # gRPC服务入口
│   ├── pb/                       # 生成的Proto代码（auto-gen）
│   │   └── platform_game/        # platform_game包的编译产物
│   │       ├── platform_game.pb.go
│   │       └── platform_game_grpc.pb.go
│   ├── platformgameservice/      # gRPC服务实现（auto-gen）
│   ├── platformgame/             # gRPC客户端（auto-gen）
│   └── internal/
│       ├── server/               # gRPC服务实现（auto-gen）
│       ├── logic/                # 业务逻辑层（auto-gen）
│       ├── svc/                  # 服务容器
│       └── config/               # 配置加载
│
├── common/                       # 公共模块 (核心)
│   ├── constant/                 # 常量定义
│   │   ├── error_code.go         # 错误码（1~30000+）
│   │   ├── business.go           # 业务常量（状态、同步操作、对象类型等）
│   │   └── messages.go           # 错误消息文本
│   ├── model/                    # 数据模型
│   │   ├── game.go               # 5大实体模型（Game/Category/Provider/Channel/Currency）
│   │   └── sync.go               # 同步相关模型（SyncDiff/SyncStats等）
│   └── response/                 # 响应处理
│       └── response.go           # 统一响应结构 {Code/Message/Data}
│
├── pkg/                          # 业务包
│   ├── game/                     # 游戏业务服务
│   │   └── game_service.go       # 游戏CRUD和状态管理
│   └── sync/                     # 同步业务服务
│       └── sync_service.go       # Preview和Run实现
│
├── Makefile                      # 构建配置（make api/make rpc生成代码）
├── go.mod                        # Go模块定义
└── README.md                     # 本文件
```

## 业务模型

### 核心5大实体

#### Level 0 - 主数据（主数据库，无依赖）

| 实体 | 表名 | 说明 |
|------|------|------|
| **Category** | game_category | 游戏分类（棋牌、电子游戏等） |
| **Provider** | game_provider | 游戏厂商（系统商、游戏提供商等） |
| **Channel** | game_channel | 游戏渠道（运营渠道、分发渠道等） |

#### Level 1 - 业务数据（依赖L0）

| 实体 | 表名 | 说明 |
|------|------|------|
| **Game** | game | 游戏基础信息，关联 Category × Provider × Channel |

#### Level 2 - 附属数据（依赖L1）

| 实体 | 表名 | 说明 |
|------|------|------|
| **Currency** | game_currency | 游戏支持的货币列表，关联 Game |

### 数据同步依赖关系

```
时间 ──────────────────────────────────────────────>

  同步L0主数据（并行同步）
  Category + Provider + Channel
            ↓
         完成L0同步
            ↓
       同步L1游戏数据
         (Game)
            ↓
       同步L2货币数据
      (Currency)
```

**关键约束**：
- L0三个实体独立，可并行同步
- L1(Game)依赖L0全部完成
- L2(Currency)依赖L1完成

## API 接口大全

### 📌 基础接口

#### Ping
```bash
GET /ping
```

### 🎮 游戏管理接口 (5个操作)

| 操作 | 方法 | 路由 | 请求体 | 返回 |
|------|------|------|--------|------|
| 创建 | POST | `/games` | GameCreateReq | GameResp |
| 获取 | GET | `/games/:id` | - | GameResp |
| 更新 | PUT | `/games/:id` | GameResp | GameResp |
| 删除 | DELETE | `/games/:id` | - | Response |
| 列表 | GET | `/games?page=1&page_size=10` | GameListReq | GameListResp |

**请求示例**：
```bash
# 创建游戏
curl -X POST http://localhost:8080/games \
  -H "Content-Type: application/json" \
  -d '{
    "code": "G001",
    "source_code": "SRC_G001",
    "name": "热血传奇",
    "category_id": 1,
    "provider_id": 2,
    "channel_id": 3,
    "logo": "https://...",
    "status": 1,
    "sort_no": 10
  }'

# 游戏列表
curl "http://localhost:8080/games?page=1&page_size=20&category_id=1"
```

### 📂 分类管理接口 (5个操作)

| 操作 | 方法 | 路由 |
|------|------|------|
| 创建 | POST | `/categories` |
| 获取 | GET | `/categories/:id` |
| 更新 | PUT | `/categories/:id` |
| 删除 | DELETE | `/categories/:id` |
| 列表 | GET | `/categories?page=1&page_size=10` |

### 🏢 厂商管理接口 (5个操作)

| 操作 | 方法 | 路由 |
|------|------|------|
| 创建 | POST | `/providers` |
| 获取 | GET | `/providers/:id` |
| 更新 | PUT | `/providers/:id` |
| 删除 | DELETE | `/providers/:id` |
| 列表 | GET | `/providers?page=1&page_size=10` |

### 🚚 渠道管理接口 (5个操作)

| 操作 | 方法 | 路由 |
|------|------|------|
| 创建 | POST | `/channels` |
| 获取 | GET | `/channels/:id` |
| 更新 | PUT | `/channels/:id` |
| 删除 | DELETE | `/channels/:id` |
| 列表 | GET | `/channels?page=1&page_size=10` |

### 🔄 数据同步接口 (2个操作)

#### 1️⃣ 同步预检查 (Preview)
```bash
POST /sync/preview

请求:
{
  "object_type": "category",  // category|provider|channel|game
  "selected_ids": [1, 2, 3],  // 可选：限定同步范围
  "selected_codes": [],        // 可选：按编码过滤
  "strict_conflict": false     // 可选：冲突是否报错
}

响应:
{
  "code": 1,
  "message": "success",
  "data": {
    "stats": {
      "remote_total": 100,     // 上游共100条
      "local_total": 95,       // 本地共95条
      "create_total": 5,       // 需要新增5条
      "update_total": 10,      // 需要更新10条
      "delete_total": 0,       // 需要删除0条
      "noop_total": 80,        // 80条无变化
      "conflict_total": 2,     // 2条冲突
      "error_total": 0         // 0条异常
    },
    "diffs": [
      {
        "object_type": "game",
        "object_id": 0,
        "object_code": "game_new_001",
        "action": "create",
        "remote_id": 500,
        "remote_code": "G_500",
        "reason": "新增游戏"
      },
      {
        "object_type": "game",
        "object_id": 5,
        "object_code": "game_001",
        "action": "update",
        "remote_id": 100,
        "remote_code": "G_100",
        "reason": "名称已变更"
      },
      {
        "object_type": "game",
        "object_id": 8,
        "object_code": "game_dup",
        "action": "conflict",
        "conflict_type": "code_duplicate",
        "reason": "编码与远端冲突"
      }
    ]
  }
}
```

#### 2️⃣ 执行同步 (Run)
```bash
POST /sync/run

请求:
{
  "object_type": "category",
  "selected_ids": [],
  "selected_codes": [],
  "auto_apply": true,          // true=真实执行，false=只预检查
  "strict_conflict": false     // true=遇冲突报错
}

响应:
{
  "code": 1,
  "message": "success",
  "data": {
    "preview": { /* 同预检查的stats和diffs */ },
    "apply": {
      "created": 5,            // 成功新增5条
      "updated": 10,           // 成功更新10条
      "deleted": 0,            // 成功删除0条
      "failed": 2,             // 失败2条
      "skipped": 2             // 跳过2条（冲突）
    }
  }
}
```

## 错误码体系

### 错误码分布

| 范围 | 类型 | 描述 |
|------|------|------|
| 1 | 成功 | 正常响应 |
| 10001-10099 | 业务错误 | 游戏/分类/厂商/渠道/币种相关 |
| 10051-10053 | 同步错误 | 同步冲突、预检查失败、执行失败 |
| 20001-20006 | 系统错误 | 数据库、网络、认证等 |
| 30001-30003 | 验证错误 | 参数验证失败 |

### 详细错误码

| 错误码 | 名称 | 含义 | 解决方案 |
|--------|------|------|---------|
| 1 | SUCCESS | 成功 | - |
| 10001 | GameNotFound | 游戏不存在 | 检查游戏ID |
| 10002 | GameAlreadyExists | 游戏已存在 | 使用不同的编码 |
| 10003 | InvalidGameStatus | 无效的游戏状态 | 状态必须为1(启用)或2(禁用) |
| 10004 | GameSyncFailed | 游戏同步失败 | 查看日志获取详细信息 |
| 10011 | CategoryNotFound | 分类不存在 | 检查分类ID |
| 10012 | CategoryAlreadyExists | 分类已存在 | 使用不同的编码 |
| 10021 | ProviderNotFound | 厂商不存在 | 检查厂商ID |
| 10022 | ProviderAlreadyExists | 厂商已存在 | 使用不同的编码 |
| 10031 | ChannelNotFound | 渠道不存在 | 检查渠道ID |
| 10032 | ChannelAlreadyExists | 渠道已存在 | 使用不同的编码 |
| 10041 | CurrencyNotFound | 币种不存在 | 检查币种ID |
| 10051 | SyncConflict | 数据冲突 | 使用strict_conflict=false跳过 |
| 10052 | SyncPreviewFailed | 同步预检查失败 | 检查系统连接 |
| 10053 | SyncExecuteFailed | 同步执行失败 | 检查auto_apply参数 |
| 20001 | DatabaseError | 数据库错误 | 检查数据库连接 |
| 20005 | InternalError | 系统内部错误 | 查看服务日志 |
| 30001 | ValidationFailed | 数据验证失败 | 检查字段类型和长度 |

## 开发指南

### 环境要求

```
Go             1.18+
Goctl          1.10.2+
Protocol Buffers  3.0+
MySQL          5.7+
```

### 快速开始

```bash
# 1. 克隆项目
git clone https://oa.98ent.com/p9/platform-game.git
cd platform-game

# 2. 下载依赖
go mod download

# 3. 生成所有代码 (API + RPC)
make gen

# 4. 查看生成结果
ls -la api/internal/handler/
ls -la rpc/internal/server/
```

### 代码生成工作流

#### 修改API定义
```bash
# 1. 编辑 api/desc/main.api
#    - 添加新的type定义
#    - 添加新的service路由
#    - 添加@handler注解

# 2. 生成新的Handler和Logic
make api

# 3. goctl会生成：
#    api/internal/handler/*.go     (新的handler)
#    api/internal/logic/*.go       (新的logic)
#    api/internal/types/types.go   (新的types)

# 4. 实现logic中的业务逻辑
```

#### 修改RPC定义
```bash
# 1. 编辑 rpc/proto/platform_game.proto
#    - 在对应的业务分区添加 message 定义
#    - 在 Service 中新增 rpc 方法
#    （参考 rpc/proto/desc/ 中的原型定义）

# 2. 生成新的Proto代码
make rpc

# 3. goctl会生成：
#    rpc/pb/platform_game/*.go      (proto编译产物)
#    rpc/internal/server/*.go       (gRPC服务实现)
#    rpc/platformgameservice/*.go   (gRPC服务骨架)

# 4. 在 rpc/internal/logic/ 中实现业务逻辑
```

### 项目分层架构

```
┌─────────────────────────────────────────┐
│         HTTP Client / gRPC Client       │
└──────────────────┬──────────────────────┘
                   ↓
┌─────────────────────────────────────────┐
│   Handler/Server (路由、验证、序列化)     │
│   ❓ 问题：请求合法吗？参数验证成功吗？   │
└──────────────────┬──────────────────────┘
                   ↓
┌─────────────────────────────────────────┐
│   Logic (业务逻辑、协调、决策)            │
│   ❓ 问题：这个请求应该怎么处理？         │
└──────────────────┬──────────────────────┘
                   ↓
┌─────────────────────────────────────────┐
│   Repository (数据访问层、缓存)          │
│   ❓ 问题：数据在哪里？怎么读写？         │
└──────────────────┬──────────────────────┘
                   ↓
┌─────────────────────────────────────────┐
│   Database/Cache/Kafka (基础设施)       │
│   MySQL / Redis / Kafka                │
└─────────────────────────────────────────┘
```

### 代码规范 (Simple Admin标准)

#### 类型命名约定

| 类型后缀 | 用途 | 示例 |
|---------|------|------|
| `Req` | 请求类型 | GameCreateReq, GameListReq |
| `Resp` | 响应类型 | GameResp, CategoryResp |
| `Info` | 列表项类型 | GameListInfo |
| 无后缀 | 数据模型 | Game, Category |

```go
// ❌ 错误命名
type GameRequest struct { }     // 不规范
type Game_Resp struct { }       // 不规范
type GetGameResponse struct { } // 不规范

// ✅ 正确命名
type GameCreateReq struct { }
type GameUpdateReq struct { }
type GameListReq struct { }
type GameResp struct { }
type GameListInfo struct { }
```

#### 字段注释规范

所有结构体字段必须有注释（用于Swagger文档生成）：

```go
type Game struct {
	// ❌ 错误：无注释
	ID   int64
	Code string
	
	// ✅ 正确：完整注释
	ID   int64  `json:"id" comment:"游戏ID"`
	Code string `json:"code" binding:"required" comment:"游戏编码"`
	Name string `json:"name" binding:"required" comment:"游戏名称"`
	
	// 可选字段示例
	Logo   string `json:"logo" comment:"游戏logo"`
	Status int    `json:"status" comment:"状态：1=启用，2=禁用"`
}
```

#### 文件命名规范

使用 snake_case 风格：

```
game_handler.go        ✅ 正确
game_logic.go          ✅ 正确
game_types.go          ✅ 正确
game_repository.go     ✅ 正确
gameHandler.go         ❌ 错误
GameLogic.go           ❌ 错误
game-types.go          ❌ 错误
```

### 常见开发任务

#### 任务1：添加新的游戏管理端点

```bash
# Step 1: 更新API定义
vi api/desc/main.api
# 添加新的type和@handler

# Step 2: 生成代码
make api

# Step 3: 实现逻辑
vi api/internal/logic/{endpoint}_logic.go
# 在Logic.{Endpoint}() 中实现业务逻辑

# Step 4: 调用外层服务(可选)
# 在logic中调用 l.svcCtx.GameService.XXX()
```

#### 任务2：添加新的gRPC服务

```bash
# Step 1: 创建proto文件
vi rpc/proto/desc/your_service.proto

# Step 2: 生成代码
make rpc

# Step 3: 实现服务
vi rpc/internal/server/yourservice/your_service_server.go
# 在 YourService.YourMethod() 中实现业务逻辑

# Step 4: 注册服务（platform_game.go）
# 将YourService注册到gRPC Server
```

## 同步流程详解

### 整体流程

```
开始同步
  ↓
─────────────────────────────────
│ Phase 1: Preview (查看变化)    │
│ ✓ 获取远程数据                 │
│ ✓ 获取本地数据                 │
│ ✓ 比较差异                     │
│ ✗ 不修改任何数据               │
│ 返回：需要做什么操作            │
─────────────────────────────────
  ↓ (用户确认)
─────────────────────────────────
│ Phase 2: Run (执行变化)        │
│ ✓ 再次Preview（确保最新）       │
│ ✓ 检查冲突                     │
│ ✓ 开始事务                     │
│ ✓ 执行Create操作               │
│ ✓ 执行Update操作               │
│ ✓ 执行Delete操作               │
│ ✓ 提交事务                     │
│ 返回：成功做了什么操作          │
─────────────────────────────────
  ↓
同步完成
```

### 冲突检测详解

| 冲突场景 | 检测方法 | 影响 | 解决方案 |
|---------|--------|------|---------|
| **编码重复** (code_duplicate) | 同一对象类型中编码重复 | 无法创建 | 检查编码唯一性 |
| **ID重复** (id_duplicate) | 上游ID在多条本地记录中 | 无法更新 | 手动合并记录 |
| **数据不匹配** (data_mismatch) | 编码相同但其他字段不同 | 需要人工确认 | 决定是否覆盖 |
| **状态不匹配** (status_mismatch) | 状态字段值不同 | 可能覆盖配置 | 审查状态变更 |

**处理冲突的参数**：
- `strict_conflict=true`: 遇冲突直接报错，不执行同步
- `strict_conflict=false`: 跳过冲突，继续同步其他数据

### Kafka事件

同步完成后会发送事件通知其他系统：

```json
{
  "event": "game.force-quit",
  "scope": "game",
  "game_ids": [1, 2, 3, 5, 8],
  "timestamp": 1704067200,
  "sync_phase": "run",
  "created": 5,
  "updated": 10,
  "deleted": 0
}
```

## 贡献指南

### 开发流程

1. **创建分支** - 从main创建feature分支
   ```bash
   git checkout -b feature/description
   ```

2. **开发功能** - 按照项目规范编写代码
   ```bash
   # 修改代码
   git add .
   git commit -m "feat: description"
   ```

3. **生成代码** - 如果修改了API或Proto
   ```bash
   make gen
   ```

4. **测试** - 运行单元测试
   ```bash
   go test ./...
   ```

5. **提交** - 推送并创建PR
   ```bash
   git push origin feature/description
   # 在GitHub/GitLab创建Pull Request
   ```

## 参考文档

📚 **业务文档**：
- [V2游戏同步速查表](./V2_GAME_SYNC_CHEATSHEET.md) - 快速参考
- [V2游戏同步综合指南](./V2_GAME_SYNC_COMPREHENSIVE.md) - 完整说明
- [V2接口迁移指南](./V2_INTERFACE_MIGRATION_GUIDE.md) - 10个迁移接口详解

📖 **外部参考**：
- [Go-Zero文档](https://go-zero.dev)
- [Protocol Buffers文档](https://developers.google.com/protocol-buffers)
- [Simple Admin规范](https://doc.renzhouu.top/)

## 许可证

Copyright (c) 98ent - All rights reserved.

内部项目，未经授权禁止使用。