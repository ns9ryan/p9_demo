package operatoradminservicelogic

import (
	"context"
	"strings"

	"oa.98ent.com/p9/platform-operator/pkg/i18nkey"
	"oa.98ent.com/p9/platform-operator/pkg/rpc/grpcerror"
	"oa.98ent.com/p9/platform-operator/rpc/internal/enterror"
	"oa.98ent.com/p9/platform-operator/rpc/internal/svc"
	"oa.98ent.com/p9/platform-operator/rpc/pb/platformoperatorrpc/adminpb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
	return &CreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// Create 创建分站管理员
func (l *CreateLogic) Create(in *adminpb.CreateAdminRequest) (*adminpb.CreateAdminResponse, error) {
	if in.OperatorId <= 0 {
		return nil, grpcerror.InvalidArgument(i18nkey.ValidationError)
	}

	// 确认分站存在
	_, err := l.svcCtx.DB.Operator.Get(l.ctx, in.OperatorId)
	if err != nil {
		return nil, enterror.Handle(l.Logger, err)
	}

	data, err := l.svcCtx.DB.OperatorAdmin.
		Create().
		SetOperatorID(in.OperatorId).
		SetUsername(strings.TrimSpace(in.Username)).
		SetPassword(in.Password).
		SetDisplayName(strings.TrimSpace(in.DisplayName)).
		SetNillableStatus(in.Status).
		Save(l.ctx)
	if err != nil {
		return nil, enterror.Handle(l.Logger, err)
	}

	return &adminpb.CreateAdminResponse{
		Id: data.ID, // 管理员ID
	}, nil
}
