package operatorprofileservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"oa.98ent.com/p9/platform-operator/pkg/i18nkey"
	"oa.98ent.com/p9/platform-operator/pkg/rpc/grpcerror"
	"oa.98ent.com/p9/platform-operator/rpc/internal/enterror"
	"oa.98ent.com/p9/platform-operator/rpc/internal/svc"
	"oa.98ent.com/p9/platform-operator/rpc/pb/operator/profile"
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

// Create 创建分站档案
func (l *CreateLogic) Create(in *profile.CreateOperatorProfileRequest) (*profile.CreateOperatorProfileResponse, error) {
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

	// 创建分站档案
	data, err := l.svcCtx.DB.OperatorProfile.
		Create().
		SetOperatorID(in.OperatorId).             // 分站ID
		SetNillableCompanyName(in.CompanyName).   // 公司名称
		SetNillableContactName(in.ContactName).   // 主要联系人名称
		SetNillableContactEmail(in.ContactEmail). // 主要联系人邮箱
		SetNillableRemark(in.Remark).             // 总网内部档案备注
		Save(l.ctx)
	if err != nil {
		// 转换Ent错误为gRPC错误
		return nil, enterror.Handle(l.Logger, err)
	}

	// 返回创建结果
	return &profile.CreateOperatorProfileResponse{
		Id: data.ID, // 档案ID
	}, nil
}
