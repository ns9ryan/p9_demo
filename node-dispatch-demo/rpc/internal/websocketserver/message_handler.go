package websocketserver

import (
	"context"
	"encoding/json"
	"fmt"

	"oa.98ent.com/p9/node-dispatch/rpc/internal/protocol"

	coderws "github.com/coder/websocket"
	"github.com/zeromicro/go-zero/core/logx"
)

// handleMessage 处理节点WebSocket业务消息
func (s *Server) handleMessage(ctx context.Context, nodeID int64, nodeCode string, messageType coderws.MessageType, data []byte) error {
	// 只处理文本消息
	if messageType != coderws.MessageText {
		return fmt.Errorf("不支持的WebSocket消息类型: %d", messageType)
	}

	// 解析通用消息
	var message protocol.Message
	if err := json.Unmarshal(data, &message); err != nil {
		return fmt.Errorf("解析WebSocket消息失败: %w", err)
	}

	// 根据消息类型处理
	switch message.Type {
	case protocol.MessageTypeTaskAck:
		return s.handleTaskAck(ctx, nodeCode, message.Data)

	case protocol.MessageTypeTaskResult:
		return s.handleTaskResult(ctx, nodeCode, message.Data)

	default:
		return fmt.Errorf("不支持的WebSocket消息类型: %s", message.Type)
	}
}

// handleTaskAck 处理任务接收确认
func (s *Server) handleTaskAck(ctx context.Context, nodeCode string, data json.RawMessage) error {
	// 解析任务确认数据
	var ack protocol.TaskAckData
	if err := json.Unmarshal(data, &ack); err != nil {
		return fmt.Errorf("解析任务确认数据失败: %w", err)
	}

	// 标记任务执行中
	if err := s.svcCtx.Task.MarkRunning(ctx, ack.TaskNo, ack.RunNo, nodeCode); err != nil {
		return fmt.Errorf("更新任务执行状态失败: %w", err)
	}

	// 记录任务确认日志
	logx.WithContext(ctx).Infow(
		"节点已确认调度任务",
		logx.Field("task_no", ack.TaskNo),
		logx.Field("run_no", ack.RunNo),
		logx.Field("node_code", nodeCode),
	)

	return nil
}

// handleTaskResult 处理任务执行结果
func (s *Server) handleTaskResult(ctx context.Context, nodeCode string, data json.RawMessage) error {
	// 解析任务结果数据
	var result protocol.TaskResultData
	if err := json.Unmarshal(data, &result); err != nil {
		return fmt.Errorf("解析任务执行结果失败: %w", err)
	}

	// 标记任务执行结果
	if err := s.svcCtx.Task.MarkResult(
		ctx,
		result.TaskNo,       // 任务编号
		result.RunNo,        // 执行序号
		nodeCode,            // 节点编号
		result.Success,      // 是否执行成功
		result.Result,       // 执行结果
		result.ErrorMessage, // 失败原因
	); err != nil {
		return fmt.Errorf("更新任务执行结果失败: %w", err)
	}

	// 记录任务结果日志
	logx.WithContext(ctx).Infow(
		"节点已返回调度任务结果",
		logx.Field("task_no", result.TaskNo),
		logx.Field("run_no", result.RunNo),
		logx.Field("node_code", nodeCode),
		logx.Field("success", result.Success),
	)

	return nil
}
