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

type UpdateOperatorAdminStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateOperatorAdminStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateOperatorAdminStatusLogic {
	return &UpdateOperatorAdminStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// UpdateOperatorAdminStatus 更新分站管理员状态
func (l *UpdateOperatorAdminStatusLogic) UpdateOperatorAdminStatus(req *types.UpdateOperatorAdminStatusRequest) (resp *types.UpdateOperatorAdminStatusResponse, err error) {
	_, err = l.svcCtx.OperatorAdminRpc.UpdateStatus(
		l.ctx,
		&adminpb.UpdateAdminStatusRequest{
			Id:     req.Id,     // 管理员ID
			Status: req.Status, // 状态: 1启用, 2停用
		},
	)
	if err != nil {
		return nil, err
	}

	return &types.UpdateOperatorAdminStatusResponse{}, nil
}
