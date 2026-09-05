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

type UpdateLanguageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateLanguageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLanguageLogic {
	return &UpdateLanguageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// UpdateLanguage 修改语言
func (l *UpdateLanguageLogic) UpdateLanguage(req *types.UpdateLanguageRequest) (resp *types.UpdateLanguageResponse, err error) {
	// 调用修改语言RPC
	_, err = l.svcCtx.LanguageRpc.Update(
		l.ctx,
		&language.UpdateLanguageRequest{
			Id:       req.Id,       // 语言ID
			NameI18N: req.NameI18n, // 多语言名称
			Status:   req.Status,   // 状态: 1启用, 2停用
		},
	)
	if err != nil {
		return nil, err
	}

	// 返回修改结果
	return &types.UpdateLanguageResponse{}, nil
}
