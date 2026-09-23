package operatorservicelogic

import (
	"context"

	"oa.98ent.com/p9/operator-base/rpc/internal/svc"
	"oa.98ent.com/p9/operator-base/rpc/pb/operatorbaserpc/operatorpb"

	"github.com/zeromicro/go-zero/core/logx"
)

type InitializeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewInitializeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InitializeLogic {
	return &InitializeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 初始化厅 operator
func (l *InitializeLogic) Initialize(in *operatorpb.InitializeOperatorRequest) (*operatorpb.InitializeOperatorResponse, error) {
	// todo: add your logic here and delete this line

	return &operatorpb.InitializeOperatorResponse{}, nil
}
