// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package operator

import (
	"context"

	"oa.98ent.com/p9/platform-operator/api/internal/svc"
	"oa.98ent.com/p9/platform-operator/api/internal/types"

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

func (l *GetOperatorLogic) GetOperator(req *types.GetOperatorRequest) (resp *types.GetOperatorResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
