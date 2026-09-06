// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package logic

import (
	"context"
	"fmt"

	"oa.98ent.com/p9/platform-game/api/internal/svc"
	"oa.98ent.com/p9/platform-game/api/internal/types"
	"oa.98ent.com/p9/platform-game/common/constant"
	"oa.98ent.com/p9/platform-game/common/logger"
	platformgame "oa.98ent.com/p9/platform-game/rpc/pb/platform_game"
)

type GameCategoryUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGameCategoryUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GameCategoryUpdateLogic {
	return &GameCategoryUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GameCategoryUpdateLogic) GameCategoryUpdate(req *types.GameCategoryUpdateReq) (resp *types.GameCategoryResp, err error) {
	// API 通过 gRPC 调用 RPC 服务进行数据库操作
	if l.svcCtx == nil || l.svcCtx.GrpcClient == nil {
		logger.Error("[API GameCategoryUpdate] gRPC client not available")
		return nil, fmt.Errorf("gRPC client not available")
	}

	// 构建 gRPC 请求
	grpcReq := &platformgame.UpdateGameCategoryRequest{
		Id:          req.ID,
		NameI18N:    req.NameI18n,
		SortNo:      req.SortNo,
		Status:      int32(req.Status),
		ForceLogout: req.ForceLogout,
	}

	// 调用 RPC 服务
	grpcResp, err := l.svcCtx.GrpcClient.GetPlatformGameServiceClient().UpdateGameCategory(l.ctx, grpcReq)
	if err != nil {
		logger.Errorf("[API GameCategoryUpdate] gRPC call failed: %v", err)
		return nil, fmt.Errorf("gRPC call failed: %s", err.Error())
	}

	if grpcResp == nil {
		logger.Error("[API GameCategoryUpdate] gRPC response is nil")
		return nil, fmt.Errorf("gRPC response is nil")
	}

	// 检查返回码
	if grpcResp.Code != constant.CodeSuccess {
		logger.Errorf("[API GameCategoryUpdate] gRPC error: Code=%d, Message=%s", grpcResp.Code, grpcResp.Message)
		return nil, fmt.Errorf("gRPC error: %s", grpcResp.Message)
	}

	// 转换响应数据
	if len(grpcResp.Data) == 0 {
		logger.Error("[API GameCategoryUpdate] gRPC response data is empty")
		return nil, fmt.Errorf("gRPC response data is empty")
	}

	categoryInfo := grpcResp.Data[0]
	resp = &types.GameCategoryResp{
		ID:           categoryInfo.Id,
		CategoryCode: categoryInfo.Code,
		NameI18n:     categoryInfo.Name,
		Status:       int16(categoryInfo.Status),
	}

	logger.Infof("[API GameCategoryUpdate] success: id=%d", req.ID)
	return resp, nil
}
