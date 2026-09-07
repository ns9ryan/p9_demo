package i18n

import (
	"context"

	"oa.98ent.com/p9/core/api/internal/convert"
	"oa.98ent.com/p9/core/api/internal/svc"
	"oa.98ent.com/p9/core/api/internal/types"
	corei18n "oa.98ent.com/p9/core/common/i18n"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateI18nLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateI18nLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateI18nLogic {
	return &UpdateI18nLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateI18nLogic) UpdateI18n(req *types.UpdateI18nReq) (resp *types.BaseMsgResp, err error) {
	_, err = l.svcCtx.Core.UpdateI18N(l.ctx, convert.UpdateI18nReq(req))
	if err != nil {
		return nil, err
	}
	corei18n.InvalidateAll()
	return &types.BaseMsgResp{Result: "success"}, nil
}
