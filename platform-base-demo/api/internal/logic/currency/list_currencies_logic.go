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

type ListCurrenciesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListCurrenciesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListCurrenciesLogic {
	return &ListCurrenciesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// ListCurrencies 获取货币管理列表
func (l *ListCurrenciesLogic) ListCurrencies(req *types.ListCurrenciesRequest) (resp *types.ListCurrenciesResponse, err error) {
	// 调用获取货币管理列表RPC
	result, err := l.svcCtx.CurrencyRpc.List(
		l.ctx,
		&currency.ListCurrenciesRequest{
			Page:     req.Page,     // 页码
			PageSize: req.PageSize, // 每页数量
			Status:   req.Status,   // 状态: 1启用, 2停用
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

	// 返回货币管理列表
	return &types.ListCurrenciesResponse{
		Total: result.Total, // 数据总数
		List:  list,         // 货币列表
	}, nil
}
