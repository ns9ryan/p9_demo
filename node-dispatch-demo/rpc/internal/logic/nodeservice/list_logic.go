package nodeservicelogic

import (
	"context"
	"strings"

	"oa.98ent.com/p9/node-dispatch/rpc/ent/node"
	"oa.98ent.com/p9/node-dispatch/rpc/internal/svc"
	"oa.98ent.com/p9/node-dispatch/rpc/pb/nodedispatchrpc/nodepb"

	"entgo.io/ent/dialect/sql"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	// 校验分页参数
	if in.Page < 1 || in.PageSize < 1 || in.PageSize > 100 {
		return nil, status.Error(codes.InvalidArgument, "invalid pagination")
	}

	// 校验节点状态
	if in.Status != nil && (*in.Status < 1 || *in.Status > 2) {
		return nil, status.Error(codes.InvalidArgument, "invalid node status")
	}

	// 获取当前在线节点快照
	onlineCodes := l.svcCtx.Connections.OnlineCodes()
	onlineSet := make(map[string]struct{}, len(onlineCodes))
	for _, nodeCode := range onlineCodes {
		onlineSet[nodeCode] = struct{}{}
	}

	// 创建节点查询
	query := l.svcCtx.DB.Node.Query()

	// 按关键字筛选
	if in.Keyword != nil {
		keyword := strings.TrimSpace(*in.Keyword)
		if keyword != "" {
			query = query.Where(
				node.Or(
					node.CodeContainsFold(keyword),
					node.NameContainsFold(keyword),
				),
			)
		}
	}

	// 按节点状态筛选
	if in.Status != nil {
		query = query.Where(node.StatusEQ(*in.Status))
	}

	// 按在线状态筛选
	if in.Online != nil {
		if *in.Online {
			// 当前没有在线节点时直接返回空列表
			if len(onlineCodes) == 0 {
				return &nodepb.ListNodesResponse{
					Total: 0,
					List:  []*nodepb.NodeInfo{},
				}, nil
			}

			query = query.Where(node.CodeIn(onlineCodes...))
		} else if len(onlineCodes) > 0 {
			query = query.Where(node.CodeNotIn(onlineCodes...))
		}
	}

	// 获取符合条件的数据总数
	total, err := query.Clone().Count(l.ctx)
	if err != nil {
		l.Logger.Errorw("获取节点总数失败", logx.Field("error", err.Error()))
		return nil, err
	}

	// 计算分页偏移量
	offset := (in.Page - 1) * in.PageSize

	// 获取当前页节点数据
	results, err := query.
		Order(
			node.ByCreatedAt(sql.OrderDesc()), // 按创建时间倒序
			node.ByID(sql.OrderDesc()),        // 创建时间相同时按ID倒序
		).
		Offset(int(offset)).
		Limit(int(in.PageSize)).
		All(l.ctx)
	if err != nil {
		l.Logger.Errorw("获取节点列表失败", logx.Field("error", err.Error()))
		return nil, err
	}

	// 转换节点列表
	list := make([]*nodepb.NodeInfo, 0, len(results))
	for _, result := range results {
		_, online := onlineSet[result.Code]
		list = append(list, toNodeInfo(result, online))
	}

	// 返回节点列表
	return &nodepb.ListNodesResponse{
		Total: int64(total), // 数据总数
		List:  list,         // 节点列表
	}, nil
}
