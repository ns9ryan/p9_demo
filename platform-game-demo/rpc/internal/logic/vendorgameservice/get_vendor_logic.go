package vendorgameservicelogic

import (
	"context"

	"oa.98ent.com/p9/platform-game/rpc/internal/svc"
	"oa.98ent.com/p9/platform-game/rpc/pb/vendors"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetVendorLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetVendorLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetVendorLogic {
	return &GetVendorLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取游戏提供商
func (l *GetVendorLogic) GetVendor(in *vendors.Empty) (*vendors.GetVendorResponse, error) {
	// todo: add your logic here and delete this line

	return &vendors.GetVendorResponse{}, nil
}
