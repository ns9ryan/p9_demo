package sync

import (
	"context"
	"sync"

	"gorm.io/gorm"
	"oa.98ent.com/p9/platform-game/common/logger"
)

// ProgressMessage 进度消息
type ProgressMessage struct {
	// 表名（如: category, provider, channel, game, currency）
	TableName string
	// 当前处理的条数
	ProcessedCount int32
	// 进度百分比（0-100）
	Progress int32
	// 检查点ID
	CheckpointID int64

	CheckpointValue string
	RemoteTotal     int32
	LocalTotal      int32

	Created int32
	Updated int32
	Deleted int32
	Failed  int32
	Skipped int32
}

// ProgressQueue 进度消息队列
type ProgressQueue struct {
	// 消息通道
	messageChan chan *ProgressMessage
	// 停止信号
	stopChan chan struct{}
	// WaitGroup 用于确保消费者优雅退出
	wg sync.WaitGroup
	// 数据库连接
	db *gorm.DB
	// 是否正在运行
	running bool
	// 锁
	mu sync.Mutex
}

const DefaultBufferSize = 1

// NewProgressQueue 创建进度消息队列
func NewProgressQueue(db *gorm.DB) *ProgressQueue {
	return &ProgressQueue{
		messageChan: make(chan *ProgressMessage, 100),
		stopChan:    make(chan struct{}),
		db:          db,
		running:     false,
	}
}

// Start 启动队列消费者
func (q *ProgressQueue) Start(ctx context.Context) {
	q.mu.Lock()
	if q.running {
		q.mu.Unlock()
		logger.Warn("[进度队列] 队列已启动，跳过重复启动")
		return
	}
	q.running = true
	q.mu.Unlock()

	logger.Info("[进度队列] 启动进度消息消费者")
	q.wg.Add(1)
	go q.consume(ctx)
}

// Stop 停止队列消费者
func (q *ProgressQueue) Stop() {
	q.mu.Lock()
	defer q.mu.Unlock()

	if !q.running {
		logger.Warn("[进度队列] 队列未启动")
		return
	}

	logger.Info("[进度队列] 停止进度消息消费者")
	close(q.stopChan)
	q.wg.Wait()
	q.running = false
	logger.Info("[进度队列] 进度消息消费者已停止")
}

// Send 发送进度消息到队列
func (q *ProgressQueue) Send(msg *ProgressMessage) error {
	if msg == nil {
		logger.Warn("[进度队列] 消息为空")
		return nil
	}

	// 非阻塞发送，防止队列满
	select {
	case q.messageChan <- msg:
		logger.Debugf("[进度队列] 发送进度消息: 表=%s, 处理数=%d, 总数=%d, 进度=%d%%",
			msg.TableName, msg.ProcessedCount, msg.RemoteTotal, msg.Progress)
		return nil
	case <-q.stopChan:
		logger.Warn("[进度队列] 队列已停止，消息发送失败")
		return nil
	default:
		// 队列满时，记录警告并继续（不阻塞发送方）
		logger.Warnf("[进度队列] 消息队列已满，跳过消息: 表=%s, 处理数=%d", msg.TableName, msg.ProcessedCount)
		return nil
	}
}

// consume 消费队列中的消息
func (q *ProgressQueue) consume(ctx context.Context) {
	defer q.wg.Done()
	logger.Info("[进度队列] 消费者开始监听消息")

	for {
		select {
		case msg := <-q.messageChan:
			if msg == nil {
				continue
			}
			q.handleProgressMessage(ctx, msg)

		case <-q.stopChan:
			// 处理停止前队列中剩余的消息
			for {
				select {
				case msg := <-q.messageChan:
					if msg != nil {
						q.handleProgressMessage(ctx, msg)
					}
				default:
					logger.Info("[进度队列] 消费者已处理完所有消息，正在退出")
					return
				}
			}

		case <-ctx.Done():
			logger.Info("[进度队列] 上下文已取消，消费者退出")
			return
		}
	}
}

// handleProgressMessage 处理进度消息
func (q *ProgressQueue) handleProgressMessage(ctx context.Context, msg *ProgressMessage) {
	if msg.CheckpointID <= 0 {
		logger.Debugf("[进度队列] 跳过处理消息（checkpointID 无效）: 表=%s, checkpointID=%d",
			msg.TableName, msg.CheckpointID)
		return
	}

	checkpointMgr := NewCheckpointManager(q.db)
	if err := checkpointMgr.UpdateCheckpointProgress(ctx, msg); err != nil {
		logger.Errorf("[进度队列] 更新检查点进度失败: %v", err)
		return
	}

	logger.Infof("[进度队列] ✓ 进度更新: 表=%s, 处理数=%d, 总数=%d, 进度=%d%%",
		msg.TableName, msg.ProcessedCount, msg.RemoteTotal, msg.LocalTotal, msg.Progress)
}

// CalculateProgress 计算进度百分比
// processedCount: 已处理条数
// totalCount: 总条数
// 返回进度百分比（0-100）
func CalculateProgress(processedCount, totalCount int64) int32 {
	if totalCount <= 0 {
		return 0
	}
	if processedCount >= totalCount {
		return 100
	}
	// 公式: processedCount / totalCount * 100
	progress := (processedCount * 100) / totalCount
	return int32(progress)
}
