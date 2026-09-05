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

type GetUserListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserListLogic {
	return &GetUserListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserListLogic) GetUserList(in *core.UserListReq) (*core.UserListResp, error) {
	list, total, err := l.svcCtx.Deps.ListUsers(l.ctx, ctxdata.ClaimsFromCtx(l.ctx), service.UserListReq{
		PageReq:     service.PageReq{Page: int(in.GetPage()), PageSize: int(in.GetPageSize())},
		Username:    in.GetUsername(),
		Mobile:      in.GetMobile(),
		Email:       in.GetEmail(),
		DisplayName: in.GetDisplayName(),
		RoleIDs:     in.GetRoleIds(),
	})
	if err != nil {
		return nil, xerr.RpcErr(err)
	}
	out := make([]*core.UserPublic, 0, len(list))
	for i := range list {
		p := service.PublicUser(&list[i], nil)
		out = append(out, logic.ToUserPublic(p))
	}
	return &core.UserListResp{List: out, Total: total}, nil
}
