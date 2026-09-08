# platform-base API 文档

`platform-base-api` 提供 P9 总网系统基础数据相关的 HTTP API，目前包含语言、时区、货币、国家地区四个模块，以及服务健康检查接口。

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

当前默认语言为：

```text
zh-CN
```

未传 `X-Lang` 时使用默认语言。

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

`debug` 仅用于开发调试，前端业务逻辑不应依赖该字段。

### 分页

管理列表接口统一使用：

| 字段          | 位置    | 必填 | 类型    | 说明              |
| ----------- | ----- | -- | ----- | --------------- |
| `page`      | query | 是  | int64 | 页码，从 `1` 开始     |
| `page_size` | query | 是  | int64 | 每页数量，范围 `1-100` |

### 公共枚举

| 字段              | 值   | 说明   |
| --------------- | --- | ---- |
| `status`        | `1` | 启用   |
| `status`        | `2` | 停用   |
| `currency_type` | `1` | 法定货币 |
| `currency_type` | `2` | 虚拟货币 |

### 多语言名称

`name_i18n` 使用语言编码作为 Key：

```json
{
  "zh-CN": "美元",
  "en-US": "US Dollar"
}
```

---

## 二、接口索引

### 公开接口

| 方法  | 路径      | 说明     |
| --- | ------- | ------ |
| GET | `/ping` | 服务健康检查 |

### 语言

| 方法   | 路径                         | 说明       |
| ---- | -------------------------- | -------- |
| POST | `/admin/language/create`   | 创建语言     |
| POST | `/admin/language/update`   | 修改语言     |
| GET  | `/admin/language/get`      | 获取语言     |
| GET  | `/admin/language/list`     | 获取语言管理列表 |
| GET  | `/admin/language/list-all` | 获取全部语言   |
| POST | `/admin/language/reorder`  | 调整语言排序   |

### 时区

| 方法   | 路径                         | 说明       |
| ---- | -------------------------- | -------- |
| POST | `/admin/timezone/create`   | 创建时区     |
| POST | `/admin/timezone/update`   | 修改时区     |
| GET  | `/admin/timezone/get`      | 获取时区     |
| GET  | `/admin/timezone/list`     | 获取时区管理列表 |
| GET  | `/admin/timezone/list-all` | 获取全部时区   |
| POST | `/admin/timezone/reorder`  | 调整时区排序   |

### 货币

| 方法   | 路径                         | 说明       |
| ---- | -------------------------- | -------- |
| POST | `/admin/currency/create`   | 创建货币     |
| POST | `/admin/currency/update`   | 修改货币     |
| GET  | `/admin/currency/get`      | 获取货币     |
| GET  | `/admin/currency/list`     | 获取货币管理列表 |
| GET  | `/admin/currency/list-all` | 获取全部货币   |
| POST | `/admin/currency/reorder`  | 调整货币排序   |

### 国家地区

| 方法   | 路径                       | 说明         |
| ---- | ------------------------ | ---------- |
| POST | `/admin/region/create`   | 创建国家地区     |
| POST | `/admin/region/update`   | 修改国家地区     |
| GET  | `/admin/region/get`      | 获取国家地区     |
| GET  | `/admin/region/list`     | 获取国家地区管理列表 |
| GET  | `/admin/region/list-all` | 获取全部国家地区   |
| POST | `/admin/region/reorder`  | 调整国家地区排序   |

---

## 三、Ping

### GET /ping

服务健康检查，不需要 JWT。

请求：

无。

响应：

```json
{
  "code": 0,
  "msg": "ok",
  "data": {}
}
```

---

## 四、语言

### LanguageInfo

| 字段          | 类型     | 说明               |
| ----------- | ------ | ---------------- |
| `id`        | int64  | 语言 ID            |
| `code`      | string | 语言编码             |
| `name_i18n` | object | 多语言名称            |
| `status`    | int64  | 状态：`1` 启用，`2` 停用 |
| `sort_no`   | int64  | 排序值，数值越小越靠前      |

### POST /admin/language/create

创建语言。

请求：

| 字段          | 位置   | 必填 | 类型     | 说明             |
| ----------- | ---- | -- | ------ | -------------- |
| `code`      | json | 是  | string | 语言编码，最大 35 个字符 |
| `name_i18n` | json | 是  | object | 多语言名称          |
| `status`    | json | 否  | int64  | `1` 启用，`2` 停用  |

