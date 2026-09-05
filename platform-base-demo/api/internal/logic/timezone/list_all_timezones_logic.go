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

type ListAllTimezonesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListAllTimezonesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListAllTimezonesLogic {
	return &ListAllTimezonesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// ListAllTimezones 获取全部时区
func (l *ListAllTimezonesLogic) ListAllTimezones(req *types.ListAllTimezonesRequest) (resp *types.ListAllTimezonesResponse, err error) {
	// 调用获取全部时区RPC
	result, err := l.svcCtx.TimezoneRpc.ListAll(
		l.ctx,
		&timezone.ListAllTimezonesRequest{
			Status: req.Status, // 状态: 1启用, 2停用
		},
	)
	if err != nil {
		return nil, err
	}

	// 转换时区列表
	list := make([]types.TimezoneInfo, 0, len(result.List))
	for _, item := range result.List {
		list = append(list, types.TimezoneInfo{
			Id:       item.Id,       // 时区ID
			Code:     item.Code,     // IANA时区编码
			NameI18n: item.NameI18N, // 多语言名称
			Status:   item.Status,   // 状态: 1启用, 2停用
			SortNo:   item.SortNo,   // 排序值
		})
	}

	// 返回全部时区
	return &types.ListAllTimezonesResponse{
		List: list, // 时区列表
	}, nil
}
