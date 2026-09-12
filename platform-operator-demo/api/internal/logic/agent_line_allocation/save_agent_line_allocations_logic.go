// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package agent_line_allocation

import (
	"context"

	"oa.98ent.com/p9/platform-operator/api/internal/svc"
	"oa.98ent.com/p9/platform-operator/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SaveAgentLineAllocationsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSaveAgentLineAllocationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveAgentLineAllocationsLogic {
	return &SaveAgentLineAllocationsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SaveAgentLineAllocationsLogic) SaveAgentLineAllocations(req *types.SaveAgentLineAllocationsRequest) (resp *types.SaveAgentLineAllocationsResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
