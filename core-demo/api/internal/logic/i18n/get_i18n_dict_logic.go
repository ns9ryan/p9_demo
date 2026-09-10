package i18n

import (
	"context"

	"oa.98ent.com/p9/core/api/internal/svc"
	"oa.98ent.com/p9/core/api/internal/types"
	"oa.98ent.com/p9/core/common/i18n"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetI18nDictLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetI18nDictLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetI18nDictLogic {
	return &GetI18nDictLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// GetI18nDict 获取多语言词条
func (l *GetI18nDictLogic) GetI18nDict(req *types.GetI18nDictReq) (resp *types.I18nDictResp, err error) {
	items := i18n.Dict(l.ctx, req.I18nCode, req.I18nGroup, req.Lang)
	if items == nil {
		items = map[string]string{}
	}

	return &types.I18nDictResp{Items: items}, nil

	// out, err := l.svcCtx.Core.GetI18NDict(l.ctx, convert.I18nDictReq(req))
	// if err != nil {
	// 	return nil, err
	// }
	// return convert.I18nDict(out), nil
}
