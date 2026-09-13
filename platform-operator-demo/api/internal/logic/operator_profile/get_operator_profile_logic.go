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

type GetOperatorProfileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetOperatorProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOperatorProfileLogic {
	return &GetOperatorProfileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// GetOperatorProfile 获取分站档案
func (l *GetOperatorProfileLogic) GetOperatorProfile(req *types.GetOperatorProfileRequest) (resp *types.GetOperatorProfileResponse, err error) {
	// 获取分站档案
	result, err := l.svcCtx.OperatorProfileRpc.Get(
		l.ctx,
		&profilepb.GetOperatorProfileRequest{
			OperatorId: req.OperatorId, // 分站ID
		},
	)
	if err != nil {
		return nil, err
	}

	// 返回分站档案信息
	return &types.GetOperatorProfileResponse{
		Profile: toOperatorProfileInfo(result.Profile), // 分站档案信息
	}, nil
}
