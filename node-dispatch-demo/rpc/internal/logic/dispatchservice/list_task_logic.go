package dispatchservicelogic

import (
	"context"
	"strings"

	"entgo.io/ent/dialect/sql"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"oa.98ent.com/p9/node-dispatch/rpc/ent"
	"oa.98ent.com/p9/node-dispatch/rpc/ent/dispatchtask"
	"oa.98ent.com/p9/node-dispatch/rpc/ent/dispatchtaskrun"
	"oa.98ent.com/p9/node-dispatch/rpc/ent/node"
	"oa.98ent.com/p9/node-dispatch/rpc/internal/svc"
	"oa.98ent.com/p9/node-dispatch/rpc/pb/nodedispatchrpc/dispatchpb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListTaskLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListTaskLogic {
	return &ListTaskLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ListTask 获取调度任务列表
func (l *ListTaskLogic) ListTask(in *dispatchpb.ListTasksRequest) (*dispatchpb.ListTasksResponse, error) {
	// ---------- 参数校验 ----------

	// 校验分页参数
	if in.Page < 1 || in.PageSize < 1 || in.PageSize > 100 {
		return nil, status.Error(codes.InvalidArgument, "invalid pagination")
	}

	// 校验任务状态
	if in.Status != nil && (*in.Status < 1 || *in.Status > 4) {
		return nil, status.Error(codes.InvalidArgument, "invalid task status")
	}

	// ---------- 查询条件 ----------

	// 创建任务查询
	query := l.svcCtx.DB.DispatchTask.Query()

	// 按关键字筛选
	if in.Keyword != nil {
		keyword := strings.TrimSpace(*in.Keyword)
		if keyword != "" {
			query = query.Where(
				dispatchtask.Or(
					dispatchtask.TaskNoContainsFold(keyword),
					dispatchtask.RequestNoContainsFold(keyword),
				),
			)
		}
	}

	// 按任务类型筛选
	if in.TaskType != nil {
		taskType := strings.TrimSpace(*in.TaskType)
		if taskType != "" {
			query = query.Where(dispatchtask.TaskTypeEQ(taskType))
		}
	}

	// 按任务状态筛选
	if in.Status != nil {
		query = query.Where(dispatchtask.StatusEQ(*in.Status))
	}

	// 按执行节点筛选
	if in.NodeCode != nil {
		nodeCode := strings.TrimSpace(*in.NodeCode)
		if nodeCode != "" {
			query = query.Where(
				dispatchtask.HasRunsWith(
					dispatchtaskrun.HasNodeWith(
						node.CodeEQ(nodeCode),
					),
				),
			)
		}
	}

	// ---------- 数据查询 ----------

	// 获取符合条件的数据总数
	total, err := query.Clone().Count(l.ctx)
	if err != nil {
		l.Logger.Errorw("获取调度任务总数失败", logx.Field("error", err.Error()))
		return nil, err
	}

	// 计算分页偏移量
	offset := (in.Page - 1) * in.PageSize

	// 获取当前页任务数据
	results, err := query.
		WithRuns(func(query *ent.DispatchTaskRunQuery) {
			query.
				WithNode().
				Order(dispatchtaskrun.ByRunNo(sql.OrderDesc()))
		}).
		Order(
			dispatchtask.ByCreatedAt(sql.OrderDesc()), // 按创建时间倒序
			dispatchtask.ByID(sql.OrderDesc()),        // 创建时间相同时按ID倒序
		).
		Offset(int(offset)).
		Limit(int(in.PageSize)).
		All(l.ctx)
	if err != nil {
		l.Logger.Errorw("获取调度任务列表失败", logx.Field("error", err.Error()))
		return nil, err
	}

	// ---------- 响应转换 ----------

	// 转换任务列表
	list := make([]*dispatchpb.TaskInfo, 0, len(results))
	for _, taskData := range results {
		var nodeCode string

		// 最近一次执行记录的节点作为任务当前执行节点
		if len(taskData.Edges.Runs) > 0 {
			runData := taskData.Edges.Runs[0]

			nodeData, nodeErr := runData.Edges.NodeOrErr()
			if nodeErr != nil {
				l.Logger.Errorw(
					"获取调度任务执行节点失败",
					logx.Field("task_no", taskData.TaskNo),
					logx.Field("run_no", runData.RunNo),
					logx.Field("error", nodeErr.Error()),
				)
				return nil, nodeErr
			}

			nodeCode = nodeData.Code
		}

		list = append(list, &dispatchpb.TaskInfo{
			TaskNo:    taskData.TaskNo,                // 调度任务编号
			RequestNo: taskData.RequestNo,             // 调用方请求编号
			Target:    taskData.Target,                // 目标服务
			TaskType:  taskData.TaskType,              // 任务类型
			Status:    taskData.Status,                // 任务状态
			NodeCode:  nodeCode,                       // 最近执行节点编码
			CreatedAt: taskData.CreatedAt.UnixMilli(), // 创建时间
			UpdatedAt: taskData.UpdatedAt.UnixMilli(), // 更新时间
		})
	}

	// 返回任务列表
	return &dispatchpb.ListTasksResponse{
		Total: int64(total), // 数据总数
		List:  list,         // 调度任务列表
	}, nil
}
