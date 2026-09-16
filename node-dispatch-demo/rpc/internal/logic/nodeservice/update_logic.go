package nodeservicelogic

import (
	"context"

	"oa.98ent.com/p9/node-dispatch/rpc/internal/svc"
	"oa.98ent.com/p9/node-dispatch/rpc/pb/nodedispatchrpc/nodepb"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLogic {
	return &UpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// Update 修改节点
func (l *UpdateLogic) Update(in *nodepb.UpdateNodeRequest) (*nodepb.UpdateNodeResponse, error) {
	// 节点ID必须大于0
	if in.Id <= 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid node id")
	}

	// 至少需要修改一个字段
	if in.Name == nil &&
		in.Status == nil &&
		in.Remark == nil {
		return nil, status.Error(codes.InvalidArgument, "no fields to update")
	}

	// 整理修改参数
	name := trimOptionalString(in.Name)
	remark := trimOptionalString(in.Remark)

	// 修改节点
	err := l.svcCtx.DB.Node.
		UpdateOneID(in.Id).
		SetNillableName(name).        // 节点名称
		SetNillableStatus(in.Status). // 节点状态: 1启用, 2停用
		SetNillableRemark(remark).    // 运维备注
		Exec(l.ctx)
	if err != nil {
		l.Logger.Errorw("修改节点失败", logx.Field("error", err.Error()))
		return nil, err
	}

	// 返回修改结果
	return &nodepb.UpdateNodeResponse{}, nil
}
