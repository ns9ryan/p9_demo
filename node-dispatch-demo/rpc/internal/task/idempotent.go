package task

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"oa.98ent.com/p9/node-dispatch/rpc/ent"
	"oa.98ent.com/p9/node-dispatch/rpc/ent/dispatchtask"
	"oa.98ent.com/p9/node-dispatch/rpc/ent/dispatchtaskrun"
)

// checkIdempotent 检查任务提交幂等性
func (s *Service) checkIdempotent(ctx context.Context, req SubmitRequest) (*ent.DispatchTask, bool, error) {
	// 根据请求编号查询已有任务
	taskData, err := s.db.DispatchTask.Query().Where(dispatchtask.RequestNoEQ(req.RequestNo)).Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, false, nil
		}

		return nil, false, fmt.Errorf("查询幂等任务失败: %w", err)
	}

	// 获取首次执行记录
	runData, err := s.db.DispatchTaskRun.
		Query().
		Where(
			dispatchtaskrun.TaskIDEQ(taskData.ID),
			dispatchtaskrun.RunNoEQ(1),
		).
		Only(ctx)
	if err != nil {
		return nil, false, err
	}

	// 获取首次执行节点
	nodeData, err := s.db.Node.Get(ctx, runData.NodeID)
	if err != nil {
		return nil, false, fmt.Errorf("查询首次任务执行节点失败: %w", err)
	}

	// 比较任务参数
	sameParams, err := equalJSON(taskData.Params, req.Params)
	if err != nil {
		return nil, false, fmt.Errorf("比较任务参数失败: %w", err)
	}

	// 相同请求编号必须保持请求内容一致
	if taskData.Target != req.Target ||
		taskData.TaskType != req.TaskType ||
		nodeData.Code != req.NodeCode ||
		!sameParams {
		return nil, false, status.Error(codes.AlreadyExists, "request_no already exists with different content")
	}

	return taskData, true, nil
}

// equalJSON 比较JSON内容是否一致
func equalJSON(left, right json.RawMessage) (bool, error) {
	leftValue, err := decodeJSON(left)
	if err != nil {
		return false, err
	}

	rightValue, err := decodeJSON(right)
	if err != nil {
		return false, err
	}

	return reflect.DeepEqual(leftValue, rightValue), nil
}

// decodeJSON 解析JSON并保留数字精度
func decodeJSON(data json.RawMessage) (any, error) {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()

	var value any
	if err := decoder.Decode(&value); err != nil {
		return nil, err
	}

	return value, nil
}
