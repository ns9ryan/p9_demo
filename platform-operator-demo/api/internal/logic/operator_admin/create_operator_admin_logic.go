// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package operator_admin

import (
	"context"
	"strings"

	apiI18nkey "oa.98ent.com/p9/platform-operator/api/internal/i18nkey"
	"oa.98ent.com/p9/platform-operator/api/internal/svc"
	"oa.98ent.com/p9/platform-operator/api/internal/types"
	"oa.98ent.com/p9/platform-operator/pkg/i18nkey"
	"oa.98ent.com/p9/platform-operator/pkg/rpc/grpcerror"
	"oa.98ent.com/p9/platform-operator/rpc/pb/platformoperatorrpc/adminpb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateOperatorAdminLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateOperatorAdminLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateOperatorAdminLogic {
	return &CreateOperatorAdminLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// CreateOperatorAdmin 创建分站管理员
func (l *CreateOperatorAdminLogic) CreateOperatorAdmin(req *types.CreateOperatorAdminRequest) (resp *types.CreateOperatorAdminResponse, err error) {
	// 分站ID必须大于0
	if req.OperatorId <= 0 {
		return nil, grpcerror.InvalidArgument(i18nkey.ValidationError)
	}

	username := strings.TrimSpace(req.Username)
	displayName := strings.TrimSpace(req.DisplayName)

	// 账号和显示名称不能为空
	if username == "" || displayName == "" {
		return nil, grpcerror.InvalidArgument(i18nkey.ValidationError)
	}

	// 明文密码长度 6-32
	if len(req.Password) < 6 || len(req.Password) > 32 {
		return nil, grpcerror.InvalidArgument(i18nkey.ValidationError)
	}

	// 校验状态
	if req.Status != nil && (*req.Status < 1 || *req.Status > 2) {
		return nil, grpcerror.InvalidArgument(i18nkey.ValidationError)
	}

	existed, err := l.svcCtx.OperatorAdminRpc.ExistsByOperatorId(
		l.ctx,
		&adminpb.ExistsAdminByOperatorIdRequest{
			OperatorId: req.OperatorId, // 分站ID
		},
	)
	if err != nil {
		return nil, err
	}
	if existed.Exists {
		return nil, grpcerror.AlreadyExists(apiI18nkey.AdminAlreadyExists)
	}

	result, err := l.svcCtx.OperatorAdminRpc.Create(
		l.ctx,
		&adminpb.CreateAdminRequest{
			OperatorId:  req.OperatorId,  // 分站ID
			Username:    req.Username,    // 账号
			Password:    req.Password,    // 密码
			DisplayName: req.DisplayName, // 显示名称
			Status:      req.Status,      // 状态: 1启用, 2停用
		},
	)
	if err != nil {
		return nil, err
	}

	return &types.CreateOperatorAdminResponse{
		Id: result.Id, // 管理员ID
	}, nil
}
