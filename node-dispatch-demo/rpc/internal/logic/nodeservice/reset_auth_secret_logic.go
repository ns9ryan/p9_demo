package nodeservicelogic

import (
	"context"

	"oa.98ent.com/p9/node-dispatch/rpc/internal/svc"
	"oa.98ent.com/p9/node-dispatch/rpc/pb/nodedispatchrpc/nodepb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ResetAuthSecretLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewResetAuthSecretLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResetAuthSecretLogic {
	return &ResetAuthSecretLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 重置节点认证密钥
func (l *ResetAuthSecretLogic) ResetAuthSecret(in *nodepb.ResetNodeAuthSecretRequest) (*nodepb.ResetNodeAuthSecretResponse, error) {
	// todo: add your logic here and delete this line

	return &nodepb.ResetNodeAuthSecretResponse{}, nil
}
