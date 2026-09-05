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

type CreateTimezoneLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateTimezoneLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateTimezoneLogic {
	return &CreateTimezoneLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// CreateTimezone 创建时区
func (l *CreateTimezoneLogic) CreateTimezone(req *types.CreateTimezoneRequest) (resp *types.CreateTimezoneResponse, err error) {
	// 调用创建时区RPC
	result, err := l.svcCtx.TimezoneRpc.Create(
		l.ctx,
		&timezone.CreateTimezoneRequest{
			Code:     req.Code,     // IANA时区编码
			NameI18N: req.NameI18n, // 多语言名称
			Status:   req.Status,   // 状态: 1启用, 2停用
		},
	)
	if err != nil {
		return nil, err
	}

	// 返回创建结果
	return &types.CreateTimezoneResponse{
		Id: result.Id, // 时区ID
	}, nil
}
