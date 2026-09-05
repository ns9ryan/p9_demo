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

type UpdateRegionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateRegionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateRegionLogic {
	return &UpdateRegionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// UpdateRegion 修改国家地区
func (l *UpdateRegionLogic) UpdateRegion(req *types.UpdateRegionRequest) (resp *types.UpdateRegionResponse, err error) {
	// 调用修改国家地区RPC
	_, err = l.svcCtx.RegionRpc.Update(
		l.ctx,
		&region.UpdateRegionRequest{
			Id:          req.Id,          // 国家或地区ID
			CallingCode: req.CallingCode, // 国际电话区号, 不包含加号
			NameI18N:    req.NameI18n,    // 多语言名称
			Status:      req.Status,      // 状态: 1启用, 2停用
		},
	)
	if err != nil {
		return nil, err
	}

	// 返回修改结果
	return &types.UpdateRegionResponse{}, nil
}
