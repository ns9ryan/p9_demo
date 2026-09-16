package dispatchcallbackservicelogic

import (
	"context"

	"oa.98ent.com/p9/node-dispatch/rpc/internal/svc"
	"oa.98ent.com/p9/node-dispatch/rpc/pb/nodedispatchrpc/callbackpb"

	"github.com/zeromicro/go-zero/core/logx"
)

type TaskResultLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTaskResultLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TaskResultLogic {
	return &TaskResultLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 回调任务结果
func (l *TaskResultLogic) TaskResult(in *callbackpb.TaskResultRequest) (*callbackpb.TaskResultResponse, error) {
	// todo: add your logic here and delete this line

	return &callbackpb.TaskResultResponse{}, nil
}
