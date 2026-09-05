// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package language

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"oa.98ent.com/p9/platform-base/api/internal/svc"
	"oa.98ent.com/p9/platform-base/api/internal/types"
	"oa.98ent.com/p9/platform-base/rpc/pb/base/language"
)

type CreateLanguageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateLanguageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLanguageLogic {
	return &CreateLanguageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// CreateLanguage 创建语言
func (l *CreateLanguageLogic) CreateLanguage(req *types.CreateLanguageRequest) (resp *types.CreateLanguageResponse, err error) {
	// 调用创建语言RPC
	result, err := l.svcCtx.LanguageRpc.Create(
		l.ctx,
		&language.CreateLanguageRequest{
			Code:     req.Code,     // 语言编码
			NameI18N: req.NameI18n, // 多语言名称
			Status:   req.Status,   // 状态: 1启用, 2停用
		},
	)
	if err != nil {
		return nil, err
	}

	// 返回创建结果
	return &types.CreateLanguageResponse{
		Id: result.Id, // 语言ID
	}, nil
}
