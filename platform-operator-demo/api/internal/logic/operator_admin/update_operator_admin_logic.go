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

type UpdateOperatorAdminLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateOperatorAdminLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateOperatorAdminLogic {
	return &UpdateOperatorAdminLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// UpdateOperatorAdmin 更新分站管理员
func (l *UpdateOperatorAdminLogic) UpdateOperatorAdmin(req *types.UpdateOperatorAdminRequest) (resp *types.UpdateOperatorAdminResponse, err error) {
	_, err = l.svcCtx.OperatorAdminRpc.Update(
		l.ctx,
		&adminpb.UpdateAdminRequest{
			Id:          req.Id,          // 管理员ID
			Username:    req.Username,    // 账号
			Password:    req.Password,    // 密码
			DisplayName: req.DisplayName, // 显示名称
			Status:      req.Status,      // 状态: 1启用, 2停用
		},
	)
	if err != nil {
		return nil, err
	}

	return &types.UpdateOperatorAdminResponse{}, nil
}
