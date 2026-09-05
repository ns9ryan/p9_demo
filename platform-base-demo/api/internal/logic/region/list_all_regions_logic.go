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

type ListAllRegionsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListAllRegionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListAllRegionsLogic {
	return &ListAllRegionsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// ListAllRegions 获取全部国家地区
func (l *ListAllRegionsLogic) ListAllRegions(req *types.ListAllRegionsRequest) (resp *types.ListAllRegionsResponse, err error) {
	// 调用获取全部国家地区RPC
	result, err := l.svcCtx.RegionRpc.ListAll(
		l.ctx,
		&region.ListAllRegionsRequest{
			Status: req.Status, // 状态: 1启用, 2停用
		},
	)
	if err != nil {
		return nil, err
	}

	// 转换国家地区列表
	list := make([]types.RegionInfo, 0, len(result.List))
	for _, item := range result.List {
		list = append(list, types.RegionInfo{
			Id:          item.Id,          // 国家或地区ID
			Code:        item.Code,        // 国家或地区编码
			CallingCode: item.CallingCode, // 国际电话区号, 不包含加号
			NameI18n:    item.NameI18N,    // 多语言名称
			Status:      item.Status,      // 状态: 1启用, 2停用
			SortNo:      item.SortNo,      // 排序值
		})
	}

	// 返回全部国家地区
	return &types.ListAllRegionsResponse{
		List: list, // 国家地区列表
	}, nil
}
