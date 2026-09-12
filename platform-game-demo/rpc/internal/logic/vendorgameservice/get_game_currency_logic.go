package vendorgameservicelogic

import (
	"context"

	"oa.98ent.com/p9/platform-game/rpc/internal/svc"
	"oa.98ent.com/p9/platform-game/rpc/pb/vendors"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetGameCurrencyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetGameCurrencyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGameCurrencyLogic {
	return &GetGameCurrencyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取游戏货币
func (l *GetGameCurrencyLogic) GetGameCurrency(in *vendors.Empty) (*vendors.GetGameCurrencyResponse, error) {
	// todo: add your logic here and delete this line

	return &vendors.GetGameCurrencyResponse{}, nil
}
