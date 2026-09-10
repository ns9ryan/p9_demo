package agentlineallocationservicelogic

import (
	"context"

	"oa.98ent.com/p9/platform-operator/rpc/internal/svc"
	"oa.98ent.com/p9/platform-operator/rpc/pb/operator/agentlineallocation"

	"github.com/zeromicro/go-zero/core/logx"
)

type AssignLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAssignLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AssignLogic {
	return &AssignLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 分配代理子线路
func (l *AssignLogic) Assign(in *agentlineallocation.AssignAgentLineRequest) (*agentlineallocation.AssignAgentLineResponse, error) {
	// todo: add your logic here and delete this line

	return &agentlineallocation.AssignAgentLineResponse{}, nil
}
