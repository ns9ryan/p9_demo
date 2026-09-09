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

type ListTimezonesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListTimezonesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListTimezonesLogic {
	return &ListTimezonesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// ListTimezones 获取时区管理列表
func (l *ListTimezonesLogic) ListTimezones(req *types.ListTimezonesRequest) (resp *types.ListTimezonesResponse, err error) {
	// 调用获取时区管理列表RPC
	result, err := l.svcCtx.TimezoneRpc.List(
		l.ctx,
		&timezone.ListTimezonesRequest{
			Page:     req.Page,     // 页码
			PageSize: req.PageSize, // 每页数量
			Status:   req.Status,   // 状态: 1启用, 2停用
		},
	)
	if err != nil {
		return nil, err
	}

	// 转换时区列表
	list := make([]types.TimezoneInfo, 0, len(result.List))
	for _, item := range result.List {
		// 获取当前语言的时区名称
		name := corei18n.TG(l.ctx, "base", item.NameKey)

		list = append(list, types.TimezoneInfo{
			Id:      item.Id,      // 时区ID
			Code:    item.Code,    // IANA时区编码
			NameKey: item.NameKey, // 名称翻译Key
			Name:    name,         // 当前语言名称
			Status:  item.Status,  // 状态: 1启用, 2停用
			SortNo:  item.SortNo,  // 排序值
		})
	}

	// 返回时区管理列表
	return &types.ListTimezonesResponse{
		Total: result.Total, // 数据总数
		List:  list,         // 时区列表
	}, nil
}
