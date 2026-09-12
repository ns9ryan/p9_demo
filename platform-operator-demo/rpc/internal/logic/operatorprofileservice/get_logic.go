package operatorprofileservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"oa.98ent.com/p9/platform-operator/pkg/i18nkey"
	"oa.98ent.com/p9/platform-operator/pkg/rpc/grpcerror"
	"oa.98ent.com/p9/platform-operator/rpc/ent"
	"oa.98ent.com/p9/platform-operator/rpc/ent/operatorprofile"
	"oa.98ent.com/p9/platform-operator/rpc/internal/enterror"
	"oa.98ent.com/p9/platform-operator/rpc/internal/svc"
	"oa.98ent.com/p9/platform-operator/rpc/pb/operator/profile"
)

type GetLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLogic {
	return &GetLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// Get 获取分站档案
func (l *GetLogic) Get(in *profile.GetOperatorProfileRequest) (*profile.GetOperatorProfileResponse, error) {
	// 分站ID必须大于0
	if in.OperatorId <= 0 {
		return nil, grpcerror.InvalidArgument(i18nkey.ValidationError)
	}

	// 确认分站存在
	_, err := l.svcCtx.DB.Operator.Get(l.ctx, in.OperatorId)
	if err != nil {
		// 转换Ent错误为gRPC错误
		return nil, enterror.Handle(l.Logger, err)
	}

	// 获取分站档案
	result, err := l.svcCtx.DB.OperatorProfile.
		Query().
		Where(operatorprofile.OperatorIDEQ(in.OperatorId)).
		Only(l.ctx)
	if err != nil {
		// 尚未创建档案时返回空
		if ent.IsNotFound(err) {
			return &profile.GetOperatorProfileResponse{
				Profile: nil, // 分站档案信息
			}, nil
		}

		// 转换Ent错误为gRPC错误
		return nil, enterror.Handle(l.Logger, err)
	}

	// 返回分站档案信息
	return &profile.GetOperatorProfileResponse{
		Profile: toOperatorProfileInfo(result), // 分站档案信息
	}, nil
}
