# oss 前端对接

Base URL（本地）：`http://127.0.0.1:8020`  
路径前缀：`/p9/oss`  
JSON 字段：`snake_case`

业务表只存 `**key**`。`url` 只用于当前展示；私有桶会过期，换链用 `GET /url/detail?key=`。  
只传 `**scene**`（场景枚举），不要传桶名。不传 `scene` 则走服务端 `default_scene`（当前为 `public`）。未在 `upload.scenes` 里列出的值会失败。

对象 key 规划：

- 有厅：`{operator_code}/{scene}/{yyyyMMdd}/{32位hex}{.ext}`  
  例：`demo/public/20260909/76d54945f0736f325ab51af53c0dfcb3.jpg`
- 无厅：`{scene}/{yyyyMMdd}/{32位hex}{.ext}`  
  例：`public/20260909/76d54945f0736f325ab51af53c0dfcb3.jpg`

---

## 统一响应

成功：HTTP **200**

```json
{ "code": 0, "msg": "ok", "data": {} }
```

失败：HTTP status **等于** `code`，没有 `data`

```json
{ "code": 400, "msg": "no file provided" }
```


| 字段     | 类型            | 说明                |
| ------ | ------------- | ----------------- |
| `code` | int           | 成功恒为 `0`          |
| `msg`  | string        | 成功恒为 `"ok"`；失败为原因 |
| `data` | object / null | 仅成功有；业务字段都在这里     |


常见失败 `code`：`400` 参数错误，`404` 不存在，`503` 未配库（档案/服务商接口）。

下面各接口的「返回」只列 `**data` 里的字段**（外层仍是 `code` / `msg` / `data`）。

---

## 一览


| 方法   | 路径                             | 给谁  | 说明                    |
| ---- | ------------------------------ | --- | --------------------- |
| POST | `/p9/oss/file/upload`     | 前端  | 小文件中转上传               |
| GET  | `/p9/oss/url/detail`      | 前端  | 用 key 换可打开的 url       |
| POST | `/p9/oss/upload/presign`  | 前端  | 大文件第 1 步：签发 PUT       |
| POST | `/p9/oss/upload/complete` | 前端  | 大文件第 3 步：确认完成         |
| POST | `/p9/oss/upload/put`      | 仅联调 | 生产不要用，应对 `put_url` 直传 |
| POST | `/p9/oss/file/list`       | 后台  | 档案列表                  |
| POST | `/p9/oss/file/delete`     | 后台  | 删对象 + 档案              |
| POST | `/p9/oss/provider/list`   | 后台  | 服务商列表                 |
| POST | `/p9/oss/provider/create` | 后台  | 新增服务商                 |
| GET  | `/p9/oss/provider/active` | 后台  | 当前生效配置                |
| GET  | `/p9/oss/provider/detail` | 后台  | 服务商详情                 |
| POST | `/p9/oss/provider/update` | 后台  | 更新服务商                 |
| POST | `/p9/oss/provider/delete` | 后台  | 删除服务商                 |
| POST | `/p9/oss/provider/enable` | 后台  | 设为默认                  |


---

## POST `/p9/oss/file/upload`

小文件经本服务中转。`Content-Type: multipart/form-data`。不要手写 Content-Type（须带 boundary）。

### 请求


| 参数      | 位置   | 必填  | 类型     | 说明                            |
| ------- | ---- | --- | ------ | ----------------------------- |
| `file`  | form | 是   | file   | 文件本体。字段名必须是 `file`            |
| `scene` | form | 否   | string | 场景枚举。空则 `default_scene`。禁止传桶名 |


### 返回 `data`


| 字段               | 类型     | 说明                              |
| ---------------- | ------ | ------------------------------- |
| `url`            | string | 当前可打开地址。公开桶长期可用；私有桶是预签名 GET，会过期 |
| `key`            | string | 对象 key，**业务表存这个**               |
| `scene`          | string | 规范化后的 scene，不必单独落库              |
| `visibility`     | string | `public` 或 `private`，由所在桶决定     |
| `url_expires_in` | int    | 秒。仅 `private`；公开不出现或为 0         |


