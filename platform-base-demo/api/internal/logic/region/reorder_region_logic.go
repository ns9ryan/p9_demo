// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package region

import (
	"context"

	"oa.98ent.com/p9/platform-base/api/internal/svc"
	"oa.98ent.com/p9/platform-base/api/internal/types"
	"oa.98ent.com/p9/platform-base/rpc/pb/base/region"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReorderRegionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewReorderRegionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReorderRegionLogic {
	return &ReorderRegionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// ReorderRegion 调整国家地区排序
func (l *ReorderRegionLogic) ReorderRegion(req *types.ReorderRegionRequest) (resp *types.ReorderRegionResponse, err error) {
	// 调用调整国家地区排序RPC
	_, err = l.svcCtx.RegionRpc.Reorder(
		l.ctx,
		&region.ReorderRegionRequest{
			Id:       req.Id,       // 需要移动的国家地区ID
			TargetId: req.TargetId, // 目标国家地区ID
		},
	)
	if err != nil {
		return nil, err
	}

	// 返回调整排序结果
	return &types.ReorderRegionResponse{}, nil
}
