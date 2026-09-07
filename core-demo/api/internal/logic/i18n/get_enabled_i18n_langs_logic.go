package i18n

import (
	"context"

	"oa.98ent.com/p9/core/api/internal/convert"
	"oa.98ent.com/p9/core/api/internal/svc"
	"oa.98ent.com/p9/core/api/internal/types"
	"oa.98ent.com/p9/core/rpc/coreclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetEnabledI18nLangsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetEnabledI18nLangsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetEnabledI18nLangsLogic {
	return &GetEnabledI18nLangsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetEnabledI18nLangsLogic) GetEnabledI18nLangs() (resp *types.I18nLangListResp, err error) {
	out, err := l.svcCtx.Core.GetEnabledI18NLangs(l.ctx, &coreclient.Empty{})
	if err != nil {
		return nil, err
	}
	return convert.I18nLangList(out), nil
}
