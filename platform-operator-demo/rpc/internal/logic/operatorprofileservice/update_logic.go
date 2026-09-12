package operatorprofileservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"oa.98ent.com/p9/platform-operator/pkg/i18nkey"
	"oa.98ent.com/p9/platform-operator/pkg/rpc/grpcerror"
	"oa.98ent.com/p9/platform-operator/rpc/ent/operatorprofile"
	"oa.98ent.com/p9/platform-operator/rpc/internal/enterror"
	"oa.98ent.com/p9/platform-operator/rpc/internal/svc"
	"oa.98ent.com/p9/platform-operator/rpc/pb/operator/profile"
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

// Update 修改分站档案
func (l *UpdateLogic) Update(in *profile.UpdateOperatorProfileRequest) (*profile.UpdateOperatorProfileResponse, error) {
	// 分站ID必须大于0
	if in.OperatorId <= 0 {
		return nil, grpcerror.InvalidArgument(i18nkey.ValidationError)
	}

	// 至少需要修改一个字段
	if in.CompanyName == nil &&
		in.ContactName == nil &&
		in.ContactEmail == nil &&
		in.Remark == nil {
		return nil, grpcerror.InvalidArgument(i18nkey.ValidationError)
	}

	// 获取分站档案
	current, err := l.svcCtx.DB.OperatorProfile.
		Query().
		Where(operatorprofile.OperatorIDEQ(in.OperatorId)).
		Only(l.ctx)
	if err != nil {
		// 转换Ent错误为gRPC错误
		return nil, enterror.Handle(l.Logger, err)
	}

	// 修改分站档案
	err = current.
		Update().
		SetNillableCompanyName(in.CompanyName).   // 公司名称
		SetNillableContactName(in.ContactName).   // 主要联系人名称
		SetNillableContactEmail(in.ContactEmail). // 主要联系人邮箱
		SetNillableRemark(in.Remark).             // 总网内部档案备注
		Exec(l.ctx)
	if err != nil {
		// 转换Ent错误为gRPC错误
		return nil, enterror.Handle(l.Logger, err)
	}

	// 返回修改结果
	return &profile.UpdateOperatorProfileResponse{}, nil
}
