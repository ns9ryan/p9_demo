package task

import (
	"context"
	"fmt"

	"oa.98ent.com/p9/node-dispatch/rpc/ent"
	"oa.98ent.com/p9/node-dispatch/rpc/ent/dispatchtask"
	"oa.98ent.com/p9/node-dispatch/rpc/ent/dispatchtaskrun"
	"oa.98ent.com/p9/node-dispatch/rpc/ent/node"
)

// ListRequest 调度任务列表请求
type ListRequest struct {
	Page     int64  // 页码
	PageSize int64  // 每页数量
	Keyword  string // 搜索关键字
	TaskType string // 任务类型
	Status   *int64 // 任务状态
	NodeCode string // 执行节点编码
}

// ListResult 调度任务列表结果
type ListResult struct {
	Total int64               // 数据总数
	List  []*ent.DispatchTask // 调度任务列表
}

// List 查询调度任务列表
func (s *Service) List(ctx context.Context, req ListRequest) (*ListResult, error) {
	// 创建任务查询
	query := s.db.DispatchTask.Query()

	// 搜索任务编号或请求编号
	if req.Keyword != "" {
		query = query.Where(
			dispatchtask.Or(
				dispatchtask.TaskNoContains(req.Keyword),
				dispatchtask.RequestNoContains(req.Keyword),
			),
		)
	}

	// 查询任务类型
	if req.TaskType != "" {
		query = query.Where(dispatchtask.TaskTypeEQ(req.TaskType))
	}

	// 查询任务状态
	if req.Status != nil {
		query = query.Where(dispatchtask.StatusEQ(*req.Status))
	}

	// 查询执行节点
	if req.NodeCode != "" {
		query = query.Where(
			dispatchtask.HasRunsWith(
				dispatchtaskrun.HasNodeWith(
					node.CodeEQ(req.NodeCode),
				),
			),
		)
	}

	// 查询总数量
	total, err := query.Count(ctx)
	if err != nil {
		return nil, fmt.Errorf("统计任务数量失败: %w", err)
	}

	// 查询任务列表
	list, err := query.
		Order(ent.Desc(dispatchtask.FieldCreatedAt)).
		Limit(int(req.PageSize)).
		Offset(int((req.Page - 1) * req.PageSize)).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("查询任务列表失败: %w", err)
	}

	return &ListResult{
		Total: int64(total), // 数据总数
		List:  list,         // 调度任务列表
	}, nil
}
