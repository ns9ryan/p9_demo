package kafka

import (
	"context"
	"errors"
	"sync"

	"github.com/segmentio/kafka-go"
	"github.com/zeromicro/go-zero/core/logx"
)

var (
	writer   *kafka.Writer
	writerMu sync.Mutex
)

// buildPool 构建 kafka writer 连接池
func buildPool() error {
	writerMu.Lock()
	defer writerMu.Unlock()

	if writer != nil {
		return nil
	}

	cfg := GetConfig()
	if cfg == nil || len(cfg.Brokers) == 0 {
		return errors.New("kafka brokers not configured")
	}

	w := &kafka.Writer{
		Addr:     kafka.TCP(cfg.Brokers...),
		Balancer: &kafka.LeastBytes{},
	}

	writer = w
	logx.Infof("kafka writer initialized with brokers: %v", cfg.Brokers)
	return nil
}

// getWriter 获取 kafka writer
func getWriter() *kafka.Writer {
	writerMu.Lock()
	defer writerMu.Unlock()
	return writer
}

// Write 发送单条消息到 kafka
func Write(ctx context.Context, message kafka.Message) error {
	config := GetConfig()
	if config == nil || len(config.Brokers) == 0 {
		logx.WithContext(ctx).Errorf("kafka config not set or brokers empty")
		return errors.New("kafka brokers not configured")
	}

	// 确保连接池已初始化
	if getWriter() == nil {
		if err := buildPool(); err != nil {
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

// WriteBatch 批量发送消息到 kafka
// 在单次调用中写入多条消息，由 kafka-go 按分区合并，减少发送往返
// 消息可携带不同的 Topic (Writer 未固定 Topic 时按 message.Topic 路由)
func WriteBatch(ctx context.Context, messages ...kafka.Message) error {
	if len(messages) == 0 {
		return nil
	}

	config := GetConfig()
	if config == nil || len(config.Brokers) == 0 {
		logx.WithContext(ctx).Errorf("kafka config not set or brokers empty")
		return errors.New("kafka brokers not configured")
	}

	// 确保连接池已初始化
	if getWriter() == nil {
		if err := buildPool(); err != nil {
			return err
		}
	}

	w := getWriter()
	if w == nil {
		return errors.New("no kafka writer available")
	}

	logx.WithContext(ctx).Debugf("sending %d kafka messages", len(messages))
	return w.WriteMessages(ctx, messages...)
}

// Produce 发送原始消息到 kafka
func Produce(ctx context.Context, topic string, key, value []byte) error {
	message := kafka.Message{
		Topic: topic,
		Key:   key,
		Value: value,
	}
	return Write(ctx, message)
}
