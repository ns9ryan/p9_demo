package public

import (
	"context"

	"oa.98ent.com/p9/core/api/internal/convert"
	"oa.98ent.com/p9/core/api/internal/svc"
	"oa.98ent.com/p9/core/api/internal/types"
	"oa.98ent.com/p9/core/rpc/coreclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type BootstrapOperatorLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBootstrapOperatorLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BootstrapOperatorLogic {
	return &BootstrapOperatorLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BootstrapOperatorLogic) BootstrapOperator(req *types.BootstrapOperatorReq) (resp *types.UserPublic, err error) {
	out, err := l.svcCtx.Core.BootstrapOperator(l.ctx, &coreclient.BootstrapOperatorReq{
		InitToken: req.InitToken, OperatorCode: req.OperatorCode, TimezoneCode: req.TimezoneCode,
		SettlementCurrencyCode: req.SettlementCurrencyCode, Username: req.Username, Password: req.Password,
		DisplayName: req.DisplayName,
	})
	if err != nil {
		return nil, err
	}
	return convert.UserPublic(out), nil
}
