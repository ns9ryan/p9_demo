package protocol

import "encoding/json"

const (
	MessageTypeTaskDispatch = "task_dispatch" // 下发任务
	MessageTypeTaskAck      = "task_ack"      // 任务接收确认
	MessageTypeTaskResult   = "task_result"   // 任务执行结果
)

// TaskDispatchData 任务下发数据
type TaskDispatchData struct {
	TaskNo   string          `json:"task_no"`   // 任务编号
	RunNo    int64           `json:"run_no"`    // 执行序号
	Target   string          `json:"target"`    // 目标服务
	TaskType string          `json:"task_type"` // 任务类型
	Params   json.RawMessage `json:"params"`    // 任务参数
}

// TaskAckData 任务接收确认数据
type TaskAckData struct {
	TaskNo string `json:"task_no"` // 任务编号
	RunNo  int64  `json:"run_no"`  // 执行序号
}

// TaskResultData 任务执行结果数据
type TaskResultData struct {
	TaskNo       string          `json:"task_no"`                 // 任务编号
	RunNo        int64           `json:"run_no"`                  // 执行序号
	Success      bool            `json:"success"`                 // 是否执行成功
	Result       json.RawMessage `json:"result,omitempty"`        // 执行结果
	ErrorMessage string          `json:"error_message,omitempty"` // 失败原因
}
