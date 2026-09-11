package vendorgameservicelogic

import (
	"context"

	"oa.98ent.com/p9/platform-game/rpc/internal/svc"
	"oa.98ent.com/p9/platform-game/rpc/pb/vendors"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetGameLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetGameLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGameLogic {
	return &GetGameLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取游戏
func (l *GetGameLogic) GetGame(in *vendors.GetGameRequest) (*vendors.GetGameResponse, error) {
	// todo: add your logic here and delete this line

	return &vendors.GetGameResponse{}, nil
}
