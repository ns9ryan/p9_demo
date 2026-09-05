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

type ListAllCurrenciesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListAllCurrenciesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListAllCurrenciesLogic {
	return &ListAllCurrenciesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// ListAllCurrencies 获取全部货币
func (l *ListAllCurrenciesLogic) ListAllCurrencies(req *types.ListAllCurrenciesRequest) (resp *types.ListAllCurrenciesResponse, err error) {
	// 调用获取全部货币RPC
	result, err := l.svcCtx.CurrencyRpc.ListAll(
		l.ctx,
		&currency.ListAllCurrenciesRequest{
			Status: req.Status, // 状态: 1启用, 2停用
		},
	)
	if err != nil {
		return nil, err
	}

	// 转换货币列表
	list := make([]types.CurrencyInfo, 0, len(result.List))
	for _, item := range result.List {
		list = append(list, types.CurrencyInfo{
			Id:           item.Id,           // 货币ID
			Code:         item.Code,         // 货币编码
			NameI18n:     item.NameI18N,     // 多语言名称
			CurrencyType: item.CurrencyType, // 货币类型: 1法定货币, 2虚拟货币
			Symbol:       item.Symbol,       // 货币符号
			AmountFactor: item.AmountFactor, // 金额换算倍率
			Status:       item.Status,       // 状态: 1启用, 2停用
			SortNo:       item.SortNo,       // 排序值
		})
	}

	// 返回全部货币
	return &types.ListAllCurrenciesResponse{
		List: list, // 货币列表
	}, nil
}
