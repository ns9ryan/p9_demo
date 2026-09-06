package logic

import (
	"context"
	"fmt"

	"oa.98ent.com/p9/platform-game/api/internal/svc"
	"oa.98ent.com/p9/platform-game/api/internal/types"
	"oa.98ent.com/p9/platform-game/common/constant"
	platformgame "oa.98ent.com/p9/platform-game/rpc/pb/platform_game"

	"github.com/zeromicro/go-zero/core/logx"
)

type ExecuteV2SqlLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewExecuteV2SqlLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ExecuteV2SqlLogic {
	return &ExecuteV2SqlLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ExecuteV2SqlLogic) ExecuteV2Sql(req *types.ExecuteV2SqlReq) (*types.ExecuteV2SqlResp, error) {
	if l.svcCtx == nil || l.svcCtx.GrpcClient == nil {
		l.Errorf("gRPC client not available")
		return nil, fmt.Errorf("gRPC client not available")
	}

	grpcReq := &platformgame.ExecuteV2SqlRequest{}
	grpcResp, err := l.svcCtx.GrpcClient.GetPlatformGameServiceClient().ExecuteV2Sql(l.ctx, grpcReq)
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
