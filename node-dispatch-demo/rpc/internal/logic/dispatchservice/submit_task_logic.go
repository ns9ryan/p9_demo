package dispatchservicelogic

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"oa.98ent.com/p9/node-dispatch/rpc/internal/svc"
	"oa.98ent.com/p9/node-dispatch/rpc/internal/task"
	"oa.98ent.com/p9/node-dispatch/rpc/pb/nodedispatchrpc/dispatchpb"
)

type SubmitTaskLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSubmitTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SubmitTaskLogic {
	return &SubmitTaskLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// SubmitTask 提交调度任务
func (l *SubmitTaskLogic) SubmitTask(in *dispatchpb.SubmitTaskRequest) (*dispatchpb.SubmitTaskResponse, error) {
	// 整理请求参数
	requestNo := strings.TrimSpace(in.RequestNo)
	target := strings.TrimSpace(in.Target)
	taskType := strings.TrimSpace(in.TaskType)
	nodeCode := strings.TrimSpace(in.NodeCode)
	params := json.RawMessage(in.Params)

	// 校验请求参数
	if requestNo == "" {
		return nil, status.Error(codes.InvalidArgument, "request_no is required")
	}
	if target == "" {
		return nil, status.Error(codes.InvalidArgument, "target is required")
	}
	if taskType == "" {
		return nil, status.Error(codes.InvalidArgument, "task_type is required")
	}
	if nodeCode == "" {
		return nil, status.Error(codes.InvalidArgument, "node_code is required")
	}
	if !json.Valid(params) {
		return nil, status.Error(codes.InvalidArgument, "params must be valid json")
	}

	// 提交调度任务
	result, err := l.svcCtx.Task.Submit(l.ctx, task.SubmitRequest{
		RequestNo: requestNo, // 调用方请求编号
		Target:    target,    // 目标服务
		TaskType:  taskType,  // 任务类型
		NodeCode:  nodeCode,  // 执行节点编码
		Params:    params,    // 任务参数
	})
	if err != nil {
		l.Logger.Errorw(
			"提交调度任务失败",
			logx.Field("request_no", requestNo),
			logx.Field("target", target),
			logx.Field("task_type", taskType),
			logx.Field("node_code", nodeCode),
			logx.Field("error", err.Error()),
		)
		return nil, err
	}

	return &dispatchpb.SubmitTaskResponse{
		TaskNo: result.TaskNo, // 调度任务编号
	}, nil
}
