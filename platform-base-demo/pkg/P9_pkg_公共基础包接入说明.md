# P9 公共基础包接入说明

## 1. 文档目的

`pkg/` 用于保存 P9 各仓库可复用的基础能力。

当前阶段暂不建立独立公共包仓库，因此每个业务仓库可以各自保存一份 `pkg/`。这些包在设计时需要保持低耦合，后续如果建立公共包仓库，应尽量做到直接迁移目录并修改 import path，而不需要重新拆分代码。

这份文档用于说明：

- `pkg/` 各目录的职责
- API 服务需要接入哪些能力
- RPC 服务需要接入哪些能力
- API / RPC 需要怎样配置和初始化
- 统一响应、参数校验、国际化、RPC 错误处理如何串起来
- 业务 Logic 应该怎样使用这些公共能力
- 接入完成后如何验证

> 示例中的 `<module>` 表示当前仓库的 Go Module，例如：
>
> `oa.98ent.com/p9/platform-base`

---

## 2. 目录结构

```text
pkg/
├─ api/                     # HTTP API 专属公共能力
│  ├─ errorhandler/         # API 统一错误处理
│  ├─ middleware/           # API 公共中间件
│  ├─ response/             # API 统一响应
│  ├─ rpcerror/             # API 调用 RPC 时的错误包装
│  └─ validate/             # API 参数校验
│
├─ rpc/                     # gRPC / RPC 专属公共能力
│  └─ grpcerror/            # RPC 统一 gRPC 错误
│
├─ database/                # 数据库通用基础能力
├─ i18n/                    # 通用国际化能力
└─ i18nkey/                 # API / RPC 共享错误 Key
```

### 2.1 目录归属原则

`pkg/api/` 中的代码离开 HTTP API 后基本没有独立使用价值，因此归 API 公共能力。

`pkg/rpc/` 中的代码依赖 gRPC / RPC 语义，因此归 RPC 公共能力。

`database`、`i18n`、`i18nkey` 不属于某一种服务类型，继续保留在 `pkg/` 根目录。

### 2.2 依赖方向

统一遵守：

```text
api/ -> pkg/
rpc/ -> pkg/
```

禁止：

```text
pkg/ -> api/
pkg/ -> rpc/
pkg/ -> api/internal/
pkg/ -> rpc/internal/
```

`pkg/` 不能依赖具体业务模块、生成代码或某个服务的 `internal` 包。

---

# 3. 各公共包职责

## 3.1 `pkg/api/response`

负责 API 统一成功和错误响应。

成功响应统一：

```json
{
  "code": 0,
  "msg": "ok",
  "data": {}
}
```

错误响应统一：

```json
{
  "code": 404,
  "msg": "数据不存在"
}
```

开发或测试环境可以额外返回：

```json
{
  "code": 500,
  "msg": "系统内部错误",
  "debug": {
    "source": "rpc",
    "rpc": "platform_base.LanguageService/Create",
    "error": "rpc error: ..."
  }
}
```

`debug` 只能用于开发和测试环境，生产环境不得返回内部错误详情。

---

## 3.2 `pkg/api/validate`

负责接入 go-zero 全局参数校验。

主要能力：

- 使用 `validate` tag 校验请求参数
- 根据 `X-Lang` 返回中文或英文校验信息
- 将参数校验错误包装为 `validate.Error`
- 供统一 `errorhandler` 判断错误类型

例如：

```go
Page int64 `form:"page" validate:"required,gte=1"`
```

请求：

```text
?page=0
```

会在进入业务 Logic 前被统一拦截。

当前参数校验翻译支持：

```text
zh / zh-CN
en / en-US
```

如果以后需要增加其他校验语言，需要在 `pkg/api/validate` 中补对应 translator 和语言映射。

---

## 3.3 `pkg/api/middleware`

目前主要提供语言中间件。

统一从请求头读取：

```text
X-Lang
```

并通过：

```go
i18n.WithLanguage(...)
```

写入当前请求 `context`。

中间件本身不决定默认语言，默认语言统一由 `pkg/i18n.Translator` 处理。

---

## 3.4 `pkg/api/rpcerror`

用于 API 作为 gRPC Client 调用 RPC 时包装错误。

主要解决两个问题：

