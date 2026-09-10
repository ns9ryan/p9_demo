// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package timezone

import (
	"context"

	corei18n "oa.98ent.com/p9/core/common/i18n"
	"oa.98ent.com/p9/platform-base/api/internal/svc"
	"oa.98ent.com/p9/platform-base/api/internal/types"
	"oa.98ent.com/p9/platform-base/rpc/pb/base/timezone"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetTimezoneLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetTimezoneLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTimezoneLogic {
	return &GetTimezoneLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// GetTimezone 获取时区
func (l *GetTimezoneLogic) GetTimezone(req *types.GetTimezoneRequest) (resp *types.GetTimezoneResponse, err error) {
	// 调用获取时区RPC
	result, err := l.svcCtx.TimezoneRpc.Get(
		l.ctx,
		&timezone.GetTimezoneRequest{
			Id: req.Id, // 时区ID
		},
	)
	if err != nil {
		return nil, err
	}

	// 获取当前语言的时区名称
	name := corei18n.TG(l.ctx, corei18n.CodePlatform, "base", result.Timezone.NameKey)

	// 返回时区信息
	return &types.GetTimezoneResponse{
		TimezoneInfo: types.TimezoneInfo{
			Id:      result.Timezone.Id,      // 时区ID
			Code:    result.Timezone.Code,    // IANA时区编码
			NameKey: result.Timezone.NameKey, // 名称翻译Key
			Name:    name,                    // 当前语言名称
			Status:  result.Timezone.Status,  // 状态: 1启用, 2停用
			SortNo:  result.Timezone.SortNo,  // 排序值
		},
	}, nil
}
