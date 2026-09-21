package operatoradminservicelogic

import (
	"context"

	"oa.98ent.com/p9/platform-operator/rpc/internal/enterror"
	"oa.98ent.com/p9/platform-operator/rpc/internal/svc"
	"oa.98ent.com/p9/platform-operator/rpc/pb/platformoperatorrpc/adminpb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLogic {
	return &UpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// Update 更新分站管理员
func (l *UpdateLogic) Update(in *adminpb.UpdateAdminRequest) (*adminpb.UpdateAdminResponse, error) {
	current, err := l.svcCtx.DB.OperatorAdmin.Get(l.ctx, in.Id)
	if err != nil {
		return nil, enterror.Handle(l.Logger, err)
	}

	err = current.
		Update().
		SetNillableUsername(trimOptionalString(in.Username)).
		SetNillablePassword(trimOptionalString(in.Password)).
		SetNillableDisplayName(trimOptionalString(in.DisplayName)).
		SetNillableStatus(in.Status).
		Exec(l.ctx)
	if err != nil {
		return nil, enterror.Handle(l.Logger, err)
	}

	return &adminpb.UpdateAdminResponse{}, nil
}
