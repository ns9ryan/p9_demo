// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package operator_domain

import (
	"context"

	"oa.98ent.com/p9/platform-operator/api/internal/svc"
	"oa.98ent.com/p9/platform-operator/api/internal/types"
	"oa.98ent.com/p9/platform-operator/rpc/pb/platformoperatorrpc/domainpb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListOperatorDomainsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListOperatorDomainsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListOperatorDomainsLogic {
	return &ListOperatorDomainsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// ListOperatorDomains 获取分站域名列表
func (l *ListOperatorDomainsLogic) ListOperatorDomains(req *types.ListOperatorDomainsRequest) (resp *types.ListOperatorDomainsResponse, err error) {
	// 获取分站域名列表
	result, err := l.svcCtx.OperatorDomainRpc.List(
		l.ctx,
		&domainpb.ListDomainsRequest{
			Page:       req.Page,       // 页码, 从1开始
			PageSize:   req.PageSize,   // 每页数量
			OperatorId: req.OperatorId, // 分站ID
			Keyword:    req.Keyword,    // 搜索关键字, 匹配域名
			DomainType: req.DomainType, // 域名类型: 1分站后台, 2代理后台, 3会员H5
			Status:     req.Status,     // 域名状态: 1启用, 2停用
		},
	)
	if err != nil {
		return nil, err
	}

	// 转换分站域名列表
	list := make([]types.OperatorDomainInfo, 0, len(result.List))
	for _, item := range result.List {
		list = append(list, toOperatorDomainInfo(item))
	}

	// 返回分站域名列表
	return &types.ListOperatorDomainsResponse{
		Total: result.Total, // 数据总数
		List:  list,         // 分站域名列表
	}, nil
}
