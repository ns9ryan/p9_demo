package dispatchservicelogic

import (
	"context"
	"strings"

	"oa.98ent.com/p9/node-dispatch/rpc/internal/svc"
	"oa.98ent.com/p9/node-dispatch/rpc/internal/task"
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
	// 整理分页参数
	page := in.GetPage()
	pageSize := in.GetPageSize()

	// 设置默认分页
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	// 整理查询参数
	keyword := strings.TrimSpace(in.GetKeyword())
	taskType := strings.TrimSpace(in.GetTaskType())
	nodeCode := strings.TrimSpace(in.GetNodeCode())

	var status *int64
	if in.Status != nil {
		value := in.GetStatus()
		status = &value
	}

	// 查询任务列表
	result, err := l.svcCtx.Task.List(l.ctx, task.ListRequest{
		Page:     page,     // 页码
		PageSize: pageSize, // 每页数量
		Keyword:  keyword,  // 搜索关键字
		TaskType: taskType, // 任务类型
		Status:   status,   // 任务状态
		NodeCode: nodeCode, // 执行节点编码
	})
	if err != nil {
		l.Logger.Errorw(
			"获取调度任务列表失败",
			logx.Field("error", err.Error()),
		)
		return nil, err
	}

	// 转换任务列表
	list := make([]*dispatchpb.TaskInfo, 0, len(result.List))

	for _, taskData := range result.List {
		list = append(list, &dispatchpb.TaskInfo{
			TaskNo:    taskData.TaskNo,                // 调度任务编号
			RequestNo: taskData.RequestNo,             // 调用方请求编号
			Target:    taskData.Target,                // 目标服务
			TaskType:  taskData.TaskType,              // 任务类型
			Status:    taskData.Status,                // 任务状态
			CreatedAt: taskData.CreatedAt.UnixMilli(), // 创建时间
			UpdatedAt: taskData.UpdatedAt.UnixMilli(), // 更新时间
		})
	}

	// 返回列表结果
	return &dispatchpb.ListTasksResponse{
		Total: result.Total, // 数据总数
		List:  list,         // 调度任务列表
	}, nil
}
