package operatoradminservicelogic

import (
	"context"

	"oa.98ent.com/p9/platform-operator/rpc/internal/enterror"
	"oa.98ent.com/p9/platform-operator/rpc/internal/svc"
	"oa.98ent.com/p9/platform-operator/rpc/pb/platformoperatorrpc/adminpb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateStatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateStatusLogic {
	return &UpdateStatusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// UpdateStatus 更新分站管理员状态
func (l *UpdateStatusLogic) UpdateStatus(in *adminpb.UpdateAdminStatusRequest) (*adminpb.UpdateAdminStatusResponse, error) {
	current, err := l.svcCtx.DB.OperatorAdmin.Get(l.ctx, in.Id)
	if err != nil {
		return nil, enterror.Handle(l.Logger, err)
	}

	err = current.
		Update().
		SetStatus(in.Status).
		Exec(l.ctx)
	if err != nil {
		return nil, enterror.Handle(l.Logger, err)
	}

	return &adminpb.UpdateAdminStatusResponse{}, nil
}
