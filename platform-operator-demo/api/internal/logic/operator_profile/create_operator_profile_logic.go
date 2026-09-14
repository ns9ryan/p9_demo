// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package operator_profile

import (
	"context"

	"oa.98ent.com/p9/platform-operator/api/internal/svc"
	"oa.98ent.com/p9/platform-operator/api/internal/types"
	"oa.98ent.com/p9/platform-operator/rpc/pb/platformoperatorrpc/profilepb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateOperatorProfileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateOperatorProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateOperatorProfileLogic {
	return &CreateOperatorProfileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// CreateOperatorProfile 创建分站档案
func (l *CreateOperatorProfileLogic) CreateOperatorProfile(req *types.CreateOperatorProfileRequest) (resp *types.CreateOperatorProfileResponse, err error) {
	// 创建分站档案
	result, err := l.svcCtx.OperatorProfileRpc.Create(
		l.ctx,
		&profilepb.CreateOperatorProfileRequest{
			OperatorId:   req.OperatorId,   // 分站ID
			CompanyName:  req.CompanyName,  // 公司名称
			ContactName:  req.ContactName,  // 主要联系人名称
			ContactEmail: req.ContactEmail, // 主要联系人邮箱
			Remark:       req.Remark,       // 总网内部档案备注
		},
	)
	if err != nil {
		return nil, err
	}

	// 返回创建结果
	return &types.CreateOperatorProfileResponse{
		Id: result.Id, // 档案ID
	}, nil
}
