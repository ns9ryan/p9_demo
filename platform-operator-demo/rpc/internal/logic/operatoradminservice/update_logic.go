package operatoradminservicelogic

import (
	"context"

	"oa.98ent.com/p9/platform-operator/pkg/i18nkey"
	"oa.98ent.com/p9/platform-operator/pkg/rpc/grpcerror"
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
	if in.Id <= 0 {
		return nil, grpcerror.InvalidArgument(i18nkey.ValidationError)
	}

	if in.Username == nil && in.Password == nil && in.DisplayName == nil && in.Status == nil {
		return nil, grpcerror.InvalidArgument(i18nkey.ValidationError)
	}

	username := trimOptionalString(in.Username)
	displayName := trimOptionalString(in.DisplayName)

	if username != nil && *username == "" {
		return nil, grpcerror.InvalidArgument(i18nkey.ValidationError)
	}

	if displayName != nil && *displayName == "" {
		return nil, grpcerror.InvalidArgument(i18nkey.ValidationError)
	}

	if in.Password != nil && (len(*in.Password) < 6 || len(*in.Password) > 32) {
		return nil, grpcerror.InvalidArgument(i18nkey.ValidationError)
	}

	if in.Status != nil && (*in.Status < 1 || *in.Status > 2) {
		return nil, grpcerror.InvalidArgument(i18nkey.ValidationError)
	}

	current, err := l.svcCtx.DB.OperatorAdmin.Get(l.ctx, in.Id)
	if err != nil {
		return nil, enterror.Handle(l.Logger, err)
	}

	err = current.
		Update().
		SetNillableUsername(username).
		SetNillablePassword(in.Password).
		SetNillableDisplayName(displayName).
		SetNillableStatus(in.Status).
		Exec(l.ctx)
	if err != nil {
		return nil, enterror.Handle(l.Logger, err)
	}

	return &adminpb.UpdateAdminResponse{}, nil
}
