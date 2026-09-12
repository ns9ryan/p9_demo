// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package execute

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/logx"
	"oa.98ent.com/p9/platform-game/api/internal/svc"
	"oa.98ent.com/p9/platform-game/api/internal/types"
	"oa.98ent.com/p9/platform-game/common/constant"
	platformgame "oa.98ent.com/p9/platform-game/rpc/pb/platform_game"
)

type ExecuteV2sqlGetLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewExecuteV2sqlGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ExecuteV2sqlGetLogic {
	return &ExecuteV2sqlGetLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ExecuteV2sqlGetLogic) ExecuteV2sqlGet() (resp *types.ExecuteV2SqlResp, err error) {
	if l.svcCtx == nil || l.svcCtx.GrpcClient == nil {
		l.Errorf("gRPC client not available")
		return nil, fmt.Errorf("gRPC client not available")
	}

	grpcReq := &platformgame.ExecuteV2SqlRequest{}
	grpcResp, err := l.svcCtx.GrpcClient.GetSyncServiceClient().ExecuteV2Sql(l.ctx, grpcReq)
	if err != nil {
		l.Errorf("ExecuteV2Sql grpc call failed: %v", err)
		return nil, err
	}
	if grpcResp == nil || grpcResp.Code != constant.CodeSuccess {
		msg := ""
		if grpcResp != nil {
			msg = grpcResp.Message
		}
		l.Errorf("ExecuteV2Sql failed: %s", msg)
		return &types.ExecuteV2SqlResp{Message: msg}, nil
	}

	return &types.ExecuteV2SqlResp{Message: grpcResp.Message}, nil
}
