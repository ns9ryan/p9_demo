// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package language

import (
	"context"

	"oa.98ent.com/p9/platform-base/api/internal/svc"
	"oa.98ent.com/p9/platform-base/api/internal/types"
	"oa.98ent.com/p9/platform-base/rpc/pb/base/language"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListLanguagesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListLanguagesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListLanguagesLogic {
	return &ListLanguagesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// ListLanguages 获取语言管理列表
func (l *ListLanguagesLogic) ListLanguages(req *types.ListLanguagesRequest) (resp *types.ListLanguagesResponse, err error) {
	// 调用获取语言管理列表RPC
	result, err := l.svcCtx.LanguageRpc.List(
		l.ctx,
		&language.ListLanguagesRequest{
			Page:     req.Page,     // 页码
			PageSize: req.PageSize, // 每页数量
			Status:   req.Status,   // 状态: 1启用, 2停用
		},
	)
	if err != nil {
		return nil, err
	}

	// 转换语言列表
	list := make([]types.LanguageInfo, 0, len(result.List))
	for _, item := range result.List {
		list = append(list, types.LanguageInfo{
			Id:       item.Id,       // 语言ID
			Code:     item.Code,     // 语言编码
			NameI18n: item.NameI18N, // 多语言名称
			Status:   item.Status,   // 状态: 1启用, 2停用
			SortNo:   item.SortNo,   // 排序值
		})
	}

	// 返回语言管理列表
	return &types.ListLanguagesResponse{
		Total: result.Total, // 数据总数
		List:  list,         // 语言列表
	}, nil
}
