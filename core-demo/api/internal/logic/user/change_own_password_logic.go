package user

import (
	"context"

	"oa.98ent.com/p9/core/api/internal/svc"
	"oa.98ent.com/p9/core/api/internal/types"
	"oa.98ent.com/p9/core/rpc/coreclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChangeOwnPasswordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewChangeOwnPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChangeOwnPasswordLogic {
	return &ChangeOwnPasswordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ChangeOwnPasswordLogic) ChangeOwnPassword(req *types.SelfPasswordReq) (resp *types.BaseMsgResp, err error) {
	_, err = l.svcCtx.Core.ChangeOwnPassword(l.ctx, &coreclient.SelfPasswordReq{OldPassword: req.OldPassword, Password: req.Password})
	if err != nil {
		return nil, err
	}
	return &types.BaseMsgResp{Result: "success"}, nil
}
