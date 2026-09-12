package domainservicelogic

import (
	"context"
	"strings"

	"entgo.io/ent/dialect/sql"
	"github.com/zeromicro/go-zero/core/logx"

	"oa.98ent.com/p9/platform-operator/pkg/i18nkey"
	"oa.98ent.com/p9/platform-operator/pkg/rpc/grpcerror"
	"oa.98ent.com/p9/platform-operator/rpc/ent/operatordomain"
	"oa.98ent.com/p9/platform-operator/rpc/internal/enterror"
	"oa.98ent.com/p9/platform-operator/rpc/internal/svc"
	"oa.98ent.com/p9/platform-operator/rpc/pb/operator/domain"
)

type ListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListLogic {
	return &ListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// List 获取分站域名管理列表
func (l *ListLogic) List(in *domain.ListDomainsRequest) (*domain.ListDomainsResponse, error) {
	// 校验分页参数
	if in.Page < 1 || in.PageSize < 1 || in.PageSize > 100 {
		return nil, grpcerror.InvalidArgument(i18nkey.ValidationError)
	}

	// 校验分站ID
	if in.OperatorId != nil && *in.OperatorId <= 0 {
		return nil, grpcerror.InvalidArgument(i18nkey.ValidationError)
	}

	// 校验域名类型
	if in.DomainType != nil && (*in.DomainType < 1 || *in.DomainType > 3) {
		return nil, grpcerror.InvalidArgument(i18nkey.ValidationError)
	}

	// 校验域名状态
	if in.Status != nil && (*in.Status < 1 || *in.Status > 2) {
		return nil, grpcerror.InvalidArgument(i18nkey.ValidationError)
	}

	// 创建分站域名查询
	query := l.svcCtx.DB.OperatorDomain.Query()

	// 按分站筛选
	if in.OperatorId != nil {
		query = query.Where(operatordomain.OperatorIDEQ(*in.OperatorId))
	}

	// 按关键字筛选
	if in.Keyword != nil {
		keyword := strings.TrimSpace(*in.Keyword)
		if keyword != "" {
			query = query.Where(operatordomain.DomainNameContainsFold(keyword))
		}
	}

	// 按域名类型筛选
	if in.DomainType != nil {
		query = query.Where(operatordomain.DomainTypeEQ(*in.DomainType))
	}

	// 按域名状态筛选
	if in.Status != nil {
		query = query.Where(operatordomain.StatusEQ(*in.Status))
	}

	// 获取符合条件的数据总数
	total, err := query.Clone().Count(l.ctx)
	if err != nil {
		// 转换Ent错误为gRPC错误
		return nil, enterror.Handle(l.Logger, err)
	}

	// 计算分页偏移量
	offset := (in.Page - 1) * in.PageSize

	// 获取当前页分站域名数据
	results, err := query.
		Order(
			operatordomain.ByCreatedAt(sql.OrderDesc()), // 按创建时间倒序
			operatordomain.ByID(sql.OrderDesc()),        // 创建时间相同时按ID倒序
		).
		Offset(int(offset)).
		Limit(int(in.PageSize)).
		All(l.ctx)
	if err != nil {
		// 转换Ent错误为gRPC错误
		return nil, enterror.Handle(l.Logger, err)
	}

	// 转换分站域名列表
	list := make([]*domain.DomainInfo, 0, len(results))
	for _, result := range results {
		list = append(list, toDomainInfo(result))
	}

	// 返回分站域名列表
	return &domain.ListDomainsResponse{
		Total: int64(total), // 数据总数
		List:  list,         // 分站域名列表
	}, nil
}
