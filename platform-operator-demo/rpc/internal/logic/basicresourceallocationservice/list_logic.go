package basicresourceallocationservicelogic

import (
	"context"

	"oa.98ent.com/p9/platform-operator/rpc/internal/svc"
	"oa.98ent.com/p9/platform-operator/rpc/pb/operator/basicresourceallocation"

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

// 获取基础资源分配列表
func (l *ListLogic) List(in *basicresourceallocation.ListBasicResourceAllocationsRequest) (*basicresourceallocation.ListBasicResourceAllocationsResponse, error) {
	// todo: add your logic here and delete this line

	return &basicresourceallocation.ListBasicResourceAllocationsResponse{}, nil
}