```json
{
  "code": 0,
  "msg": "ok",
  "data": {
    "url": "http://127.0.0.1:9200/oss-bucket/public/20260831/76d54945f0736f325ab51af53c0dfcb3.jpg",
    "key": "public/20260831/76d54945f0736f325ab51af53c0dfcb3.jpg",
    "scene": "public",
    "visibility": "public"
  }
}
```

失败示例：`{"code":400,"msg":"no file provided"}`、`file too large`、`file extension not allowed`、`invalid scene`。

---

## GET `/p9/oss/url/detail`

用已存的 `key` 换当前可打开的 `url`（私有链过期后走这条）。需登录：`Authorization: Bearer <access_token>`。

已配档案库时，只给本租户档案换链；别人的 key 返回 404。

### 请求


| 参数    | 位置    | 必填  | 类型     | 说明                          |
| ----- | ----- | --- | ------ | --------------------------- |
| `Authorization` | header | 配了 `core_rpc` 时必填 | string | `Bearer <access_token>`；并走 Casbin（启动 `RegisterCatalog` 写入 `sys_api`，超管自动有权） |
| `key` | query | 推荐  | string | 对象 key                      |
| `url` | query | 否   | string | 旧访问地址；与 `key` 同时传以 `key` 为准 |


### 返回 `data`

与上传相同：`url` / `key` / `scene` / `visibility` / `url_expires_in`。

```json
{
  "code": 0,
  "msg": "ok",
  "data": {
    "url": "http://127.0.0.1:9200/oss-bucket/public/20260831/76d54945f0736f325ab51af53c0dfcb3.jpg",
    "key": "public/20260831/76d54945f0736f325ab51af53c0dfcb3.jpg",
    "scene": "public",
    "visibility": "public"
  }
}
```

失败：`url or key is required`、`invalid key`、`invalid scene`、未登录 `401`、本租户无此档案 `404`。

管理后台业务列表不要对每行打本接口。公开图用 `GET /provider/active` 缓存 `file_url`、`path_style`、桶名，本地按 `path_style` 拼：

```text
true  → {file_url}/{bucket}/{key}     # MinIO
false → {file_url}/{key}              # CDN / 虚拟主机
```

私有桶必须预签名，不能拼；需要时再走本接口或批量换链。

---

## 大文件直传（前端生产流程）

不要把文件字节再打到 oss。

1. `POST /upload/presign` 拿到 `put_url`、`headers`、`key`
2. **PUT `put_url`**（打到 MinIO/OSS，不是本服务）。`headers` 原样带上，body 是文件本身
3. `POST /upload/complete` 提交 `key`，拿最终 `url` / `key` 并入库

第 2 步存储返回 HTTP 200/204，不是本服务的 `{code,msg}`。

---

## POST `/p9/oss/upload/presign`

`Content-Type: application/json`

### 请求


| 字段         | 必填  | 类型     | 说明                       |
| ---------- | --- | ------ | ------------------------ |
| `filename` | 是   | string | 带真实扩展名，如 `clip.mp4`      |
| `size`     | 是   | int    | 即将 PUT 的字节数（`file.size`） |
| `scene`    | 否   | string | 场景枚举。空则 `default_scene`  |


```json
{ "filename": "photo.jpg", "size": 204800, "scene": "public" }
```

### 返回 `data`


| 字段               | 类型     | 说明                                         |
| ---------------- | ------ | ------------------------------------------ |
| `put_url`        | string | 第 2 步 PUT 地址。不要落库、不要用浏览器打开                 |
| `method`         | string | 固定 `PUT`                                   |
| `headers`        | object | 必须原样带上，一般含 `Content-Type`、`Content-Length` |
| `key`            | string | 第 3 步提交；业务表存 complete 返回的 key              |
| `url`            | string | 当前可打开地址；此时对象可能还未上传，以 complete 为准           |
| `expires_in`     | int    | `put_url` 有效秒数                             |
| `scene`          | string | 规范化后的 scene                                |
| `visibility`     | string | `public` 或 `private`                       |
| `url_expires_in` | int    | 读链有效秒数；仅 private                           |


