// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package basic_resource_allocation

import (
	"context"

	"oa.98ent.com/p9/platform-operator/api/internal/svc"
	"oa.98ent.com/p9/platform-operator/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListBasicResourceAllocationsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListBasicResourceAllocationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListBasicResourceAllocationsLogic {
	return &ListBasicResourceAllocationsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListBasicResourceAllocationsLogic) ListBasicResourceAllocations(req *types.ListBasicResourceAllocationsRequest) (resp *types.ListBasicResourceAllocationsResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
