package websocketserver

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"oa.98ent.com/p9/node-dispatch/rpc/ent"
	"oa.98ent.com/p9/node-dispatch/rpc/ent/dispatchtask"
	"oa.98ent.com/p9/node-dispatch/rpc/ent/dispatchtaskrun"
	"oa.98ent.com/p9/node-dispatch/rpc/internal/websocket"

	coderws "github.com/coder/websocket"
	"github.com/zeromicro/go-zero/core/logx"
)

// handleMessage 处理节点WebSocket业务消息
func (s *Server) handleMessage(ctx context.Context, nodeID int64, nodeCode string, messageType coderws.MessageType, data []byte) error {
	// 只处理文本消息
	if messageType != coderws.MessageText {
		return fmt.Errorf("不支持的WebSocket消息类型: %d", messageType)
	}

	// 解析业务消息
	var message websocket.Message
	if err := json.Unmarshal(data, &message); err != nil {
		return fmt.Errorf("解析WebSocket消息失败: %w", err)
	}

	// 根据消息类型处理
	switch message.Type {
	case websocket.MessageTypeTaskAck:
		return s.handleTaskAck(ctx, nodeID, nodeCode, message.Data)

	default:
		return fmt.Errorf("不支持的WebSocket业务消息类型: %s", message.Type)
	}
}

// handleTaskAck 处理任务接收确认
func (s *Server) handleTaskAck(ctx context.Context, nodeID int64, nodeCode string, data json.RawMessage) error {
	// ---------- 消息解析 ----------

	// 解析任务确认数据
	var ack websocket.TaskAckData
	if err := json.Unmarshal(data, &ack); err != nil {
		return fmt.Errorf("解析任务确认数据失败: %w", err)
	}

	// 校验任务编号
	taskNo := strings.TrimSpace(ack.TaskNo)
	if taskNo == "" {
		return fmt.Errorf("任务编号不能为空")
	}

	// ---------- 任务查询 ----------

	// 查询调度任务
	taskData, err := s.svcCtx.DB.DispatchTask.
		Query().
		Where(dispatchtask.TaskNoEQ(taskNo)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return fmt.Errorf("调度任务不存在: %s", taskNo)
		}

		return fmt.Errorf("查询调度任务失败: %w", err)
	}

	// 查询当前节点最近一次执行记录
	runData, err := s.svcCtx.DB.DispatchTaskRun.
		Query().
		Where(
			dispatchtaskrun.TaskIDEQ(taskData.ID),
			dispatchtaskrun.NodeIDEQ(nodeID),
		).
		Order(ent.Desc(dispatchtaskrun.FieldID)).
		First(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return fmt.Errorf("调度任务执行记录不存在: %s", taskNo)
		}

		return fmt.Errorf("查询调度任务执行记录失败: %w", err)
	}

	// ---------- 幂等处理 ----------

	// 已经处理过的确认直接忽略, 不重复修改开始时间
	if taskData.Status != 1 || runData.Status != 1 {
		return nil
	}

	// ---------- 状态更新 ----------

	// 开启事务
	tx, err := s.svcCtx.DB.Tx(ctx)
	if err != nil {
		return fmt.Errorf("开启任务确认事务失败: %w", err)
	}

	// 标记调度任务执行中
	err = tx.DispatchTask.
		UpdateOneID(taskData.ID).
		SetStatus(2).
		Exec(ctx)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("更新调度任务状态失败: %w", err)
	}

	// 标记任务执行记录执行中
	err = tx.DispatchTaskRun.
		UpdateOneID(runData.ID).
		SetStatus(2).
		SetStartedAt(time.Now()).
		Exec(ctx)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("更新调度任务执行状态失败: %w", err)
	}

	// 提交事务
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("提交任务确认事务失败: %w", err)
	}

	logx.WithContext(ctx).Infow("节点已确认调度任务", logx.Field("task_no", taskNo), logx.Field("node_code", nodeCode))

	return nil
}
