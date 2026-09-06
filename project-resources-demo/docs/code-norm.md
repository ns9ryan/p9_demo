# 开发技术规范

---

## 1. 量化限制


| 项                            | 上限                                       | 说明                                                        |
| ---------------------------- | ---------------------------------------- | --------------------------------------------------------- |
| 函数体行数（不含签名与空行后的注释块）          | **80**                                   | Handler **≤ 40**；`main` / `MustInit` 编排可到 100             |
| 函数参数个数                       | **5**                                    | 含 `ctx`、接收者不算。超过则改为 `XxxReq` / `Options` 结构体              |
| 命名返回值个数                      | **3**                                    | 常见 `(T, error)` / `(T, int64, error)`。禁止 `(a, b, c, err)` |
| `if` / `for` / `switch` 嵌套深度 | **4**                                    | 超了用提前返回、抽函数                                               |
| 单文件行数                        | **800**                                  | 生成代码、纯路由表可到 1000；再大按资源拆文件                                 |
| 单行长度                         | **120**                                  | import、URL、DDL 字符串除外                                      |
| 圈复杂度（单个函数）                   | **15**                                   | `switch` 多分支按路径计；优先表驱动                                    |
| 函数内局部变量                      | **12**                                   | 过多说明该拆                                                    |
| 结构体导出字段（单个 DTO）              | **20**                                   | 再多分子结构（分页、筛选、主体）                                          |
| slice/map 字面量元素（业务常量表除外）     | 建议 **≤ 16** 后抽 `var xxx = []T{...}` 到文件级 |                                                           |
| 包的对外导出函数（单 package）          | 建议 **≤ 30**                              | 根包 `rbacx` 只 re-export，实现放子包                              |


下面每条对应上表一行：**怎么数**、**符合**、**违反及改法**。

### 1.1 函数体行数 ≤ 80（Handler ≤ 40）

**怎么数：** 从 `{` 下一行到配对 `}` 上一行，计可执行语句与控制结构（含闭包）。签名、函数上方注释、函数体内的空行不计入。禁止多句挤一行规避。

**符合（Handler 薄，远小于 40 行）：**

```go
func (h *H) userCreate(w http.ResponseWriter, r *http.Request) {
	var req service.CreateUserReq
	if err := decodeJSON(r, &req); err != nil {
		fail(w, service.BadRequest("invalid json"))
		return
	}
	u, err := h.D.CreateUser(r.Context(), claims(r), req)
	if err != nil {
		fail(w, err)
		return
	}
	ok(w, service.PublicUser(u, nil))
}
```

**违反：** 一个 `CreateUser` 里校验 + 查重 + 哈希 + 写用户 + 绑角色 + 写 Casbin + 拼响应，体长 120 行。

**改法：** 拆成 `validateCreateUser` → `insertUser` → `bindRoles`，编排函数本身 ≤ 80 行。`MustInit` / `main` 只做装配，上限 100 行。

---



### 1.2 函数参数个数 ≤ 5

**怎么数：** 签名里逗号分隔的形参。**接收者** `(d *Deps)` **不算。** `ctx` 算 1 个。可变参数 `ids ...int64` 算 1 个。

**符合（3 个：ctx、claims、req）：**

```go
func (d *Deps) CreateUser(ctx context.Context, claims *scope.Claims, req CreateUserReq) (*model.User, error)
```

**违反（远超 5 个）：**

```go
func CreateUser(ctx context.Context, db *gorm.DB, e *casbin.Enforcer,
	username, password, display, mobile, email string, status int16, operatorID *int64, roleIDs []int64)
```

**改法：** 业务字段进 `CreateUserReq`；`db`/`enforcer` 放 `*Deps` 接收者。函数选项最多 3 个 `func(*Opt)`，或改成一个 `Opt` 结构体。

---



### 1.3 返回值个数 ≤ 3

**怎么数：** 括号里的结果个数。`(T, error)` 为 2；`(list, total, error)` 为 3。

**符合：**

```go
func (d *Deps) ListUsers(ctx context.Context, claims *scope.Claims, page PageReq) ([]model.User, int64, error)
func (d *Deps) Login(ctx context.Context, req LoginReq) (*LoginResult, error)
```

（`Login` 已是 3 个返回值，不能再加 `roles`。）

**违反：**

```go
func Login(...) (token string, expire int64, user UserPublic, roles []string, err error) // 5 个
```

**改法：** `token+expire` 合成 `TokenInfo`；`user` 内带 `role_codes`。需要第 4 个结果时用结构体，例如 `type LoginResult struct { Token *TokenInfo; User UserPublic }`。

---



### 1.4 嵌套深度 ≤ 4

**怎么数：** `if` / `for` / `switch` / `select` 层层包裹。函数体为 0；每深入一层 +1。`else if` 与外层 `if` 同级，不算更深一层。

**符合（深度 2：for → if）：**

