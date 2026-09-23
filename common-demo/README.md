# p9/common

Go 模块 `oa.98ent.com/p9/common`，供 p9 业务服务（如 [core](https://oa.98ent.com/p9/core)）共用的 HTTP/RPC 基础能力：统一响应、错误码、多语言、请求上下文、追踪与错误日志。

要求 **Go 1.26**。

```bash
go get oa.98ent.com/p9/common@latest
```

## 包一览

| 包 | 作用 |
| --- | --- |
| [`i18n`](i18n) | 语言解析、进程内文案翻译、数据库词典缓存 |
| [`xerr`](xerr) | 带 HTTP 状态 / 文案 key / 堆栈的错误，与 gRPC status 互转 |
| [`response`](response) | go-zero `httpx` 统一信封；把解析错误、gRPC 错误转成 `xerr` |
| [`ctxdata`](ctxdata) | Claims、Token、ClientIP 写入 context 与 gRPC metadata |
| [`errorlog`](errorlog) | 5xx 错误采集：上下文暂存 + 异步上报 |
| [`tracing`](tracing) | 把错误/堆栈写进 OpenTelemetry span（Jaeger 标红） |
| [`utils`](utils) | 客户端 IP、IP 白名单、ID 解析 |
| [`modules/cache`](modules/cache) | 进程内 TTL 缓存（词典 `TG` 用 2 分钟档） |

## HTTP 接入

业务 API 启动时：加载 locale JSON → `SetupHTTPX` → 中间件写入语言和 IP。core-api / promo-api 的典型顺序：

```go
trans, err := i18n.New(c.I18n, localeFS) // localeFS 为业务自己的 embed.FS
logx.Must(err)

server.Use(func(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        lang := i18n.ParseLang(r.Header.Get("X-Lang"))
        next(w, r.WithContext(i18n.WithLang(r.Context(), lang)))
    }
})

response.SetupHTTPX(trans, i18n.CodePlatform, isDebug)
```

成功与失败信封：

```json
{ "code": 0, "msg": "成功", "data": {} }
{ "code": 400, "msg": "参数错误" }
```

`isDebug` 为 true 时错误体会多 `debug.error` / `debug.stack` / `debug.cause`（仅 Dev/Test 打开）。

`code` 参数是站点 `i18n_code`（`platform` / `operator` / `promo` 等），失败文案先走 `Translator`，无模板参数时再按 `error` 组查数据库词典 `i18n.TG`。

## 多语言

两条通道：

| 函数 | 来源 | 用途 |
| --- | --- | --- |
| `Translator.T` / `Tf` | 业务 embed 的 `locale/*.json` | 信封 `msg`、登录/角色等进程内文案 |
| `i18n.Dict` / `TG` | `SetDictLoader` + 2 分钟缓存 | 菜单标题、接口说明、前端词条 |

语言码：`zh-CN` / `zh-HK` / `en-US`。`ParseLang` 读 `X-Lang`（或 Accept-Language 第一段）；空则默认 `en-US`。

```go
i18n.SetDictLoader(func(ctx context.Context, code, group, lang string) (map[string]string, error) {
    // 通常调 core-rpc GetI18nDict
    return dict, nil
})

title := i18n.TG(ctx, i18n.CodePlatform, i18n.GroupMenu, "route.dashboard")
```

写词典后调 `i18n.InvalidateAll()`。站点编码：`CodeByPartnerMode("on")` → `operator`，否则 `platform`。

文案 key 常量在 [`i18n/msg.go`](i18n/msg.go)（如 `i18n.InvalidParam`）。locale JSON 由**各业务仓库**提供，本模块不内置语言文件。

`Translator` 配置：

```yaml
I18n:
  DefaultLanguage: zh-CN
```

## 错误

业务返回 `*xerr.Error`，`Message` 用 `i18n` 的 key，不要写死中文：

```go
return xerr.BadRequest(i18n.InvalidParam)
return xerr.BadRequestWith(i18n.ParamRequired, map[string]any{"Field": "id"})
```

RPC 侧用 `xerr.RpcErr(err)` 转 gRPC status（HTTP 状态映射到 gRPC code；5xx 通过 `errdetails.DebugInfo` 带 cause/stack）。API 的 `response.FromError` 再转回 HTTP：`InvalidArgument`→400、`Unauthenticated`→401、`PermissionDenied`→403、`NotFound`→404；其余 4xx/5xx 可直接用对应数字 code。

`response` 还会把 go-zero/json 解析失败收成 400，例如缺字段 → `paramRequired`，类型不对 → `paramTypeMismatch`，非法 JSON → `invalidParam`。

## 请求上下文

JWT 校验后把身份放进 context，并写入 outgoing gRPC metadata（`x-user-id`、`x-operator-code`、`x-client-ip` 等）。RPC 入站无盒子时从 incoming metadata 还原。

```go
ctx = ctxdata.WithClaims(ctx, claims)
ctx = ctxdata.WithClientIP(ctx, utils.ClientIP(r))
c := ctxdata.ClaimsFromCtx(ctx)
```

`SkipTenant` / `SkipSoftDelete` 给 Ent 拦截器关闭分站过滤或软删过滤。

## 错误日志与追踪

- `errorlog.WithBag` 在请求开始挂袋子；`response.SetupHTTPX` 对 5xx 调 `errorlog.Note`。中间件再 `Report` 异步写入（实现 `errorlog.Recorder`）。
- `tracing.Error(ctx, err)` 把错误和堆栈打到当前 span。RPC 可挂 `tracing.UnaryServerInterceptor()`；HTTP 5xx 用 `tracing.HTTPError`。

## 工具

[`utils`](utils)：`ClientIP`（`X-Forwarded-For` / `RemoteAddr`）、`NormalizeIP`、`IPAllowed`（精确 IP 或 CIDR）、`ResolveIDs`。

[`modules/cache`](modules/cache)：预置 1s / 15s / 30s / 2min / 10min / 1d / 1w。`GetC` 在过期前半个 TTL 后台刷新。词典走 `cache.TwoMinuteCache`。

## 测试

```bash
go test ./...
```
