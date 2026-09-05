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

type ReorderCurrencyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewReorderCurrencyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReorderCurrencyLogic {
	return &ReorderCurrencyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// ReorderCurrency 调整货币排序
func (l *ReorderCurrencyLogic) ReorderCurrency(req *types.ReorderCurrencyRequest) (resp *types.ReorderCurrencyResponse, err error) {
	// 调用调整货币排序RPC
	_, err = l.svcCtx.CurrencyRpc.Reorder(
		l.ctx,
		&currency.ReorderCurrencyRequest{
			Id:       req.Id,       // 需要移动的货币ID
			TargetId: req.TargetId, // 目标货币ID
		},
	)
	if err != nil {
		return nil, err
	}

	// 返回调整排序结果
	return &types.ReorderCurrencyResponse{}, nil
}
