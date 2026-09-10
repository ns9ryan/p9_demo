# platform-base API 文档

`platform-base-api` 提供 P9 总网系统基础数据相关的 HTTP API。

当前包含：

- 时区
- 货币
- 国家地区
- 服务健康检查

语言管理已从 `platform-base` 移除，统一由 Core 管理。语言相关接口请参考同目录的 [core-api.md](./core-api.md)。

---

## 一、接口约定

### 服务地址

默认开发端口：

```text
18001
```

实际 Host 和域名以部署环境为准。

### 鉴权

除 `GET /ping` 外，当前所有 `/admin/*` 接口均需要后台登录鉴权和接口权限校验。

请求头：

```http
Authorization: Bearer <access_token>
```

Token 获取、刷新和过期处理规则见同目录的 [core-api.md](./core-api.md)。

### 多语言

请求语言通过 `X-Lang` Header 指定：

```http
X-Lang: zh-CN
```

当前使用的语言编码：

| 语言编码 | 说明 |
| --- | --- |
| `zh-CN` | 简体中文 |
| `zh-HK` | 繁体中文 |
| `en-US` | 英文 |

默认语言为 `zh-CN`。

未传 `X-Lang` 时使用默认语言。

### 名称字段

时区、货币、国家地区统一返回：

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `name_key` | string | 名称翻译 Key，长期正式字段 |
| `name` | string | 当前请求语言对应的名称，开发阶段用于调试 |

例如：

```json
{
  "name_key": "timezone.asia_tokyo.name",
  "name": "日本标准时间"
}
```

当前前端可以直接使用 `name` 展示。

`name` 为开发阶段过渡字段，后续前端多语言接入稳定后可能移除。长期应以 `name_key` 作为多语言名称的稳定引用。

### 基础数据创建

当前 `platform-base` 不提供时区、货币、国家地区的创建接口。

这些数据属于系统基础数据，由后端初始化维护。前端管理页面只提供：

- 修改允许变更的字段
- 查询详情
- 分页列表
- 全部列表
- 调整排序

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

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `code` | int | 成功固定为 `0` |
| `msg` | string | 成功固定为 `ok` |
| `data` | object | 接口业务数据 |

没有业务数据返回的接口，当前 `data` 为：

```json
{}
```

### 错误响应

失败时 `code` 使用对应 HTTP 状态码：

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
    "rpc": "platformbase.TimezoneService/Update",
    "error": "..."
  }
}
```

`debug` 仅用于开发调试，前端业务逻辑不应依赖该字段。

### 分页

管理列表接口统一使用：

| 字段 | 位置 | 必填 | 类型 | 说明 |
| --- | --- | --- | --- | --- |
| `page` | query | 是 | int64 | 页码，从 `1` 开始 |
| `page_size` | query | 是 | int64 | 每页数量，范围 `1-100` |

### 公共枚举

| 字段 | 值 | 说明 |
| --- | --- | --- |
| `status` | `1` | 启用 |
| `status` | `2` | 停用 |
| `currency_type` | `1` | 法定货币 |
| `currency_type` | `2` | 虚拟货币 |

---

## 二、接口索引

### 公开接口

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| GET | `/ping` | 服务健康检查 |

### 时区

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| POST | `/admin/timezone/update` | 修改时区 |
| GET | `/admin/timezone/get` | 获取时区 |
| GET | `/admin/timezone/list` | 获取时区管理列表 |
| GET | `/admin/timezone/list-all` | 获取全部时区 |
| POST | `/admin/timezone/reorder` | 调整时区排序 |

### 货币

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| POST | `/admin/currency/update` | 修改货币 |
| GET | `/admin/currency/get` | 获取货币 |
| GET | `/admin/currency/list` | 获取货币管理列表 |
| GET | `/admin/currency/list-all` | 获取全部货币 |
| POST | `/admin/currency/reorder` | 调整货币排序 |

### 国家地区

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| POST | `/admin/region/update` | 修改国家地区 |
| GET | `/admin/region/get` | 获取国家地区 |
| GET | `/admin/region/list` | 获取国家地区管理列表 |
| GET | `/admin/region/list-all` | 获取全部国家地区 |
| POST | `/admin/region/reorder` | 调整国家地区排序 |

---

## 三、Ping

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

## 四、时区

时区属于系统基础数据，不提供创建接口。

### TimezoneInfo

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `id` | int64 | 时区 ID |
| `code` | string | IANA 时区编码 |
| `name_key` | string | 名称翻译 Key |
| `name` | string | 当前语言名称，开发阶段用于调试 |
| `status` | int64 | 状态：`1` 启用，`2` 停用 |
| `sort_no` | int64 | 排序值，数值越小越靠前 |

示例：

```json
{
  "id": 1,
  "code": "Asia/Tokyo",
  "name_key": "timezone.asia_tokyo.name",
  "name": "日本标准时间",
  "status": 1,
  "sort_no": 10
}
```

### POST /admin/timezone/update

修改时区。

当前只允许修改 `status`。

| 字段 | 位置 | 必填 | 类型 | 说明 |
| --- | --- | --- | --- | --- |
| `id` | json | 是 | int64 | 时区 ID，大于 `0` |
| `status` | json | 是 | int64 | `1` 启用，`2` 停用 |

说明：`code`、`name_key` 不允许通过接口修改。时区当前只有 `status` 一个可修改字段，因此请求必须传 `status`。

示例：

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

### GET /admin/timezone/get

Query 示例：

```text
/admin/timezone/get?id=1
```

| 字段 | 位置 | 必填 | 类型 | 说明 |
| --- | --- | --- | --- | --- |
| `id` | query | 是 | int64 | 时区 ID，大于 `0` |

响应：

```json
{
  "code": 0,
  "msg": "ok",
  "data": {
    "id": 1,
    "code": "Asia/Tokyo",
    "name_key": "timezone.asia_tokyo.name",
    "name": "日本标准时间",
    "status": 1,
    "sort_no": 10
  }
}
```

### GET /admin/timezone/list

Query 示例：

```text
/admin/timezone/list?page=1&page_size=20&status=1
```

| 字段 | 位置 | 必填 | 类型 | 说明 |
| --- | --- | --- | --- | --- |
| `page` | query | 是 | int64 | 页码，从 `1` 开始 |
| `page_size` | query | 是 | int64 | 每页数量，范围 `1-100` |
| `status` | query | 否 | int64 | `1` 启用，`2` 停用 |

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
        "code": "Asia/Tokyo",
        "name_key": "timezone.asia_tokyo.name",
        "name": "日本标准时间",
        "status": 1,
        "sort_no": 10
      }
    ]
  }
}
```