1. 保存具体失败的 RPC 方法
2. 保留原始 gRPC Status，不破坏 `codes.NotFound`、`codes.Internal` 等状态识别

例如原始 RPC 方法：

```text
/platform_base.LanguageService/Create
```

包装后可以保存：

```text
platform_base.LanguageService/Create
```

后续 `errorhandler` 可以在调试信息和日志中准确定位失败的 RPC 方法。

> 一个底层 `zrpc.Client` 只需要注册一次 `rpcerror.UnaryClientInterceptor`。

---

## 3.5 `pkg/api/errorhandler`

负责 API 全局错误处理。

主要处理：

- API 参数校验错误
- JSON / Form / Path 等请求解析错误
- RPC gRPC 错误
- gRPC Code -> HTTP Status
- i18n Key -> 用户可读文案
- 5xx 内部错误日志
- 开发 / 测试环境 debug 信息

当前主要映射关系：

| gRPC Code | HTTP Status |
|---|---:|
| InvalidArgument | 400 |
| FailedPrecondition | 400 |
| OutOfRange | 400 |
| Unauthenticated | 401 |
| PermissionDenied | 403 |
| NotFound | 404 |
| AlreadyExists | 409 |
| Aborted | 409 |
| ResourceExhausted | 429 |
| Canceled | 408 |
| Unimplemented | 501 |
| Unavailable | 503 |
| DeadlineExceeded | 504 |
| 其他 | 500 |

---

## 3.6 `pkg/rpc/grpcerror`

供 RPC Logic 创建统一 gRPC 错误。

业务 Logic 不直接返回 HTTP 错误，也不返回已经翻译好的中文或英文文案。

推荐：

```go
return nil, grpcerror.InvalidArgument(i18nkey.ValidationError)
```

```go
return nil, grpcerror.NotFound(i18nkey.DataNotFound)
```

RPC 返回：

```text
gRPC Code + i18n Key
```

API 再负责：

```text
gRPC Code
-> HTTP Status

i18n Key
-> 当前语言文案
```

---

## 3.7 `pkg/i18n`

负责通用国际化能力。

主要包括：

- `WithLanguage`
- `LanguageFromContext`
- `Translator`
- 加载 JSON 语言资源
- 默认语言回退
- 根据 context 翻译 i18n Key

它不属于 API 或 RPC 专属能力，因此保留在 `pkg/i18n`。

---

## 3.8 `pkg/i18nkey`

统一定义跨 API / RPC 使用的错误 Key。

例如：

```go
const (
    InvalidRequest  = "common.invalid_request"
    DataNotFound    = "common.data_not_found"
    ConstraintError = "common.constraint_error"
    ValidationError = "common.validation_error"
    DatabaseError   = "common.database_error"
)
```

RPC 使用这些 Key 返回错误，API 使用相同 Key 翻译，因此不能把它放到 API 或 RPC 单侧目录。

---

## 3.9 `pkg/database`

负责数据库通用基础能力。

例如：

- PostgreSQL 配置
- DSN
- `sql.DB` 连接池
- Ping
- Ent SQL Driver

目前 P9 的正常调用链是：

```text
API
-> RPC
-> Ent
-> PostgreSQL
```

因此当前通常只有 RPC 会使用 `pkg/database`。

但 `database` 本身不依赖 RPC，未来 Job、Consumer、迁移工具等也可能使用，所以继续保留在：

```text
pkg/database/
```

不要移动到：

```text
pkg/rpc/database/
```

---

# 4. API 接入步骤

下面以一个同时包含 API 和 RPC 的仓库为例。

---

## 第一步：复制 `pkg/`

将公共 `pkg/` 放到仓库根目录：

```text
<repo>/
├─ api/
├─ rpc/
├─ pkg/
└─ go.mod
```

复制后先执行：

```bash
go mod tidy
```

让 Go 自动整理公共包需要的依赖。

不要直接复制其他仓库的业务代码或旧 Module import。

---

## 第二步：准备 API 多语言资源

语言资源属于具体 API 应用，不放到 `pkg/`。

建议：

```text
api/internal/locales/
├─ embed.go
└─ locale/
   ├─ zh-CN.json
   └─ en-US.json
```

`embed.go`：

