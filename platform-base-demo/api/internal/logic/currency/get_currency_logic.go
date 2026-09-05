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

type GetCurrencyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetCurrencyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCurrencyLogic {
	return &GetCurrencyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// GetCurrency 获取货币
func (l *GetCurrencyLogic) GetCurrency(req *types.GetCurrencyRequest) (resp *types.GetCurrencyResponse, err error) {
	// 调用获取货币RPC
	result, err := l.svcCtx.CurrencyRpc.Get(
		l.ctx,
		&currency.GetCurrencyRequest{
			Id: req.Id, // 货币ID
		},
	)
	if err != nil {
		return nil, err
	}

	// 返回货币信息
	return &types.GetCurrencyResponse{
		Currency: types.CurrencyInfo{
			Id:           result.Currency.Id,           // 货币ID
			Code:         result.Currency.Code,         // 货币编码
			NameI18n:     result.Currency.NameI18N,     // 多语言名称
			CurrencyType: result.Currency.CurrencyType, // 货币类型: 1法定货币, 2虚拟货币
			Symbol:       result.Currency.Symbol,       // 货币符号
			AmountFactor: result.Currency.AmountFactor, // 金额换算倍率
			Status:       result.Currency.Status,       // 状态: 1启用, 2停用
			SortNo:       result.Currency.SortNo,       // 排序值
		},
	}, nil
}
