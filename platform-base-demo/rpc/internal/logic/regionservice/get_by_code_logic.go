package regionservicelogic

import (
	"context"

	"oa.98ent.com/p9/platform-base/rpc/internal/svc"
	"oa.98ent.com/p9/platform-base/rpc/pb/base/region"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetByCodeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetByCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetByCodeLogic {
	return &GetByCodeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 按编码获取国家地区
func (l *GetByCodeLogic) GetByCode(in *region.GetRegionByCodeRequest) (*region.GetRegionByCodeResponse, error) {
	// todo: add your logic here and delete this line

	return &region.GetRegionByCodeResponse{}, nil
}
