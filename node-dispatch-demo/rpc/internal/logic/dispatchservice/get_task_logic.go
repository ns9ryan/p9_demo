package dispatchservicelogic

import (
	"context"
	"fmt"
	"strings"

	"oa.98ent.com/p9/node-dispatch/rpc/ent"
	"oa.98ent.com/p9/node-dispatch/rpc/ent/dispatchtask"
	"oa.98ent.com/p9/node-dispatch/rpc/ent/dispatchtaskrun"
	"oa.98ent.com/p9/node-dispatch/rpc/internal/svc"
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
	// ---------- 参数整理 ----------

	// 整理查询编号
	taskNo := strings.TrimSpace(in.GetTaskNo())
	requestNo := strings.TrimSpace(in.GetRequestNo())

	// 任务编号和请求编号必须且只能提供一个
	if taskNo == "" && requestNo == "" {
		return nil, status.Error(codes.InvalidArgument, "task_no or request_no is required")
	}
	if taskNo != "" && requestNo != "" {
		return nil, status.Error(codes.InvalidArgument, "task_no and request_no cannot be provided together")
	}

	// ---------- 任务查询 ----------

	// 创建任务查询
	query := l.svcCtx.DB.DispatchTask.Query()

	// 根据指定编号查询任务
	if taskNo != "" {
		query = query.Where(dispatchtask.TaskNoEQ(taskNo))
	} else {
		query = query.Where(dispatchtask.RequestNoEQ(requestNo))
	}

	// 获取调度任务
	taskData, err := query.Only(l.ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, status.Error(codes.NotFound, "task not found")
		}

		l.Logger.Errorw(
			"查询调度任务失败",
			logx.Field("task_no", taskNo),
			logx.Field("request_no", requestNo),
			logx.Field("error", err.Error()),
		)
		return nil, err
	}

	// ---------- 执行记录查询 ----------

	// 获取任务全部执行记录及执行节点
	runList, err := l.svcCtx.DB.DispatchTaskRun.
		Query().
		Where(dispatchtaskrun.TaskIDEQ(taskData.ID)).
		WithNode().
		Order(ent.Asc(dispatchtaskrun.FieldRunNo)).
		All(l.ctx)
	if err != nil {
		l.Logger.Errorw(
			"查询调度任务执行记录失败",
			logx.Field("task_no", taskData.TaskNo),
			logx.Field("error", err.Error()),
		)
		return nil, err
	}

	// 调度任务创建后必须至少存在一条执行记录
	if len(runList) == 0 {
		err = fmt.Errorf("调度任务执行记录不存在: %s", taskData.TaskNo)
		l.Logger.Errorw(
			"调度任务执行记录不存在",
			logx.Field("task_no", taskData.TaskNo),
			logx.Field("error", err.Error()),
		)
		return nil, err
	}

	// ---------- 响应转换 ----------

	// 转换任务执行记录
	runs := make([]*dispatchpb.TaskRunInfo, 0, len(runList))
	var nodeCode string

	for _, runData := range runList {
		// 获取执行节点
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
			runInfo.FinishedAt = new(runData.FinishedAt.UnixMilli())
		}

		runs = append(runs, runInfo)
	}

	// 返回任务详情
	return &dispatchpb.GetTaskResponse{
		Task: &dispatchpb.TaskInfo{
			TaskNo:    taskData.TaskNo,                // 调度任务编号
			RequestNo: taskData.RequestNo,             // 调用方请求编号
			Target:    taskData.Target,                // 目标服务
			TaskType:  taskData.TaskType,              // 任务类型
			Status:    taskData.Status,                // 任务状态
			NodeCode:  nodeCode,                       // 最近执行节点编码
			CreatedAt: taskData.CreatedAt.UnixMilli(), // 创建时间
			UpdatedAt: taskData.UpdatedAt.UnixMilli(), // 更新时间
		},
		Params: string(taskData.Params), // 任务参数
		Runs:   runs,                    // 执行记录
	}, nil
}
