package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"oa.98ent.com/p9/node-agent/internal/protocol"

	coderws "github.com/coder/websocket"
	"github.com/zeromicro/go-zero/core/logx"
)

// handleMessage 处理调度中心WebSocket消息
func (c *Client) handleMessage(ctx context.Context, conn *coderws.Conn, messageType coderws.MessageType, data []byte) error {
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
	case protocol.MessageTypeTaskDispatch:
		return c.handleTaskDispatch(ctx, conn, message.Data)

	default:
		return fmt.Errorf("不支持的WebSocket业务消息类型: %s", message.Type)
	}
}

// handleTaskDispatch 处理任务下发消息
func (c *Client) handleTaskDispatch(ctx context.Context, conn *coderws.Conn, data json.RawMessage) error {
	// 解析任务下发数据
	var taskData protocol.TaskDispatchData
	if err := json.Unmarshal(data, &taskData); err != nil {
		return fmt.Errorf("解析任务下发数据失败: %w", err)
	}

	// 校验任务基本信息
	if strings.TrimSpace(taskData.TaskNo) == "" {
		return fmt.Errorf("任务编号不能为空")
	}
	if taskData.RunNo <= 0 {
		return fmt.Errorf("任务执行序号必须大于0")
	}
	if strings.TrimSpace(taskData.Target) == "" {
		return fmt.Errorf("任务目标不能为空")
	}
	if strings.TrimSpace(taskData.TaskType) == "" {
		return fmt.Errorf("任务类型不能为空")
	}

	logger := logx.WithContext(ctx)
	logger.Infow(
		"收到调度任务",
		logx.Field("task_no", taskData.TaskNo),
		logx.Field("run_no", taskData.RunNo),
		logx.Field("target", taskData.Target),
		logx.Field("task_type", taskData.TaskType),
	)

	// 确认任务已经接收
	if err := c.sendTaskAck(ctx, conn, taskData); err != nil {
		return err
	}

	// 执行调度任务
	result, err := c.task.Execute(ctx, taskData.Target, taskData.TaskType, taskData.Params)
	if err != nil {
		// 返回任务执行失败结果
		if sendErr := c.sendTaskResult(ctx, conn, taskData, false, nil, err.Error()); sendErr != nil {
			return sendErr
		}

		logger.Errorw(
			"调度任务执行失败",
			logx.Field("task_no", taskData.TaskNo),
			logx.Field("run_no", taskData.RunNo),
			logx.Field("error", err.Error()),
		)
		return nil
	}

	// 返回任务执行成功结果
	if err = c.sendTaskResult(ctx, conn, taskData, true, result, ""); err != nil {
		return err
	}

	logger.Infow("调度任务执行成功", logx.Field("task_no", taskData.TaskNo), logx.Field("run_no", taskData.RunNo))

	return nil
}

// sendTaskAck 发送任务接收确认
func (c *Client) sendTaskAck(ctx context.Context, conn *coderws.Conn, taskData protocol.TaskDispatchData) error {
	// 编码任务接收确认数据
	data, err := json.Marshal(protocol.TaskAckData{
		TaskNo: taskData.TaskNo, // 任务编号
		RunNo:  taskData.RunNo,  // 执行序号
	})
	if err != nil {
		return fmt.Errorf("编码任务接收确认数据失败: %w", err)
	}

	// 发送任务接收确认
	if err = c.sendMessage(ctx, conn, protocol.Message{
		Type: protocol.MessageTypeTaskAck, // 消息类型
		Data: data,                        // 确认数据
	}); err != nil {
		return fmt.Errorf("发送任务接收确认失败: %w", err)
	}

	return nil
}

// sendTaskResult 发送任务执行结果
func (c *Client) sendTaskResult(
	ctx context.Context,
	conn *coderws.Conn,
	taskData protocol.TaskDispatchData,
	success bool,
	result json.RawMessage,
	errorMessage string,
) error {
	// 编码任务执行结果
	data, err := json.Marshal(protocol.TaskResultData{
		TaskNo:       taskData.TaskNo, // 任务编号
		RunNo:        taskData.RunNo,  // 执行序号
		Success:      success,         // 是否执行成功
		Result:       result,          // 执行结果
		ErrorMessage: errorMessage,    // 失败原因
	})
	if err != nil {
		return fmt.Errorf("编码任务执行结果失败: %w", err)
	}

	// 发送任务执行结果
	if err = c.sendMessage(ctx, conn, protocol.Message{
		Type: protocol.MessageTypeTaskResult, // 消息类型
		Data: data,                           // 执行结果
	}); err != nil {
		return fmt.Errorf("发送任务执行结果失败: %w", err)
	}

	return nil
}
