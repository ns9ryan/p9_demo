package task

import (
	"context"
	"fmt"

	"oa.98ent.com/p9/node-dispatch/rpc/ent"
	"oa.98ent.com/p9/node-dispatch/rpc/ent/dispatchtask"
	"oa.98ent.com/p9/node-dispatch/rpc/ent/dispatchtaskrun"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetRequest 获取调度任务请求
type GetRequest struct {
	TaskNo    string // 调度任务编号
	RequestNo string // 调用方请求编号
}

// GetResult 获取调度任务结果
type GetResult struct {
	Task *ent.DispatchTask      // 调度任务
	Runs []*ent.DispatchTaskRun // 执行记录
}

// Get 获取调度任务
func (s *Service) Get(ctx context.Context, req GetRequest) (*GetResult, error) {
	// 创建任务查询
	query := s.db.DispatchTask.Query()

	// 根据指定编号查询任务
	if req.TaskNo != "" {
		query = query.Where(dispatchtask.TaskNoEQ(req.TaskNo))
	} else {
		query = query.Where(dispatchtask.RequestNoEQ(req.RequestNo))
	}

	// 获取调度任务
	taskData, err := query.Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, status.Error(codes.NotFound, "task not found")
		}

		return nil, fmt.Errorf("查询调度任务失败: %w", err)
	}

	// 获取全部执行记录及执行节点
	runList, err := s.db.DispatchTaskRun.
		Query().
		Where(dispatchtaskrun.TaskIDEQ(taskData.ID)).
		WithNode().
		Order(ent.Asc(dispatchtaskrun.FieldRunNo)).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("查询任务执行记录失败: %w", err)
	}

	// 调度任务必须至少存在一条执行记录
	if len(runList) == 0 {
		return nil, fmt.Errorf("调度任务执行记录不存在: %s", taskData.TaskNo)
	}

	// 返回任务详情
	return &GetResult{
		Task: taskData, // 调度任务
		Runs: runList,  // 执行记录
	}, nil
}
