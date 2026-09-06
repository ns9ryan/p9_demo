// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

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

type GameCategoryGetLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGameCategoryGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GameCategoryGetLogic {
	return &GameCategoryGetLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GameCategoryGetLogic) GameCategoryGet(req *types.GameCategoryGetReq) (resp *types.GameCategoryResp, err error) {
	l.Infof("[API GameCategoryGet] received req: id=%d", req.ID)

	if l.svcCtx == nil || l.svcCtx.GrpcClient == nil {
		l.Error("[API GameCategoryGet] gRPC client not available")
		return nil, fmt.Errorf("gRPC client not available")
	}

	grpcReq := &platformgame.GetGameCategoryRequest{
		Id: req.ID,
	}

	client := l.svcCtx.GrpcClient.GetPlatformGameServiceClient()
	l.Infof("[API GameCategoryGet] calling gRPC service: GetGameCategory with id=%d", grpcReq.Id)

	grpcResp, err := client.GetGameCategory(l.ctx, grpcReq)
	if err != nil {
		l.Errorf("[API GameCategoryGet] gRPC call failed: %v", err)
		return nil, err
	}

	l.Infof("[API GameCategoryGet] gRPC response received: grpcResp=%+v", grpcResp)

	if grpcResp == nil {
		l.Error("[API GameCategoryGet] gRPC response is nil")
		return nil, fmt.Errorf("gRPC response is nil")
	}

	l.Infof("[API GameCategoryGet] gRPC response: Code=%d, Message=%s, DataLen=%d",
		grpcResp.Code, grpcResp.Message, len(grpcResp.Data))

	if len(grpcResp.Data) == 0 {
		l.Errorf("[API GameCategoryGet] gRPC response data is empty, Code=%d, Message=%s",
			grpcResp.Code, grpcResp.Message)
		return nil, fmt.Errorf("gRPC response data is empty: %s", grpcResp.Message)
	}

	if grpcResp.Code != constant.CodeSuccess {
		l.Errorf("[API GameCategoryGet] gRPC error: Code=%d, Message=%s", grpcResp.Code, grpcResp.Message)
		return nil, fmt.Errorf("gRPC error: %s", grpcResp.Message)
	}

	item := grpcResp.Data[0]
	resp = &types.GameCategoryResp{
		ID:           item.Id,
		CategoryCode: item.Code,
		NameI18n:     item.Name,
		Status:       int16(item.Status),
	}

	l.Infof("[API GameCategoryGet] success: id=%d", item.Id)
	return resp, nil
}
