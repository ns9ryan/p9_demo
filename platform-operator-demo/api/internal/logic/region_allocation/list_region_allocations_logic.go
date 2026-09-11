// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package region_allocation

import (
	"context"

	"oa.98ent.com/p9/platform-operator/api/internal/svc"
	"oa.98ent.com/p9/platform-operator/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListRegionAllocationsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListRegionAllocationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListRegionAllocationsLogic {
	return &ListRegionAllocationsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListRegionAllocationsLogic) ListRegionAllocations(req *types.ListRegionAllocationsRequest) (resp *types.ListRegionAllocationsResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
