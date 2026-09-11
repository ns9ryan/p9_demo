// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package operator_profile

import (
	"context"

	"oa.98ent.com/p9/platform-operator/api/internal/svc"
	"oa.98ent.com/p9/platform-operator/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetOperatorProfileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetOperatorProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOperatorProfileLogic {
	return &GetOperatorProfileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetOperatorProfileLogic) GetOperatorProfile(req *types.GetOperatorProfileRequest) (resp *types.GetOperatorProfileResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
