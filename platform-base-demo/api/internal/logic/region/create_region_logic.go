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

type CreateRegionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateRegionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateRegionLogic {
	return &CreateRegionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// CreateRegion 创建国家地区
func (l *CreateRegionLogic) CreateRegion(req *types.CreateRegionRequest) (resp *types.CreateRegionResponse, err error) {
	// 调用创建国家地区RPC
	result, err := l.svcCtx.RegionRpc.Create(
		l.ctx,
		&region.CreateRegionRequest{
			Code:        req.Code,        // 国家或地区编码
			CallingCode: req.CallingCode, // 国际电话区号, 不包含加号
			NameI18N:    req.NameI18n,    // 多语言名称
			Status:      req.Status,      // 状态: 1启用, 2停用
		},
	)
	if err != nil {
		return nil, err
	}

	// 返回创建结果
	return &types.CreateRegionResponse{
		Id: result.Id, // 国家或地区ID
	}, nil
}
