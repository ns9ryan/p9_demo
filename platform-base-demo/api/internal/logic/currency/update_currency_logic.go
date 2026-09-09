// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package currency

import (
	"context"

	"oa.98ent.com/p9/platform-base/api/internal/svc"
	"oa.98ent.com/p9/platform-base/api/internal/types"
	"oa.98ent.com/p9/platform-base/rpc/pb/base/currency"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateCurrencyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateCurrencyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateCurrencyLogic {
	return &UpdateCurrencyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// UpdateCurrency 修改货币
func (l *UpdateCurrencyLogic) UpdateCurrency(req *types.UpdateCurrencyRequest) (resp *types.UpdateCurrencyResponse, err error) {
	// 调用修改货币RPC
	_, err = l.svcCtx.CurrencyRpc.Update(
		l.ctx,
		&currency.UpdateCurrencyRequest{
			Id:     req.Id,     // 货币ID
			Symbol: req.Symbol, // 货币符号
			Status: req.Status, // 状态: 1启用, 2停用
		},
	)
	if err != nil {
		return nil, err
	}

	// 返回修改结果
	return &types.UpdateCurrencyResponse{}, nil
}
