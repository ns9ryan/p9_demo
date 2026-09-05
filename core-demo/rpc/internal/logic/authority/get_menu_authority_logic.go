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

type GetMenuAuthorityLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetMenuAuthorityLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMenuAuthorityLogic {
	return &GetMenuAuthorityLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetMenuAuthorityLogic) GetMenuAuthority(in *core.RoleIdReq) (*core.MenuAuthResp, error) {
	ids, err := l.svcCtx.Deps.GetMenuAuthority(l.ctx, ctxdata.ClaimsFromCtx(l.ctx), utils.GetEffectiveId(in.RoleId, in.Id))
	if err != nil {
		return nil, xerr.RpcErr(err)
	}
	if ids == nil {
		ids = []int64{}
	}
	return &core.MenuAuthResp{MenuIds: ids}, nil
}
