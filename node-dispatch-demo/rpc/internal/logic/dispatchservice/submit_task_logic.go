package dispatchservicelogic

import (
	"context"

	"oa.98ent.com/p9/node-dispatch/rpc/internal/svc"
	"oa.98ent.com/p9/node-dispatch/rpc/pb/nodedispatchrpc/dispatchpb"

	"github.com/zeromicro/go-zero/core/logx"
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

// 提交调度任务
func (l *SubmitTaskLogic) SubmitTask(in *dispatchpb.SubmitTaskRequest) (*dispatchpb.SubmitTaskResponse, error) {
	// todo: add your logic here and delete this line

	return &dispatchpb.SubmitTaskResponse{}, nil
}
