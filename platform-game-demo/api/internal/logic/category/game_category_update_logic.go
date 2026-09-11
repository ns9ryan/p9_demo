// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package category

import (
	"context"
	"fmt"

	"oa.98ent.com/p9/platform-game/api/internal/svc"
	"oa.98ent.com/p9/platform-game/api/internal/types"
	"oa.98ent.com/p9/platform-game/common/constant"
	platformgame "oa.98ent.com/p9/platform-game/rpc/pb/platform_game"

	"github.com/bytedance/gopkg/util/logger"
	"github.com/zeromicro/go-zero/core/logx"
)

type GameCategoryUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGameCategoryUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GameCategoryUpdateLogic {
	return &GameCategoryUpdateLogic{
		Logger: logx.WithContext(ctx),
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
	grpcResp, err := l.svcCtx.GrpcClient.GetGameCategoryServiceClient().UpdateGameCategory(l.ctx, grpcReq)
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
	if grpcResp.Data == nil {
		logger.Error("[API GameCategoryUpdate] gRPC response data is empty")
		return nil, fmt.Errorf("gRPC response data is empty")
	}

	resp = &types.GameCategoryResp{
		ID:           grpcResp.Data.Id,
		CategoryCode: grpcResp.Data.CategoryCode,
		NameI18n:     grpcResp.Data.NameI18N,
		Status:       int16(grpcResp.Data.Status),
	}

	logger.Infof("[API GameCategoryUpdate] success: id=%d", req.ID)
	return resp, nil
}
