package dispatchservicelogic

import (
	"context"
	"strings"

	"oa.98ent.com/p9/node-dispatch/rpc/internal/svc"
	"oa.98ent.com/p9/node-dispatch/rpc/internal/task"
	"oa.98ent.com/p9/node-dispatch/rpc/pb/nodedispatchrpc/dispatchpb"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GetTaskLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTaskLogic {
	return &GetTaskLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetTask 获取调度任务
func (l *GetTaskLogic) GetTask(in *dispatchpb.GetTaskRequest) (*dispatchpb.GetTaskResponse, error) {
	// 整理查询参数
	taskNo := strings.TrimSpace(in.GetTaskNo())
	requestNo := strings.TrimSpace(in.GetRequestNo())

	// 校验查询参数
	if taskNo == "" && requestNo == "" {
		return nil, status.Error(codes.InvalidArgument, "task_no or request_no is required")
	}
	if taskNo != "" && requestNo != "" {
		return nil, status.Error(codes.InvalidArgument, "task_no and request_no cannot be provided together")
	}

	// 查询调度任务
	taskResult, err := l.svcCtx.Task.Get(l.ctx, task.GetRequest{
		TaskNo:    taskNo,    // 调度任务编号
		RequestNo: requestNo, // 调用方请求编号
	})
	if err != nil {
		l.Logger.Errorw(
			"获取调度任务失败",
			logx.Field("task_no", taskNo),
			logx.Field("request_no", requestNo),
			logx.Field("error", err.Error()),
		)
		return nil, err
	}

	// 校验查询结果
	if taskResult == nil || taskResult.Task == nil {
		l.Logger.Errorw(
			"获取调度任务结果为空",
			logx.Field("task_no", taskNo),
			logx.Field("request_no", requestNo),
		)
		return nil, status.Error(codes.Internal, "task result is empty")
	}

	taskData := taskResult.Task

	// 转换执行记录
	runs := make([]*dispatchpb.TaskRunInfo, 0, len(taskResult.Runs))
	var nodeCode string

	for _, runData := range taskResult.Runs {
		// 获取执行节点
		nodeData, err := runData.Edges.NodeOrErr()
		if err != nil {
			l.Logger.Errorw(
				"获取任务执行节点失败",
				logx.Field("task_no", taskData.TaskNo),
				logx.Field("run_no", runData.RunNo),
				logx.Field("error", err.Error()),
			)
			return nil, err
		}

		// 保存执行节点编码
		nodeCode = nodeData.Code

		// 最近一条执行记录的节点作为任务当前执行节点
		nodeCode = nodeData.Code

		// 创建执行记录响应
		runInfo := &dispatchpb.TaskRunInfo{
			RunNo:    runData.RunNo,  // 执行序号
			NodeCode: nodeData.Code,  // 执行节点编码
			Status:   runData.Status, // 执行状态
		}

		// 设置执行结果
		if len(runData.Result) > 0 {
			runInfo.Result = new(string(runData.Result))
		}

		// 设置失败原因
		if runData.ErrorMessage != nil {
			runInfo.ErrorMessage = new(*runData.ErrorMessage)
		}

		// 设置开始执行时间
		if runData.StartedAt != nil {
			runInfo.StartedAt = new(runData.StartedAt.UnixMilli())
		}

		// 设置执行结束时间
		if runData.FinishedAt != nil {
			runInfo.StartedAt = new(runData.StartedAt.UnixMilli())
		}

		runs = append(runs, runInfo)
	}

	// 返回调度任务信息
	return &dispatchpb.GetTaskResponse{
		Task: &dispatchpb.TaskInfo{
			TaskNo:    taskData.TaskNo,                // 调度任务编号
			RequestNo: taskData.RequestNo,             // 调用方请求编号
			Target:    taskData.Target,                // 目标服务
			TaskType:  taskData.TaskType,              // 任务类型
			Status:    taskData.Status,                // 任务状态
			NodeCode:  nodeCode,                       // 执行节点编码
			CreatedAt: taskData.CreatedAt.UnixMilli(), // 创建时间
			UpdatedAt: taskData.UpdatedAt.UnixMilli(), // 更新时间
		},
		Params: string(taskData.Params), // 任务参数
		Runs:   runs,                    // 执行记录
	}, nil
}
