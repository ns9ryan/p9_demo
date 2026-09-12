# integration 前端对接

Base URL（本地）：`http://127.0.0.1:18005`  
路径前缀：`/p9/integration`（IP/GEO 模块为 `/p9/integration/ipgeo`）  
JSON 字段：`snake_case`  
全部 HTTP 要登录：`Authorization: Bearer <access_token>`（core `CheckToken`），并走 Casbin（`Enforce`）。启动时 `RegisterCatalog` 把本服务 path 写入 `sys_api` 并给各厅 `super_admin` 补权；未授权角色会 `403`。未登录 `401`。

查询顺序：合法公网按当前租户策略选出启用源（未配则用服务商默认），按 `priority` 从小到大；每个源先查缓存（KV，键 `name|IP`，整包 JSON；默认表 `ip_geo_cache`，命中续期），miss 再打该源。改策略不必清缓存：停用的源不会再被读到。单 IP 失败不拖垮整批。服务商连接配置无管理接口；租户策略走 `/policy/*`。机器可读契约：[swagger.json](swagger.json)（由 `api/desc` 生成，`make swagger` / `make -f Makefile.win swagger`；已包真实外壳 `{code,msg,data}`）。

其它服务请走 gRPC，见 [RPC.md](RPC.md)。配置与本地启动见 [OPS.md](OPS.md)。服务说明见 [integration.md](integration.md)。

---

## 统一响应

成功：HTTP **200**

```json
{ "code": 0, "msg": "ok", "data": {} }
```

失败：HTTP status **等于** `code`，**没有** `data`

```json
{ "code": 400, "msg": "max ips is 100" }
```

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `code` | int | 成功恒为 `0` |
| `msg` | string | 成功恒为 `"ok"`；失败为原因 |
| `data` | object / null | 仅成功有；业务字段都在这里。无载荷可为 `null` |

常见失败 `code`：`400` 参数错误，`401` 未登录，`404` 不存在。

下面各接口的「返回」只列 **`data` 里的字段**（外层仍是 `code` / `msg` / `data`）。

---

## 一览

| 方法 | 路径 | 给谁 | 说明 |
| --- | --- | --- | --- |
| POST | `/p9/integration/ipgeo/ip/search` | 业务 / 后台 | 批量查 IP 地理信息 |
| POST | `/p9/integration/ipgeo/policy/list` | 后台 | 当前租户策略列表 |
| POST | `/p9/integration/ipgeo/policy/save` | 后台 | 保存当前租户策略 |
| POST | `/p9/integration/ipgeo/policy/default/save` | 后台 | 保存默认策略 |
| POST | `/p9/integration/ipgeo/local/list` | 后台 | 本地 IP 列表 |
| POST | `/p9/integration/ipgeo/local/create` | 后台 | 新增本地 IP |
| POST | `/p9/integration/ipgeo/local/update` | 后台 | 更新本地 IP |
| POST | `/p9/integration/ipgeo/local/delete` | 后台 | 删除本地 IP |
| GET | `/p9/integration/ipgeo/local/detail` | 后台 | 本地 IP 详情 |

`Content-Type: application/json`（有 body 的 POST）。`Authorization: Bearer <access_token>`。

---

## POST `/p9/integration/ipgeo/ip/search`

批量查询。最多 **100** 条；超过整批失败。无效 IP 仍占一位，只回 `query`。私网 / 回环不打三方：`country` 为 `私有IP`，`source` 为 `private`。某条源失败继续下一条，整批仍 200。

名称类字段按 `lang` 从多语言里挑（缺则 `en`，再任意一个）。

合法公网按当前租户策略选出启用源（未配则用 yaml 种子默认），按 `priority` 依次查；每个源先读缓存（KV，键 `name|IP`；默认 `ip_geo_cache.payload`），miss 再 Fetch。停用某源则不再读它的缓存。不收 `operator_code`。租户取自 JWT（其它服务走 gRPC metadata `x-operator-code`）。流程图见 [ip_search_pfd.md](ip_search_pfd.md)。

### 请求

| 字段 | 必填 | 类型 | 说明 |
| --- | --- | --- | --- |
| `ips` | 是 | string[] | IPv4 / IPv6；最多 100 |
| `lang` | 否 | string | 空则 `en`。仅：`en` `zh-CN` `pt-BR` `ru` `fr` `ja` `es` `de` |

```json
{ "ips": ["8.8.8.8", "1.1.1.1"], "lang": "zh-CN" }
```

### 返回 `data`

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `list` | array | 与 `ips` 等长；不会是 `null` |

`list[]`：

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `query` | string | 请求里的原始 IP |
| `continent` | string | 大洲名（随 `lang`） |
| `continent_code` | string | 如 `NA` |
| `country` | string | 国家名（随 `lang`） |
| `country_code` | string | 如 `US` |
| `province` | string | 省 / 州 |
| `province_iso_code` | string | 如 `CA` |
| `city` | string | 城市 |
| `zip` | string | 邮编（与 `postal` 同源） |
| `postal` | string | 邮编 |
| `lat` | float | 纬度 |
| `lon` | float | 经度 |
| `timezone` | string | 时区 |
| `isp` | string | ISP |
| `org` | string | 组织 |
| `as` | string | ASN 文本 |
| `source` | string | 命中的服务商 `name`，或 `private` |

