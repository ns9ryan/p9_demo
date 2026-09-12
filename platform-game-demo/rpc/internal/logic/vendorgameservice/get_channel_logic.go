package vendorgameservicelogic

import (
	"context"

	"oa.98ent.com/p9/platform-game/rpc/internal/svc"
	"oa.98ent.com/p9/platform-game/rpc/pb/vendors"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetChannelLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetChannelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetChannelLogic {
	return &GetChannelLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取游戏渠道
func (l *GetChannelLogic) GetChannel(in *vendors.Empty) (*vendors.GetChannelResponse, error) {
	// todo: add your logic here and delete this line

	return &vendors.GetChannelResponse{}, nil
}
