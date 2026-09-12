package operatorservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"oa.98ent.com/p9/platform-operator/pkg/i18nkey"
	"oa.98ent.com/p9/platform-operator/pkg/rpc/grpcerror"
	"oa.98ent.com/p9/platform-operator/rpc/internal/enterror"
	"oa.98ent.com/p9/platform-operator/rpc/internal/svc"
	"oa.98ent.com/p9/platform-operator/rpc/pb/operator/operator"
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

// Update 修改分站
func (l *UpdateLogic) Update(in *operator.UpdateOperatorRequest) (*operator.UpdateOperatorResponse, error) {
	// 分站ID必须大于0
	if in.Id <= 0 {
		return nil, grpcerror.InvalidArgument(i18nkey.ValidationError)
	}

	// 至少需要修改一个字段
	if in.Name == nil &&
		in.TimezoneCode == nil &&
		in.SettlementCurrencyCode == nil &&
		in.Status == nil &&
		in.Remark == nil {
		return nil, grpcerror.InvalidArgument(i18nkey.ValidationError)
	}

	// 获取当前分站
	current, err := l.svcCtx.DB.Operator.Get(l.ctx, in.Id)
	if err != nil {
		// 转换Ent错误为gRPC错误
		return nil, enterror.Handle(l.Logger, err)
	}

	// 判断是否修改时区或结算币种
	timezoneChanged := in.TimezoneCode != nil && *in.TimezoneCode != current.TimezoneCode
	currencyChanged := in.SettlementCurrencyCode != nil &&
		*in.SettlementCurrencyCode != current.SettlementCurrencyCode

	// 发布中或已经发布过的分站不能修改时区和结算币种
	if (current.PublishStatus == 2 || current.PublishedAt != nil) &&
		(timezoneChanged || currencyChanged) {
		return nil, grpcerror.InvalidArgument(i18nkey.ValidationError)
	}

	// 修改分站
	err = current.
		Update().
		SetNillableName(in.Name).                                     // 分站名称
		SetNillableTimezoneCode(in.TimezoneCode).                     // 时区编码
		SetNillableSettlementCurrencyCode(in.SettlementCurrencyCode). // 结算币种编码
		SetNillableStatus(in.Status).                                 // 分站状态: 1正常, 2暂停, 3关闭
		SetNillableRemark(in.Remark).                                 // 内部备注
		Exec(l.ctx)
	if err != nil {
		// 转换Ent错误为gRPC错误
		return nil, enterror.Handle(l.Logger, err)
	}

	// 返回修改结果
	return &operator.UpdateOperatorResponse{}, nil
}
