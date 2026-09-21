// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package operator_admin

import (
	"context"

	"oa.98ent.com/p9/platform-operator/api/internal/svc"
	"oa.98ent.com/p9/platform-operator/api/internal/types"
	"oa.98ent.com/p9/platform-operator/rpc/pb/platformoperatorrpc/adminpb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ResetOperatorAdminPasswordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewResetOperatorAdminPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResetOperatorAdminPasswordLogic {
	return &ResetOperatorAdminPasswordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// ResetOperatorAdminPassword 重置分站管理员密码
func (l *ResetOperatorAdminPasswordLogic) ResetOperatorAdminPassword(req *types.ResetOperatorAdminPasswordRequest) (resp *types.ResetOperatorAdminPasswordResponse, err error) {
	// 重置密码
	_, err = l.svcCtx.OperatorAdminRpc.ResetPassword(
		l.ctx,
		&adminpb.ResetAdminPasswordRequest{
			Id:       req.Id,       // 管理员ID
			Password: req.Password, // 新密码
		},
	)
	if err != nil {
		return nil, err
	}

	return &types.ResetOperatorAdminPasswordResponse{}, nil
}
