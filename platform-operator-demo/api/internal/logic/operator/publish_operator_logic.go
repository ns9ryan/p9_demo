// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package operator

import (
	"context"

	"oa.98ent.com/p9/platform-operator/api/internal/svc"
	"oa.98ent.com/p9/platform-operator/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PublishOperatorLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPublishOperatorLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublishOperatorLogic {
	return &PublishOperatorLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PublishOperatorLogic) PublishOperator(req *types.PublishOperatorRequest) (resp *types.PublishOperatorResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