```json
{
  "list": [
    {
      "query": "8.8.8.8",
      "continent": "北美洲",
      "continent_code": "NA",
      "country": "美国",
      "country_code": "US",
      "province": "",
      "province_iso_code": "",
      "city": "",
      "zip": "",
      "lat": 37.751,
      "lon": -97.822,
      "timezone": "",
      "isp": "",
      "org": "",
      "as": "",
      "postal": "",
      "source": "maxmind"
    }
  ]
}
```

失败：`max ips is 100`、`unsupported lang`。

---

## 策略

服务商连接配置来自 yaml 种子，后台不管理。每个租户可覆盖某个服务商的 `is_enabled`、`priority`。登录账号的 `operator_code` 从 JWT 取，不收请求体。该租户没有策略行时，用 `ip_geo_provider` 上的默认启用/优先级。

`tag` 存在服务商表：`在线`（`driver=http`）、`本地`（`driver=local`）、`离线`（`mmdb`、`ip2region`）。列表只回策略字段 + `name` + `tag`，不回连接密钥。

`page` 不适用：一次返回全部服务商，按有效 `priority`、`id` 升序。

### 列表字段

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `provider_id` | int | 全局服务商 id，保存时回传 |
| `name` | string | 服务商名称 |
| `is_enabled` | bool | 当前租户是否启用；未配策略时为默认 |
| `priority` | int | 当前租户优先级，越小越先 |
| `tag` | string | `在线` / `本地` / `离线` |

## POST `/p9/integration/ipgeo/policy/list`

无请求体。列出全部服务商的当前租户策略（没有则用默认）。

### 返回 `data`

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `list` | array | 不会是 `null` |
| `total` | int | 条数 |

```json
{
  "list": [
    { "provider_id": 1, "name": "local", "is_enabled": true, "priority": 0, "tag": "本地" },
    { "provider_id": 2, "name": "maxmind", "is_enabled": true, "priority": 1, "tag": "在线" },
    { "provider_id": 3, "name": "geoip-city", "is_enabled": true, "priority": 2, "tag": "离线" }
  ],
  "total": 3
}
```

## POST `/p9/integration/ipgeo/policy/save`

批量保存**当前租户**对服务商的 `is_enabled` / `priority`。不存在则创建，存在则更新。不改连接配置，也不改其它租户。同一 `provider_id` 不能出现多次，否则整批失败。

```json
{
  "list": [
    { "provider_id": 1, "is_enabled": true, "priority": 0 },
    { "provider_id": 2, "is_enabled": false, "priority": 10 }
  ]
}
```

| 字段 | 必填 | 类型 | 说明 |
| --- | --- | --- | --- |
| `list` | 是 | array | 不能为空 |
| `list[].provider_id` | 是 | int | 全局服务商 id |
| `list[].is_enabled` | 是 | bool | 该租户是否启用 |
| `list[].priority` | 是 | int | 该租户优先级，越小越先 |

### 返回 `data`

与 list 相同：全部服务商（已套当前租户策略）。失败：`policy list is required`、`max policies is 100`、`duplicate provider_id`、`invalid provider`、`provider not found`。

## POST `/p9/integration/ipgeo/policy/default/save`

批量保存**默认**启用/优先级，写入 `ip_geo_provider`。租户没配策略时用这些值。已有租户策略的行不受影响。请求体与 `/policy/save` 的 `list` 相同。

### 返回 `data`

全部服务商的**默认** `is_enabled` / `priority`（未套租户策略），字段同 list。

---

## 本地 IP 表

`driver=local` 的数据源。全局一份，不按租户隔离。是否启用、优先级看当前租户策略（没有则用 yaml 种子默认）。表中没有该 IP 时，查询继续走后面的源。增删改会清掉该 IP 的查询缓存。

`page` 从 1；`page_size` 默认 20、最多 100；≤0 用默认，不报错。

### 字段

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `id` | int | 更新/删除必填 |
| `ip` | string | 必填，IPv4/IPv6，全局唯一 |
| `continent` / `continent_code` | string | 大洲 |
| `country` / `country_code` | string | 国家 |
| `province` / `province_iso_code` | string | 省/州 |
| `city` | string | 城市 |
| `zip` / `postal` | string | 邮编，写入同一列；只填其中一个即可 |
| `lat` / `lon` | float | 经纬度 |
| `timezone` `isp` `org` `as` | string | 其它 |

## POST `/p9/integration/ipgeo/local/list`

```json
{ "page": 1, "page_size": 20, "ip": "8.8" }
```

`ip` 可选，子串过滤。返回 `data.list` + `data.total`。

## POST `/p9/integration/ipgeo/local/create`

```json
{ "ip": "8.8.8.8", "country": "中国", "country_code": "CN", "city": "深圳" }
```

失败：`ip is required`、`invalid ip`、`ip already exists`。

## POST `/p9/integration/ipgeo/local/update`

必须带 `id`，其余字段同 create。

## POST `/p9/integration/ipgeo/local/delete`

```json
{ "id": 1 }
```

### 返回 `data`

```json
{ "result": "success" }
```

## GET `/p9/integration/ipgeo/local/detail`

`GET /p9/integration/ipgeo/local/detail?id=1`

失败：`local ip not found`。

---

## gRPC（其它服务）

业务进程查 IP 见 [RPC.md](RPC.md)（`Integration`）。管理后台仍用本文 HTTP（查询 + 策略 + 本地 IP）。
