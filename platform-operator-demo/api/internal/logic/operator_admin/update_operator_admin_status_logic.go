// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package operator_admin

import (
	"context"

	"oa.98ent.com/p9/common/xerr"
	apiI18nkey "oa.98ent.com/p9/platform-operator/api/internal/i18nkey"
	"oa.98ent.com/p9/platform-operator/api/internal/svc"
	"oa.98ent.com/p9/platform-operator/api/internal/types"
	"oa.98ent.com/p9/platform-operator/pkg/i18nkey"
	"oa.98ent.com/p9/platform-operator/rpc/pb/platformoperatorrpc/adminpb"
	"oa.98ent.com/p9/platform-operator/rpc/pb/platformoperatorrpc/operatorpb"

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
	// 获取管理员信息
	admin, err := l.svcCtx.OperatorAdminRpc.Get(l.ctx, &adminpb.GetAdminRequest{
		Id: req.Id,
	})
	if err != nil {
		return nil, err
	}
	if admin == nil || admin.Admin == nil {
		return nil, xerr.NotFound(i18nkey.DataNotFound)
	}

	// 判断分站状态，已发布则不能更新
	operator, err := l.svcCtx.OperatorRpc.Get(l.ctx, &operatorpb.GetOperatorRequest{
		Id: admin.Admin.OperatorId,
	})
	if err != nil {
		return nil, err
	}
	if operator == nil || operator.Operator == nil {
		return nil, xerr.NotFound(i18nkey.DataNotFound)
	}
	if operator.Operator.PublishStatus != 1 {
		return nil, xerr.BadRequest(apiI18nkey.ForbiddenUpdOpHasPublished)
	}

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
