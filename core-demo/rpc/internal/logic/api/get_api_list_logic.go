package api

import (
	"context"

	"oa.98ent.com/p9/core/common/xerr"
	"oa.98ent.com/p9/core/rpc/internal/logic"
	"oa.98ent.com/p9/core/rpc/internal/svc"
	"oa.98ent.com/p9/core/rpc/service"
	"oa.98ent.com/p9/core/rpc/types/core"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetApiListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetApiListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetApiListLogic {
	return &GetApiListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetApiListLogic) GetApiList(in *core.ApiListReq) (*core.ApiListResp, error) {
	list, total, err := l.svcCtx.Deps.ListAPIs(l.ctx, service.APIListReq{
		PageReq:     service.PageReq{Page: int(in.GetPage()), PageSize: int(in.GetPageSize())},
		Path:        in.GetPath(),
		Method:      in.GetMethod(),
		APIGroup:    in.GetApiGroup(),
		ServiceName: in.GetServiceName(),
		Description: in.GetDescription(),
	})
	if err != nil {
		return nil, xerr.RpcErr(err)
	}
	out := make([]*core.ApiInfo, 0, len(list))
	for _, a := range list {
		out = append(out, logic.ToApiInfo(a))
	}
	return &core.ApiListResp{List: out, Total: total}, nil
}
