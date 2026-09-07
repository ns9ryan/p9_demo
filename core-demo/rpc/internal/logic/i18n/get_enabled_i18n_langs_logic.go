package i18n

import (
	"context"

	"oa.98ent.com/p9/core/common/xerr"
	"oa.98ent.com/p9/core/rpc/internal/logic"
	"oa.98ent.com/p9/core/rpc/internal/svc"
	"oa.98ent.com/p9/core/rpc/types/core"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetEnabledI18nLangsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetEnabledI18nLangsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetEnabledI18nLangsLogic {
	return &GetEnabledI18nLangsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetEnabledI18nLangsLogic) GetEnabledI18NLangs(_ *core.Empty) (*core.I18NLangListResp, error) {
	list, err := l.svcCtx.Deps.ListEnabledI18nLangs(l.ctx)
	if err != nil {
		return nil, xerr.RpcErr(err)
	}
	out := make([]*core.I18NLangInfo, 0, len(list))
	for _, row := range list {
		out = append(out, logic.ToI18nLangInfo(row))
	}
	return &core.I18NLangListResp{List: out, Total: int64(len(out))}, nil
}
