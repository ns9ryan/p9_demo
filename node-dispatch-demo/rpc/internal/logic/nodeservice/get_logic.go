package nodeservicelogic

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"oa.98ent.com/p9/node-dispatch/rpc/internal/svc"
	"oa.98ent.com/p9/node-dispatch/rpc/pb/nodedispatchrpc/nodepb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLogic {
	return &GetLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// Get 获取节点
func (l *GetLogic) Get(in *nodepb.GetNodeRequest) (*nodepb.GetNodeResponse, error) {
	// 节点ID必须大于0
	if in.Id <= 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid node id")
	}

	// 获取节点
	data, err := l.svcCtx.DB.Node.Get(l.ctx, in.Id)
	if err != nil {
		l.Logger.Errorw("获取节点失败", logx.Field("error", err.Error()))
		return nil, err
	}

	// TODO Connection Manager 完成后获取节点真实在线状态
	online := false

	// 返回节点信息
	return &nodepb.GetNodeResponse{
		Node: toNodeInfo(data, online), // 节点信息
	}, nil
}
