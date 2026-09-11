package regionallocationservicelogic

import (
	"context"

	"oa.98ent.com/p9/platform-operator/rpc/internal/svc"
	"oa.98ent.com/p9/platform-operator/rpc/pb/operator/regionallocation"

	"github.com/zeromicro/go-zero/core/logx"
)

type SaveLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSaveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveLogic {
	return &SaveLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 保存经营地区分配
func (l *SaveLogic) Save(in *regionallocation.SaveRegionAllocationsRequest) (*regionallocation.SaveRegionAllocationsResponse, error) {
	// todo: add your logic here and delete this line

	return &regionallocation.SaveRegionAllocationsResponse{}, nil
}
