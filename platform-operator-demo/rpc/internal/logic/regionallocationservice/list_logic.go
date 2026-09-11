package regionallocationservicelogic

import (
	"context"

	"oa.98ent.com/p9/platform-operator/rpc/internal/svc"
	"oa.98ent.com/p9/platform-operator/rpc/pb/operator/regionallocation"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListLogic {
	return &ListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取经营地区分配列表
func (l *ListLogic) List(in *regionallocation.ListRegionAllocationsRequest) (*regionallocation.ListRegionAllocationsResponse, error) {
	// todo: add your logic here and delete this line

	return &regionallocation.ListRegionAllocationsResponse{}, nil
}
