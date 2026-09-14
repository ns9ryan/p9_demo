// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package region_allocation

import (
	"context"

	"oa.98ent.com/p9/platform-operator/api/internal/svc"
	"oa.98ent.com/p9/platform-operator/api/internal/types"
	"oa.98ent.com/p9/platform-operator/rpc/pb/platformoperatorrpc/regionallocationpb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListRegionAllocationsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListRegionAllocationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListRegionAllocationsLogic {
	return &ListRegionAllocationsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// ListRegionAllocations 获取经营地区分配列表
func (l *ListRegionAllocationsLogic) ListRegionAllocations(req *types.ListRegionAllocationsRequest) (resp *types.ListRegionAllocationsResponse, err error) {
	// 获取分站当前经营地区分配关系
	result, err := l.svcCtx.RegionAllocationRpc.List(
		l.ctx,
		&regionallocationpb.ListRegionAllocationsRequest{
			OperatorId: req.OperatorId, // 分站ID
		},
	)
	if err != nil {
		return nil, err
	}

	// 转换经营地区分配列表
	list := make([]types.RegionAllocationInfo, 0, len(result.List))
	for _, item := range result.List {
		list = append(list, types.RegionAllocationInfo{
			RegionCode:  item.RegionCode, // 国家地区编码
			AllocatedAt: item.CreatedAt,  // 分配时间, Unix毫秒时间戳
		})
	}

	// 返回经营地区分配列表
	return &types.ListRegionAllocationsResponse{
		List: list, // 已分配经营地区列表
	}, nil
}
