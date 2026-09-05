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

type ListAllLanguagesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListAllLanguagesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListAllLanguagesLogic {
	return &ListAllLanguagesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// ListAllLanguages 获取全部语言
func (l *ListAllLanguagesLogic) ListAllLanguages(req *types.ListAllLanguagesRequest) (resp *types.ListAllLanguagesResponse, err error) {
	// 调用获取全部语言RPC
	result, err := l.svcCtx.LanguageRpc.ListAll(
		l.ctx,
		&language.ListAllLanguagesRequest{
			Status: req.Status, // 状态: 1启用, 2停用
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

	// 返回全部语言
	return &types.ListAllLanguagesResponse{
		List: list, // 语言列表
	}, nil
}
