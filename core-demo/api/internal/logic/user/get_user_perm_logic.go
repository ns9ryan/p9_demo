package user

import (
	"context"

	"oa.98ent.com/p9/core/api/internal/convert"
	"oa.98ent.com/p9/core/api/internal/svc"
	"oa.98ent.com/p9/core/api/internal/types"
	"oa.98ent.com/p9/core/rpc/coreclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserPermLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserPermLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserPermLogic {
	return &GetUserPermLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserPermLogic) GetUserPerm() (resp *types.PermResp, err error) {
	out, err := l.svcCtx.Core.GetUserPerm(l.ctx, &coreclient.Empty{})
	if err != nil {
		return nil, err
	}
	return convert.PermResp(out), nil
}
