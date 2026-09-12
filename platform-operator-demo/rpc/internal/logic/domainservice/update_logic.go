package domainservicelogic

import (
	"context"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"

	"oa.98ent.com/p9/platform-operator/pkg/i18nkey"
	"oa.98ent.com/p9/platform-operator/pkg/rpc/grpcerror"
	"oa.98ent.com/p9/platform-operator/rpc/internal/enterror"
	"oa.98ent.com/p9/platform-operator/rpc/internal/svc"
	"oa.98ent.com/p9/platform-operator/rpc/pb/operator/domain"
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

// Update 修改分站域名
func (l *UpdateLogic) Update(in *domain.UpdateDomainRequest) (*domain.UpdateDomainResponse, error) {
	// 域名ID必须大于0
	if in.Id <= 0 {
		return nil, grpcerror.InvalidArgument(i18nkey.ValidationError)
	}

	// 至少需要修改一个字段
	if in.DomainName == nil &&
		in.DomainType == nil &&
		in.Status == nil &&
		in.Remark == nil {
		return nil, grpcerror.InvalidArgument(i18nkey.ValidationError)
	}

	// 获取当前分站域名
	current, err := l.svcCtx.DB.OperatorDomain.Get(l.ctx, in.Id)
	if err != nil {
		// 转换Ent错误为gRPC错误
		return nil, enterror.Handle(l.Logger, err)
	}

	// 标准化域名
	var domainName *string
	if in.DomainName != nil {
		value := strings.ToLower(strings.TrimSpace(*in.DomainName))
		domainName = &value
	}

	// 修改分站域名
	err = current.
		Update().
		SetNillableDomainName(domainName).    // 域名
		SetNillableDomainType(in.DomainType). // 域名类型: 1分站后台, 2代理后台, 3会员H5
		SetNillableStatus(in.Status).         // 域名状态: 1启用, 2停用
		SetNillableRemark(in.Remark).         // 总网内部备注
		Exec(l.ctx)
	if err != nil {
		// 转换Ent错误为gRPC错误
		return nil, enterror.Handle(l.Logger, err)
	}

	// 返回修改结果
	return &domain.UpdateDomainResponse{}, nil
}
