// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package operator_profile

import (
	"context"

	"oa.98ent.com/p9/platform-operator/api/internal/svc"
	"oa.98ent.com/p9/platform-operator/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateOperatorProfileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateOperatorProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateOperatorProfileLogic {
	return &UpdateOperatorProfileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateOperatorProfileLogic) UpdateOperatorProfile(req *types.UpdateOperatorProfileRequest) (resp *types.UpdateOperatorProfileResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
