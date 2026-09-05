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

type GetLanguageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetLanguageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLanguageLogic {
	return &GetLanguageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// GetLanguage 获取语言
func (l *GetLanguageLogic) GetLanguage(req *types.GetLanguageRequest) (resp *types.GetLanguageResponse, err error) {
	// 调用获取语言RPC
	result, err := l.svcCtx.LanguageRpc.Get(
		l.ctx,
		&language.GetLanguageRequest{
			Id: req.Id, // 语言ID
		},
	)
	if err != nil {
		return nil, err
	}

	// 返回语言信息
	return &types.GetLanguageResponse{
		Language: types.LanguageInfo{
			Id:       result.Language.Id,       // 语言ID
			Code:     result.Language.Code,     // 语言编码
			NameI18n: result.Language.NameI18N, // 多语言名称
			Status:   result.Language.Status,   // 状态: 1启用, 2停用
			SortNo:   result.Language.SortNo,   // 排序值
		},
	}, nil
}