```go
package locales

import "embed"

// FS API多语言资源
//
//go:embed locale/*.json
var FS embed.FS
```

中文示例：

```json
{
  "common": {
    "invalid_request": "请求参数格式错误",
    "data_not_found": "数据不存在",
    "constraint_error": "数据约束冲突",
    "validation_error": "数据校验失败",
    "database_error": "数据库操作失败",
    "internal_error": "系统内部错误",
    "too_many_requests": "请求过于频繁",
    "service_unavailable": "服务暂不可用",
    "request_timeout": "请求超时"
  }
}
```

英文示例：

```json
{
  "common": {
    "invalid_request": "Invalid request",
    "data_not_found": "Data not found",
    "constraint_error": "Data constraint conflict",
    "validation_error": "Validation failed",
    "database_error": "Database operation failed",
    "internal_error": "Internal server error",
    "too_many_requests": "Too many requests",
    "service_unavailable": "Service unavailable",
    "request_timeout": "Request timeout"
  }
}
```

如果某种语言没有 API 文案资源，不需要创建空文件，Translator 会回退到默认语言。

---

## 第三步：增加 API Config

示例：

```go
package config

import (
    "<module>/pkg/i18n"

    "github.com/zeromicro/go-zero/rest"
    "github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
    rest.RestConf

    // RPC配置
    XxxRpc zrpc.RpcClientConf

    // 国际化配置
    I18n i18n.Config
}
```

`XxxRpc` 按当前仓库实际 RPC 名称修改。

例如 `platform-base`：

```go
PlatformBaseRpc zrpc.RpcClientConf
```

---

## 第四步：增加 API YAML 配置

示例：

```yaml
# API服务配置
Name: xxx-api
Host: 0.0.0.0
Port: 18001
Mode: dev

# RPC配置
XxxRpc:
  Etcd:
    Hosts:
      - 127.0.0.1:12379
    Key: xxx-rpc

# 国际化配置
I18n:
  DefaultLanguage: zh-CN
```

不要再额外增加一套 `Debug: true/false` 配置。

是否返回 debug 信息直接根据 go-zero `Mode` 判断：

```text
dev  -> 返回
test -> 返回
pro  -> 不返回
```

---

## 第五步：初始化 Translator 和 Language Middleware

在 API `ServiceContext` 中增加：

```go
Trans *i18n.Translator
Language rest.Middleware
```

初始化：

```go
trans, err := i18n.New(c.I18n, locales.FS)
logx.Must(err)
```

```go
Language: middleware.NewLanguageMiddleware().Handle,
```

示例 import：

```go
"<module>/api/internal/locales"
"<module>/pkg/api/middleware"
"<module>/pkg/i18n"
```

---

## 第六步：创建底层 RPC Client

如果一个 RPC 进程里定义多个 protobuf Service：

```text
PingService
LanguageService
TimezoneService
CurrencyService
RegionService
```

不要为每个 Service 创建一个底层 `zrpc.Client`。

只创建一个：

```go
rpcClient := zrpc.MustNewClient(
    c.XxxRpc,
    zrpc.WithUnaryClientInterceptor(rpcerror.UnaryClientInterceptor),
)
```

然后分别创建生成的 Service Client：

```go
PingRpc: pingservice.NewPingService(rpcClient),
LanguageRpc: languageservice.NewLanguageService(rpcClient),
TimezoneRpc: timezoneservice.NewTimezoneService(rpcClient),
```

这样：

```text
一个 zrpc.Client
-> 多个 generated Service Client
```

不会因为多个 protobuf Service 建立多套底层 RPC Client。

---

## 第七步：完整 ServiceContext 示例

```go
type ServiceContext struct {
    Config config.Config

    PingRpc     pingservice.PingService
    LanguageRpc languageservice.LanguageService

    Trans    *i18n.Translator
    Language rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {
    // 创建翻译器
    trans, err := i18n.New(c.I18n, locales.FS)
    logx.Must(err)

    // 创建底层RPC客户端
    rpcClient := zrpc.MustNewClient(
        c.XxxRpc,
        zrpc.WithUnaryClientInterceptor(rpcerror.UnaryClientInterceptor),
    )

    return &ServiceContext{
        Config: c,

        PingRpc:     pingservice.NewPingService(rpcClient),
        LanguageRpc: languageservice.NewLanguageService(rpcClient),

        Trans:    trans,
        Language: middleware.NewLanguageMiddleware().Handle,
    }
}
```

