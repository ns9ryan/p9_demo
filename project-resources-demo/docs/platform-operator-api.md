# platform-operator API 文档

`platform-operator-api` 提供 P9 总网侧分站管理相关 HTTP API。

当前已实现并纳入本文档的功能：

- 分站管理
- 分站档案
- 分站域名
- 分站管理员
- 基础资源分配汇总
- 语言分配
- 经营地区分配
- 代理子线路分配
- 服务健康检查

游戏资源分配目前仍是接口骨架，由其他模块或同事继续实现，本文档暂不作为前端对接依据。

---

## 一、接口约定

### 服务地址

默认开发端口：

```text
18002
```

实际 Host、域名和网关地址以部署环境为准。

例如直接访问开发服务器时：

```text
http://<host>:18002
```

### 鉴权

除 `GET /ping` 外，当前所有 `/admin/*` 接口均使用：

```text
Jwt, ActionLog, Authority
```

前端请求需要携带后台登录 Token：

```http
Authorization: Bearer <access_token>
```

Token 获取、刷新和过期处理规则参考同目录的 [core-api.md](./core-api.md)。

### 多语言

请求语言通过 `X-Lang` Header 指定：

```http
X-Lang: zh-CN
```

当前默认语言：

```text
zh-CN
```

错误文案等会根据当前请求语言返回。

### 请求格式

POST 接口统一使用 JSON：

```http
Content-Type: application/json
```

GET 接口参数通过 Query String 传递。

### 成功响应

成功响应统一格式：

```json
{
  "code": 0,
  "msg": "ok",
  "data": {}
}
```


| 字段     | 类型     | 说明         |
| ------ | ------ | ---------- |
| `code` | int    | 成功固定为 `0`  |
| `msg`  | string | 成功固定为 `ok` |
| `data` | object | 接口业务数据     |


没有业务数据返回的接口，当前 `data` 为：

```json
{}
```



### 错误响应

失败时 `code` 使用对应 HTTP 状态码，例如：

```json
{
  "code": 400,
  "msg": "参数校验失败"
}
```

开发和测试环境可能额外返回 `debug`：

```json
{
  "code": 400,
  "msg": "参数校验失败",
  "debug": {
    "source": "api",
    "error": "..."
  }
}
```

RPC 调用错误时可能包含：

```json
{
  "code": 400,
  "msg": "参数校验失败",
  "debug": {
    "source": "rpc",
    "rpc": "...",
    "error": "..."
  }
}
```

`debug` 仅用于开发调试，前端业务逻辑不应依赖该字段。

### 分页

分页列表统一使用：


| 字段          | 位置    | 必填  | 类型    | 说明              |
| ----------- | ----- | --- | ----- | --------------- |
| `page`      | query | 是   | int64 | 页码，从 `1` 开始     |
| `page_size` | query | 是   | int64 | 每页数量，范围 `1-100` |




### 时间字段

本文档中的时间字段统一为：

```text
Unix 毫秒时间戳
```

例如：

```text
1757836800000
```



### 公共枚举



#### 分站创建状态 `creation_status`


| 值   | 说明  |
| --- | --- |
| `1` | 草稿  |
| `2` | 已完成 |




#### 分站发布状态 `publish_status`


| 值   | 说明   |
| --- | ---- |
| `1` | 未发布  |
| `2` | 发布中  |
| `3` | 已发布  |
| `4` | 发布失败 |




#### 分站状态 `status`


| 值   | 说明  |
| --- | --- |
| `1` | 正常  |
| `2` | 暂停  |
| `3` | 关闭  |




#### 域名类型 `domain_type`


| 值   | 说明    |
| --- | ----- |
| `1` | 分站后台  |
| `2` | 代理后台  |
| `3` | 会员 H5 |




#### 域名状态 `status`


| 值   | 说明  |
| --- | --- |
| `1` | 启用  |
| `2` | 停用  |




#### 管理员状态 `status`


| 值   | 说明  |
| --- | --- |
| `1` | 启用  |
| `2` | 停用  |




#### 代理子线路编码


| 编码                | 说明     |
| ----------------- | ------ |
| `CASH_PRODUCTION` | 现金正式线路 |
| `CASH_DEMO`       | 现金试玩线路 |
| `CASH_TEST`       | 现金测试线路 |
| `CREDIT_DEMO`     | 信誉试玩线路 |
| `CREDIT_TEST`     | 信誉测试线路 |




### 关联基础数据

创建或修改分站时使用的时区、结算币种，以及经营地区分配所使用的国家地区，都来自 `platform-base`。

前端可参考：

- [platform-base API 文档](./platform-base-api.md)

语言分配使用 Core 当前启用的语言，前端可参考：

- [core API 文档](./core-api.md)

常用数据来源：


| 数据    | 建议接口                                                             |
| ----- | ---------------------------------------------------------------- |
| 时区    | `GET /admin/timezone/list-all?status=1`                          |
| 结算币种  | `GET /admin/currency/list-all?status=1`                          |
| 国家地区  | `GET /admin/region/list-all?status=1`                            |
| 启用语言  | `GET /admin/i18n/lang/enabled`                                   |
| 代理子线路 | `GET /admin/operator/agent-line-allocation/list?operator_id=...` |


---



## 二、当前开发状态说明

当前公开代码中有几个前端对接时需要特别注意的状态：

