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

// ListItem 调度任务列表项
type ListItem struct {
	Task     *ent.DispatchTask // 调度任务
	NodeCode string            // 当前执行节点编码
}

// ListResult 调度任务列表结果
type ListResult struct {
	Total int64       // 数据总数
	List  []*ListItem // 调度任务列表
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

	// 按任务类型筛选
	if req.TaskType != "" {
		query = query.Where(dispatchtask.TaskTypeEQ(req.TaskType))
	}

	// 按任务状态筛选
	if req.Status != nil {
		query = query.Where(dispatchtask.StatusEQ(*req.Status))
	}

	// 按执行节点筛选
	if req.NodeCode != "" {
		query = query.Where(
			dispatchtask.HasRunsWith(
				dispatchtaskrun.HasNodeWith(node.CodeEQ(req.NodeCode)),
			),
		)
	}

	// 查询数据总数
	total, err := query.Clone().Count(ctx)
	if err != nil {
		return nil, fmt.Errorf("统计任务数量失败: %w", err)
	}

	// 查询当前页任务
	taskList, err := query.
		Order(ent.Desc(dispatchtask.FieldCreatedAt), ent.Desc(dispatchtask.FieldID)).
		Limit(int(req.PageSize)).
		Offset(int((req.Page - 1) * req.PageSize)).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("查询任务列表失败: %w", err)
	}

	// 当前页没有数据时直接返回
	if len(taskList) == 0 {
		return &ListResult{
			Total: int64(total),  // 数据总数
			List:  []*ListItem{}, // 调度任务列表
		}, nil
	}

	// 收集当前页任务ID
	taskIDs := make([]int64, 0, len(taskList))
	for _, taskData := range taskList {
		taskIDs = append(taskIDs, taskData.ID)
	}

	// 查询当前页任务的执行记录
	runList, err := s.db.DispatchTaskRun.
		Query().
		Where(dispatchtaskrun.TaskIDIn(taskIDs...)).
		WithNode().
		Order(
			ent.Asc(dispatchtaskrun.FieldTaskID),
			ent.Desc(dispatchtaskrun.FieldRunNo),
		).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("查询任务执行记录失败: %w", err)
	}

	// 获取每个任务最新一次执行节点
	nodeCodeMap := make(map[int64]string, len(taskList))
	for _, runData := range runList {
		if _, exists := nodeCodeMap[runData.TaskID]; exists {
			continue
		}

		nodeData, err := runData.Edges.NodeOrErr()
		if err != nil {
			return nil, fmt.Errorf("获取任务执行节点失败: %w", err)
		}

		nodeCodeMap[runData.TaskID] = nodeData.Code
	}

	// 组装列表结果
	list := make([]*ListItem, 0, len(taskList))
	for _, taskData := range taskList {
		list = append(list, &ListItem{
			Task:     taskData,                 // 调度任务
			NodeCode: nodeCodeMap[taskData.ID], // 当前执行节点编码
		})
	}

	return &ListResult{
		Total: int64(total), // 数据总数
		List:  list,         // 调度任务列表
	}, nil
}
