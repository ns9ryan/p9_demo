// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package basic_resource_allocation

import (
	"context"

	"oa.98ent.com/p9/platform-operator/api/internal/svc"
	"oa.98ent.com/p9/platform-operator/api/internal/types"
	"oa.98ent.com/p9/platform-operator/rpc/pb/platformoperatorrpc/basicresourceallocationpb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListBasicResourceAllocationsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListBasicResourceAllocationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListBasicResourceAllocationsLogic {
	return &ListBasicResourceAllocationsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// ListBasicResourceAllocations 获取基础资源分配列表
func (l *ListBasicResourceAllocationsLogic) ListBasicResourceAllocations(req *types.ListBasicResourceAllocationsRequest) (resp *types.ListBasicResourceAllocationsResponse, err error) {
	// 获取基础资源分配列表
	result, err := l.svcCtx.BasicResourceAllocationRpc.List(
		l.ctx,
		&basicresourceallocationpb.ListBasicResourceAllocationsRequest{
			Page:     req.Page,     // 页码, 从1开始
			PageSize: req.PageSize, // 每页数量
			Keyword:  req.Keyword,  // 搜索关键字, 匹配分站编码或名称
		},
	)
	if err != nil {
		return nil, err
	}

	// 转换基础资源分配列表
	list := make([]types.BasicResourceAllocationInfo, 0, len(result.List))
	for _, item := range result.List {
		list = append(list, types.BasicResourceAllocationInfo{
			OperatorId:     item.OperatorId,     // 分站ID
			OperatorCode:   item.OperatorCode,   // 分站业务编码
			OperatorName:   item.OperatorName,   // 分站名称
			LanguageCount:  item.LanguageCount,  // 已分配语言数量
			RegionCount:    item.RegionCount,    // 已分配经营地区数量
			AgentLineCount: item.AgentLineCount, // 已分配代理子线路数量
		})
	}

	// 返回基础资源分配列表
	return &types.ListBasicResourceAllocationsResponse{
		Total: result.Total, // 数据总数
		List:  list,         // 基础资源分配列表
	}, nil
}
