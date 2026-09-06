// API使用示例

## 1. Ping接口

```bash
curl http://localhost:8080/ping
```

## 2. 游戏管理

### 创建游戏
```bash
curl -X POST http://localhost:8080/games \
  -H "Content-Type: application/json" \
  -d '{
    "code": "GAME001",
    "source_code": "SOURCE_GAME001",
    "name": "热血传奇",
    "category_id": 1,
    "provider_id": 2,
    "channel_id": 3,
    "logo": "https://example.com/logo.png",
    "status": 1,
    "sort_no": 10
  }'
```

### 获取游戏
```bash
curl http://localhost:8080/games/1
```

### 更新游戏
```bash
curl -X PUT http://localhost:8080/games/1 \
  -H "Content-Type: application/json" \
  -d '{
    "id": 1,
    "code": "GAME001",
    "source_code": "SOURCE_GAME001",
    "name": "热血传奇-更新",
    "category_id": 1,
    "provider_id": 2,
    "channel_id": 3,
    "status": 1,
    "sort_no": 20
  }'
```

### 删除游戏
```bash
curl -X DELETE http://localhost:8080/games/1
```

### 游戏列表
```bash
# 基础列表
curl http://localhost:8080/games

# 分页查询
curl "http://localhost:8080/games?page=1&page_size=20"

# 按分类过滤
curl "http://localhost:8080/games?category_id=1"

# 按厂商过滤
curl "http://localhost:8080/games?provider_id=2"

# 按渠道过滤
curl "http://localhost:8080/games?channel_id=3"

# 综合过滤
curl "http://localhost:8080/games?page=1&page_size=20&category_id=1&provider_id=2"
```

## 3. 分类管理

### 创建分类
```bash
curl -X POST http://localhost:8080/categories \
  -H "Content-Type: application/json" \
  -d '{
    "code": "CAT001",
    "source_code": "SOURCE_CAT001",
    "name": "棋牌类",
    "source_name": "Poker",
    "status": 1,
    "sort_no": 10
  }'
```

### 分类列表
```bash
curl "http://localhost:8080/categories?page=1&page_size=20"
```

## 4. 厂商管理

### 创建厂商
```bash
curl -X POST http://localhost:8080/providers \
  -H "Content-Type: application/json" \
  -d '{
    "code": "PROV001",
    "source_code": "SOURCE_PROV001",
    "name": "AG电子游艺",
    "source_name": "AG Gaming",
    "logo": "https://example.com/provider.png",
    "status": 1,
    "sort_no": 10
  }'
```

### 厂商列表
```bash
curl "http://localhost:8080/providers?page=1&page_size=20"
```

## 5. 渠道管理

### 创建渠道
```bash
curl -X POST http://localhost:8080/channels \
  -H "Content-Type: application/json" \
  -d '{
    "code": "CHAN001",
    "source_code": "SOURCE_CHAN001",
    "name": "官方渠道",
    "source_name": "Official",
    "status": 1,
    "sort_no": 10
  }'
```

### 渠道列表
```bash
curl "http://localhost:8080/channels?page=1&page_size=20"
```

## 6. 数据同步

### 同步预检查（Preview）
```bash
# 预检查分类同步
curl -X POST http://localhost:8080/sync/preview \
  -H "Content-Type: application/json" \
  -d '{
    "object_type": "category",
    "strict_conflict": false
  }'

# 预检查游戏同步（限定ID范围）
curl -X POST http://localhost:8080/sync/preview \
  -H "Content-Type: application/json" \
  -d '{
    "object_type": "game",
    "selected_ids": [1, 2, 3],
    "strict_conflict": false
  }'

# 预检查游戏同步（按编码过滤）
curl -X POST http://localhost:8080/sync/preview \
  -H "Content-Type: application/json" \
  -d '{
    "object_type": "game",
    "selected_codes": ["GAME001", "GAME002"],
    "strict_conflict": false
  }'
```

### 执行同步（Run）
```bash
# 真正执行分类同步
curl -X POST http://localhost:8080/sync/run \
  -H "Content-Type: application/json" \
  -d '{
    "object_type": "category",
    "auto_apply": true,
    "strict_conflict": false
  }'

# 预览游戏同步（不修改数据）
curl -X POST http://localhost:8080/sync/run \
  -H "Content-Type: application/json" \
  -d '{
    "object_type": "game",
    "auto_apply": false,
    "strict_conflict": false
  }'
```

## 7. 响应格式说明

### 成功响应
```json
{
  "code": 1,
  "message": "success",
  "data": {
    "id": 1,
    "code": "GAME001",
    "name": "热血传奇",
    ...
  }
}
```

### 错误响应
```json
{
  "code": 10001,
  "message": "游戏不存在",
  "data": null
}
```

### 列表响应
```json
{
  "code": 1,
  "message": "success",
  "data": {
    "list": [
      { "id": 1, "code": "GAME001", ... },
      { "id": 2, "code": "GAME002", ... }
    ],
    "total": 100
  }
}
```

## 8. gRPC调用示例

### 使用grpcurl
```bash
# 列出所有服务
grpcurl -plaintext localhost:9080 list

# 调用Ping服务
grpcurl -plaintext -d "{}" localhost:9080 game.PingService/Ping

# 调用创建游戏
grpcurl -plaintext -d @game.json localhost:9080 game.GameService/Create
```

game.json:
```json
{
  "code": "GAME001",
  "source_code": "SOURCE_GAME001",
  "name": "热血传奇",
  "category_id": 1,
  "provider_id": 2,
  "channel_id": 3
}
```