示例：

```json
{
  "code": "zh-CN",
  "name_i18n": {
    "zh-CN": "简体中文",
    "en-US": "Simplified Chinese"
  },
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

### POST /admin/language/update

修改语言。

当前接口不提供语言编码 `code` 的修改。

请求：

| 字段          | 位置   | 必填 | 类型     | 说明            |
| ----------- | ---- | -- | ------ | ------------- |
| `id`        | json | 是  | int64  | 语言 ID，大于 `0`  |
| `name_i18n` | json | 否  | object | 多语言名称，不传时不修改  |
| `status`    | json | 否  | int64  | `1` 启用，`2` 停用 |

示例：

```json
{
  "id": 1,
  "name_i18n": {
    "zh-CN": "简体中文",
    "en-US": "Simplified Chinese"
  },
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

### GET /admin/language/get

根据 ID 获取语言。

Query 示例：

```text
/admin/language/get?id=1
```

请求：

| 字段   | 位置    | 必填 | 类型    | 说明           |
| ---- | ----- | -- | ----- | ------------ |
| `id` | query | 是  | int64 | 语言 ID，大于 `0` |

响应：

```json
{
  "code": 0,
  "msg": "ok",
  "data": {
    "language": {
      "id": 1,
      "code": "zh-CN",
      "name_i18n": {
        "zh-CN": "简体中文",
        "en-US": "Simplified Chinese"
      },
      "status": 1,
      "sort_no": 10
    }
  }
}
```

### GET /admin/language/list

获取分页语言管理列表。

Query 示例：

```text
/admin/language/list?page=1&page_size=20&status=1
```

请求：

| 字段          | 位置    | 必填 | 类型    | 说明              |
| ----------- | ----- | -- | ----- | --------------- |
| `page`      | query | 是  | int64 | 页码，从 `1` 开始     |
| `page_size` | query | 是  | int64 | 每页数量，范围 `1-100` |
| `status`    | query | 否  | int64 | `1` 启用，`2` 停用   |

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
        "code": "zh-CN",
        "name_i18n": {
          "zh-CN": "简体中文",
          "en-US": "Simplified Chinese"
        },
        "status": 1,
        "sort_no": 10
      }
    ]
  }
}
```

### GET /admin/language/list-all

不分页获取全部语言，可按状态筛选。

Query 示例：

```text
/admin/language/list-all?status=1
```

请求：

| 字段       | 位置    | 必填 | 类型    | 说明            |
| -------- | ----- | -- | ----- | ------------- |
| `status` | query | 否  | int64 | `1` 启用，`2` 停用 |

响应：

```json
{
  "code": 0,
  "msg": "ok",
  "data": {
    "list": [
      {
        "id": 1,
        "code": "zh-CN",
        "name_i18n": {
          "zh-CN": "简体中文",
          "en-US": "Simplified Chinese"
        },
        "status": 1,
        "sort_no": 10
      }
    ]
  }
}
```

### POST /admin/language/reorder

调整语言排序。

请求：

| 字段          | 位置   | 必填 | 类型    | 说明                |
| ----------- | ---- | -- | ----- | ----------------- |
| `id`        | json | 是  | int64 | 需要移动的语言 ID，大于 `0` |
| `target_id` | json | 是  | int64 | 目标语言 ID，大于 `0`    |

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

## 五、时区

### TimezoneInfo

| 字段          | 类型     | 说明               |
| ----------- | ------ | ---------------- |
| `id`        | int64  | 时区 ID            |
| `code`      | string | IANA 时区编码        |
| `name_i18n` | object | 多语言名称            |
| `status`    | int64  | 状态：`1` 启用，`2` 停用 |
| `sort_no`   | int64  | 排序值，数值越小越靠前      |

### POST /admin/timezone/create

创建时区。

请求：

| 字段          | 位置   | 必填 | 类型     | 说明                  |
| ----------- | ---- | -- | ------ | ------------------- |
| `code`      | json | 是  | string | IANA 时区编码，最大 64 个字符 |
| `name_i18n` | json | 是  | object | 多语言名称               |
| `status`    | json | 否  | int64  | `1` 启用，`2` 停用       |

示例：

```json
{
  "code": "Asia/Tokyo",
  "name_i18n": {
    "zh-CN": "东京",
    "en-US": "Tokyo"
  },
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

### POST /admin/timezone/update

修改时区。

当前接口不提供时区编码 `code` 的修改。

请求：

| 字段          | 位置   | 必填 | 类型     | 说明            |
| ----------- | ---- | -- | ------ | ------------- |
| `id`        | json | 是  | int64  | 时区 ID，大于 `0`  |
| `name_i18n` | json | 否  | object | 多语言名称，不传时不修改  |
| `status`    | json | 否  | int64  | `1` 启用，`2` 停用 |

示例：

```json
{
  "id": 1,
  "name_i18n": {
    "zh-CN": "东京",
    "en-US": "Tokyo"
  },
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

### GET /admin/timezone/get

根据 ID 获取时区。

Query 示例：

```text
/admin/timezone/get?id=1
```

请求：

| 字段   | 位置    | 必填 | 类型    | 说明           |
| ---- | ----- | -- | ----- | ------------ |
| `id` | query | 是  | int64 | 时区 ID，大于 `0` |

响应：

```json
{
  "code": 0,
  "msg": "ok",
  "data": {
    "timezone": {
      "id": 1,
      "code": "Asia/Tokyo",
      "name_i18n": {
        "zh-CN": "东京",
        "en-US": "Tokyo"
      },
      "status": 1,
      "sort_no": 10
    }
  }
}
```

### GET /admin/timezone/list

获取分页时区管理列表。

Query 示例：

```text
/admin/timezone/list?page=1&page_size=20&status=1
```

请求：

| 字段          | 位置    | 必填 | 类型    | 说明              |
| ----------- | ----- | -- | ----- | --------------- |
| `page`      | query | 是  | int64 | 页码，从 `1` 开始     |
| `page_size` | query | 是  | int64 | 每页数量，范围 `1-100` |
| `status`    | query | 否  | int64 | `1` 启用，`2` 停用   |

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
        "name_i18n": {
          "zh-CN": "东京",
          "en-US": "Tokyo"
        },
        "status": 1,
        "sort_no": 10
      }
    ]
  }
}
```

### GET /admin/timezone/list-all

不分页获取全部时区，可按状态筛选。

Query 示例：

```text
/admin/timezone/list-all?status=1
```

请求：

| 字段       | 位置    | 必填 | 类型    | 说明            |
| -------- | ----- | -- | ----- | ------------- |
| `status` | query | 否  | int64 | `1` 启用，`2` 停用 |

响应：

```json
{
  "code": 0,
  "msg": "ok",
  "data": {
    "list": [
      {
        "id": 1,
        "code": "Asia/Tokyo",
        "name_i18n": {
          "zh-CN": "东京",
          "en-US": "Tokyo"
        },
        "status": 1,
        "sort_no": 10
      }
    ]
  }
}
```

### POST /admin/timezone/reorder

调整时区排序。

请求：

| 字段          | 位置   | 必填 | 类型    | 说明                |
| ----------- | ---- | -- | ----- | ----------------- |
| `id`        | json | 是  | int64 | 需要移动的时区 ID，大于 `0` |
| `target_id` | json | 是  | int64 | 目标时区 ID，大于 `0`    |

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

## 六、货币

### CurrencyInfo

| 字段              | 类型     | 说明                          |
| --------------- | ------ | --------------------------- |
| `id`            | int64  | 货币 ID                       |
| `code`          | string | 货币编码                        |
| `name_i18n`     | object | 多语言名称                       |
| `currency_type` | int64  | 货币类型：`1` 法定货币，`2` 虚拟货币      |
| `symbol`        | string | 货币符号                        |
| `amount_factor` | int64  | 金额换算倍率，例如 `USD=100`、`VND=1` |
| `status`        | int64  | 状态：`1` 启用，`2` 停用            |
| `sort_no`       | int64  | 排序值，数值越小越靠前                 |

### POST /admin/currency/create

创建货币。

请求：

| 字段              | 位置   | 必填 | 类型     | 说明                |
| --------------- | ---- | -- | ------ | ----------------- |
| `code`          | json | 是  | string | 货币编码，最大 16 个字符    |
| `name_i18n`     | json | 是  | object | 多语言名称             |
| `currency_type` | json | 是  | int64  | `1` 法定货币，`2` 虚拟货币 |
| `symbol`        | json | 是  | string | 货币符号，最大 16 个字符    |
| `amount_factor` | json | 是  | int64  | 金额换算倍率，大于 `0`     |
| `status`        | json | 否  | int64  | `1` 启用，`2` 停用     |

示例：

```json
{
  "code": "USD",
  "name_i18n": {
    "zh-CN": "美元",
    "en-US": "US Dollar"
  },
  "currency_type": 1,
  "symbol": "$",
  "amount_factor": 100,
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

### POST /admin/currency/update

修改货币。

当前接口不提供以下字段的修改：

* `code`
* `currency_type`
* `amount_factor`

请求：

| 字段          | 位置   | 必填 | 类型     | 说明             |
| ----------- | ---- | -- | ------ | -------------- |
| `id`        | json | 是  | int64  | 货币 ID，大于 `0`   |
| `name_i18n` | json | 否  | object | 多语言名称，不传时不修改   |
| `symbol`    | json | 否  | string | 货币符号，最大 16 个字符 |
| `status`    | json | 否  | int64  | `1` 启用，`2` 停用  |

示例：

```json
{
  "id": 1,
  "name_i18n": {
    "zh-CN": "美元",
    "en-US": "US Dollar"
  },
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

根据 ID 获取货币。

Query 示例：

```text
/admin/currency/get?id=1
```

请求：

| 字段   | 位置    | 必填 | 类型    | 说明           |
| ---- | ----- | -- | ----- | ------------ |
| `id` | query | 是  | int64 | 货币 ID，大于 `0` |

响应：

```json
{
  "code": 0,
  "msg": "ok",
  "data": {
    "currency": {
      "id": 1,
      "code": "USD",
      "name_i18n": {
        "zh-CN": "美元",
        "en-US": "US Dollar"
      },
      "currency_type": 1,
      "symbol": "$",
      "amount_factor": 100,
      "status": 1,
      "sort_no": 10
    }
  }
}
```

### GET /admin/currency/list

获取分页货币管理列表。

Query 示例：

```text
/admin/currency/list?page=1&page_size=20&status=1
```

请求：

| 字段          | 位置    | 必填 | 类型    | 说明              |
| ----------- | ----- | -- | ----- | --------------- |
| `page`      | query | 是  | int64 | 页码，从 `1` 开始     |
| `page_size` | query | 是  | int64 | 每页数量，范围 `1-100` |
| `status`    | query | 否  | int64 | `1` 启用，`2` 停用   |

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
        "code": "USD",
        "name_i18n": {
          "zh-CN": "美元",
          "en-US": "US Dollar"
        },
        "currency_type": 1,
        "symbol": "$",
        "amount_factor": 100,
        "status": 1,
        "sort_no": 10
      }
    ]
  }
}
```

### GET /admin/currency/list-all

不分页获取全部货币，可按状态筛选。

Query 示例：

```text
/admin/currency/list-all?status=1
```

请求：

| 字段       | 位置    | 必填 | 类型    | 说明            |
| -------- | ----- | -- | ----- | ------------- |
| `status` | query | 否  | int64 | `1` 启用，`2` 停用 |

响应：

```json
{
  "code": 0,
  "msg": "ok",
  "data": {
    "list": [
      {
        "id": 1,
        "code": "USD",
        "name_i18n": {
          "zh-CN": "美元",
          "en-US": "US Dollar"
        },
        "currency_type": 1,
        "symbol": "$",
        "amount_factor": 100,
        "status": 1,
        "sort_no": 10
      }
    ]
  }
}
```

### POST /admin/currency/reorder

调整货币排序。

请求：

| 字段          | 位置   | 必填 | 类型    | 说明                |
| ----------- | ---- | -- | ----- | ----------------- |
| `id`        | json | 是  | int64 | 需要移动的货币 ID，大于 `0` |
| `target_id` | json | 是  | int64 | 目标货币 ID，大于 `0`    |

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

## 七、国家地区

### RegionInfo

| 字段             | 类型     | 说明               |
| -------------- | ------ | ---------------- |
| `id`           | int64  | 国家或地区 ID         |
| `code`         | string | 国家或地区编码          |
| `calling_code` | string | 国际电话区号，不包含 `+`   |
| `name_i18n`    | object | 多语言名称            |
| `status`       | int64  | 状态：`1` 启用，`2` 停用 |
| `sort_no`      | int64  | 排序值，数值越小越靠前      |

### POST /admin/region/create

创建国家地区。

请求：

| 字段             | 位置   | 必填 | 类型     | 说明                      |
| -------------- | ---- | -- | ------ | ----------------------- |
| `code`         | json | 是  | string | 国家或地区编码，固定 2 个字符        |
| `calling_code` | json | 是  | string | 国际电话区号，不包含 `+`，最大 3 个字符 |
| `name_i18n`    | json | 是  | object | 多语言名称                   |
| `status`       | json | 否  | int64  | `1` 启用，`2` 停用           |

示例：

```json
{
  "code": "JP",
  "calling_code": "81",
  "name_i18n": {
    "zh-CN": "日本",
    "en-US": "Japan"
  },
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

### POST /admin/region/update

修改国家地区。

当前接口不提供国家地区编码 `code` 的修改。

请求：

| 字段             | 位置   | 必填 | 类型     | 说明                      |
| -------------- | ---- | -- | ------ | ----------------------- |
| `id`           | json | 是  | int64  | 国家或地区 ID，大于 `0`         |
| `calling_code` | json | 否  | string | 国际电话区号，不包含 `+`，最大 3 个字符 |
| `name_i18n`    | json | 否  | object | 多语言名称，不传时不修改            |
| `status`       | json | 否  | int64  | `1` 启用，`2` 停用           |

示例：

```json
{
  "id": 1,
  "calling_code": "81",
  "name_i18n": {
    "zh-CN": "日本",
    "en-US": "Japan"
  },
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

根据 ID 获取国家地区。

Query 示例：

```text
/admin/region/get?id=1
```

请求：

| 字段   | 位置    | 必填 | 类型    | 说明              |
| ---- | ----- | -- | ----- | --------------- |
| `id` | query | 是  | int64 | 国家或地区 ID，大于 `0` |

响应：

```json
{
  "code": 0,
  "msg": "ok",
  "data": {
    "region": {
      "id": 1,
      "code": "JP",
      "calling_code": "81",
      "name_i18n": {
        "zh-CN": "日本",
        "en-US": "Japan"
      },
      "status": 1,
      "sort_no": 10
    }
  }
}
```

### GET /admin/region/list

获取分页国家地区管理列表。

Query 示例：

```text
/admin/region/list?page=1&page_size=20&status=1
```

请求：

| 字段          | 位置    | 必填 | 类型    | 说明              |
| ----------- | ----- | -- | ----- | --------------- |
| `page`      | query | 是  | int64 | 页码，从 `1` 开始     |
| `page_size` | query | 是  | int64 | 每页数量，范围 `1-100` |
| `status`    | query | 否  | int64 | `1` 启用，`2` 停用   |

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
        "code": "JP",
        "calling_code": "81",
        "name_i18n": {
          "zh-CN": "日本",
          "en-US": "Japan"
        },
        "status": 1,
        "sort_no": 10
      }
    ]
  }
}
```

### GET /admin/region/list-all

不分页获取全部国家地区，可按状态筛选。

Query 示例：

```text
/admin/region/list-all?status=1
```

请求：

| 字段       | 位置    | 必填 | 类型    | 说明            |
| -------- | ----- | -- | ----- | ------------- |
| `status` | query | 否  | int64 | `1` 启用，`2` 停用 |

响应：

```json
{
  "code": 0,
  "msg": "ok",
  "data": {
    "list": [
      {
        "id": 1,
        "code": "JP",
        "calling_code": "81",
        "name_i18n": {
          "zh-CN": "日本",
          "en-US": "Japan"
        },
        "status": 1,
        "sort_no": 10
      }
    ]
  }
}
```

### POST /admin/region/reorder

调整国家地区排序。

请求：

| 字段          | 位置   | 必填 | 类型    | 说明                  |
| ----------- | ---- | -- | ----- | ------------------- |
| `id`        | json | 是  | int64 | 需要移动的国家地区 ID，大于 `0` |
| `target_id` | json | 是  | int64 | 目标国家地区 ID，大于 `0`    |

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
