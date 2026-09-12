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

// Create 创建分站域名
func (l *CreateLogic) Create(in *domain.CreateDomainRequest) (*domain.CreateDomainResponse, error) {
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

	// 标准化域名
	domainName := strings.ToLower(strings.TrimSpace(in.DomainName))

	// 创建分站域名
	data, err := l.svcCtx.DB.OperatorDomain.
		Create().
		SetOperatorID(in.OperatorId). // 分站ID
		SetDomainName(domainName).    // 域名
		SetDomainType(in.DomainType). // 域名类型: 1分站后台, 2代理后台, 3会员H5
		SetNillableStatus(in.Status). // 域名状态: 1启用, 2停用
		SetNillableRemark(in.Remark). // 总网内部备注
		Save(l.ctx)
	if err != nil {
		// 转换Ent错误为gRPC错误
		return nil, enterror.Handle(l.Logger, err)
	}

	// 返回创建结果
	return &domain.CreateDomainResponse{
		Id: data.ID, // 域名ID
	}, nil
}