1. `POST /admin/operator/publish` 已经暴露 HTTP 接口，但当前 RPC 的发布业务逻辑仍是 TODO。现阶段调用会返回成功结构，但不会真正执行完整发布流程，前端暂时不要把成功响应视为“分站已经完成发布”。
2. `POST /admin/operator/delete` 已经完整提供 HTTP API、权限目录和 RPC 删除逻辑。当前只允许删除 `publish_status = 1` 未发布或 `publish_status = 4` 发布失败的分站。删除时会先按分站清理管理员账号，再在 `platform-operator` 本地事务中清理语言、经营地区、代理子线路、域名、档案和分站本身。`platform-game` 游戏资源的跨服务清理目前仍保留 TODO，待对应服务提供 RPC 后接入。
3. `game_allocation` 当前仍是空请求/空响应的接口骨架，本文档暂不纳入正式对接接口。
4. `POST /admin/operator/complete` 当前实现主要用于把 `creation_status` 从 `1` 更新为 `2`。当前后端不会在该接口中统一检查档案、域名、语言、地区、代理子线路等配置是否全部完成，前端创建向导仍应按照页面流程控制调用时机。

---



## 三、接口索引

- [服务健康检查](#get-ping)
  - [GET /ping](#get-ping)
- [分站管理](#operator)
  - [POST /admin/operator/create](#post-adminoperatorcreate)
  - [POST /admin/operator/update](#post-adminoperatorupdate)
  - [GET /admin/operator/get](#get-adminoperatorget)
  - [GET /admin/operator/list](#get-adminoperatorlist)
  - [POST /admin/operator/complete](#post-adminoperatorcomplete)
  - [POST /admin/operator/publish](#post-adminoperatorpublish)
  - [POST /admin/operator/delete](#post-adminoperatordelete)
- [分站档案](#operator-profile)
  - [POST /admin/operator/profile/create](#post-adminoperatorprofilecreate)
  - [POST /admin/operator/profile/update](#post-adminoperatorprofileupdate)
  - [GET /admin/operator/profile/get](#get-adminoperatorprofileget)
- [分站域名](#operator-domain)
  - [POST /admin/operator/domain/create](#post-adminoperatordomaincreate)
  - [POST /admin/operator/domain/update](#post-adminoperatordomainupdate)
  - [GET /admin/operator/domain/get](#get-adminoperatordomainget)
  - [GET /admin/operator/domain/list](#get-adminoperatordomainlist)
  - [POST /admin/operator/domain/delete](#post-adminoperatordomaindelete)
- [分站管理员](#operator-admin)
  - [POST /admin/operator/admin/create](#post-adminoperatoradmincreate)
  - [GET /admin/operator/admin/list](#get-adminoperatoradminlist)
  - [POST /admin/operator/admin/update](#post-adminoperatoradminupdate)
  - [POST /admin/operator/admin/resetPassword](#post-adminoperatoradminresetpassword)
  - [POST /admin/operator/admin/updateStatus](#post-adminoperatoradminupdatestatus)
- [基础资源分配](#basic-resource-allocation)
  - [GET /admin/operator/basic-resource-allocation/list](#get-adminoperatorbasic-resource-allocationlist)
  - [GET /admin/operator/language-allocation/list](#get-adminoperatorlanguage-allocationlist)
  - [POST /admin/operator/language-allocation/save](#post-adminoperatorlanguage-allocationsave)
  - [GET /admin/operator/region-allocation/list](#get-adminoperatorregion-allocationlist)
  - [POST /admin/operator/region-allocation/save](#post-adminoperatorregion-allocationsave)
  - [GET /admin/operator/agent-line-allocation/list](#get-adminoperatoragent-line-allocationlist)
  - [POST /admin/operator/agent-line-allocation/save](#post-adminoperatoragent-line-allocationsave)

---



## 四、Ping



### GET /ping

服务健康检查，不需要 JWT。

请求：无。

响应：

```json
{
  "code": 0,
  "msg": "ok",
  "data": {}
}
```

---



## 五、分站管理



### OperatorInfo


| 字段                         | 类型            | 说明                                    |
| -------------------------- | ------------- | ------------------------------------- |
| `id`                       | int64         | 分站 ID                                 |
| `code`                     | string        | 分站全局唯一业务编码，由后端创建时生成                   |
| `name`                     | string        | 分站名称                                  |
| `timezone_code`            | string        | IANA 时区编码                             |
| `settlement_currency_code` | string        | 结算币种编码                                |
| `creation_status`          | int64         | 创建状态：`1` 草稿，`2` 已完成                   |
| `publish_status`           | int64         | 发布状态：`1` 未发布，`2` 发布中，`3` 已发布，`4` 发布失败 |
| `status`                   | int64         | 分站状态：`1` 正常，`2` 暂停，`3` 关闭             |
| `remark`                   | string / null | 总网内部备注                                |
| `published_at`             | int64 / null  | 首次发布成功时间，Unix 毫秒时间戳                   |
| `created_at`               | int64         | 创建时间，Unix 毫秒时间戳                       |
| `updated_at`               | int64         | 更新时间，Unix 毫秒时间戳                       |


示例：

```json
{
  "id": 1,
  "code": "OP_8D7091378B244D89A51FB102251489F1",
  "name": "A01",
  "timezone_code": "Asia/Tokyo",
  "settlement_currency_code": "USD",
  "creation_status": 1,
  "publish_status": 1,
  "status": 1,
  "remark": "内部测试分站",
  "published_at": null,
  "created_at": 1757836800000,
  "updated_at": 1757836800000
}
```



### POST /admin/operator/create

创建分站基础信息。

新建分站默认：

```text
creation_status = 1
publish_status  = 1
status          = 1  // 未传 status 时
```

分站业务编码 `code` 由后端自动生成，前端不需要传。

#### 请求参数


| 字段                         | 位置   | 必填  | 类型     | 说明                   |
| -------------------------- | ---- | --- | ------ | -------------------- |
| `name`                     | json | 是   | string | 分站名称，非空，最大 100 个字符   |
| `timezone_code`            | json | 是   | string | 时区编码，最大 64 个字符       |
| `settlement_currency_code` | json | 是   | string | 结算币种编码，最大 16 个字符     |
| `status`                   | json | 否   | int64  | `1` 正常，`2` 暂停，`3` 关闭 |
| `remark`                   | json | 否   | string | 总网内部备注，最大 1000 个字符   |


请求示例：

```json
{
  "name": "A01",
  "timezone_code": "Asia/Tokyo",
  "settlement_currency_code": "USD",
  "status": 1,
  "remark": "内部测试分站"
}
```



#### 业务规则

- `timezone_code` 必须是 `platform-base` 中存在且已启用的时区。
- `settlement_currency_code` 必须是 `platform-base` 中存在且已启用的货币。
- 结算币种编码会自动去除首尾空格并转换为大写。
- 创建成功后返回系统生成的分站 ID 和业务编码。

响应：

```json
{
  "code": 0,
  "msg": "ok",
  "data": {
    "id": 1,
    "code": "OP_8D7091378B244D89A51FB102251489F1"
  }
}
```



### POST /admin/operator/update

修改分站基础信息。

#### 请求参数


| 字段                         | 位置   | 必填  | 类型     | 说明                   |
| -------------------------- | ---- | --- | ------ | -------------------- |
| `id`                       | json | 是   | int64  | 分站 ID，大于 `0`         |
| `name`                     | json | 否   | string | 分站名称，非空，最大 100 个字符   |
| `timezone_code`            | json | 否   | string | 时区编码，最大 64 个字符       |
| `settlement_currency_code` | json | 否   | string | 结算币种编码，最大 16 个字符     |
| `status`                   | json | 否   | int64  | `1` 正常，`2` 暂停，`3` 关闭 |
| `remark`                   | json | 否   | string | 总网内部备注，最大 1000 个字符   |


除 `id` 外，至少需要传一个需要修改的字段。

请求示例：

```json
{
  "id": 1,
  "name": "A01 新名称",
  "status": 2,
  "remark": "临时暂停"
}
```



#### 业务规则

- 修改时区时，新时区必须存在且处于启用状态。
- 修改结算币种时，新币种必须存在且处于启用状态。
- 发布中，或者已经有过成功发布时间的分站，不允许再修改 `timezone_code` 和 `settlement_currency_code`。
- 其他允许修改的字段仍按照后端当前规则处理。

响应：

```json
{
  "code": 0,
  "msg": "ok",
  "data": {}
}
```



### GET /admin/operator/get

按分站 ID 获取分站详情。

Query 示例：

```text
/admin/operator/get?id=1
```



#### 请求参数


| 字段   | 位置    | 必填  | 类型    | 说明           |
| ---- | ----- | --- | ----- | ------------ |
| `id` | query | 是   | int64 | 分站 ID，大于 `0` |


响应：

```json
{
  "code": 0,
  "msg": "ok",
  "data": {
    "id": 1,
    "code": "OP_8D7091378B244D89A51FB102251489F1",
    "name": "A01",
    "timezone_code": "Asia/Tokyo",
    "settlement_currency_code": "USD",
    "creation_status": 1,
    "publish_status": 1,
    "status": 1,
    "remark": "内部测试分站",
    "published_at": null,
    "created_at": 1757836800000,
    "updated_at": 1757836800000
  }
}
```



### GET /admin/operator/list

获取分站管理列表。

Query 示例：

```text
/admin/operator/list?page=1&page_size=20&keyword=A01&creation_status=1&publish_status=1&status=1
```



#### 请求参数


| 字段                | 位置    | 必填  | 类型     | 说明                               |
| ----------------- | ----- | --- | ------ | -------------------------------- |
| `page`            | query | 是   | int64  | 页码，从 `1` 开始                      |
| `page_size`       | query | 是   | int64  | 每页数量，范围 `1-100`                  |
| `keyword`         | query | 否   | string | 搜索关键字，匹配分站编码或名称，最大 100 个字符       |
| `creation_status` | query | 否   | int64  | `1` 草稿，`2` 已完成                   |
| `publish_status`  | query | 否   | int64  | `1` 未发布，`2` 发布中，`3` 已发布，`4` 发布失败 |
| `status`          | query | 否   | int64  | `1` 正常，`2` 暂停，`3` 关闭             |


列表按创建时间倒序，同一创建时间下按 ID 倒序。

响应：

```json
{
  "code": 0,
  "msg": "ok",
  "data": {
    "total": 1,
    "list": [
      {
        "id": 1,
        "code": "OP_8D7091378B244D89A51FB102251489F1",
        "name": "A01",
        "timezone_code": "Asia/Tokyo",
        "settlement_currency_code": "USD",
        "creation_status": 1,
        "publish_status": 1,
        "status": 1,
        "remark": "内部测试分站",
        "published_at": null,
        "created_at": 1757836800000,
        "updated_at": 1757836800000
      }
    ]
  }
}
```



### POST /admin/operator/complete

完成分站创建，将创建状态更新为：

```text
creation_status = 2
```



#### 请求参数


| 字段   | 位置   | 必填  | 类型    | 说明           |
| ---- | ---- | --- | ----- | ------------ |
| `id` | json | 是   | int64 | 分站 ID，大于 `0` |


请求示例：

```json
{
  "id": 1
}
```

说明：

- 已经是 `creation_status = 2` 时再次调用会直接返回成功。
- 当前后端实现不会在此接口统一校验档案、域名和基础资源是否完整，前端应在创建向导流程完成后再调用。

响应：

```json
{
  "code": 0,
  "msg": "ok",
  "data": {}
}
```



### POST /admin/operator/publish

发布分站。

> 当前接口已经暴露，但发布 RPC 业务逻辑仍是 TODO。当前调用会返回成功结构，但不会真正完成正式发布流程。前端暂时不要依赖该接口判断发布是否成功。



#### 请求参数


| 字段   | 位置   | 必填  | 类型    | 说明           |
| ---- | ---- | --- | ----- | ------------ |
| `id` | json | 是   | int64 | 分站 ID，大于 `0` |


请求示例：

```json
{
  "id": 1
}
```

当前响应：

```json
{
  "code": 0,
  "msg": "ok",
  "data": {}
}
```



### POST /admin/operator/delete

删除分站。

#### 请求参数


| 字段   | 位置   | 必填  | 类型    | 说明           |
| ---- | ---- | --- | ----- | ------------ |
| `id` | json | 是   | int64 | 分站 ID，大于 `0` |


请求示例：

```json
{
  "id": 1
}
```



#### 删除条件

只有下面两种发布状态允许删除：


| `publish_status` | 说明   | 是否允许删除 |
| ---------------- | ---- | ------ |
| `1`              | 未发布  | 是      |
| `2`              | 发布中  | 否      |
| `3`              | 已发布  | 否      |
| `4`              | 发布失败 | 是      |


因此前端列表中的删除按钮建议只在：

```text
publish_status == 1 || publish_status == 4
```

时允许操作。

#### 删除范围

当前 `platform-operator` RPC 会在同一个本地数据库事务中依次删除：

```text
分站管理员账号（按 operator_id 删除）
→ 语言分配
→ 经营地区分配
→ 代理子线路分配
→ 分站域名
→ 分站档案
→ 分站
```

管理员账号会先通过分站管理员服务按 `operator_id` 清理；其余数据在 `platform-operator` 本地事务中删除。任意一步失败时，本地事务不会提交。

> `platform-game` 游戏资源属于其他服务。当前 API Logic 已预留清理位置，但对应跨服务清理 RPC 尚未接入。在该 TODO 完成前，删除接口只能保证 `platform-operator` 自己的数据清理完整。

响应：

```json
{
  "code": 0,
  "msg": "ok",
  "data": {}
}
```

---



## 六、分站档案

一个分站最多只有一条档案记录。

前端建议流程：

```text
GET /admin/operator/profile/get
        ↓
profile == null
        → 调用 create
profile != null
        → 调用 update
```



### OperatorProfileInfo


| 字段              | 类型            | 说明              |
| --------------- | ------------- | --------------- |
| `id`            | int64         | 档案 ID           |
| `operator_id`   | int64         | 分站 ID           |
| `company_name`  | string / null | 公司名称            |
| `contact_name`  | string / null | 主要联系人名称         |
| `contact_email` | string / null | 主要联系人邮箱         |
| `remark`        | string / null | 总网内部档案备注        |
| `created_at`    | int64         | 创建时间，Unix 毫秒时间戳 |
| `updated_at`    | int64         | 更新时间，Unix 毫秒时间戳 |




### POST /admin/operator/profile/create

创建分站档案。

#### 请求参数


| 字段              | 位置   | 必填  | 类型     | 说明                      |
| --------------- | ---- | --- | ------ | ----------------------- |
| `operator_id`   | json | 是   | int64  | 分站 ID，大于 `0`            |
| `company_name`  | json | 否   | string | 公司名称，最大 200 个字符         |
| `contact_name`  | json | 否   | string | 主要联系人名称，最大 100 个字符      |
| `contact_email` | json | 否   | string | 主要联系人邮箱，邮箱格式，最大 255 个字符 |
| `remark`        | json | 否   | string | 总网内部档案备注，最大 1000 个字符    |


请求示例：

```json
{
  "operator_id": 1,
  "company_name": "Example Company",
  "contact_name": "Ryan",
  "contact_email": "ryan@example.com",
  "remark": "内部档案备注"
}
```

响应：

```json
{
  "code": 0,
  "msg": "ok",
  "data": {
    "id": 1
  }
}
```



### POST /admin/operator/profile/update

修改分站档案。

#### 请求参数


| 字段              | 位置   | 必填  | 类型     | 说明                      |
| --------------- | ---- | --- | ------ | ----------------------- |
| `operator_id`   | json | 是   | int64  | 分站 ID，大于 `0`            |
| `company_name`  | json | 否   | string | 公司名称，最大 200 个字符         |
| `contact_name`  | json | 否   | string | 主要联系人名称，最大 100 个字符      |
| `contact_email` | json | 否   | string | 主要联系人邮箱，邮箱格式，最大 255 个字符 |
| `remark`        | json | 否   | string | 总网内部档案备注，最大 1000 个字符    |


除 `operator_id` 外，至少需要传一个需要修改的字段。

请求示例：

```json
{
  "operator_id": 1,
  "contact_name": "Ryan Chen",
  "contact_email": "ryan.chen@example.com"
}
```

响应：

```json
{
  "code": 0,
  "msg": "ok",
  "data": {}
}
```



### GET /admin/operator/profile/get

获取分站档案。

Query 示例：

```text
/admin/operator/profile/get?operator_id=1
```



#### 请求参数


| 字段            | 位置    | 必填  | 类型    | 说明           |
| ------------- | ----- | --- | ----- | ------------ |
| `operator_id` | query | 是   | int64 | 分站 ID，大于 `0` |


已创建档案时：

```json
{
  "code": 0,
  "msg": "ok",
  "data": {
    "profile": {
      "id": 1,
      "operator_id": 1,
      "company_name": "Example Company",
      "contact_name": "Ryan",
      "contact_email": "ryan@example.com",
      "remark": "内部档案备注",
      "created_at": 1757836800000,
      "updated_at": 1757836800000
    }
  }
}
```

分站存在但尚未创建档案时：

```json
{
  "code": 0,
  "msg": "ok",
  "data": {
    "profile": null
  }
}
```

---



## 七、分站域名

域名只传主机名，不包含协议、路径和端口。

正确示例：

```text
bo.example.com
agent.example.com
m.example.com
```

不要传：

```text
https://bo.example.com
bo.example.com:443
https://bo.example.com/login
```

后端会自动去除首尾空格并转换为小写。

### OperatorDomainInfo


| 字段            | 类型            | 说明                          |
| ------------- | ------------- | --------------------------- |
| `id`          | int64         | 域名 ID                       |
| `operator_id` | int64         | 分站 ID                       |
| `domain_name` | string        | 域名，不包含协议和端口                 |
| `domain_type` | int64         | `1` 分站后台，`2` 代理后台，`3` 会员 H5 |
| `status`      | int64         | `1` 启用，`2` 停用               |
| `remark`      | string / null | 总网内部备注                      |
| `created_at`  | int64         | 创建时间，Unix 毫秒时间戳             |
| `updated_at`  | int64         | 更新时间，Unix 毫秒时间戳             |




### 域名约束

当前后端约束包括：

- 同一个分站不能重复保存同一个域名。
- 同一个启用中的域名不能同时绑定多个分站。
- 同一个分站的同一种域名类型只能存在一个启用中的域名。
- `domain_name` 必须是合法 FQDN。
- 创建时不传 `status`，默认按启用状态处理。



### POST /admin/operator/domain/create

创建分站域名。

#### 请求参数


| 字段            | 位置   | 必填  | 类型     | 说明                          |
| ------------- | ---- | --- | ------ | --------------------------- |
| `operator_id` | json | 是   | int64  | 分站 ID，大于 `0`                |
| `domain_name` | json | 是   | string | FQDN，最大 253 个字符，不含协议和端口     |
| `domain_type` | json | 是   | int64  | `1` 分站后台，`2` 代理后台，`3` 会员 H5 |
| `status`      | json | 否   | int64  | `1` 启用，`2` 停用               |
| `remark`      | json | 否   | string | 总网内部备注，最大 1000 个字符          |


请求示例：

```json
{
  "operator_id": 1,
  "domain_name": "bo.example.com",
  "domain_type": 1,
  "status": 1,
  "remark": "分站后台域名"
}
```

响应：

```json
{
  "code": 0,
  "msg": "ok",
  "data": {
    "id": 1
  }
}
```



### POST /admin/operator/domain/update

修改分站域名。

#### 请求参数


| 字段            | 位置   | 必填  | 类型     | 说明                          |
| ------------- | ---- | --- | ------ | --------------------------- |
| `id`          | json | 是   | int64  | 域名 ID，大于 `0`                |
| `domain_name` | json | 否   | string | FQDN，最大 253 个字符             |
| `domain_type` | json | 否   | int64  | `1` 分站后台，`2` 代理后台，`3` 会员 H5 |
| `status`      | json | 否   | int64  | `1` 启用，`2` 停用               |
| `remark`      | json | 否   | string | 总网内部备注，最大 1000 个字符          |


请求示例：

```json
{
  "id": 1,
  "status": 2,
  "remark": "暂时停用"
}
```

响应：

```json
{
  "code": 0,
  "msg": "ok",
  "data": {}
}
```



### GET /admin/operator/domain/get

按域名 ID 获取详情。

Query 示例：

```text
/admin/operator/domain/get?id=1
```



#### 请求参数


| 字段   | 位置    | 必填  | 类型    | 说明           |
| ---- | ----- | --- | ----- | ------------ |
| `id` | query | 是   | int64 | 域名 ID，大于 `0` |


响应：

```json
{
  "code": 0,
  "msg": "ok",
  "data": {
    "id": 1,
    "operator_id": 1,
    "domain_name": "bo.example.com",
    "domain_type": 1,
    "status": 1,
    "remark": "分站后台域名",
    "created_at": 1757836800000,
    "updated_at": 1757836800000
  }
}
```



### GET /admin/operator/domain/list

获取分站域名分页列表。

Query 示例：

```text
/admin/operator/domain/list?page=1&page_size=20&operator_id=1&domain_type=1&status=1&keyword=example
```



#### 请求参数


| 字段            | 位置    | 必填  | 类型     | 说明                          |
| ------------- | ----- | --- | ------ | --------------------------- |
| `page`        | query | 是   | int64  | 页码，从 `1` 开始                 |
| `page_size`   | query | 是   | int64  | 每页数量，范围 `1-100`             |
| `operator_id` | query | 否   | int64  | 分站 ID，大于 `0`                |
| `keyword`     | query | 否   | string | 域名关键字，最大 253 个字符            |
| `domain_type` | query | 否   | int64  | `1` 分站后台，`2` 代理后台，`3` 会员 H5 |
| `status`      | query | 否   | int64  | `1` 启用，`2` 停用               |


响应：

```json
{
  "code": 0,
  "msg": "ok",
  "data": {
    "total": 1,
    "list": [
      {
        "id": 1,
        "operator_id": 1,
        "domain_name": "bo.example.com",
        "domain_type": 1,
        "status": 1,
        "remark": "分站后台域名",
        "created_at": 1757836800000,
        "updated_at": 1757836800000
      }
    ]
  }
}
```



### POST /admin/operator/domain/delete

删除分站域名。

#### 请求参数


| 字段   | 位置   | 必填  | 类型    | 说明           |
| ---- | ---- | --- | ----- | ------------ |
| `id` | json | 是   | int64 | 域名 ID，大于 `0` |


请求示例：

```json
{
  "id": 1
}
```

响应：

```json
{
  "code": 0,
  "msg": "ok",
  "data": {}
}
```

---



## 八、分站管理员

管理员账号挂在具体分站下。同一个分站内 `username` 唯一。列表和详情均不返回密码。

明文密码长度 `6-32`，入库前由服务端哈希；表字段长度留给哈希值，前端不要按 64 位明文理解。

可通过更新接口修改账号、密码、显示名称和启停状态。当前不提供按 ID 删除单条账号、详情查询接口。删除分站时会按 `operator_id` 一并清理该分站下全部管理员账号，无需单独调用删除接口。

### OperatorAdminInfo


| 字段              | 类型     | 说明              |
| --------------- | ------ | --------------- |
| `id`            | int64  | 管理员 ID          |
| `operator_id`   | int64  | 分站 ID           |
| `operator_name` | string | 分站名称            |
| `username`      | string | 账号              |
| `display_name`  | string | 显示名称            |
| `status`        | int64  | `1` 启用，`2` 停用   |
| `created_at`    | int64  | 创建时间，Unix 毫秒时间戳 |
| `updated_at`    | int64  | 更新时间，Unix 毫秒时间戳 |




### 管理员约束

- 同一个分站下账号不可重复。
- 创建时不传 `status`，默认按启用状态处理。
- 列表、创建响应均不包含 `password`。



### POST /admin/operator/admin/create

创建分站管理员。

#### 请求参数


| 字段             | 位置   | 必填  | 类型     | 说明              |
| -------------- | ---- | --- | ------ | --------------- |
| `operator_id`  | json | 是   | int64  | 分站 ID，大于 `0`    |
| `username`     | json | 是   | string | 账号，最大 64 个字符    |
| `password`     | json | 是   | string | 明文密码，长度 `6-32`  |
| `display_name` | json | 是   | string | 显示名称，最大 100 个字符 |
| `status`       | json | 否   | int64  | `1` 启用，`2` 停用   |


请求示例：

```json
{
  "operator_id": 1,
  "username": "opadmin",
  "password": "Passw0rd",
  "display_name": "分站管理员",
  "status": 1
}
```

响应：

```json
{
  "code": 0,
  "msg": "ok",
  "data": {
    "id": 1
  }
}
```



### GET /admin/operator/admin/list

获取分站管理员分页列表。

Query 示例：

```text
/admin/operator/admin/list?page=1&page_size=20&operator_id=1&status=1&keyword=admin
```



#### 请求参数


| 字段            | 位置    | 必填  | 类型     | 说明                   |
| ------------- | ----- | --- | ------ | -------------------- |
| `page`        | query | 是   | int64  | 页码，从 `1` 开始          |
| `page_size`   | query | 是   | int64  | 每页数量，范围 `1-100`      |
| `operator_id` | query | 否   | int64  | 分站 ID，大于 `0`         |
| `keyword`     | query | 否   | string | 匹配账号或显示名称，最大 100 个字符 |
| `status`      | query | 否   | int64  | `1` 启用，`2` 停用        |


响应：

```json
{
  "code": 0,
  "msg": "ok",
  "data": {
    "total": 1,
    "list": [
      {
        "id": 1,
        "operator_id": 1,
        "operator_name": "示例分站",
        "username": "opadmin",
        "display_name": "分站管理员",
        "status": 1,
        "created_at": 1757836800000,
        "updated_at": 1757836800000
      }
    ]
  }
}
```



### POST /admin/operator/admin/update

更新分站管理员。至少需要传入 `username`、`password`、`display_name`、`status` 中的一项。不修改所属分站。

#### 请求参数


| 字段             | 位置   | 必填  | 类型     | 说明              |
| -------------- | ---- | --- | ------ | --------------- |
| `id`           | json | 是   | int64  | 管理员 ID，大于 `0`   |
| `username`     | json | 否   | string | 账号，最大 64 个字符    |
| `password`     | json | 否   | string | 明文密码，长度 `6-32`  |
| `display_name` | json | 否   | string | 显示名称，最大 100 个字符 |
| `status`       | json | 否   | int64  | `1` 启用，`2` 停用   |


请求示例：

```json
{
  "id": 1,
  "username": "opadmin",
  "password": "Passw0rd",
  "display_name": "分站管理员",
  "status": 1
}
```

响应：

```json
{
  "code": 0,
  "msg": "ok",
  "data": {}
}
```



### POST /admin/operator/admin/resetPassword

重置分站管理员密码。

#### 请求参数


| 字段         | 位置   | 必填  | 类型     | 说明              |
| ---------- | ---- | --- | ------ | --------------- |
| `id`       | json | 是   | int64  | 管理员 ID，大于 `0`   |
| `password` | json | 是   | string | 新明文密码，长度 `6-32` |


请求示例：

```json
{
  "id": 1,
  "password": "NewPass1"
}
```

响应：

```json
{
  "code": 0,
  "msg": "ok",
  "data": {}
}
```



### POST /admin/operator/admin/updateStatus

更新分站管理员启停状态。

#### 请求参数


| 字段       | 位置   | 必填  | 类型    | 说明            |
| -------- | ---- | --- | ----- | ------------- |
| `id`     | json | 是   | int64 | 管理员 ID，大于 `0` |
| `status` | json | 是   | int64 | `1` 启用，`2` 停用 |


请求示例：

```json
{
  "id": 1,
  "status": 2
}
```

响应：

```json
{
  "code": 0,
  "msg": "ok",
  "data": {}
}
```

---



## 九、基础资源分配

基础资源分配页面当前由三类资源组成：

```text
基础资源分配
├─ 语言
├─ 经营地区
└─ 代理子线路
```

同时提供一个汇总列表，用于展示每个分站当前已经分配的资源数量。

### 保存接口说明

语言、经营地区、代理子线路三个 `/save` 接口都按“当前最终选择结果”保存。

例如分站原来分配：

```text
A, B, C
```

请求保存：

```text
A, C, D
```

保存后最终结果就是：

```text
A, C, D
```

也就是说：

- 不再传的旧数据会被取消分配。
- 新传的数据会新增分配。
- 重复编码会被后端去重。
- 传空数组 `[]` 表示清空当前这一类全部分配。

因此前端保存时应提交当前页面完整选中集合，不是只提交本次新增项。

### GET /admin/operator/basic-resource-allocation/list

获取各分站基础资源分配汇总列表。

#### BasicResourceAllocationInfo


| 字段                 | 类型     | 说明         |
| ------------------ | ------ | ---------- |
| `operator_id`      | int64  | 分站 ID      |
| `operator_code`    | string | 分站业务编码     |
| `operator_name`    | string | 分站名称       |
| `language_count`   | int64  | 已分配语言数量    |
| `region_count`     | int64  | 已分配经营地区数量  |
| `agent_line_count` | int64  | 已分配代理子线路数量 |


Query 示例：

```text
/admin/operator/basic-resource-allocation/list?page=1&page_size=20&keyword=A01
```



#### 请求参数


| 字段          | 位置    | 必填  | 类型     | 说明                   |
| ----------- | ----- | --- | ------ | -------------------- |
| `page`      | query | 是   | int64  | 页码，从 `1` 开始          |
| `page_size` | query | 是   | int64  | 每页数量，范围 `1-100`      |
| `keyword`   | query | 否   | string | 匹配分站编码或名称，最大 100 个字符 |


响应：

```json
{
  "code": 0,
  "msg": "ok",
  "data": {
    "total": 1,
    "list": [
      {
        "operator_id": 1,
        "operator_code": "OP_8D7091378B244D89A51FB102251489F1",
        "operator_name": "A01",
        "language_count": 2,
        "region_count": 3,
        "agent_line_count": 2
      }
    ]
  }
}
```

---



## 十、语言分配



### LanguageAllocationInfo


| 字段              | 类型     | 说明              |
| --------------- | ------ | --------------- |
| `language_code` | string | 语言编码            |
| `allocated_at`  | int64  | 分配时间，Unix 毫秒时间戳 |




### GET /admin/operator/language-allocation/list

获取某个分站当前已经分配的语言。

Query 示例：

```text
/admin/operator/language-allocation/list?operator_id=1
```



#### 请求参数


| 字段            | 位置    | 必填  | 类型    | 说明           |
| ------------- | ----- | --- | ----- | ------------ |
| `operator_id` | query | 是   | int64 | 分站 ID，大于 `0` |


响应：

```json
{
  "code": 0,
  "msg": "ok",
  "data": {
    "list": [
      {
        "language_code": "zh-CN",
        "allocated_at": 1757836800000
      },
      {
        "language_code": "en-US",
        "allocated_at": 1757836800000
      }
    ]
  }
}
```

说明：

- 此接口只返回当前已分配语言。
- 前端可使用 Core 的 `GET /admin/i18n/lang/enabled` 获取当前可选择的启用语言，再与本接口结果组合展示。



### POST /admin/operator/language-allocation/save

保存分站当前语言分配。

#### 请求参数


| 字段               | 位置   | 必填  | 类型       | 说明                           |
| ---------------- | ---- | --- | -------- | ---------------------------- |
| `operator_id`    | json | 是   | int64    | 分站 ID，大于 `0`                 |
| `language_codes` | json | 否   | string[] | 当前最终选中的语言编码，每项非空，单项最大 16 个字符 |


请求示例：

```json
{
  "operator_id": 1,
  "language_codes": [
    "zh-CN",
    "en-US"
  ]
}
```

清空全部语言分配：

```json
{
  "operator_id": 1,
  "language_codes": []
}
```



#### 业务规则

- 传入语言必须属于 Core 当前启用语言。
- 语言编码比较时不区分大小写，后端最终保存 Core 返回的标准语言编码。
- 保存采用全量覆盖语义。

响应：

```json
{
  "code": 0,
  "msg": "ok",
  "data": {}
}
```

---



## 十一、经营地区分配



### RegionAllocationInfo


| 字段             | 类型     | 说明              |
| -------------- | ------ | --------------- |
| `region_code`  | string | 国家地区编码          |
| `allocated_at` | int64  | 分配时间，Unix 毫秒时间戳 |




### GET /admin/operator/region-allocation/list

获取某个分站当前已经分配的经营地区。

Query 示例：

```text
/admin/operator/region-allocation/list?operator_id=1
```



#### 请求参数


| 字段            | 位置    | 必填  | 类型    | 说明           |
| ------------- | ----- | --- | ----- | ------------ |
| `operator_id` | query | 是   | int64 | 分站 ID，大于 `0` |


响应：

```json
{
  "code": 0,
  "msg": "ok",
  "data": {
    "list": [
      {
        "region_code": "JP",
        "allocated_at": 1757836800000
      },
      {
        "region_code": "TH",
        "allocated_at": 1757836800000
      }
    ]
  }
}
```

说明：

- 此接口只返回当前已分配国家地区。
- 前端可使用 `platform-base` 的 `GET /admin/region/list-all?status=1` 获取当前可选择的启用国家地区。



### POST /admin/operator/region-allocation/save

保存分站当前经营地区分配。

#### 请求参数


| 字段             | 位置   | 必填  | 类型       | 说明                       |
| -------------- | ---- | --- | -------- | ------------------------ |
| `operator_id`  | json | 是   | int64    | 分站 ID，大于 `0`             |
| `region_codes` | json | 否   | string[] | 当前最终选中的国家地区编码，每项固定 2 个字符 |


请求示例：

```json
{
  "operator_id": 1,
  "region_codes": [
    "JP",
    "TH"
  ]
}
```

清空全部经营地区分配：

```json
{
  "operator_id": 1,
  "region_codes": []
}
```



#### 业务规则

- 传入地区必须属于 `platform-base` 当前启用的国家地区。
- 国家地区编码会自动去除首尾空格并转换为大写。
- 保存采用全量覆盖语义。

响应：

```json
{
  "code": 0,
  "msg": "ok",
  "data": {}
}
```

---



## 十二、代理子线路分配

代理子线路目前不依赖基础数据表，后端使用固定编码。

当前支持：


| 编码                | 说明     |
| ----------------- | ------ |
| `CASH_PRODUCTION` | 现金正式线路 |
| `CASH_DEMO`       | 现金试玩线路 |
| `CASH_TEST`       | 现金测试线路 |
| `CREDIT_DEMO`     | 信誉试玩线路 |
| `CREDIT_TEST`     | 信誉测试线路 |




### AgentLineAllocationInfo


| 字段                | 类型           | 说明          |
| ----------------- | ------------ | ----------- |
| `agent_line_code` | string       | 代理子线路编码     |
| `allocated`       | bool         | 是否已经分配给当前分站 |
| `allocated_at`    | int64 / null | 分配时间，未分配时为空 |




### GET /admin/operator/agent-line-allocation/list

获取某个分站的代理子线路分配状态。

与语言、经营地区接口不同，此接口会返回全部支持的代理子线路，并通过 `allocated` 标识当前是否已分配。

Query 示例：

```text
/admin/operator/agent-line-allocation/list?operator_id=1
```



#### 请求参数


| 字段            | 位置    | 必填  | 类型    | 说明           |
| ------------- | ----- | --- | ----- | ------------ |
| `operator_id` | query | 是   | int64 | 分站 ID，大于 `0` |


响应示例：

```json
{
  "code": 0,
  "msg": "ok",
  "data": {
    "list": [
      {
        "agent_line_code": "CASH_PRODUCTION",
        "allocated": true,
        "allocated_at": 1757836800000
      },
      {
        "agent_line_code": "CASH_DEMO",
        "allocated": true,
        "allocated_at": 1757836800000
      },
      {
        "agent_line_code": "CASH_TEST",
        "allocated": false,
        "allocated_at": null
      },
      {
        "agent_line_code": "CREDIT_DEMO",
        "allocated": false,
        "allocated_at": null
      },
      {
        "agent_line_code": "CREDIT_TEST",
        "allocated": false,
        "allocated_at": null
      }
    ]
  }
}
```



### POST /admin/operator/agent-line-allocation/save

保存分站当前代理子线路分配。

#### 请求参数


| 字段                 | 位置   | 必填  | 类型       | 说明                         |
| ------------------ | ---- | --- | -------- | -------------------------- |
| `operator_id`      | json | 是   | int64    | 分站 ID，大于 `0`               |
| `agent_line_codes` | json | 否   | string[] | 当前最终选中的代理子线路编码，每项最大 32 个字符 |


请求示例：

```json
{
  "operator_id": 1,
  "agent_line_codes": [
    "CASH_PRODUCTION",
    "CASH_DEMO"
  ]
}
```

清空全部代理子线路分配：

```json
{
  "operator_id": 1,
  "agent_line_codes": []
}
```



#### 业务规则

- 只允许传当前后端支持的固定代理子线路编码。
- 编码必须使用上表中的标准值。
- 重复编码会自动去重。
- 保存采用全量覆盖语义。

响应：

```json
{
  "code": 0,
  "msg": "ok",
  "data": {}
}
```

---



## 十三、前端页面对接建议

当前接口和页面可以按下面方式对应：


| 页面           | 主要接口                                                                                                          |
| ------------ | ------------------------------------------------------------------------------------------------------------- |
| 分站列表         | `GET /admin/operator/list`、`POST /admin/operator/delete`                                                      |
| 创建分站-基础信息    | `POST /admin/operator/create`、`POST /admin/operator/update`                                                   |
| 创建分站-档案信息    | `GET /admin/operator/profile/get`、`POST /admin/operator/profile/create`、`POST /admin/operator/profile/update` |
| 创建分站-域名信息    | `/admin/operator/domain/*`                                                                                    |
| 创建完成         | `POST /admin/operator/complete`                                                                               |
| 分站详情         | `GET /admin/operator/get` + 档案/域名/资源分配接口                                                                      |
| 域名管理         | `/admin/operator/domain/*`                                                                                    |
| 管理员账号        | `/admin/operator/admin/*`                                                                                     |
| 基础资源分配列表     | `GET /admin/operator/basic-resource-allocation/list`                                                          |
| 基础资源分配-语言    | `/admin/operator/language-allocation/*`                                                                       |
| 基础资源分配-经营地区  | `/admin/operator/region-allocation/*`                                                                         |
| 基础资源分配-代理子线路 | `/admin/operator/agent-line-allocation/*`                                                                     |


基础资源分配页面的三个 Tab：

```text
基础资源分配
├─ 语言
├─ 经营地区
└─ 代理子线路
```

不需要拆成三个独立左侧菜单页面。