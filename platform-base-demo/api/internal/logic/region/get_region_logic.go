// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package region

import (
	"context"

	corei18n "oa.98ent.com/p9/core/common/i18n"
	"oa.98ent.com/p9/platform-base/api/internal/svc"
	"oa.98ent.com/p9/platform-base/api/internal/types"
	"oa.98ent.com/p9/platform-base/rpc/pb/base/region"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRegionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetRegionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRegionLogic {
	return &GetRegionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// GetRegion 获取国家地区
func (l *GetRegionLogic) GetRegion(req *types.GetRegionRequest) (resp *types.GetRegionResponse, err error) {
	// 调用获取国家地区RPC
	result, err := l.svcCtx.RegionRpc.Get(
		l.ctx,
		&region.GetRegionRequest{
			Id: req.Id, // 国家或地区ID
		},
	)
	if err != nil {
		return nil, err
	}

	// 获取当前语言的国家地区名称
	name := corei18n.TG(l.ctx, corei18n.CodePlatform, "base", result.Region.NameKey)

	// 返回国家地区信息
	return &types.GetRegionResponse{
		RegionInfo: types.RegionInfo{
			Id:          result.Region.Id,          // 国家或地区ID
			Code:        result.Region.Code,        // 国家或地区编码
			CallingCode: result.Region.CallingCode, // 国际电话区号, 不包含加号
			NameKey:     result.Region.NameKey,     // 名称翻译Key
			Name:        name,                      // 当前语言名称
			Status:      result.Region.Status,      // 状态: 1启用, 2停用
			SortNo:      result.Region.SortNo,      // 排序值
		},
	}, nil
}
