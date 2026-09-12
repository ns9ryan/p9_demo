package currencyservicelogic

import (
	"context"

	"oa.98ent.com/p9/platform-base/rpc/internal/svc"
	"oa.98ent.com/p9/platform-base/rpc/pb/base/currency"

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

// 按编码获取货币
func (l *GetByCodeLogic) GetByCode(in *currency.GetCurrencyByCodeRequest) (*currency.GetCurrencyByCodeResponse, error) {
	// todo: add your logic here and delete this line

	return &currency.GetCurrencyByCodeResponse{}, nil
}
