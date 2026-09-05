package user

import (
	"context"

	"oa.98ent.com/p9/core/api/internal/svc"
	"oa.98ent.com/p9/core/api/internal/types"
	"oa.98ent.com/p9/core/rpc/coreclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type LogoutLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLogoutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LogoutLogic {
	return &LogoutLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LogoutLogic) Logout(req *types.LogoutReq) (resp *types.BaseMsgResp, err error) {
	_, err = l.svcCtx.Core.Logout(l.ctx, &coreclient.LogoutReq{RefreshToken: req.RefreshToken})
	if err != nil {
		return nil, err
	}
	return &types.BaseMsgResp{
		Result: "success",
	}, nil
}