---

## 第八步：在 API main 中注册公共能力

API 启动入口需要接入：

```go
v, err := validate.New(c.I18n.DefaultLanguage)
logx.Must(err)

// 注册全局参数校验器
httpx.SetValidator(v)
```

创建 ServiceContext 后：

```go
ctx := svc.NewServiceContext(c)
```

根据运行模式决定是否返回 debug：

```go
debug := c.Mode == service.DevMode || c.Mode == service.TestMode
```

注册统一错误处理：

```go
httpx.SetErrorHandlerCtx(
    errorhandler.New(ctx.Trans, debug).Handle,
)
```

注册统一成功响应：

```go
httpx.SetOkHandler(response.Ok)
```

注册语言中间件：

```go
server.Use(ctx.Language)
```

完整关键顺序：

```text
加载配置
↓
创建 Validator
↓
注册 Validator
↓
创建 HTTP Server
↓
创建 ServiceContext
↓
注册 ErrorHandler
↓
注册 OkHandler
↓
注册 Language Middleware
↓
注册 Handler
↓
启动 Server
```

---

# 5. API 业务 Logic 编写规则

API Logic 保持薄。

推荐职责：

```text
API Request
↓
转换 Proto Request
↓
调用 RPC
↓
RPC error 直接 return
↓
转换 RPC Response
↓
返回 API Response
```

示例：

```go
func (l *CreateLogic) Create(req *types.CreateRequest) (*types.CreateResponse, error) {
    result, err := l.svcCtx.LanguageRpc.Create(
        l.ctx,
        &language.CreateLanguageRequest{
            Code:     req.Code,
            NameI18N: req.NameI18n,
            Status:   req.Status,
        },
    )
    if err != nil {
        return nil, err
    }

    return &types.CreateResponse{
        Id: result.Id,
    }, nil
}
```

API Logic 不应该重复：

- 拼统一成功响应
- 拼统一错误响应
- 翻译 i18n
- 判断 gRPC Code
- 写 HTTP Status
- 记录 RPC 方法名

这些统一由公共层处理。

---

# 6. API 完整错误链

RPC 业务 Logic：

```go
return nil, grpcerror.NotFound(i18nkey.DataNotFound)
```

产生：

```text
gRPC Code = NotFound
Message = common.data_not_found
```

API RPC Client：

```text
rpcerror.UnaryClientInterceptor
↓
保存具体RPC方法
↓
保留原始gRPC Status
```

API Logic：

```go
if err != nil {
    return nil, err
}
```

全局 ErrorHandler：

```text
NotFound
↓
HTTP 404
↓
common.data_not_found
↓
Translator
↓
数据不存在
```

最终：

```json
{
  "code": 404,
  "msg": "数据不存在"
}
```

---

# 7. RPC 接入步骤

## 第一步：使用 `pkg/rpc/grpcerror`

RPC 业务错误统一返回 gRPC Status。

例如：

```go
import (
    "<module>/pkg/i18nkey"
    "<module>/pkg/rpc/grpcerror"
)
```

```go
if in.Id <= 0 {
    return nil, grpcerror.InvalidArgument(i18nkey.ValidationError)
}
```

查询不到：

```go
return nil, grpcerror.NotFound(i18nkey.DataNotFound)
```

内部错误：

```go
return nil, grpcerror.Internal(i18nkey.DatabaseError)
```

RPC 不负责翻译成中文或英文。

---

## 第二步：Ent 错误转换留在具体 RPC 项目

Ent generated package 与具体仓库强绑定，因此类似：

```text
rpc/internal/enterror/
```

继续保留在业务仓库，不放入 `pkg/`。

职责例如：

```text
Ent NotFound
-> gRPC NotFound

Ent Validation
-> gRPC InvalidArgument

Ent Constraint
-> 对应统一错误

未知数据库错误
-> gRPC Internal
```

这样 `pkg/` 不会依赖：

```text
<module>/rpc/ent
```

保持未来可迁移性。

---

## 第三步：数据库接入

RPC Config 示例：

