package utils

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/zeromicro/go-zero/core/logx"
)

type KafkaPayloadBasis struct {
	Stage    int    `json:"stage"`
	Producer string `json:"producer"`
}

type ForceQuitEvent struct {
	Basis     KafkaPayloadBasis `json:"basis"`             // 基础信息
	Scope     string            `json:"scope"`             // 作用域：game_category/game_provider/game_channel/game
	PartnerID int64             `json:"partner_id"`        // 分站ID，0表示总网事件
	CatID     int64             `json:"cat_id,omitempty"`  // ScopeGameCategory 携带
	VenID     int64             `json:"ven_id,omitempty"`  // ScopeGameVendor 携带
	ChanID    int64             `json:"chan_id,omitempty"` // ScopeGameChannel 携带
	GameIDs   []int64           `json:"game_ids"`          // 游戏ID列表
	QuitAt    time.Time         `json:"quit_at"`           // 踢线时间
}

var (
	writer   *kafka.Writer
	writerMu sync.Mutex
)

const (
	TopicGameForceQuit = "game.force-quit"
)

func buildPool(brokers []string) error {
	writerMu.Lock()
	defer writerMu.Unlock()

	if writer != nil {
		return nil
	}

	if len(brokers) == 0 {
		return errors.New("kafka brokers not configured")
	}

	w := &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Balancer: &kafka.LeastBytes{},
	}

	writer = w
	logx.Infof("kafka writer initialized with brokers: %v", brokers)
	return nil
}

// getWriter 获取 kafka writer
func getWriter() *kafka.Writer {
	writerMu.Lock()
	defer writerMu.Unlock()
	return writer
}

func Write(ctx context.Context, brokers []string, message kafka.Message) error {
	if len(brokers) == 0 {
		logx.WithContext(ctx).Errorf("kafka brokers not configured")
		return errors.New("kafka brokers not configured")
	}

	// 确保连接池已初始化
	if getWriter() == nil {
		if err := buildPool(brokers); err != nil {
			return err
		}
	}

	w := getWriter()
	if w == nil {
		return errors.New("no kafka writer available")
	}

	logx.WithContext(ctx).Debugf("sending kafka message to topic=%s", message.Topic)
	return w.WriteMessages(ctx, message)
}

// SendForceQuitEvent 发送游戏强制踢线事件
func SendForceQuitEvent(ctx context.Context, brokers []string, payload *ForceQuitEvent) error {
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

	return Write(ctx, brokers, message)
}
