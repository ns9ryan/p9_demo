package protocol

import "encoding/json"

// Message WebSocket通用消息
type Message struct {
	Type string          `json:"type"` // 消息类型
	Data json.RawMessage `json:"data"` // 消息数据
}
