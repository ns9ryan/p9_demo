# p9 core 说明文档

## 一、HTTP API

core-api 默认 `http://192.168.0.15:18000`，前缀 `/admin`。JSON 字段以 `[api/desc](../api/desc)` 为准。

`PartnerMode` 在 `rpc/etc/core.yaml`：`"on"` 分分站，`"off"` 平台。`on` 时 `operator_id` 只信 JWT，忽略请求体里的分站字段。

## 接口索引

- [约定](#约定)
- [公开（无 JWT）](#公开无-jwt)
  - [POST /admin/login](#post-adminlogin)
  - [POST /admin/refresh](#post-adminrefresh)
- [仅 JWT](#仅-jwt)
  - [POST /admin/logout](#post-adminlogout)
  - [POST /admin/logout/all](#post-adminlogoutall)
  - [GET /admin/user/info](#get-adminuserinfo)
  - [GET /admin/user/perm](#get-adminuserperm)
  - [GET /admin/menu/role](#get-adminmenurole)
  - [POST /admin/user/password/self](#post-adminuserpasswordself)
- [JWT + Casbin](#jwt--casbin)
  - [分站（](#分站on-使用)`on` [使用）](#分站on-使用)
    - [GET /admin/operator/self](#get-adminoperatorself)
    - [POST /admin/operator/update](#post-adminoperatorupdate)
  - [用户](#用户)
    - [POST /admin/user/create](#post-adminusercreate)
    - [POST /admin/user/update](#post-adminuserupdate)
    - [POST /admin/user/delete](#post-adminuserdelete)
    - [POST /admin/user/list](#post-adminuserlist)
    - [GET /admin/user/detail](#get-adminuserdetail)
    - [POST /admin/user/password](#post-adminuserpassword)
    - [POST /admin/user/roles](#post-adminuserroles)
  - [角色](#角色)
    - [POST /admin/role/create](#post-adminrolecreate)
    - [POST /admin/role/update](#post-adminroleupdate)
    - [POST /admin/role/delete](#post-adminroledelete)
    - [POST /admin/role/list](#post-adminrolelist)
    - [GET /admin/role/detail](#get-adminroledetail)
  - [菜单（全局，不分分站）](#菜单全局不分分站)
    - [POST /admin/menu/create](#post-adminmenucreate)
    - [POST /admin/menu/update](#post-adminmenuupdate)
    - [POST /admin/menu/delete](#post-adminmenudelete)
    - [POST /admin/menu/list](#post-adminmenulist)
  - [API（全局，不分分站）](#api全局不分分站)
    - [POST /admin/api/create](#post-adminapicreate)
    - [POST /admin/api/update](#post-adminapiupdate)
    - [POST /admin/api/delete](#post-adminapidelete)
    - [POST /admin/api/list](#post-adminapilist)
  - [授权](#授权)
    - [POST /admin/authority/menu/update](#post-adminauthoritymenuupdate)
    - [POST /admin/authority/menu/role](#post-adminauthoritymenurole)
    - [POST /admin/authority/api/update](#post-adminauthorityapiupdate)
    - [POST /admin/authority/api/role](#post-adminauthorityapirole)
- [二、RPC 初始化](#二rpc-初始化)
  - [CreatePlatformAdmin / bootstrapAdmin](#1-createplatformadmin--bootstrapadmin-创建总网超级管理员账号)
  - [CreateOperatorAdmin / bootstrapOperator](#2-createoperatoradmin--bootstrapoperator-创建分站超级管理员账号)
- [三、菜单 / API 自动迁移（RegisterCatalog）](#三菜单--api-自动迁移registercatalog)
  - [示例：promo-api](#示例promo-api)

---



## 约定



### 响应结构体

成功：

```json
{ "code": 0, "msg": "ok", "data": {} }
```

失败时 HTTP 状态码与 `code` 相同，无 `data`：

```json
{ "code": 401, "msg": "未登录" }
```

access 过期（前端用 refresh 后续请求）：

```json
{ "code": 498, "msg": "登录已过期" }
```

无 `returns` 的接口成功时 `data` 为 `null`。时间字段为 Unix 秒。下文示例均为完整 HTTP body。

### 鉴权

登录后请求头：`Authorization: Bearer <access_token>`。


| 分组           | 中间件                        |
| ------------ | -------------------------- |
| 公开           | 无                          |
| 仅 JWT        | JWT（`CheckToken`）          |
| JWT + Casbin | JWT + Authority（`Enforce`） |

### 多语言

请求头：`X-Lang: zh-CN`（缺省）或 `en-US`。`zh*` 归一为 `zh-CN`，`en*` 为 `en-US`。

内置菜单 `title`、角色 `role_name`、接口 `description` 在库中存 i18n key（如 `route.dashboard`、`role.superAdmin`），响应按当前语言翻译。自定义名称没有对应词条时原样返回。信封 `msg`（成功/失败）同样翻译。下文示例默认 `zh-CN`。

### 分页


| 字段          | 说明                         |
| ----------- | -------------------------- |
| `page`      | 从 1 起，缺省 1                 |
| `page_size` | 用户列表缺省 20，角色列表缺省 50，上限 100 |




### 枚举


| 字段                                       | 值                    |
| ---------------------------------------- | -------------------- |
| `status`                                 | `1` 启用，`2` 停用        |
| `menu_type`                              | `0` 目录，`1` 菜单，`2` 按钮 |
| `hide_menu` / `disabled` / `is_required` | `0` 否，`1` 是          |




### 错误码


| code | 含义                                                               |
| ---- | ---------------------------------------------------------------- |
| 400  | 参数错误                                                             |
| 401  | 未登录 / token 无效 / 需重新登录                                           |
| 403  | 无权限或资源被禁用                                                        |
| 404  | 不存在                                                              |
| 498  | access token 过期：用 `refresh_token` 调 `POST /admin/refresh` 后重试原请求 |
| 500  | 内部错误                                                             |


refresh token 过期仍返回 **401**，不要用 498 刷新，避免死循环。

### 响应字段说明

信封：


| 字段     | 类型                    | 说明                   |
| ------ | --------------------- | -------------------- |
| `code` | int                   | `0` 成功，失败为 HTTP 状态码  |
| `msg`  | string                | 提示文案                 |
| `data` | object / array / null | 业务数据；无返回体的接口为 `null` |


`data.token`（登录 / 刷新）：


| 字段               | 类型     | 说明                              |
| ---------------- | ------ | ------------------------------- |
| `access_token`   | string | 业务请求 Bearer，过期返回 498            |
| `refresh_token`  | string | 仅用于 refresh / logout，不能当 Bearer |
| `expire`         | int64  | access 过期时间，Unix 秒              |
| `refresh_expire` | int64  | refresh 过期时间，Unix 秒             |


`UserPublic`：


| 字段               | 类型       | 说明                        |
| ---------------- | -------- | ------------------------- |
| `id`             | int64    | 用户主键                      |
| `user_code`      | string   | 用户对外编码                    |
| `username`       | string   | 登录名                       |
| `display_name`   | string   | 显示名                       |
| `operator_id`    | int64    | 所属分站；平台用户可省略或为 0          |
| `is_super_admin` | bool     | 是否超管；创建用户接口不可设为 true      |
| `status`         | int32    | `1` 启用，`2` 停用             |
| `role_codes`     | string[] | 角色编码列表                    |
| `home_path`      | string   | 登录后首页，默认 `/dashboard`     |
| `created_at`     | int64    | 创建时间，Unix 秒               |
| `last_login_at`  | int64    | 最后登录时间，Unix 秒；从未登录可省略或为 0 |
| `mobile`         | string   | 手机号；未填写可省略或为空             |
| `email`          | string   | 邮箱；未填写可省略或为空              |


`RoleInfo`：


| 字段            | 类型     | 说明                             |
| ------------- | ------ | ------------------------------ |
| `id`          | int64  | 角色主键                           |
| `operator_id` | int64  | 所属分站；平台角色可省略或为 0               |
| `role_code`   | string | 角色编码，分站内唯一                     |
| `role_name`   | string | 角色名称；内置为 i18n key，响应已按语言翻译     |
| `description` | string | 备注                             |
| `status`      | int32  | `1` 启用，`2` 停用；停用会清 Casbin      |
| `is_system`   | bool   | 内置角色（如 `super_admin`）不可删、不可改编码 |
| `sort_no`     | int32  | 排序，越小越前                        |
| `created_at`  | int64  | 创建时间，Unix 秒                    |
| `updated_at`  | int64  | 更新时间，Unix 秒                    |


`MenuInfo` / `MenuNode`：


| 字段           | 类型         | 说明                                 |
| ------------ | ---------- | ---------------------------------- |
| `id`         | int64      | 菜单主键                               |
| `parent_id`  | int64      | 父菜单 id，根为 `0`                      |
| `menu_type`  | int32      | `0` 目录，`1` 菜单，`2` 按钮               |
| `path`       | string     | 前端路由 path                          |
| `name`       | string     | 全局唯一标识，`RegisterCatalog` 按此 upsert |
| `component`  | string     | 前端组件路径                             |
| `redirect`   | string     | 重定向                                |
| `title`      | string     | 显示标题；内置为 i18n key，响应已按语言翻译     |
| `icon`       | string     | 图标                                 |
| `permission` | string     | 按钮权限码；动态菜单只含空字符串项                  |
| `hide_menu`  | int32      | `1` 侧栏隐藏                           |
| `sort`       | int32      | 排序，越小越前                            |
| `disabled`   | int32      | `1` 停用（仅 MenuInfo）                 |
| `created_at` | int64      | 创建时间，Unix 秒（仅 MenuInfo）            |
| `updated_at` | int64      | 更新时间，Unix 秒（仅 MenuInfo）            |
| `children`   | MenuNode[] | 子节点（仅 MenuNode）                    |


`ApiInfo`：


| 字段             | 类型     | 说明                                  |
| -------------- | ------ | ----------------------------------- |
| `id`           | int64  | 接口主键                                |
| `description`  | string | 接口说明；内置为 i18n key，响应已按语言翻译        |
| `api_group`    | string | 分组，如 `user`                         |
| `method`       | string | HTTP 方法，限 GET/POST/PUT/PATCH/DELETE |
| `path`         | string | 以 `/` 开头；与 method 组成唯一键             |
| `is_required`  | int32  | `1` 分配 API 权限时强制带上                  |
| `service_name` | string | 所属服务，如 `core-api`                   |
| `created_at`   | int64  | 创建时间，Unix 秒                         |
| `updated_at`   | int64  | 更新时间，Unix 秒                         |


`OperatorInfo`：


| 字段                         | 类型     | 说明                   |
| -------------------------- | ------ | -------------------- |
| `id`                       | int64  | 分站主键                 |
| `operator_code`            | string | 分站代码，不可改             |
| `timezone_code`            | string | 时区，如 `Asia/Shanghai` |
| `settlement_currency_code` | string | 结算币种，如 `CNY`         |
| `status`                   | int32  | `1` 启用，`2` 停用        |
| `required_config_version`  | int32  | 要求完成的配置版本            |
| `completed_config_version` | int32  | 已完成的配置版本             |
| `config_completed_at`      | int64  | 配置完成时间，Unix 秒，未完成可省略 |
| `created_at`               | int64  | 创建时间，Unix 秒          |
| `updated_at`               | int64  | 更新时间，Unix 秒          |


列表类 `data`：


| 字段      | 类型    | 说明              |
| ------- | ----- | --------------- |
| `list`  | array | 当前页数据，元素见上面对应类型 |
| `total` | int64 | 总条数             |


其它 `data`：


| 字段            | 类型       | 说明                                      |
| ------------- | -------- | --------------------------------------- |
| `permissions` | string[] | 当前用户按钮权限码（`GET /user/perm`）             |
| `menu_ids`    | int64[]  | 角色已分配菜单 id（`POST /authority/menu/role`） |
| `path`        | string   | 接口路径（API 授权项）                           |
| `method`      | string   | HTTP 方法（API 授权项）                        |


---



## 公开（无 JWT）



### POST /admin/login

`off`：`username` + `password`。`on`：再加 `operator_code`。

**请求**


| 字段              | 位置   | 必填   | 类型     | 说明                                 |
| --------------- | ---- | ---- | ------ | ---------------------------------- |
| `username`      | json | 是    | string | 登录用户名                              |
| `password`      | json | 是    | string | 登录密码                               |
| `operator_code` | json | on 是 | string | 分站代码；`PartnerMode=on` 必填，`off` 不要传 |


```json
{
  "username": "admin",
  "password": "Admin@123",
  "operator_code": "demo"
}
```

**响应**

```json
{
  "code": 0,
  "msg": "ok",
  "data": {
    "token": {
      "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
      "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
      "expire": 1710000000,
      "refresh_expire": 1710600000
    },
    "user": {
      "id": 1,
      "user_code": "a1b2c3d4e5f6",
      "username": "admin",
      "display_name": "分站超管",
      "operator_id": 1,
      "is_super_admin": true,
      "status": 1,
      "role_codes": ["super_admin"],
      "home_path": "/dashboard"
    }
  }
}
```



### POST /admin/refresh

校验 refresh 密钥 + 用户 salt，下发新 token 对并拉黑旧 refresh。过期返回 401。

**请求**


| 字段              | 位置   | 必填  | 类型     | 说明                           |
| --------------- | ---- | --- | ------ | ---------------------------- |
| `refresh_token` | json | 是   | string | 登录下发的 refresh；过期返回 401 需重新登录 |


```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**响应**

```json
{
  "code": 0,
  "msg": "ok",
  "data": {
    "token": {
      "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
      "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
      "expire": 1710086400,
      "refresh_expire": 1710686400
    },
    "user": {
      "id": 1,
      "user_code": "a1b2c3d4e5f6",
      "username": "admin",
      "display_name": "分站超管",
      "operator_id": 1,
      "is_super_admin": true,
      "status": 1,
      "role_codes": ["super_admin"],
      "home_path": "/dashboard"
    }
  }
}
```

---



## 仅 JWT

需 Header：`Authorization: Bearer <access_token>`。

### POST /admin/logout

拉黑当前 access 与该 refresh（需 Redis）。不改 salt。

**请求**


| 字段              | 位置   | 必填  | 类型     | 说明                                        |
| --------------- | ---- | --- | ------ | ----------------------------------------- |
| `refresh_token` | json | 是   | string | 要拉黑的 refresh；当前 access 从 Authorization 读取 |


```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**响应**

```json
{ "code": 0, "msg": "ok", "data": { "result": "success" } }
```



### POST /admin/logout/all

轮换该用户 `salt`，全部 access / refresh 失效。无请求体。

**响应**

```json
{ "code": 0, "msg": "ok", "data": { "result": "success" } }
```



### GET /admin/user/info

当前用户。无请求参数。

**响应**

```json
{
  "code": 0,
  "msg": "ok",
  "data": {
    "id": 1,
    "user_code": "a1b2c3d4e5f6",
    "username": "admin",
    "display_name": "分站超管",
    "operator_id": 1,
    "is_super_admin": true,
    "status": 1,
    "role_codes": ["super_admin"],
    "home_path": "/dashboard"
  }
}
```



### GET /admin/user/perm

本角色菜单中 `permission != ""` 的按钮码。

**响应**

```json
{
  "code": 0,
  "msg": "ok",
  "data": {
    "permissions": ["user:create", "role:create", "menu:create", "api:create"]
  }
}
```



### GET /admin/menu/role

动态菜单：`permission == ""` 的项组树。

**响应**

```json
{
  "code": 0,
  "msg": "ok",
  "data": [
    {
      "id": 1,
      "parent_id": 0,
      "menu_type": 1,
      "path": "/dashboard",
      "name": "Dashboard",
      "component": "dashboard/index",
      "redirect": "",
      "title": "工作台",
      "icon": "",
      "permission": "",
      "hide_menu": 0,
      "sort": 1,
      "children": []
    },
    {
      "id": 10,
      "parent_id": 0,
      "menu_type": 0,
      "path": "/system",
      "name": "System",
      "component": "",
      "redirect": "",
      "title": "系统管理",
      "icon": "",
      "permission": "",
      "hide_menu": 0,
      "sort": 10,
      "children": [
        {
          "id": 11,
          "parent_id": 10,
          "menu_type": 1,
          "path": "/system/user",
          "name": "User",
          "component": "system/user/index",
          "redirect": "",
          "title": "用户管理",
          "icon": "",
          "permission": "",
          "hide_menu": 0,
          "sort": 11,
          "children": []
        }
      ]
    }
  ]
}
```



### POST /admin/user/password/self

改自己密码并轮换自己 salt。

**请求**


| 字段             | 位置   | 必填  | 类型     | 说明                          |
| -------------- | ---- | --- | ------ | --------------------------- |
| `old_password` | json | 是   | string | 当前密码                        |
| `password`     | json | 是   | string | 新密码；成功后轮换 salt，旧 token 全部失效 |


```json
{
  "old_password": "Admin@123",
  "password": "Admin@456"
}
```

**响应**

```json
{ "code": 0, "msg": "ok", "data": { "result": "success" } }
```

---



## JWT + Casbin

需 JWT + Casbin。Header：`Authorization: Bearer <access_token>`。

### 分站（`on` 使用）

无 `operator/create|list|delete`。开分站只用 bootstrap / `CreateOperatorAdmin`。

#### GET /admin/operator/self

当前 JWT 对应分站。无请求参数。

**响应**

```json
{
  "code": 0,
  "msg": "ok",
  "data": {
    "id": 1,
    "operator_code": "demo",
    "timezone_code": "Asia/Shanghai",
    "settlement_currency_code": "CNY",
    "status": 1,
    "required_config_version": 0,
    "completed_config_version": 0,
    "created_at": 1700000000,
    "updated_at": 1700000000
  }
}
```



#### POST /admin/operator/update

不可改 `operator_code`。

**请求**


| 字段                         | 位置   | 必填  | 类型     | 说明                                       |
| -------------------------- | ---- | --- | ------ | ---------------------------------------- |
| `timezone_code`            | json | 否   | string | 时区，如 `Asia/Shanghai`；不可改 `operator_code` |
| `settlement_currency_code` | json | 否   | string | 结算币种，如 `CNY`                             |


```json
{
  "timezone_code": "UTC",
  "settlement_currency_code": "USD"
}
```

**响应**

```json
{ "code": 0, "msg": "ok", "data": { "result": "success" } }
```



### 用户

`off` 写 `operator_id=NULL`；`on` 写 JWT 分站。禁止设 `is_super_admin`。

#### POST /admin/user/create

**请求**


| 字段             | 位置   | 必填  | 类型      | 说明                     |
| -------------- | ---- | --- | ------- | ---------------------- |
| `username`     | json | 是   | string  | 登录名，同分站内唯一             |
| `password`     | json | 是   | string  | 初始密码                   |
| `display_name` | json | 是   | string  | 显示名                    |
| `mobile`       | json | 否   | string  | 手机号                    |
| `email`        | json | 否   | string  | 邮箱                     |
| `status`       | json | 否   | int32   | `1` 启用 / `2` 停用，缺省 `1` |
| `role_ids`     | json | 否   | int64[] | 绑定角色 id，须与用户同分站        |


```json
{
  "username": "alice",
  "password": "Alice@123",
  "display_name": "运营",
  "mobile": "13800000000",
  "email": "alice@example.com",
  "status": 1,
  "role_ids": [2]
}
```

**响应**

```json
{
  "code": 0,
  "msg": "ok",
  "data": {
    "id": 2,
    "user_code": "f6e5d4c3b2a1",
    "username": "alice",
    "display_name": "运营",
    "operator_id": 1,
    "is_super_admin": false,
    "status": 1,
    "role_codes": ["editor"],
    "home_path": "/dashboard"
  }
}
```



#### POST /admin/user/update

不可改 `operator_id` / `is_super_admin` / `salt`。

**请求**


| 字段             | 位置   | 必填  | 类型     | 说明              |
| -------------- | ---- | --- | ------ | --------------- |
| `id`           | json | 是   | int64  | 要更新的用户主键        |
| `display_name` | json | 否   | string | 显示名             |
| `mobile`       | json | 否   | string | 手机号             |
| `email`        | json | 否   | string | 邮箱              |
| `status`       | json | 否   | int32  | `1` 启用 / `2` 停用 |


```json
{
  "id": 2,
  "display_name": "运营主管",
  "mobile": "13900000000",
  "email": "alice2@example.com",
  "status": 1
}
```

**响应**

```json
{ "code": 0, "msg": "ok", "data": { "result": "success" } }
```



#### POST /admin/user/delete

软删。不可删自己、不可删 root。`id` 与 `ids` 至少一个，`ids` 优先。

**请求**


| 字段    | 位置   | 必填  | 类型      | 说明                        |
| ----- | ---- | --- | ------- | ------------------------- |
| `id`  | json | 否   | int64   | 单个用户 id；与 `ids` 同时传时忽略本字段 |
| `ids` | json | 否   | int64[] | 批量用户 id，优先于 `id`          |


```json
{ "ids": [2, 3] }
```

或：

```json
{ "id": 2 }
```

**响应**

```json
{ "code": 0, "msg": "ok", "data": { "result": "success" } }
```



#### POST /admin/user/list

`off` 全量；`on` 仅本分站。

**请求**


| 字段            | 位置   | 必填  | 类型      | 说明                                      |
| ------------- | ---- | --- | ------- | --------------------------------------- |
| `page`        | json | 否   | int32   | 页码，从 1 起，缺省 1                           |
| `page_size`   | json | 否   | int32   | 每页条数，缺省 20，上限 100                       |
| `username`    | json | 否   | string  | 用户名，模糊匹配；空则忽略                           |
| `mobile`      | json | 否   | string  | 手机号，模糊匹配；空则忽略                           |
| `email`       | json | 否   | string  | 邮箱，模糊匹配；空则忽略                            |
| `display_name` | json | 否   | string  | 显示名称，模糊匹配；空则忽略                          |
| `role_ids`    | json | 否   | []int64 | 角色多选，命中任一角色即可；空则忽略。多条件与其它字段同时生效（AND） |


```json
{
  "page": 1,
  "page_size": 20,
  "username": "ali",
  "mobile": "138",
  "email": "@example.com",
  "display_name": "运营",
  "role_ids": [2, 3]
}
```

**响应**

```json
{
  "code": 0,
  "msg": "ok",
  "data": {
    "list": [
      {
        "id": 1,
        "user_code": "a1b2c3d4e5f6",
        "username": "admin",
        "display_name": "分站超管",
        "operator_id": 1,
        "is_super_admin": true,
        "status": 1,
        "role_codes": ["super_admin"],
        "home_path": "/dashboard",
        "created_at": 1710000000,
        "last_login_at": 1710003600
      },
      {
        "id": 2,
        "user_code": "f6e5d4c3b2a1",
        "username": "alice",
        "display_name": "运营",
        "operator_id": 1,
        "is_super_admin": false,
        "status": 1,
        "role_codes": ["editor"],
        "home_path": "/dashboard",
        "created_at": 1710001200,
        "last_login_at": 0
      }
    ],
    "total": 2
  }
}
```



#### GET /admin/user/detail

Query：`/admin/user/detail?id=2`

**请求**


| 字段   | 位置    | 必填  | 类型    | 说明                               |
| ---- | ----- | --- | ----- | -------------------------------- |
| `id` | query | 是   | int64 | 用户主键，如 `/admin/user/detail?id=2` |


**响应**

```json
{
  "code": 0,
  "msg": "ok",
  "data": {
    "id": 2,
    "user_code": "f6e5d4c3b2a1",
    "username": "alice",
    "display_name": "运营",
    "operator_id": 1,
    "is_super_admin": false,
    "status": 1,
    "role_codes": ["editor"],
    "home_path": "/dashboard",
    "created_at": 1710001200,
    "last_login_at": 1710003600,
    "mobile": "13800000000",
    "email": "alice@example.com"
  }
}
```



#### POST /admin/user/password

改他人密码并轮换对方 salt。

**请求**


| 字段         | 位置   | 必填  | 类型     | 说明               |
| ---------- | ---- | --- | ------ | ---------------- |
| `user_id`  | json | 是   | int64  | 被改密的用户主键         |
| `password` | json | 是   | string | 新密码；成功后轮换对方 salt |


```json
{
  "user_id": 2,
  "password": "Alice@456"
}
```

**响应**

```json
{ "code": 0, "msg": "ok", "data": { "result": "success" } }
```



#### POST /admin/user/roles

用户与角色的 `operator_id` 必须同为 NULL 或同值。

**请求**


| 字段         | 位置   | 必填  | 类型      | 说明                  |
| ---------- | ---- | --- | ------- | ------------------- |
| `user_id`  | json | 是   | int64   | 用户主键                |
| `role_ids` | json | 是   | int64[] | 全量替换该用户角色；角色须与用户同分站 |


```json
{
  "user_id": 2,
  "role_ids": [2, 3]
}
```

**响应**

```json
{ "code": 0, "msg": "ok", "data": { "result": "success" } }
```



### 角色

接口不能把 `is_system` 设为 true。停用角色会清掉该 `(role_code, v1)` 的 Casbin 策略。`v1`：`off` 为 `""`，`on` 为 `operator_id` 十进制字符串。

#### POST /admin/role/create

**请求**


| 字段            | 位置   | 必填  | 类型     | 说明                       |
| ------------- | ---- | --- | ------ | ------------------------ |
| `role_code`   | json | 是   | string | 角色编码，分站内唯一；不可设为系统保留码冒充超管 |
| `role_name`   | json | 是   | string | 角色名称                     |
| `description` | json | 否   | string | 备注                       |
| `status`      | json | 否   | int32  | `1` 启用 / `2` 停用，缺省 `1`   |
| `sort_no`     | json | 否   | int32  | 排序，越小越前                  |


```json
{
  "role_code": "editor",
  "role_name": "运营",
  "description": "运营角色",
  "status": 1,
  "sort_no": 10
}
```

**响应**

```json
{
  "code": 0,
  "msg": "ok",
  "data": {
    "id": 2,
    "operator_id": 1,
    "role_code": "editor",
    "role_name": "运营",
    "description": "运营角色",
    "status": 1,
    "is_system": false,
    "sort_no": 10,
    "created_at": 1700000100,
    "updated_at": 1700000100
  }
}
```



#### POST /admin/role/update

内置角色的 `role_code` / `is_system` 不可改。

**请求**


| 字段            | 位置   | 必填  | 类型     | 说明                                       |
| ------------- | ---- | --- | ------ | ---------------------------------------- |
| `id`          | json | 是   | int64  | 要更新的角色主键                                 |
| `role_name`   | json | 否   | string | 角色名称；内置角色的 `role_code` / `is_system` 不可改 |
| `description` | json | 否   | string | 备注                                       |
| `status`      | json | 否   | int32  | `1` 启用 / `2` 停用；停用会清 Casbin              |
| `sort_no`     | json | 否   | int32  | 排序，越小越前                                  |


```json
{
  "id": 2,
  "role_name": "高级运营",
  "description": "可看用户列表",
  "status": 1,
  "sort_no": 20
}
```

**响应**

```json
{ "code": 0, "msg": "ok", "data": { "result": "success" } }
```



#### POST /admin/role/delete

`is_system` 不可删；有用户绑定则拒绝；按 `(role_code, v1)` 清 Casbin。

**请求**


| 字段    | 位置   | 必填  | 类型      | 说明                        |
| ----- | ---- | --- | ------- | ------------------------- |
| `id`  | json | 否   | int64   | 单个角色 id；与 `ids` 同时传时忽略本字段 |
| `ids` | json | 否   | int64[] | 批量角色 id，优先于 `id`          |


```json
{ "ids": [2] }
```

**响应**

```json
{ "code": 0, "msg": "ok", "data": { "result": "success" } }
```



#### POST /admin/role/list

**请求**


| 字段          | 位置   | 必填  | 类型    | 说明                                          |
| ----------- | ---- | --- | ----- | ------------------------------------------- |
| `page`      | json | 否   | int32 | 页码，从 1 起，缺省 1                               |
| `page_size` | json | 否   | int32 | 每页条数，缺省 50，上限 100                           |
| `role_name` | json | 否   | string | 模糊匹配 `role_name` **或** `role_code`；空则忽略 |


```json
{ "page": 1, "page_size": 50, "role_name": "admin" }
```

**响应**

```json
{
  "code": 0,
  "msg": "ok",
  "data": {
    "list": [
      {
        "id": 1,
        "operator_id": 1,
        "role_code": "super_admin",
        "role_name": "超级管理员",
        "description": "",
        "status": 1,
        "is_system": true,
        "sort_no": 0,
        "created_at": 1700000000,
        "updated_at": 1700000000
      },
      {
        "id": 2,
        "operator_id": 1,
        "role_code": "editor",
        "role_name": "运营",
        "description": "运营角色",
        "status": 1,
        "is_system": false,
        "sort_no": 10,
        "created_at": 1700000100,
        "updated_at": 1700000100
      }
    ],
    "total": 2
  }
}
```



#### GET /admin/role/detail

Query：`/admin/role/detail?id=2`

**请求**


| 字段   | 位置    | 必填  | 类型    | 说明                               |
| ---- | ----- | --- | ----- | -------------------------------- |
| `id` | query | 是   | int64 | 角色主键，如 `/admin/role/detail?id=2` |


**响应**

```json
{
  "code": 0,
  "msg": "ok",
  "data": {
    "id": 2,
    "operator_id": 1,
    "role_code": "editor",
    "role_name": "运营",
    "description": "运营角色",
    "status": 1,
    "is_system": false,
    "sort_no": 10,
    "created_at": 1700000100,
    "updated_at": 1700000100
  }
}
```



### 菜单（全局，不分分站）

`name` 唯一。创建后补授权给各分站 `super_admin`。超管角色不可减菜单。

#### POST /admin/menu/create

**请求**


| 字段           | 位置   | 必填  | 类型     | 说明                                    |
| ------------ | ---- | --- | ------ | ------------------------------------- |
| `name`       | json | 是   | string | 全局唯一标识，后续 `RegisterCatalog` 按此 upsert |
| `title`      | json | 是   | string | 显示标题                                  |
| `menu_type`  | json | 是   | int32  | `0` 目录，`1` 菜单，`2` 按钮                  |
| `parent_id`  | json | 否   | int64  | 父菜单 id，须已存在；根省略为 `0`                  |
| `path`       | json | 否   | string | 前端路由 path                             |
| `component`  | json | 否   | string | 前端组件路径                                |
| `redirect`   | json | 否   | string | 重定向地址                                 |
| `icon`       | json | 否   | string | 图标                                    |
| `permission` | json | 否   | string | 按钮权限码，如 `user:export`；菜单/目录留空         |
| `hide_menu`  | json | 否   | int32  | `1` 侧栏隐藏，缺省 `0`                       |
| `sort`       | json | 否   | int32  | 排序，越小越前                               |
| `disabled`   | json | 否   | int32  | `1` 停用，缺省 `0`                         |


```json
{
  "parent_id": 11,
  "menu_type": 2,
  "name": "UserExport",
  "title": "导出用户",
  "permission": "user:export",
  "sort": 112
}
```

**响应**

```json
{
  "code": 0,
  "msg": "ok",
  "data": {
    "id": 1112,
    "parent_id": 11,
    "menu_type": 2,
    "path": "",
    "name": "UserExport",
    "component": "",
    "redirect": "",
    "title": "导出用户",
    "icon": "",
    "permission": "user:export",
    "hide_menu": 0,
    "sort": 112,
    "disabled": 0,
    "created_at": 1700000200,
    "updated_at": 1700000200
  }
}
```



#### POST /admin/menu/update

按 `id` 部分更新；禁止挂到自己的子孙。

**请求**


| 字段           | 位置   | 必填  | 类型     | 说明                   |
| ------------ | ---- | --- | ------ | -------------------- |
| `id`         | json | 是   | int64  | 要更新的菜单主键             |
| `parent_id`  | json | 否   | int64  | 新父菜单 id，不可挂到自己的子孙    |
| `menu_type`  | json | 否   | int32  | `0` 目录，`1` 菜单，`2` 按钮 |
| `path`       | json | 否   | string | 前端路由 path            |
| `name`       | json | 否   | string | 全局唯一标识               |
| `component`  | json | 否   | string | 前端组件路径               |
| `redirect`   | json | 否   | string | 重定向地址                |
| `title`      | json | 否   | string | 显示标题                 |
| `icon`       | json | 否   | string | 图标                   |
| `permission` | json | 否   | string | 按钮权限码                |
| `hide_menu`  | json | 否   | int32  | `1` 侧栏隐藏             |
| `sort`       | json | 否   | int32  | 排序，越小越前              |
| `disabled`   | json | 否   | int32  | `1` 停用               |


```json
{
  "id": 1112,
  "title": "导出用户 Excel",
  "sort": 113,
  "hide_menu": 0
}
```

**响应**

```json
{ "code": 0, "msg": "ok", "data": { "result": "success" } }
```



#### POST /admin/menu/delete

有子节点则拒绝；先清 `sys_role_menu`。

**请求**


| 字段    | 位置   | 必填  | 类型      | 说明                        |
| ----- | ---- | --- | ------- | ------------------------- |
| `id`  | json | 否   | int64   | 单个菜单 id；与 `ids` 同时传时忽略本字段 |
| `ids` | json | 否   | int64[] | 批量菜单 id，优先于 `id`          |


```json
{ "id": 1112 }
```

**响应**

```json
{ "code": 0, "msg": "ok", "data": { "result": "success" } }
```



#### POST /admin/menu/list

授权页勾选，全量未停用菜单。无请求体。

**响应**

```json
{
  "code": 0,
  "msg": "ok",
  "data": [
    {
      "id": 10,
      "parent_id": 0,
      "menu_type": 0,
      "path": "/system",
      "name": "System",
      "component": "",
      "redirect": "",
      "title": "系统管理",
      "icon": "",
      "permission": "",
      "hide_menu": 0,
      "sort": 10,
      "disabled": 0,
      "created_at": 1700000000,
      "updated_at": 1700000000
    },
    {
      "id": 11,
      "parent_id": 10,
      "menu_type": 1,
      "path": "/system/user",
      "name": "User",
      "component": "system/user/index",
      "redirect": "",
      "title": "用户管理",
      "icon": "",
      "permission": "",
      "hide_menu": 0,
      "sort": 11,
      "disabled": 0,
      "created_at": 1700000000,
      "updated_at": 1700000000
    }
  ]
}
```



### API（全局，不分分站）

`(method, path)` 唯一。改 path/method 或删除会同步 Casbin。创建后补授权给各分站 `super_admin`。

#### POST /admin/api/create

**请求**


| 字段             | 位置   | 必填  | 类型     | 说明                                  |
| -------------- | ---- | --- | ------ | ----------------------------------- |
| `method`       | json | 是   | string | HTTP 方法，限 GET/POST/PUT/PATCH/DELETE |
| `path`         | json | 是   | string | 接口路径，须以 `/` 开头；与 method 组成唯一键       |
| `description`  | json | 否   | string | 接口说明，授权页展示                          |
| `api_group`    | json | 否   | string | 分组名，如 `user`                        |
| `is_required`  | json | 否   | int32  | `1` 表示分配 API 权限时强制带上，缺省 `0`         |
| `service_name` | json | 否   | string | 所属服务，如 `core-api`                   |


```json
{
  "description": "导出用户",
  "api_group": "user",
  "method": "POST",
  "path": "/admin/user/export",
  "is_required": 0,
  "service_name": "core-api"
}
```

**响应**

```json
{
  "code": 0,
  "msg": "ok",
  "data": {
    "id": 30,
    "description": "导出用户",
    "api_group": "user",
    "method": "POST",
    "path": "/admin/user/export",
    "is_required": 0,
    "service_name": "core-api",
    "created_at": 1700000300,
    "updated_at": 1700000300
  }
}
```



#### POST /admin/api/update

**请求**


| 字段             | 位置   | 必填  | 类型     | 说明                               |
| -------------- | ---- | --- | ------ | -------------------------------- |
| `id`           | json | 是   | int64  | 要更新的接口主键                         |
| `description`  | json | 否   | string | 接口说明                             |
| `api_group`    | json | 否   | string | 分组名                              |
| `method`       | json | 否   | string | HTTP 方法；改 method/path 会同步 Casbin |
| `path`         | json | 否   | string | 接口路径，须以 `/` 开头                   |
| `is_required`  | json | 否   | int32  | `1` 分配权限时强制带上                    |
| `service_name` | json | 否   | string | 所属服务                             |


```json
{
  "id": 30,
  "description": "导出用户 Excel",
  "api_group": "user"
}
```

**响应**

```json
{ "code": 0, "msg": "ok", "data": { "result": "success" } }
```



#### POST /admin/api/delete

同时删除对应 Casbin 策略。

**请求**


| 字段    | 位置   | 必填  | 类型      | 说明                        |
| ----- | ---- | --- | ------- | ------------------------- |
| `id`  | json | 否   | int64   | 单个接口 id；与 `ids` 同时传时忽略本字段 |
| `ids` | json | 否   | int64[] | 批量接口 id，优先于 `id`          |


```json
{ "ids": [30] }
```

**响应**

```json
{ "code": 0, "msg": "ok", "data": { "result": "success" } }
```



#### POST /admin/api/list

无请求体。

**响应**

```json
{
  "code": 0,
  "msg": "ok",
  "data": [
    {
      "id": 5,
      "description": "后台用户列表",
      "api_group": "user",
      "method": "POST",
      "path": "/admin/user/list",
      "is_required": 0,
      "service_name": "core-api",
      "created_at": 1700000000,
      "updated_at": 1700000000
    },
    {
      "id": 6,
      "description": "后台用户详情",
      "api_group": "user",
      "method": "GET",
      "path": "/admin/user/detail",
      "is_required": 0,
      "service_name": "core-api",
      "created_at": 1700000000,
      "updated_at": 1700000000
    }
  ]
}
```



### 授权

`role_id` 与 `id` 二选一（`role_id` 优先）。角色必须本租户。

#### POST /admin/authority/menu/update

`is_system` 超管角色禁止减菜单。

**请求**


| 字段         | 位置   | 必填  | 类型      | 说明                  |
| ---------- | ---- | --- | ------- | ------------------- |
| `role_id`  | json | 是   | int64   | 角色主键，必须本租户          |
| `menu_ids` | json | 否   | int64[] | 全量替换该角色菜单；超管角色禁止减菜单 |


```json
{
  "role_id": 2,
  "menu_ids": [1, 10, 11]
}
```

**响应**

```json
{ "code": 0, "msg": "ok", "data": { "result": "success" } }
```



#### POST /admin/authority/menu/role

**请求**


| 字段        | 位置   | 必填  | 类型    | 说明                     |
| --------- | ---- | --- | ----- | ---------------------- |
| `role_id` | json | 否   | int64 | 角色主键，与 `id` 二选一，优先本字段  |
| `id`      | json | 否   | int64 | 角色主键别名，`role_id` 未传时使用 |


```json
{ "role_id": 2 }
```

**响应**

```json
{
  "code": 0,
  "msg": "ok",
  "data": {
    "menu_ids": [1, 10, 11]
  }
}
```



#### POST /admin/authority/api/update

按 `(role_code, v1)` 替换；合并 `is_required`；超管角色禁止减 API。

**请求**


| 字段              | 位置   | 必填  | 类型       | 说明                     |
| --------------- | ---- | --- | -------- | ---------------------- |
| `role_id`       | json | 是   | int64    | 角色主键，必须本租户             |
| `data`          | json | 否   | object[] | 全量替换该角色 API；超管禁止减 API  |
| `data[].path`   | json | 是   | string   | 接口 path，须已在 sys_api 登记 |
| `data[].method` | json | 是   | string   | HTTP 方法                |


```json
{
  "role_id": 2,
  "data": [
    { "path": "/admin/user/list", "method": "POST" },
    { "path": "/admin/user/detail", "method": "GET" }
  ]
}
```

**响应**

```json
{ "code": 0, "msg": "ok", "data": { "result": "success" } }
```



#### POST /admin/authority/api/role

**请求**


| 字段        | 位置   | 必填  | 类型    | 说明                     |
| --------- | ---- | --- | ----- | ---------------------- |
| `role_id` | json | 否   | int64 | 角色主键，与 `id` 二选一，优先本字段  |
| `id`      | json | 否   | int64 | 角色主键别名，`role_id` 未传时使用 |


```json
{ "role_id": 2 }
```

**响应**

```json
{
  "code": 0,
  "msg": "ok",
  "data": [
    { "path": "/admin/user/list", "method": "POST" },
    { "path": "/admin/user/detail", "method": "GET" }
  ]
}
```

---



## 二、RPC 初始化



### 1. CreatePlatformAdmin / `bootstrapAdmin` 创建总网超级管理员账号

`PartnerMode=off`。实现：`[rpc/bootstrap/bootstrap.go](../rpc/bootstrap/bootstrap.go)` `CreatePlatformAdmin`。

**请求**（RPC `BootstrapAdminReq`）


| 字段             | 位置  | 必填  | 类型     | 说明                           |
| -------------- | --- | --- | ------ | ---------------------------- |
| `init_token`   | rpc | 是   | string | 初始化口令，须与 yaml `InitToken` 一致 |
| `username`     | rpc | 是   | string | 平台超管登录名                      |
| `password`     | rpc | 是   | string | 平台超管密码                       |
| `display_name` | rpc | 否   | string | 显示名，缺省为 username             |


```json
{
  "init_token": "change-me-init-token",
  "username": "admin",
  "password": "Admin@123",
  "display_name": "平台超管"
}
```

**响应**

```json
{
  "id": 1,
  "user_code": "a1b2c3d4e5f6",
  "username": "admin",
  "display_name": "平台超管",
  "is_super_admin": true,
  "status": 1,
  "role_codes": ["super_admin"],
  "home_path": "/dashboard"
}
```



### 2. CreateOperatorAdmin / `bootstrapOperator` 创建分站超级管理员账号

`PartnerMode=on`。实现：`CreateOperatorAdmin`。分站已存在则复用，再创建本分站超管。

**请求**（RPC `BootstrapOperatorReq`）


| 字段                         | 位置  | 必填  | 类型     | 说明                           |
| -------------------------- | --- | --- | ------ | ---------------------------- |
| `init_token`               | rpc | 是   | string | 初始化口令，须与 yaml `InitToken` 一致 |
| `operator_code`            | rpc | 是   | string | 分站代码，已存在则复用该分站               |
| `timezone_code`            | rpc | 否   | string | 时区，缺省 `UTC`                  |
| `settlement_currency_code` | rpc | 否   | string | 结算币种，缺省 `USD`                |
| `username`                 | rpc | 是   | string | 本分站超管登录名                     |
| `password`                 | rpc | 是   | string | 本分站超管密码                      |
| `display_name`             | rpc | 否   | string | 显示名，缺省为 username             |


```json
{
  "init_token": "change-me-init-token",
  "operator_code": "demo",
  "timezone_code": "Asia/Shanghai",
  "settlement_currency_code": "CNY",
  "username": "admin",
  "password": "Admin@123",
  "display_name": "分站超管"
}
```

**响应**

```json
{
  "id": 1,
  "user_code": "a1b2c3d4e5f6",
  "username": "admin",
  "display_name": "分站超管",
  "operator_id": 1,
  "is_super_admin": true,
  "status": 1,
  "role_codes": ["super_admin"],
  "home_path": "/dashboard"
}
```

HTTP 对应：`POST /admin/bootstrap/admin`、`POST /admin/bootstrap/operator`（`init_token` 走 Header `X-Init-Token`）。

---



## 三、菜单 / API 自动迁移（RegisterCatalog）

业务服务**不要**在 core-rpc 里写死 path。各 HTTP 进程启动时调 RPC `registerCatalog`：

- 菜单按 `name` upsert；`parent_name` 在本批全部写入后再挂父子
- API 按 `(method, path)` upsert
- 新建或更新后给各分站 `super_admin` 补菜单与 Casbin（增量，不替换已有授权）

core-api 启动时注册系统管理菜单和 `/admin/user|role|menu|api|authority|operator/*`（见 `[api/internal/catalog/catalog.go](../api/internal/catalog/catalog.go)`）。若先 bootstrap 再启 HTTP，重启一次即可写入。

### 示例：promo-api

`[example/promo-api/main.go](../example/promo-api/main.go)` 等价请求：

`RegisterCatalogReq`：


| 字段      | 必填  | 类型       | 说明                            |
| ------- | --- | -------- | ----------------------------- |
| `menus` | 否   | object[] | 要注册的菜单，按 `name` upsert        |
| `apis`  | 否   | object[] | 要注册的 API，按 method+path upsert |


`menus[]`（`RegisterMenuReq`）：


| 字段            | 必填  | 类型     | 说明                    |
| ------------- | --- | ------ | --------------------- |
| `name`        | 是   | string | 全局唯一标识                |
| `title`       | 是   | string | 显示标题                  |
| `menu_type`   | 是   | int32  | `0` 目录，`1` 菜单，`2` 按钮  |
| `path`        | 否   | string | 前端路由 path             |
| `component`   | 否   | string | 前端组件路径                |
| `redirect`    | 否   | string | 重定向                   |
| `icon`        | 否   | string | 图标                    |
| `permission`  | 否   | string | 按钮权限码                 |
| `hide_menu`   | 否   | int32  | `1` 侧栏隐藏              |
| `sort`        | 否   | int32  | 排序，越小越前               |
| `disabled`    | 否   | int32  | `1` 停用                |
| `parent_name` | 否   | string | 父菜单的 `name`，本批写入后再挂父子 |


`apis[]`（`CreateApiReq`）：


| 字段             | 必填  | 类型     | 说明                 |
| -------------- | --- | ------ | ------------------ |
| `method`       | 是   | string | HTTP 方法            |
| `path`         | 是   | string | 以 `/` 开头           |
| `description`  | 否   | string | 接口说明               |
| `api_group`    | 否   | string | 分组                 |
| `is_required`  | 否   | int32  | `1` 分配权限时强制带上      |
| `service_name` | 否   | string | 所属服务，如 `promo-api` |


```json
{
  "menus": [
    {
      "name": "PromoCenter",
      "title": "route.promoCenter",
      "menu_type": 0,
      "path": "/promo",
      "sort": 20
    },
    {
      "name": "PromoActivityList",
      "title": "route.promoActivityList",
      "menu_type": 1,
      "path": "/promo/activity/list",
      "component": "promo/activity/list",
      "parent_name": "PromoCenter",
      "sort": 21
    }
  ],
  "apis": [
    {
      "description": "api.promoList",
      "api_group": "promo",
      "method": "GET",
      "path": "/admin/promo/list",
      "service_name": "promo-api"
    }
  ]
}
```

```go
_, err := cli.RegisterCatalog(ctx, req)
```

仍可单独调 `RegisterApi`。