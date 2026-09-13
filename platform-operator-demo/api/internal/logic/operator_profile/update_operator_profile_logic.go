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

type UpdateOperatorProfileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateOperatorProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateOperatorProfileLogic {
	return &UpdateOperatorProfileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// UpdateOperatorProfile 修改分站档案
func (l *UpdateOperatorProfileLogic) UpdateOperatorProfile(req *types.UpdateOperatorProfileRequest) (resp *types.UpdateOperatorProfileResponse, err error) {
	// 修改分站档案
	_, err = l.svcCtx.OperatorProfileRpc.Update(
		l.ctx,
		&profilepb.UpdateOperatorProfileRequest{
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

	// 返回修改结果
	return &types.UpdateOperatorProfileResponse{}, nil
}
