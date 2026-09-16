package nodeservicelogic

import (
	"context"

	"oa.98ent.com/p9/node-dispatch/rpc/internal/svc"
	"oa.98ent.com/p9/node-dispatch/rpc/pb/nodedispatchrpc/nodepb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLogic {
	return &UpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 修改节点
func (l *UpdateLogic) Update(in *nodepb.UpdateNodeRequest) (*nodepb.UpdateNodeResponse, error) {
	// todo: add your logic here and delete this line

	return &nodepb.UpdateNodeResponse{}, nil
}
