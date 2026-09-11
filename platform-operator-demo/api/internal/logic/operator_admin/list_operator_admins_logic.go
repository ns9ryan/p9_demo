// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package operator_admin

import (
	"context"

	"oa.98ent.com/p9/platform-operator/api/internal/svc"
	"oa.98ent.com/p9/platform-operator/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListOperatorAdminsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListOperatorAdminsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListOperatorAdminsLogic {
	return &ListOperatorAdminsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListOperatorAdminsLogic) ListOperatorAdmins(req *types.ListOperatorAdminsRequest) (resp *types.ListOperatorAdminsResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
