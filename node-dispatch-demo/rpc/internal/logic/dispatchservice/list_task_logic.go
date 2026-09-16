package dispatchservicelogic

import (
	"context"

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

// 获取调度任务列表
func (l *ListTaskLogic) ListTask(in *dispatchpb.ListTasksRequest) (*dispatchpb.ListTasksResponse, error) {
	// todo: add your logic here and delete this line

	return &dispatchpb.ListTasksResponse{}, nil
}
