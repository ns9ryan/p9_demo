// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package region_allocation

import (
	"context"
	"strings"

	"oa.98ent.com/p9/common/xerr"
	"oa.98ent.com/p9/platform-base/rpc/pb/platformbaserpc/regionpb"
	"oa.98ent.com/p9/platform-operator/api/internal/i18nkey"
	"oa.98ent.com/p9/platform-operator/api/internal/svc"
	"oa.98ent.com/p9/platform-operator/api/internal/types"
	"oa.98ent.com/p9/platform-operator/rpc/pb/platformoperatorrpc/regionallocationpb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SaveRegionAllocationsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSaveRegionAllocationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveRegionAllocationsLogic {
	return &SaveRegionAllocationsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// SaveRegionAllocations 保存经营地区分配
func (l *SaveRegionAllocationsLogic) SaveRegionAllocations(req *types.SaveRegionAllocationsRequest) (resp *types.SaveRegionAllocationsResponse, err error) {
	regionCodes := make([]string, 0, len(req.RegionCodes))

	if len(req.RegionCodes) > 0 {
		// 获取全部启用的国家地区
		result, err := l.svcCtx.RegionRpc.ListAll(
			l.ctx,
			&regionpb.ListAllRegionsRequest{
				Status: new(int64(1)), // 只获取启用的国家地区
			},
		)
		if err != nil {
			return nil, err
		}

		// 建立可用国家地区编码索引
		regionMap := make(map[string]string, len(result.List))
		for _, region := range result.List {
			regionMap[strings.ToUpper(region.Code)] = region.Code
		}

		// 校验并转换为Platform Base标准国家地区编码
		for _, regionCode := range req.RegionCodes {
			code := strings.ToUpper(strings.TrimSpace(regionCode))

			standardCode, ok := regionMap[code]
			if !ok {
				return nil, xerr.BadRequest(i18nkey.RegionUnavailable)
			}

			regionCodes = append(regionCodes, standardCode)
		}
	}

	// 保存分站经营地区分配
	_, err = l.svcCtx.RegionAllocationRpc.Save(
		l.ctx,
		&regionallocationpb.SaveRegionAllocationsRequest{
			OperatorId:  req.OperatorId, // 分站ID
			RegionCodes: regionCodes,    // 当前分配的国家地区编码
		},
	)
	if err != nil {
		return nil, err
	}

	// 返回保存结果
	return &types.SaveRegionAllocationsResponse{}, nil
}
