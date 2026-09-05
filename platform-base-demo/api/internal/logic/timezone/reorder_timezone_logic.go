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

type ReorderTimezoneLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewReorderTimezoneLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReorderTimezoneLogic {
	return &ReorderTimezoneLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// ReorderTimezone 调整时区排序
func (l *ReorderTimezoneLogic) ReorderTimezone(req *types.ReorderTimezoneRequest) (resp *types.ReorderTimezoneResponse, err error) {
	// 调用调整时区排序RPC
	_, err = l.svcCtx.TimezoneRpc.Reorder(
		l.ctx,
		&timezone.ReorderTimezoneRequest{
			Id:       req.Id,       // 需要移动的时区ID
			TargetId: req.TargetId, // 目标时区ID
		},
	)
	if err != nil {
		return nil, err
	}

	// 返回调整排序结果
	return &types.ReorderTimezoneResponse{}, nil
}
