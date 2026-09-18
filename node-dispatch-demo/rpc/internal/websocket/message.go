package websocket

import "encoding/json"

const (
	MessageTypeTaskDispatch = "task_dispatch" // 下发任务
	MessageTypeTaskAck      = "task_ack"      // 任务接收确认
)

// Message WebSocket业务消息
type Message struct {
	Type string          `json:"type"` // 消息类型
	Data json.RawMessage `json:"data"` // 消息数据
}

// TaskDispatchData 任务下发数据
type TaskDispatchData struct {
	TaskNo   string          `json:"task_no"`   // 任务编号
	TaskType string          `json:"task_type"` // 任务类型
	Params   json.RawMessage `json:"params"`    // 任务参数
}

// TaskAckData 任务接收确认数据
type TaskAckData struct {
	TaskNo string `json:"task_no"` // 任务编号
}
