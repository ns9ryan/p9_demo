package utils

import (
	"encoding/json"
	"time"

	"oa.98ent.com/p9/platform-game/common/model"
)

// MustMarshalJSON 将JSONMap序列化为JSON字符串
func MustMarshalJSON(data model.JSONMap) []byte {
	if data == nil {
		return []byte("{}")
	}
	bytes, _ := json.Marshal(data)
	return bytes
}

// MustParseJSON 将JSON字符串解析为JSONMap
func MustParseJSON(data []byte) model.JSONMap {
	var result model.JSONMap
	_ = json.Unmarshal(data, &result)
	if result == nil {
		result = make(model.JSONMap)
	}
	return result
}

// TimeToUnixTimestamp 将 *time.Time 转换为 *int64（Unix时间戳，单位：秒）
// 如果输入为 nil，则返回 nil
func TimeToUnixTimestamp(t *time.Time) *int64 {
	if t == nil {
		return nil
	}
	ts := t.Unix()
	return &ts
}

// JSON 将任意数据序列化为JSON字符串
func JSON(data interface{}) string {
	bytes, _ := json.Marshal(data)
	return string(bytes)
}

func StringToJSON(data string, v interface{}) error {
	return json.Unmarshal([]byte(data), v)
}