```go
type Config struct {
    zrpc.RpcServerConf

    DatabaseConf database.DatabaseConf
}
```

YAML：

```yaml
DatabaseConf:
  Host: 127.0.0.1
  Port: 5432
  DBName: p9_platform
  Username: root
  Password: root
  SSLMode: disable
```

ServiceContext：

```go
driver, err := c.DatabaseConf.NewDriver()
logx.Must(err)

entOpts := []ent.Option{
    ent.Log(logx.Info),
    ent.Driver(driver),
}

if c.Mode == service.DevMode || c.Mode == service.TestMode {
    entOpts = append(entOpts, ent.Debug())
}

db := ent.NewClient(entOpts...)
```

数据库连接生命周期由具体 RPC 服务负责。

例如启动入口：

```go
ctx := svc.NewServiceContext(c)
defer ctx.DB.Close()
```

---

## 第四步：数据库迁移属于具体项目

`pkg/database` 只负责数据库连接和 Driver，不负责某个项目的 Ent Schema。

例如：

```text
rpc/internal/svc/migrate.go
```

继续属于具体 RPC 项目。

---

# 8. 多 protobuf Service 注意事项

一个 RPC 进程可以包含多个 protobuf Service。

例如：

```text
platform-base-rpc
├─ PingService
├─ LanguageService
├─ TimezoneService
├─ CurrencyService
└─ RegionService
```

API 侧：

```text
一个 zrpc.Client
-> 多个 generated Service Client
```

RPC Server 侧需要逐个注册：

```go
base.RegisterPingServiceServer(...)
base.RegisterLanguageServiceServer(...)
base.RegisterTimezoneServiceServer(...)
```

`goctl` 对已经存在的 RPC 主入口不会在后续生成时自动补充新的 Service 注册。

因此新增 protobuf Service 后，需要检查并手动维护 RPC 主入口的：

```go
RegisterXXXServiceServer(...)
```

不要通过删除主入口重新生成来解决，因为主入口通常已经包含数据库迁移、资源关闭等项目自己的启动逻辑。

---

# 9. i18n 接入规则

## 9.1 请求语言

API 统一使用：

```text
X-Lang
```

例如：

```text
X-Lang: zh-CN
```

```text
X-Lang: en-US
```

不要在不同模块分别使用其他 Header。

---

## 9.2 默认语言

默认语言放配置：

```yaml
I18n:
  DefaultLanguage: zh-CN
```

Middleware 只负责把请求语言写入 context。

Translator 负责：

```text
请求未传语言
或对应语言资源不可用
↓
回退 DefaultLanguage
```

---

## 9.3 错误 Key 由 RPC 返回

RPC 推荐：

```go
grpcerror.InvalidArgument(i18nkey.ValidationError)
```

不要：

```go
errors.New("数据校验失败")
```

也不要在 RPC 根据语言翻译。

语言翻译统一放 API。

---

# 10. 参数校验规则

API DSL / Request Struct 负责声明请求约束：

```go
Id int64 `json:"id" validate:"required,gt=0"`
```

```go
Status *int64 `json:"status,optional" validate:"omitempty,oneof=1 2"`
```

API Validator 负责提供友好的请求入口校验。

RPC 对真正影响数据正确性的规则仍然需要自己校验或由 Ent / DB 兜底。

不能因为 API 已经校验，就认为 RPC 一定不会被直接调用。

推荐理解：

```text
API validate
-> 对外请求第一层友好校验

RPC Logic
-> 业务规则校验

Ent / DB
-> 数据完整性兜底
```

---

# 11. Debug 规则

不要单独维护：

```yaml
Debug: true
```

统一使用 go-zero：

```yaml
Mode: dev
```

判断：

```go
debug := c.Mode == service.DevMode || c.Mode == service.TestMode
```

开发 / 测试：

```json
{
  "code": 500,
  "msg": "系统内部错误",
  "debug": {
    "source": "rpc",
    "rpc": "...",
    "error": "..."
  }
}
```

生产：

```json
{
  "code": 500,
  "msg": "系统内部错误"
}
```

生产环境不得向前端暴露：

- SQL 错误
- RPC 内部错误
- 数据库信息
- 调用栈
- 内部服务信息

---

# 12. 接入完成检查清单

## API

确认：

