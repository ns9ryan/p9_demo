// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package channel

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/logx"
	"oa.98ent.com/p9/platform-game/api/internal/logic"
	"oa.98ent.com/p9/platform-game/api/internal/svc"
	"oa.98ent.com/p9/platform-game/api/internal/types"
	"oa.98ent.com/p9/platform-game/common/constant"
	platformgame "oa.98ent.com/p9/platform-game/rpc/pb/platform_game"
)

type GameChannelGetLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGameChannelGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GameChannelGetLogic {
	return &GameChannelGetLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GameChannelGetLogic) GameChannelGet(req *types.GameChannelGetReq) (resp *types.GameChannelResp, err error) {
	l.Infof("[API GameChannelGet] received req: id=%d", req.ID)

	// 检查 gRPC 客户端是否可用
	if l.svcCtx == nil || l.svcCtx.GrpcClient == nil {
		l.Error("[API GameChannelGet] gRPC client not available")
		return nil, fmt.Errorf("gRPC client not available")
	}

	// 构建 gRPC 请求
	grpcReq := &platformgame.GetGameChannelRequest{
		Id: req.ID,
	}

	// 调用 RPC 服务
	client := l.svcCtx.GrpcClient.GetGameChannelServiceClient()
	grpcResp, err := client.GetGameChannel(l.ctx, grpcReq)
	if err != nil {
		l.Errorf("[API GameChannelGet] gRPC call failed: %v", err)
		return nil, err
	}

	// 检�?gRPC 响应
	if grpcResp == nil || grpcResp.Data == nil {
		l.Error("[API GameChannelGet] gRPC response is nil")
		return nil, fmt.Errorf("gRPC response is nil")
	}

	// 转换响应
	if grpcResp.Code != constant.CodeSuccess {
		l.Errorf("[API GameChannelGet] gRPC error: Code=%d, Message=%s", grpcResp.Code, grpcResp.Message)
		return nil, fmt.Errorf("gRPC error: %s", grpcResp.Message)
	}

	item := grpcResp.Data
	resp = logic.ChannelProtoToResponse(item)

	l.Infof("[API GameChannelGet] success: id=%d", item.Id)
	return resp, nil
}
