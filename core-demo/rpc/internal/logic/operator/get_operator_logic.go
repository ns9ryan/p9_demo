package operator

import (
	"context"

	"oa.98ent.com/p9/core/common/ctxdata"
	"oa.98ent.com/p9/core/common/xerr"
	"oa.98ent.com/p9/core/rpc/internal/logic"
	"oa.98ent.com/p9/core/rpc/internal/svc"
	"oa.98ent.com/p9/core/rpc/types/core"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetOperatorLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetOperatorLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOperatorLogic {
	return &GetOperatorLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetOperatorLogic) GetOperator(in *core.Empty) (*core.OperatorInfo, error) {
	op, err := l.svcCtx.Deps.GetOperatorSelf(l.ctx, ctxdata.ClaimsFromCtx(l.ctx))
	if err != nil {
		return nil, xerr.RpcErr(err)
	}
	return logic.ToOperatorInfo(op), nil
}
