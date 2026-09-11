// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package game_allocation

import (
	"context"

	"oa.98ent.com/p9/platform-operator/api/internal/svc"
	"oa.98ent.com/p9/platform-operator/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListGameAllocationsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListGameAllocationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListGameAllocationsLogic {
	return &ListGameAllocationsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListGameAllocationsLogic) ListGameAllocations(req *types.ListGameAllocationsRequest) (resp *types.ListGameAllocationsResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
