package kafka

import (
	"context"
	"encoding/json"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/zeromicro/go-zero/core/logx"
)

// SendForceQuitEvent 发送游戏强制踢线事件
func SendForceQuitEvent(ctx context.Context, resourceID int64, scope string) error {
	payload := &ForceQuitEvent{
		Scope:     scope,
		PartnerID: 0, // 0 表示总网事件，影响所有分站
		GameIDs:   []int64{resourceID},
		QuitAt:    time.Now(),
	}

	// 序列化消息
	data, err := json.Marshal(payload)
	if err != nil {
		logx.WithContext(ctx).Errorf("序列化 ForceQuitEvent 失败: %v", err)
		return err
	}

	// 使用 Write 函数发送消息
	message := kafka.Message{
		Topic: TopicGameForceQuit,
		Value: data,
	}

	return Write(ctx, message)
}
