# core

使用 **go-zero RPC + HTTP API** 的后台权限中心：登录、JWT、Casbin、厅（租户）隔离、菜单/API/角色 CRUD。业务服务（如 `example/promo-api`）通过 RPC `CheckToken` / `Enforce` 鉴权，并用 `common/entmixin` 做厅隔离。

## 目录

```
common/           可被 RPC、API、业务服务共用
  ctxdata/        JWT claims（Value + gRPC outgoing/incoming metadata）与 SkipTenant / SkipSoftDelete
  entmixin/       Time / SoftDelete / Tenant Mixin
  response/       {code, msg, data}
  middleware/     JWT + Casbin（调 Core RPC）
rpc/              Core 权限 RPC
  desc/*.proto    按资源拆分，goctls 合并为 rpc/core.proto
  ent/            权限表 schema（go generate）
  service/        业务实现
api/              后台 HTTP（/admin/*），只转发 Core RPC
  desc/all.api    入口；desc/core/*.api 按资源拆分
  internal/catalog  启动时 RegisterCatalog（系统菜单 + 后台 API + 默认语言）
example/promo-api 业务示例：GET /admin/promo/list
```



## 运行

依赖本机 Postgres、Redis（DSN 见 yaml，与 rbacx-v2 example 相同）。

```bash
# 终端 1：权限 RPC
cd rpc && go run . -f etc/core.yaml

# 终端 2：后台 HTTP
cd api && go run . -f etc/core-api.yaml

# 终端 3（可选）：业务 API
cd example/promo-api && go run . -f etc/promo-api.yaml
```

Partner 模式在 `rpc/etc/core.yaml` 的 `PartnerMode`：`"on"` 分厅，`"off"` 平台。

## 初始化与登录(使用rpc方式初始化)

`on` 模式：

```bash
curl -s localhost:8889/admin/bootstrap/operator \
  -H 'X-Init-Token: change-me-init-token' \
  -H 'Content-Type: application/json' \
  -d '{"operator_code":"demo","username":"admin","password":"Admin@123","display_name":"分站超管"}'

curl -s localhost:8889/admin/login \
  -H 'Content-Type: application/json' \
  -d '{"username":"admin","password":"Admin@123","operator_code":"demo"}'
```

`off` 模式把 `PartnerMode` 改为 `"off"`，改调 `POST /admin/bootstrap/admin`（body 无 `operator_code`），登录也不带厅代码。

之后请求头带 `Authorization: Bearer <access_token>`。响应统一为 `{ "code": 0, "msg": "ok", "data": ... }`，失败 `code` 为 HTTP 状态码。

## 接口

后台 HTTP 的请求参数与响应见 [docs/api.md](docs/api.md)。

- 公开：`/login` `/refresh` `/bootstrap/admin` `/bootstrap/operator`
- 仅 JWT：`/logout` `/logout/all` `/user/info` `/user/perm` `/menu/role` `/user/password/self`
- JWT + Casbin：用户/角色/菜单/API/授权、厅 `operator/self|update`

业务接口示例：`GET /admin/promo/list`（promo-api 启动时一次 `RegisterCatalog` 写入菜单「优惠中心 / 活动列表」和对应 API）。

## 业务服务怎么接

1. HTTP 中间件用 `common/middleware.JWT` + `Authority`，内部调 Core `CheckToken` / `Enforce`。JWT 会 `ctxdata.WithClaims` / `WithRawToken`，同时写入 gRPC outgoing metadata。
2. 调 RPC 时传入 HTTP 的 `ctx` 即可，gRPC 会自动带上 metadata。不要挂 Unary Client/Server Interceptor。RPC logic / Ent mixin 用 `ctxdata.ClaimsFromCtx`（无 Value 时从 incoming metadata 还原）。
3. 业务 Ent schema 嵌入 `entmixin.TimeMixin` + `entmixin.TenantMixin`（`operator_id` 可空）。`claims.OperatorID != 0` 且未 `ctxdata.SkipTenant` 时自动按厅过滤；创建时自动盖章。
4. `import _ "your/module/ent/runtime"`，启动时 `Schema.Create`。
5. 菜单和需鉴权的 HTTP path 不要写进 core-rpc。各 HTTP 服务启动时调 RPC `RegisterCatalog`（菜单按 `name` upsert，API 按 method+path upsert，多语言按 group+key+lang upsert，支持的语言按 `lang` 幂等插入，并给各厅 `super_admin` 补授权）。`core-api` 注册系统管理菜单、`/admin/user|role|menu|api|i18n|authority|operator/*` 和默认语言 `zh-CN` / `zh-HK` / `en-US`；`example/promo-api` 注册「优惠中心 / 活动列表」和 `GET /admin/promo/list`。若先 bootstrap 再启对应 HTTP 服务，重启一次即可写入。仍可单独调 `RegisterApi`。



## 生成代码

按 goctls 规范：RPC 写 `rpc/desc/*.proto`（`// group: xxx` + lowerCamelCase rpc），API 写 `api/desc/all.api` 及其 import。改完后：

```bash
make gen-rpc
make gen-api

cd rpc/ent && go generate
cd example/promo-api/ent && go generate
```

`goctls` v1.12.6、protoc 需在 PATH 中。`PROJECT_STYLE=go_zero`。RPC 生成会把 desc 合并进 `rpc/core.proto`，logic 按 group 落在 `rpc/internal/logic/<group>/`，client 在 `rpc/coreclient`。

登录仍用 `*http.Request` 取 UserAgent / 回退 IP：`api/internal/logic/public/login_logic.go`。刷新 / 预览走生成代码 `NewXxxLogic(r.Context(), svcCtx)`，客户端 IP 由全局 `middleware.ClientIP`（与 `middleware.I18n` 同级）写入上下文。使用 core JWT 的业务 API 同样要 `server.Use(middleware.ClientIP)`；JWT 仅在 ctx 里还没有 IP 时回填。重新 `make gen-api` 后若 login handler 被覆盖，把 `NewLoginLogic(r.Context(), svcCtx)` 改回 `NewLoginLogic(r, svcCtx)`。
