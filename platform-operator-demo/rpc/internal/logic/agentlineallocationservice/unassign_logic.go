package agentlineallocationservicelogic

import (
	"context"

	"oa.98ent.com/p9/platform-operator/rpc/internal/svc"
	"oa.98ent.com/p9/platform-operator/rpc/pb/operator/agentlineallocation"

	"github.com/zeromicro/go-zero/core/logx"
)

type UnassignLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUnassignLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnassignLogic {
	return &UnassignLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 取消代理子线路分配
func (l *UnassignLogic) Unassign(in *agentlineallocation.UnassignAgentLineRequest) (*agentlineallocation.UnassignAgentLineResponse, error) {
	// todo: add your logic here and delete this line

	return &agentlineallocation.UnassignAgentLineResponse{}, nil
}