```json
{
  "code": 0,
  "msg": "ok",
  "data": {
    "put_url": "http://127.0.0.1:9200/oss-bucket/public/20260825/ab12cd.jpg?X-Amz-...",
    "method": "PUT",
    "headers": { "Content-Type": "image/jpeg", "Content-Length": "204800" },
    "key": "public/20260825/ab12cd.jpg",
    "url": "http://127.0.0.1:9200/oss-bucket/public/20260825/ab12cd.jpg",
    "expires_in": 600,
    "scene": "public",
    "visibility": "public"
  }
}
```

失败：`filename is required`、`size is required`、`file too large`、`file extension not allowed`、`invalid scene`。

---

## POST `/p9/oss/upload/put`

仅 Swagger / 联调：本服务代为写入。**生产前端不要调用。**

`Content-Type: multipart/form-data`

### 请求


| 参数      | 位置   | 必填  | 类型     | 说明                    |
| ------- | ---- | --- | ------ | --------------------- |
| `key`   | form | 是   | string | presign 返回的 key       |
| `file`  | form | 是   | file   | 文件本体                  |
| `scene` | form | 否   | string | 须与第 1 步一致；不传则从 key 解析 |


### 返回 `data`

与上传相同。此步**不入库**，仍须 complete。

---

## POST `/p9/oss/upload/complete`

`Content-Type: application/json`

### 请求


| 字段      | 必填  | 类型     | 说明                         |
| ------- | --- | ------ | -------------------------- |
| `key`   | 是   | string | presign 返回的 key            |
| `scene` | 否   | string | 不传则从 key 第一段解析；若传须与第 1 步一致 |


```json
{ "key": "public/20260825/ab12cd.jpg" }
```

### 返回 `data`

与上传相同：`url` / `key` / `scene` / `visibility` / `url_expires_in`。

失败：`key is required`、`object not found: PUT put_url first, then complete`、`invalid scene`。

---

## POST `/p9/oss/file/list`

档案列表。需 PostgreSQL。`Content-Type: application/json`

### 请求


| 字段          | 必填  | 类型     | 说明                      |
| ----------- | --- | ------ | ----------------------- |
| `scene`     | 否   | string | 按 scene 过滤              |
| `page`      | 否   | int    | 从 1 起，默认 1；`<=0` 当 1    |
| `page_size` | 否   | int    | 默认 20，最大 100；`<=0` 当 20 |


```json
{ "scene": "public", "page": 1, "page_size": 20 }
```

### 返回 `data`


| 字段      | 类型    | 说明                     |
| ------- | ----- | ---------------------- |
| `list`  | array | 当前页。空为 `[]`，不会是 `null` |
| `total` | int   | 符合条件的总条数，不是本页条数        |


`list[]`：


| 字段               | 类型     | 说明                   |
| ---------------- | ------ | -------------------- |
| `id`             | int    | 档案 id                |
| `key`            | string | 对象 key               |
| `scene`          | string | 上传场景                 |
| `bucket`         | string | 所在桶                  |
| `name`           | string | 原始文件名                |
| `size`           | int    | 字节                   |
| `content_type`   | string | MIME                 |
| `url`            | string | 当前可打开地址              |
| `visibility`     | string | `public` 或 `private` |
| `url_expires_in` | int    | 秒；仅 private          |
| `created_at`     | int    | unix 毫秒              |


```json
{
  "code": 0,
  "msg": "ok",
  "data": {
    "list": [
      {
        "id": 1,
        "key": "public/20260831/76d54945f0736f325ab51af53c0dfcb3.jpg",
        "scene": "public",
        "bucket": "oss-bucket",
        "name": "a.jpg",
        "size": 20480,
        "content_type": "image/jpeg",
        "url": "http://127.0.0.1:9200/oss-bucket/public/20260831/76d54945f0736f325ab51af53c0dfcb3.jpg",
        "visibility": "public",
        "created_at": 1756636800000
      }
    ],
    "total": 1
  }
}
```

失败：`{"code":503,"msg":"database not configured"}`。

---

## POST `/p9/oss/file/delete`

