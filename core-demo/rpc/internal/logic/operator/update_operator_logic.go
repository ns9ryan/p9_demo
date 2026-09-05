package operator

import (
	"context"

	"oa.98ent.com/p9/core/common/ctxdata"
	"oa.98ent.com/p9/core/common/xerr"
	"oa.98ent.com/p9/core/rpc/internal/svc"
	"oa.98ent.com/p9/core/rpc/service"
	"oa.98ent.com/p9/core/rpc/types/core"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateOperatorLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateOperatorLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateOperatorLogic {
	return &UpdateOperatorLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateOperatorLogic) UpdateOperator(in *core.UpdateOperatorReq) (*core.Empty, error) {
	err := l.svcCtx.Deps.UpdateOperatorSelf(l.ctx, ctxdata.ClaimsFromCtx(l.ctx), service.UpdateOperatorReq{
		TimezoneCode: in.TimezoneCode, SettlementCurrencyCode: in.SettlementCurrencyCode,
	})
	if err != nil {
		return nil, xerr.RpcErr(err)
	}
	return &core.Empty{}, nil
}
