package public

import (
	"context"

	"oa.98ent.com/p9/core/api/internal/convert"
	"oa.98ent.com/p9/core/api/internal/svc"
	"oa.98ent.com/p9/core/api/internal/types"
	"oa.98ent.com/p9/core/rpc/coreclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type IssuePreviewTokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewIssuePreviewTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IssuePreviewTokenLogic {
	return &IssuePreviewTokenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *IssuePreviewTokenLogic) IssuePreviewToken(req *types.IssuePreviewTokenReq) (resp *types.IssuePreviewTokenResp, err error) {
	out, err := l.svcCtx.Core.IssuePreviewToken(l.ctx, &coreclient.IssuePreviewTokenReq{
		OperatorCode: req.OperatorCode,
	})
	if err != nil {
		return nil, err
	}
	return convert.IssuePreviewTokenResp(out), nil
}