- [ ] 已复制需要的 `pkg/`
- [ ] 已执行 `go mod tidy`
- [ ] 已建立 `api/internal/locales`
- [ ] 已配置 `I18n.DefaultLanguage`
- [ ] 已配置 RPC Client
- [ ] Translator 初始化成功
- [ ] 一个底层 RPC Client 只创建一次
- [ ] RPC Client 已注册 `rpcerror.UnaryClientInterceptor`
- [ ] 已注册 `httpx.SetValidator`
- [ ] 已注册 `httpx.SetErrorHandlerCtx`
- [ ] 已注册 `httpx.SetOkHandler`
- [ ] 已注册 Language Middleware
- [ ] API Logic 不手工拼统一响应
- [ ] API Logic RPC 错误直接向上返回
- [ ] `Mode: pro` 时响应不存在 `debug`

## RPC

确认：

- [ ] 已使用 `pkg/rpc/grpcerror`
- [ ] 已使用统一 `i18nkey`
- [ ] RPC 不返回已经翻译好的中文或英文错误
- [ ] Ent 错误转换保留在项目 `rpc/internal`
- [ ] 数据库通过 `pkg/database` 初始化
- [ ] Ent Client 生命周期正确关闭
- [ ] 需要自动迁移时由具体项目负责
- [ ] 新增 protobuf Service 后已检查 RPC 主入口注册

---

# 13. 推荐测试场景

## 成功响应

请求一个正常接口，确认：

```json
{
  "code": 0,
  "msg": "ok",
  "data": {}
}
```

---

## API 参数校验

传入非法参数，例如：

```text
page=0
```

确认：

- HTTP Status = 400
- 响应结构统一
- `X-Lang: zh-CN` 返回中文校验信息
- `X-Lang: en-US` 返回英文校验信息

---

## RPC NotFound

RPC 返回：

```go
grpcerror.NotFound(i18nkey.DataNotFound)
```

确认 API：

```text
HTTP 404
```

中文：

```json
{
  "code": 404,
  "msg": "数据不存在"
}
```

英文：

```json
{
  "code": 404,
  "msg": "Data not found"
}
```

---

## RPC Internal

模拟 RPC 内部错误。

开发环境确认：

```json
{
  "code": 500,
  "msg": "系统内部错误",
  "debug": {
    "source": "rpc",
    "rpc": "...",
    "error": "..."
  }
}
```

生产环境确认：

```json
{
  "code": 500,
  "msg": "系统内部错误"
}
```

---

## RPC 不可用

停止 RPC 服务后调用 API。

预期：

```text
HTTP 503
```

响应：

```json
{
  "code": 503,
  "msg": "服务暂不可用"
}
```

---

# 14. 后续迁移公共仓库原则

以后如果建立真正的公共包仓库，目标是将当前：

```text
<repo>/pkg/
```

整体迁移到公共仓库。

迁移时原则上只需要：

1. 移动公共包
2. 修改 Module / import path
3. 各业务仓库删除本地重复 `pkg`
4. `go mod tidy`

如果某个 `pkg` 目录在迁移时需要大量引用业务仓库的 `internal` 包才能工作，说明当前公共包边界设计有问题，应先拆除业务依赖，再迁移。

---

# 15. 总体调用链

最终 API / RPC 基础链路：

```text
HTTP Request
    ↓
Language Middleware
    ↓
API参数解析
    ↓
Validator
    ↓
API Logic
    ↓
一个底层 zrpc.Client
    ↓
rpcerror Interceptor
    ↓
Generated Service Client
    ↓
RPC Logic
    ↓
业务校验
    ↓
Ent / PostgreSQL
```

成功：

```text
RPC Response
↓
API Logic
↓
response.Ok
↓
code=0, msg=ok
```

失败：

```text
RPC gRPC Error
↓
rpcerror
↓
API Logic直接返回
↓
errorhandler
├─ gRPC Code -> HTTP Status
├─ i18n Key -> 当前语言文案
├─ 5xx -> 内部日志
└─ dev/test -> debug
↓
统一 ErrorResponse
```

这套公共能力的目标是让业务代码只关注业务本身，不在每个 Handler / Logic 中重复编写参数错误、响应包装、gRPC 状态转换、国际化和调试信息处理。
