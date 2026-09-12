// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package operator_admin

import (
	"context"

	"oa.98ent.com/p9/platform-operator/api/internal/svc"
	"oa.98ent.com/p9/platform-operator/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateOperatorAdminLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateOperatorAdminLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateOperatorAdminLogic {
	return &CreateOperatorAdminLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateOperatorAdminLogic) CreateOperatorAdmin(req *types.CreateOperatorAdminRequest) (resp *types.CreateOperatorAdminResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
