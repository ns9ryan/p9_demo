package operatoradminservicelogic

import (
	"context"

	"oa.98ent.com/p9/platform-operator/rpc/internal/enterror"
	"oa.98ent.com/p9/platform-operator/rpc/internal/svc"
	"oa.98ent.com/p9/platform-operator/rpc/pb/platformoperatorrpc/adminpb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ResetPasswordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewResetPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResetPasswordLogic {
	return &ResetPasswordLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ResetPassword 重置分站管理员密码
func (l *ResetPasswordLogic) ResetPassword(in *adminpb.ResetAdminPasswordRequest) (*adminpb.ResetAdminPasswordResponse, error) {
	current, err := l.svcCtx.DB.OperatorAdmin.Get(l.ctx, in.Id)
	if err != nil {
		return nil, enterror.Handle(l.Logger, err)
	}

	err = current.
		Update().
		SetPassword(in.Password).
		Exec(l.ctx)
	if err != nil {
		return nil, enterror.Handle(l.Logger, err)
	}

	return &adminpb.ResetAdminPasswordResponse{}, nil
}
