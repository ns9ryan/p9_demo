package operator

import (
	"context"

	"oa.98ent.com/p9/core/api/internal/convert"
	"oa.98ent.com/p9/core/api/internal/svc"
	"oa.98ent.com/p9/core/api/internal/types"
	"oa.98ent.com/p9/core/rpc/coreclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetOperatorLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetOperatorLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOperatorLogic {
	return &GetOperatorLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetOperatorLogic) GetOperator() (resp *types.OperatorInfo, err error) {
	out, err := l.svcCtx.Core.GetOperator(l.ctx, &coreclient.Empty{})
	if err != nil {
		return nil, err
	}
	return convert.OperatorInfo(out), nil
}
