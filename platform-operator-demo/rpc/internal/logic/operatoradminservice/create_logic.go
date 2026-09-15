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
	// 分站ID必须大于0
	if in.OperatorId <= 0 {
		return nil, grpcerror.InvalidArgument(i18nkey.ValidationError)
	}

	username := strings.TrimSpace(in.Username)
	displayName := strings.TrimSpace(in.DisplayName)

	// 账号和显示名称不能为空
	if username == "" || displayName == "" {
		return nil, grpcerror.InvalidArgument(i18nkey.ValidationError)
	}

	// 明文密码长度 6-32
	if len(in.Password) < 6 || len(in.Password) > 32 {
		return nil, grpcerror.InvalidArgument(i18nkey.ValidationError)
	}

	// 校验状态
	if in.Status != nil && (*in.Status < 1 || *in.Status > 2) {
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
		SetUsername(username).
		SetPassword(in.Password).
		SetDisplayName(displayName).
		SetNillableStatus(in.Status).
		Save(l.ctx)
	if err != nil {
		return nil, enterror.Handle(l.Logger, err)
	}

	return &adminpb.CreateAdminResponse{
		Id: data.ID, // 管理员ID
	}, nil
}
