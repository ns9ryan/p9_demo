package nodeservicelogic

import (
	"context"

	"oa.98ent.com/p9/node-dispatch/rpc/internal/svc"
	"oa.98ent.com/p9/node-dispatch/rpc/pb/nodedispatchrpc/nodepb"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

// ResetAuthSecret 重置节点认证密钥
func (l *ResetAuthSecretLogic) ResetAuthSecret(in *nodepb.ResetNodeAuthSecretRequest) (*nodepb.ResetNodeAuthSecretResponse, error) {
	// 节点ID必须大于0
	if in.Id <= 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid node id")
	}

	// 生成新的节点认证密钥
	authSecret, authSecretHash, err := generateAuthSecret()
	if err != nil {
		l.Logger.Errorw("生成节点认证密钥失败", logx.Field("error", err.Error()))
		return nil, err
	}

	// 更新节点认证密钥哈希
	err = l.svcCtx.DB.Node.
		UpdateOneID(in.Id).
		SetAuthSecretHash(authSecretHash). // 节点认证密钥哈希
		Exec(l.ctx)
	if err != nil {
		l.Logger.Errorw("重置节点认证密钥失败", logx.Field("error", err.Error()))
		return nil, err
	}

	// TODO Connection Manager 完成后主动断开当前节点连接, 强制使用新密钥重新认证

	// 返回新的节点认证密钥
	return &nodepb.ResetNodeAuthSecretResponse{
		AuthSecret: authSecret, // 新节点认证密钥
	}, nil
}