```go
for _, id := range ids {
	u, err := d.mustTenantUser(ctx, claims, id)
	if err != nil {
		return err
	}
	if u.IsSuperAdmin {
		return Forbidden("cannot delete root user")
	}
	_ = d.DB.WithContext(ctx).Model(u).Update("deleted_at", nowPtr())
}
```

**违反（深度 5：if mode → if claims → for → if → if）：**

```go
if d.Mode == ModeOn {
	if claims != nil {
		for _, id := range req.IDs {
			if id > 0 {
				if u, err := d.get(id); err == nil {
					if u.OperatorID != nil && *u.OperatorID == claims.OperatorID {
						// ...
					}
				}
			}
		}
	}
}
```

**改法：** 提前返回（`if claims == nil { return Unauthorized(...) }`），循环体抽 `deleteOne(ctx, claims, id)`。

---



### 1.5 单文件行数 ≤ 800

**怎么数：** `wc -l` 含空行与注释。goctl 生成、纯路由表可到 **1000**。

**符合：** `service/user.go` 按用户能力单文件，约 280 行。

**违反：** `handler/register.go` 把全部 Handler + `ok`/`fail` 堆在一起超过 500 行。

**改法：** 按资源拆 `handler/user.go`、`handler/role.go`、`handler/http.go`（`ok`/`fail`），`register.go` 只留 `Register`。

---



### 1.6 单行长度 ≤ 120

**怎么数：** 字符数（非字节）。`import` 路径、URL、DDL 整段字符串除外。

**符合：**

```go
q := OperatorScope(d.DB.WithContext(ctx).Model(&model.User{}), d.Mode, operatorID(claims))
```

**违反：**

```go
err := d.DB.WithContext(ctx).Where("deleted_at IS NULL AND LOWER(username) = LOWER(?) AND operator_id = ? AND status = ?", req.Username, op.ID, model.StatusNormal).First(&u).Error
```

**改法：** 链式换行，每行一个 `Where` / 参数，缩进对齐。

---



### 1.7 圈复杂度 ≤ 15

**怎么数：** 函数入口为 1，每个 `if` / `else if` / `for` / `case` / `&&` / `||` 的独立决策 +1。15 条互斥 `case` 的纯分发可以，但不要在每个 `case` 里再套 `if`。

**符合：** `fail` 只做类型断言 + 写 JSON，复杂度约 1～2。登录校验拆成「取厅 → 取用户 → 验密 → 签 Token」，每段 < 8。

**违反：** 一个 `Login` 里 10+ 个 `if err`、模式分支、厅停用、用户停用、无角色、IP 解析全揉在一起，路径数 > 15。

**改法：** 表驱动（`switch d.Mode` 只负责选查询条件）；重复的 `if err != nil` 留在小函数里，编排函数本身分支变少。

---



### 1.8 函数内局部变量 ≤ 12

**怎么数：** 函数体内 `:=` / `var` 声明的标识符（不含短语句 `if err :=` 可与外层 `err` 复用）。接收者、参数不算。

**符合：**

```go
func (d *Deps) GetUser(ctx context.Context, claims *scope.Claims, id int64) (*model.User, []string, error) {
	u, err := d.mustTenantUser(ctx, claims, id)
	if err != nil {
		return nil, nil, err
	}
	codes, err := d.RoleCodesOfUser(ctx, u.ID)
	return u, codes, err
}
```

局部量：`u`、`err`、`codes`（`err` 复用），约 3 个。

**违反：** 一个函数里 `hash` `salt` `oid` `exist` `role` `menus` `policies` `dom` `now` `ip` `token` `exp` `pub` … 超过 12 个。

**改法：** 建用户抽 `createUserRow(...)`，发 Token 抽 `SignToken`，让编排函数只留 5～8 个局部量。

---



### 1.9 单个 DTO 导出字段 ≤ 20

**怎么数：** 结构体里 **导出**（大写）且带 `json` 的字段。嵌套结构体按各自文件计数，不把子结构字段加到父上。

**符合：** `UserPublic` 约 9 个字段；`CreateUserReq` 约 8 个。

**违反：**

```go
type UserUpdateReq struct {
	ID, OperatorID, Status int64
	Username, Password, DisplayName, Mobile, Email, Salt, UserCode string
	IsSuperAdmin, IsSystem bool
	RoleIDs []int64
	Timezone, Currency, Remark, Avatar, Dept, LastIP string
	Page, PageSize int // 更新和分页混在一个结构体
}
```

**改法：** 更新用 `UpdateUserReq`（只含允许改的指针字段）；分页用 `PageReq`；禁止把 `salt` / `is_super_admin` 放进 DTO。

---



### 1.10 slice / map 字面量元素建议 ≤ 16

**怎么数：** 函数体内 `[]T{ a, b, ... }` / `map[K]V{ ... }` 的元素个数。文件级 `var builtin = []Item{...}`（如 `internalapis.Builtin`）不受限，但单项仍应可读。

**符合：** Handler 里组 2～3 条测试路由；`APIAuthReq.Data` 运行时再 append。

**违反：** 在 `Register` 函数中间写 25 条 `rest.Route{...}` 字面量，导致函数行数、复杂度一起爆。

