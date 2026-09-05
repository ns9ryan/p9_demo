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

type CreateApiLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateApiLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateApiLogic {
	return &CreateApiLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateApiLogic) CreateApi(in *core.CreateApiReq) (*core.ApiInfo, error) {
	row, err := l.svcCtx.Deps.CreateAPI(l.ctx, service.CreateAPIReq{
		Description: in.Description, APIGroup: in.ApiGroup, Method: in.Method, Path: in.Path,
		IsRequired: int16(in.IsRequired), ServiceName: in.ServiceName,
	})
	if err != nil {
		return nil, xerr.RpcErr(err)
	}
	return logic.ToApiInfo(*row), nil
}
