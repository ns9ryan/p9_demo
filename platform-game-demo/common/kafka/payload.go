package kafka

import (
	"encoding/json"
	"time"
)

// ForceQuitEvent 游戏强制踢线事件
type ForceQuitEvent struct {
	Scope     string    `json:"scope"`      // 作用域：game_category/game_provider/game_channel/game
	PartnerID int64     `json:"partner_id"` // 分站ID，0表示总网事件
	GameIDs   []int64   `json:"game_ids"`   // 游戏ID列表
	QuitAt    time.Time `json:"quit_at"`    // 踢线时间
}

// Payload 将事件序列化为 JSON 字节
func (e *ForceQuitEvent) Payload() []byte {
	data, _ := json.Marshal(e)
	return data
}
