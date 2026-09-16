package nodeservicelogic

import (
	"context"

	"oa.98ent.com/p9/node-dispatch/rpc/internal/svc"
	"oa.98ent.com/p9/node-dispatch/rpc/pb/nodedispatchrpc/nodepb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListLogic {
	return &ListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// List 获取节点列表
func (l *ListLogic) List(in *nodepb.ListNodesRequest) (*nodepb.ListNodesResponse, error) {
	// TODO Connection Manager 完成后实现节点在线状态获取及 online 筛选
	// online 筛选必须在数据库分页前转换为节点编码条件, 保证 total 和分页结果正确

	return &nodepb.ListNodesResponse{}, nil
}
