// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package timezone

import (
	"context"

	"oa.98ent.com/p9/platform-base/api/internal/svc"
	"oa.98ent.com/p9/platform-base/api/internal/types"
	"oa.98ent.com/p9/platform-base/rpc/pb/base/timezone"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateTimezoneLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateTimezoneLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateTimezoneLogic {
	return &UpdateTimezoneLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// UpdateTimezone 修改时区
func (l *UpdateTimezoneLogic) UpdateTimezone(req *types.UpdateTimezoneRequest) (resp *types.UpdateTimezoneResponse, err error) {
	// 调用修改时区RPC
	_, err = l.svcCtx.TimezoneRpc.Update(
		l.ctx,
		&timezone.UpdateTimezoneRequest{
			Id:     req.Id,     // 时区ID
			Status: req.Status, // 状态: 1启用, 2停用
		},
	)
	if err != nil {
		return nil, err
	}

	// 返回修改结果
	return &types.UpdateTimezoneResponse{}, nil
}
