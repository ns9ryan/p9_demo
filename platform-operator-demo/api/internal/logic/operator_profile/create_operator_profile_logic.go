// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package operator_profile

import (
	"context"

	"oa.98ent.com/p9/platform-operator/api/internal/svc"
	"oa.98ent.com/p9/platform-operator/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateOperatorProfileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateOperatorProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateOperatorProfileLogic {
	return &CreateOperatorProfileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateOperatorProfileLogic) CreateOperatorProfile(req *types.CreateOperatorProfileRequest) (resp *types.CreateOperatorProfileResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
