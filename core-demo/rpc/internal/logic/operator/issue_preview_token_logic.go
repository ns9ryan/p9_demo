package operator

import (
	"context"

	"oa.98ent.com/p9/core/common/xerr"
	"oa.98ent.com/p9/core/rpc/internal/svc"
	"oa.98ent.com/p9/core/rpc/service"
	"oa.98ent.com/p9/core/rpc/types/core"

	"github.com/zeromicro/go-zero/core/logx"
)

type IssuePreviewTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewIssuePreviewTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IssuePreviewTokenLogic {
	return &IssuePreviewTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *IssuePreviewTokenLogic) IssuePreviewToken(in *core.IssuePreviewTokenReq) (*core.IssuePreviewTokenResp, error) {
	res, err := l.svcCtx.Deps.IssuePreviewToken(l.ctx, service.IssuePreviewTokenReq{
		OperatorCode: in.OperatorCode,
	})
	if err != nil {
		return nil, xerr.RpcErr(err)
	}
	return &core.IssuePreviewTokenResp{
		AccessToken:  res.AccessToken,
		Expire:       res.Expire,
		OperatorCode: res.OperatorCode,
		HomePath:     res.HomePath,
	}, nil
}
