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

type CreateCurrencyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateCurrencyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateCurrencyLogic {
	return &CreateCurrencyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// CreateCurrency 创建货币
func (l *CreateCurrencyLogic) CreateCurrency(req *types.CreateCurrencyRequest) (resp *types.CreateCurrencyResponse, err error) {
	// 调用创建货币RPC
	result, err := l.svcCtx.CurrencyRpc.Create(
		l.ctx,
		&currency.CreateCurrencyRequest{
			Code:         req.Code,         // 货币编码
			NameI18N:     req.NameI18n,     // 多语言名称
			CurrencyType: req.CurrencyType, // 货币类型: 1法定货币, 2虚拟货币
			Symbol:       req.Symbol,       // 货币符号
			AmountFactor: req.AmountFactor, // 金额换算倍率
			Status:       req.Status,       // 状态: 1启用, 2停用
		},
	)
	if err != nil {
		return nil, err
	}

	// 返回创建结果
	return &types.CreateCurrencyResponse{
		Id: result.Id, // 货币ID
	}, nil
}
