package authority

import (
	"context"

	"oa.98ent.com/p9/core/common/ctxdata"
	"oa.98ent.com/p9/core/common/utils"
	"oa.98ent.com/p9/core/common/xerr"
	"oa.98ent.com/p9/core/rpc/internal/svc"
	"oa.98ent.com/p9/core/rpc/types/core"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetApiAuthorityLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetApiAuthorityLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetApiAuthorityLogic {
	return &GetApiAuthorityLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetApiAuthorityLogic) GetApiAuthority(in *core.RoleIdReq) (*core.ApiAuthResp, error) {
	list, err := l.svcCtx.Deps.GetAPIAuthority(l.ctx, ctxdata.ClaimsFromCtx(l.ctx), utils.GetEffectiveId(in.RoleId, in.Id))
	if err != nil {
		return nil, xerr.RpcErr(err)
	}
	out := make([]*core.ApiAuthItem, 0, len(list))
	for _, it := range list {
		out = append(out, &core.ApiAuthItem{Path: it.Path, Method: it.Method})
	}
	return &core.ApiAuthResp{List: out}, nil
}