**改法：** 提到包级 `var jwtOnlyRoutes = []rest.Route{...}`，或按资源拆 `userRoutes()` 返回切片（每个小函数内部元素仍建议 ≤ 16）。

---



### 1.11 单 package 导出函数建议 ≤ 30

**怎么数：** 该包 `func Xxx` / `func (T) Xxx` 中 **大写** 的个数（不含测试文件）。

---



## 2. 接口风格与路由



### 2.1 风格

本项目后台 **不是 REST**，采用 **动作路径**：

```text
POST /admin/{资源}/{动作}     # 写、列表、授权
GET  /admin/{资源}/{动作}     # 当前用户、详情、树
```

例：`POST /admin/user/create`、`POST /admin/user/list`、`GET /admin/user/detail`。

### 2.2 HTTP 方法约定


| 方法     | 用途                                |
| ------ | --------------------------------- |
| `POST` | 创建、更新、删除、列表（可带筛选 body）、授权保存、登录、登出 |
| `GET`  | 无 body 的只读：当前用户、详情、菜单树、权限码、厅信息    |


- 列表统一 `POST` + JSON 分页，避免超长 query。
- 详情 `GET`，主键用 query：`?id=1`（正整数）。
- 删除 body：`{"id":1}` 或 `{"ids":[1,2]}`。



### 2.3 鉴权分组（注册时必须拆开）


| 组            | 中间件                         | 例子                                           |
| ------------ | --------------------------- | -------------------------------------------- |
| 公开           | 无                           | `login`、`refresh`、`bootstrap/*`              |
| 仅 JWT        | `MiddlewareJWT`             | `user/info`、`menu/role`、`user/perm`、`logout` |
| JWT + Casbin | JWT + `MiddlewareAuthority` | 用户/角色 CRUD、授权、业务 `/admin/promo/*`            |




### 2.4 Header


| Header                           | 用途                           |
| -------------------------------- | ---------------------------- |
| `Authorization: Bearer <token>`  | 后台 JWT                       |
| `Content-Type: application/json` | 有 body 的请求                   |
| `X-Forwarded-For`                | 可选；登录写 `last_login_ip` 时取第一段 |
| `X-Lang`                         | 语言 ｜                         |


---



## 3. 接口响应格式

与当前 `handler/register.go` 的 `ok` / `fail` 一致。业务项目后台应对齐，避免前端两套解析。

### 3.1 成功

- HTTP **200**
- Body：

```json
{
  "code": 0,
  "msg": "ok",
  "data": {}
}
```


| 字段     | 类型                    | 规则                                         |
| ------ | --------------------- | ------------------------------------------ |
| `code` | number                | 成功 **恒为** `0`（不要用 HTTP 200 里的 200）         |
| `msg`  | string                | 成功固定 `"ok"`                                |
| `data` | object / array / null | 无载荷时 `data` 可为 `null` 或省略；**不要**把业务字段平铺到根上 |


**列表：**

```json
{
  "code": 0,
  "msg": "ok",
  "data": {
    "list": [],
    "total": 0
  }
}
```

- `list` 必为数组（空则 `[]`，禁止 `null`）。
- `total` 为满足条件的总数，不是当前页条数。
- 分页入参：`page` 从 **1** 起；`page_size` 默认 20（角色列表可 50），**最大 100**。`page<=0`、`page_size<=0` 时用默认值，不报错。



### 3.2 失败

- HTTP status **等于** body 里的 `code`（4xx/5xx）。
- Body **不返回** `data`：

```json
{
  "code": 401,
  "msg": "unauthorized"
}
```


| HTTP / code | 含义                               | 前端          |
| ----------- | -------------------------------- | ----------- |
| 400         | 参数错误、模式不符、重复 bootstrap、角色仍绑定用户等  | 展示 `msg`    |
| 401         | 未登录、Token 无效、salt 不匹配、黑名单、用户/厅停用 | 清 Token，跳登录 |
| 403         | Casbin 未通过、不可删 root/系统角色、跨厅伪造    | 提示无权限       |
| 404         | 本租户下资源不存在（`on` 下他厅数据也当 404，防探测）  | —           |
| 500         | 未分类错误                            | 不展示内部 SQL   |


---



## 4. 命名与代码风格



### 4.1 命名


| 对象     | 规则                  | 例                                           |
| ------ | ------------------- | ------------------------------------------- |
| 包名     | 小写、单单词              | `service` `casbinx`                         |
| 导出类型   | 驼峰，避免 `XxxInfoData` | `CreateUserReq` `UserPublic`                |
| JSON   | `snake_case`        | `role_code` `page_size` `access_token`      |
| DB 列   | `snake_case`        | 与 schema.md 一致                              |
| 布尔     | `is_` / `has_` 前缀   | `is_super_admin` `is_system`                |
| 接口路径动作 | 小写动词                | `list` `detail`                             |
| 测试     | `TestXxx_场景`        | `TestOnOperatorAdminLoginScopeAndIsolation` |


