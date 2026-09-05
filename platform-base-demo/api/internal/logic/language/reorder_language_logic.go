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

type ReorderLanguageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewReorderLanguageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReorderLanguageLogic {
	return &ReorderLanguageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// ReorderLanguage 调整语言排序
func (l *ReorderLanguageLogic) ReorderLanguage(req *types.ReorderLanguageRequest) (resp *types.ReorderLanguageResponse, err error) {
	// 调用调整语言排序RPC
	_, err = l.svcCtx.LanguageRpc.Reorder(
		l.ctx,
		&language.ReorderLanguageRequest{
			Id:       req.Id,       // 需要移动的语言ID
			TargetId: req.TargetId, // 目标语言ID
		},
	)
	if err != nil {
		return nil, err
	}

	// 返回调整排序结果
	return &types.ReorderLanguageResponse{}, nil
}
