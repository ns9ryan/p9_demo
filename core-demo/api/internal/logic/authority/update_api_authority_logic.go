package authority

import (
	"context"

	"oa.98ent.com/p9/core/api/internal/convert"
	"oa.98ent.com/p9/core/api/internal/svc"
	"oa.98ent.com/p9/core/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateApiAuthorityLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateApiAuthorityLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateApiAuthorityLogic {
	return &UpdateApiAuthorityLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateApiAuthorityLogic) UpdateApiAuthority(req *types.ApiAuthReq) (resp *types.BaseMsgResp, err error) {
	_, err = l.svcCtx.Core.UpdateApiAuthority(l.ctx, convert.ApiAuthReq(req))
	if err != nil {
		return nil, err
	}
	return &types.BaseMsgResp{Result: "success"}, nil
}