### GET /admin/timezone/list-all

不分页获取全部时区，可按状态筛选。

```text
/admin/timezone/list-all?status=1
```

| 字段 | 位置 | 必填 | 类型 | 说明 |
| --- | --- | --- | --- | --- |
| `status` | query | 否 | int64 | `1` 启用，`2` 停用 |

响应结构与分页列表中的 `list` 项一致，不返回 `total`。

### POST /admin/timezone/reorder

| 字段 | 位置 | 必填 | 类型 | 说明 |
| --- | --- | --- | --- | --- |
| `id` | json | 是 | int64 | 需要移动的时区 ID，大于 `0` |
| `target_id` | json | 是 | int64 | 目标时区 ID，大于 `0` |

`id` 和 `target_id` 不能相同。

示例：

```json
{
  "id": 3,
  "target_id": 1
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

## 五、货币

货币属于系统基础数据，不提供创建接口。

### CurrencyInfo

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `id` | int64 | 货币 ID |
| `code` | string | 货币编码 |
| `name_key` | string | 名称翻译 Key |
| `name` | string | 当前语言名称，开发阶段用于调试 |
| `currency_type` | int64 | 货币类型：`1` 法定货币，`2` 虚拟货币 |
| `symbol` | string | 货币符号 |
| `amount_factor` | int64 | 金额换算倍率，例如 `USD=100`、`VND=1` |
| `status` | int64 | 状态：`1` 启用，`2` 停用 |
| `sort_no` | int64 | 排序值，数值越小越靠前 |

### POST /admin/currency/update

当前允许修改 `symbol`、`status`。

| 字段 | 位置 | 必填 | 类型 | 说明 |
| --- | --- | --- | --- | --- |
| `id` | json | 是 | int64 | 货币 ID，大于 `0` |
| `symbol` | json | 否 | string | 货币符号，最大 16 个字符 |
| `status` | json | 否 | int64 | `1` 启用，`2` 停用 |

说明：

- `symbol` 和 `status` 至少传一个
- `symbol` 传入时不能为空字符串
- `code`、`name_key`、`currency_type`、`amount_factor` 不允许通过接口修改

示例：

```json
{
  "id": 1,
  "symbol": "$",
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

### GET /admin/currency/get

```text
/admin/currency/get?id=1
```

| 字段 | 位置 | 必填 | 类型 | 说明 |
| --- | --- | --- | --- | --- |
| `id` | query | 是 | int64 | 货币 ID，大于 `0` |

响应：

```json
{
  "code": 0,
  "msg": "ok",
  "data": {
    "id": 1,
    "code": "USD",
    "name_key": "currency.usd.name",
    "name": "美元",
    "currency_type": 1,
    "symbol": "$",
    "amount_factor": 100,
    "status": 1,
    "sort_no": 10
  }
}
```

### GET /admin/currency/list

```text
/admin/currency/list?page=1&page_size=20&status=1
```

| 字段 | 位置 | 必填 | 类型 | 说明 |
| --- | --- | --- | --- | --- |
| `page` | query | 是 | int64 | 页码，从 `1` 开始 |
| `page_size` | query | 是 | int64 | 每页数量，范围 `1-100` |
| `status` | query | 否 | int64 | `1` 启用，`2` 停用 |

响应中的 `list` 项结构为 `CurrencyInfo`，并返回 `total`。

### GET /admin/currency/list-all

不分页获取全部货币，可按状态筛选。

```text
/admin/currency/list-all?status=1
```

| 字段 | 位置 | 必填 | 类型 | 说明 |
| --- | --- | --- | --- | --- |
| `status` | query | 否 | int64 | `1` 启用，`2` 停用 |

响应中的 `list` 项结构为 `CurrencyInfo`，不返回 `total`。

### POST /admin/currency/reorder

| 字段 | 位置 | 必填 | 类型 | 说明 |
| --- | --- | --- | --- | --- |
| `id` | json | 是 | int64 | 需要移动的货币 ID，大于 `0` |
| `target_id` | json | 是 | int64 | 目标货币 ID，大于 `0` |

`id` 和 `target_id` 不能相同。

示例：

```json
{
  "id": 3,
  "target_id": 1
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

## 六、国家地区

国家地区属于系统基础数据，不提供创建接口。

### RegionInfo

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `id` | int64 | 国家或地区 ID |
| `code` | string | 国家或地区编码 |
| `calling_code` | string | 国际电话区号，不包含 `+` |
| `name_key` | string | 名称翻译 Key |
| `name` | string | 当前语言名称，开发阶段用于调试 |
| `status` | int64 | 状态：`1` 启用，`2` 停用 |
| `sort_no` | int64 | 排序值，数值越小越靠前 |

### POST /admin/region/update

当前允许修改 `calling_code`、`status`。

| 字段 | 位置 | 必填 | 类型 | 说明 |
| --- | --- | --- | --- | --- |
| `id` | json | 是 | int64 | 国家或地区 ID，大于 `0` |
| `calling_code` | json | 否 | string | 国际电话区号，不包含 `+`，最多 3 位数字 |
| `status` | json | 否 | int64 | `1` 启用，`2` 停用 |

说明：

- `calling_code` 和 `status` 至少传一个
- `calling_code` 使用字符串传递，例如 `"81"`
- `code`、`name_key` 不允许通过接口修改

示例：

```json
{
  "id": 1,
  "calling_code": "81",
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

### GET /admin/region/get

```text
/admin/region/get?id=1
```

| 字段 | 位置 | 必填 | 类型 | 说明 |
| --- | --- | --- | --- | --- |
| `id` | query | 是 | int64 | 国家或地区 ID，大于 `0` |

响应：

```json
{
  "code": 0,
  "msg": "ok",
  "data": {
    "id": 1,
    "code": "JP",
    "calling_code": "81",
    "name_key": "region.jp.name",
    "name": "日本",
    "status": 1,
    "sort_no": 10
  }
}
```

### GET /admin/region/list

```text
/admin/region/list?page=1&page_size=20&status=1
```

| 字段 | 位置 | 必填 | 类型 | 说明 |
| --- | --- | --- | --- | --- |
| `page` | query | 是 | int64 | 页码，从 `1` 开始 |
| `page_size` | query | 是 | int64 | 每页数量，范围 `1-100` |
| `status` | query | 否 | int64 | `1` 启用，`2` 停用 |

响应中的 `list` 项结构为 `RegionInfo`，并返回 `total`。

### GET /admin/region/list-all

不分页获取全部国家地区，可按状态筛选。

```text
/admin/region/list-all?status=1
```

| 字段 | 位置 | 必填 | 类型 | 说明 |
| --- | --- | --- | --- | --- |
| `status` | query | 否 | int64 | `1` 启用，`2` 停用 |

响应中的 `list` 项结构为 `RegionInfo`，不返回 `total`。

### POST /admin/region/reorder

| 字段 | 位置 | 必填 | 类型 | 说明 |
| --- | --- | --- | --- | --- |
| `id` | json | 是 | int64 | 需要移动的国家地区 ID，大于 `0` |
| `target_id` | json | 是 | int64 | 目标国家地区 ID，大于 `0` |

`id` 和 `target_id` 不能相同。

示例：

```json
{
  "id": 3,
  "target_id": 1
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

## 七、前端对接重点

1. `platform-base` 已删除语言管理接口，语言统一由 Core 管理。
2. 时区、货币、国家地区都不再提供创建接口，前端不要展示“创建”按钮。
3. `name_i18n` 已删除，统一改为 `name_key`。
4. 当前接口同时返回 `name`，用于开发阶段展示和调试。
5. `name_key` 是长期稳定字段，前端多语言正式接入后应以它作为名称翻译引用。
6. 时区只允许修改 `status`。
7. 货币只允许修改 `symbol`、`status`。
8. 国家地区只允许修改 `calling_code`、`status`。
9. `code`、`name_key`、`currency_type`、`amount_factor` 等系统基础属性不通过管理 API 修改。
10. 列表统一按 `sort_no` 展示，数值越小越靠前。
