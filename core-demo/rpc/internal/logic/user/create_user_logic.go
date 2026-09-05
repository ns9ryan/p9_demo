package user

import (
	"context"

	"oa.98ent.com/p9/core/common/ctxdata"
	"oa.98ent.com/p9/core/common/xerr"
	"oa.98ent.com/p9/core/rpc/internal/logic"
	"oa.98ent.com/p9/core/rpc/internal/svc"
	"oa.98ent.com/p9/core/rpc/service"
	"oa.98ent.com/p9/core/rpc/types/core"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateUserLogic {
	return &CreateUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateUserLogic) CreateUser(in *core.CreateUserReq) (*core.UserPublic, error) {
	u, err := l.svcCtx.Deps.CreateUser(l.ctx, ctxdata.ClaimsFromCtx(l.ctx), service.CreateUserReq{
		Username: in.Username, Password: in.Password, DisplayName: in.DisplayName,
		Mobile: in.Mobile, Email: in.Email, Status: int16(in.Status), RoleIDs: in.RoleIds,
	})
	if err != nil {
		return nil, xerr.RpcErr(err)
	}
	return logic.ToUserPublic(service.PublicUser(u, nil)), nil
}