删存储对象并删档案行。需 PostgreSQL。`Content-Type: application/json`

### 请求


| 字段    | 必填  | 类型     | 说明     |
| ----- | --- | ------ | ------ |
| `key` | 是   | string | 对象 key |


```json
{ "key": "public/20260831/76d54945f0736f325ab51af53c0dfcb3.jpg" }
```

### 返回 `data`


| 字段           | 类型     | 说明                   |
| ------------ | ------ | -------------------- |
| `url`        | string | 删除后为空                |
| `key`        | string | 被删的 key              |
| `scene`      | string | 场景                   |
| `visibility` | string | `public` 或 `private` |


失败：`key is required`、`503` 未配库。

---

## 服务商（后台配置，前端业务页一般不用）

`secret_key` 列表/详情不返回。创建不要传 `id`、`source`。未配库返回 `503`。

### Provider `data` 字段


| 字段              | 类型     | 说明                                                   |
| --------------- | ------ | ---------------------------------------------------- |
| `id`            | int    | 创建后返回                                                |
| `name`          | string | 名称                                                   |
| `is_default`    | bool   | `true` 时立即覆盖运行时                                      |
| `driver`        | string | `s3` / `minio`                                       |
| `endpoint`      | string | S3 API 地址                                            |
| `file_url`      | string | 公开读 URL 前缀                                           |
| `access_key`    | string | 访问密钥                                                 |
| `secret_key`    | string | 仅创建/更新可写；列表不返回；更新时空则保留                               |
| `path_style`    | bool   | 公开链怎么拼。`true`：`{file_url}/{bucket}/{key}`（MinIO）；`false`：`{file_url}/{key}` |
| `region`        | string | 如 `default`                                          |
| `buckets`       | object | 桶名 → `{ "mode": "public|private", "sign_ttl": 600 }` |
| `source`        | string | 仅响应：`yaml` / `etcd` / `provider:名称`                  |


---

## POST `/p9/oss/provider/list`

无请求参数。

### 返回 `data`


| 字段      | 类型    | 说明                         |
| ------- | ----- | -------------------------- |
| `list`  | array | Provider 数组，无 `secret_key` |
| `total` | int   | 条数                         |


---

## POST `/p9/oss/provider/create`

请求 body：上表（不含 `id` / `source`）。`is_default=true` 立即生效。

### 返回 `data`

单个 Provider。

---

## GET `/p9/oss/provider/active`

无请求参数。没有默认服务商时返回进程 `oss`（`source` 为 `yaml` 或 `etcd`），与实际上传使用的驱动一致。

### 返回 `data`

单个 Provider，含 `source`。

---

## GET `/p9/oss/provider/detail`

### 请求


| 参数   | 位置    | 必填  | 类型  | 说明  |
| ---- | ----- | --- | --- | --- |
| `id` | query | 是   | int | 正整数 |


### 返回 `data`

单个 Provider。失败：`invalid id`、`provider not found`。

---

## POST `/p9/oss/provider/update`

### 请求

body 同创建，**必须带 `id`**。`secret_key` 空则不改。

### 返回 `data`

单个 Provider。

---

## POST `/p9/oss/provider/delete`

### 请求


| 字段    | 必填  | 类型    | 说明           |
| ----- | --- | ----- | ------------ |
| `id`  | 二选一 | int   | 单个 id        |
| `ids` | 二选一 | int[] | 目前按第一个 id 删除 |


```json
{ "id": 1 }
```

### 返回

成功仅 `{"code":0,"msg":"ok"}`，无 `data`。若删的是默认项，运行时退回 yaml/etcd。

---

## POST `/p9/oss/provider/enable`

### 请求


| 字段   | 必填  | 类型  | 说明      |
| ---- | --- | --- | ------- |
| `id` | 是   | int | 要启用的服务商 |


```json
{ "id": 1 }
```

### 返回 `data`

单个 Provider，`is_default` 为 `true`。

---

## gRPC（其它服务）

其它服务新建租户见 [RPC.md](RPC.md)（`InitTenant`，body 传 `operator_code`）。改服务商走本文 HTTP `/provider/*`。前端仍用本文 HTTP。