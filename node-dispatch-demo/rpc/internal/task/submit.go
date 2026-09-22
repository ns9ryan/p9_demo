package task

import (
	"context"
	"encoding/json"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"oa.98ent.com/p9/node-dispatch/rpc/ent"
	"oa.98ent.com/p9/node-dispatch/rpc/ent/node"
	"oa.98ent.com/p9/node-dispatch/rpc/internal/protocol"

	"github.com/google/uuid"
)

// SubmitRequest 提交调度任务请求
type SubmitRequest struct {
	RequestNo string          // 调用方请求编号
	Target    string          // 目标服务
	TaskType  string          // 任务类型
	NodeCode  string          // 执行节点编码
	Params    json.RawMessage // 任务参数
}

// SubmitResponse 提交调度任务响应
type SubmitResponse struct {
	TaskNo string // 调度任务编号
}

// Submit 提交调度任务
func (s *Service) Submit(ctx context.Context, req SubmitRequest) (*SubmitResponse, error) {
	// 检查任务幂等
	taskData, exists, err := s.checkIdempotent(ctx, req)
	if err != nil {
		return nil, err
	}
	if exists {
		return &SubmitResponse{
			TaskNo: taskData.TaskNo, // 调度任务编号
		}, nil
	}

	// 获取执行节点
	nodeData, err := s.db.Node.
		Query().
		Where(node.CodeEQ(req.NodeCode)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, status.Error(codes.NotFound, "node not found")
		}

		return nil, fmt.Errorf("查询执行节点失败: %w", err)
	}

	// 检查节点状态
	if nodeData.Status != 1 {
		return nil, status.Error(codes.FailedPrecondition, "node is disabled")
	}

	// 检查节点连接状态
	if !s.connections.IsOnline(req.NodeCode) {
		return nil, status.Error(codes.Unavailable, "node is offline")
	}

	// 开启事务
	tx, err := s.db.Tx(ctx)
	if err != nil {
		return nil, fmt.Errorf("开启事务失败: %w", err)
	}

	// 创建调度任务
	taskData, err = s.createTask(ctx, tx, req)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	// 创建首次执行记录
	runData, err := s.createRun(ctx, tx, taskData.ID, nodeData.ID, 1)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	// 提交事务
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("提交事务失败: %w", err)
	}

	// 下发任务
	if err = s.dispatch(ctx, taskData.TaskNo, runData.RunNo, req); err != nil {
		// 标记任务下发失败
		s.markDispatchFailed(ctx, taskData.ID, runData.ID, err.Error())

		return nil, err
	}

	// 返回任务编号
	return &SubmitResponse{
		TaskNo: taskData.TaskNo, // 调度任务编号
	}, nil
}

// createTask 创建调度任务
func (s *Service) createTask(ctx context.Context, tx *ent.Tx, req SubmitRequest) (*ent.DispatchTask, error) {
	taskData, err := tx.DispatchTask.
		Create().
		SetTaskNo(uuid.NewString()). // 调度任务编号
		SetRequestNo(req.RequestNo). // 调用方请求编号
		SetTarget(req.Target).       // 目标服务
		SetTaskType(req.TaskType).   // 任务类型
		SetParams(req.Params).       // 任务参数
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("创建调度任务失败: %w", err)
	}

	return taskData, nil
}

// createRun 创建任务执行记录
func (s *Service) createRun(ctx context.Context, tx *ent.Tx, taskID int64, nodeID int64, runNo int64) (*ent.DispatchTaskRun, error) {
	runData, err := tx.DispatchTaskRun.
		Create().
		SetTaskID(taskID). // 调度任务ID
		SetNodeID(nodeID). // 执行节点ID
		SetRunNo(runNo).   // 执行序号
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("创建任务执行记录失败: %w", err)
	}

	return runData, nil
}

// dispatch 下发调度任务
func (s *Service) dispatch(ctx context.Context, taskNo string, runNo int64, req SubmitRequest) error {
	// 编码任务数据
	data, err := json.Marshal(protocol.TaskDispatchData{
		TaskNo:   taskNo,       // 调度任务编号
		RunNo:    runNo,        // 执行序号
		Target:   req.Target,   // 目标服务
		TaskType: req.TaskType, // 任务类型
		Params:   req.Params,   // 任务参数
	})
	if err != nil {
		return fmt.Errorf("编码任务数据失败: %w", err)
	}

	// 发送任务消息
	if err = s.connections.Send(ctx, req.NodeCode, protocol.Message{
		Type: protocol.MessageTypeTaskDispatch, // 消息类型
		Data: data,                             // 消息数据
	}); err != nil {
		return fmt.Errorf("下发任务失败: %w", err)
	}

	return nil
}
