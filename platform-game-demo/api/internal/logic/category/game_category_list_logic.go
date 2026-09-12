// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package category

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

type GameCategoryListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGameCategoryListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GameCategoryListLogic {
	return &GameCategoryListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GameCategoryListLogic) GameCategoryList(req *types.GameCategoryListReq) (resp *types.GameCategoryListResp, err error) {
	l.Infof("[API GameCategoryList] received req: page=%d, page_size=%d", req.Page, req.PageSize)

	if l.svcCtx == nil || l.svcCtx.GrpcClient == nil {
		l.Error("[API GameCategoryList] gRPC client not available")
		return nil, fmt.Errorf("gRPC client not available")
	}

	grpcReq := &platformgame.GetGameCategoryListRequest{
		Page:         int32(req.Page),
		PageSize:     int32(req.PageSize),
		CategoryCode: req.CategoryCode,
		Name:         req.Name,
		Status:       int32(req.Status),
		IsDeleted:    int32(req.IsDeleted),
		SortBy:       req.SortBy,
		SortOrder:    req.SortOrder,
	}

	client := l.svcCtx.GrpcClient.GetGameCategoryServiceClient()
	grpcResp, err := client.GetGameCategoryList(l.ctx, grpcReq)
	if err != nil {
		l.Errorf("[API GameCategoryList] gRPC call failed: %v", err)
		return nil, err
	}

	if grpcResp == nil {
		l.Error("[API GameCategoryList] gRPC response is nil")
		return nil, fmt.Errorf("gRPC response is nil")
	}

	if grpcResp.Code != constant.CodeSuccess {
		l.Errorf("[API GameCategoryList] gRPC error: Code=%d, Message=%s", grpcResp.Code, grpcResp.Message)
		return nil, fmt.Errorf("gRPC error: %s", grpcResp.Message)
	}

	items := make([]types.GameCategoryResp, 0, len(grpcResp.Data))
	for _, item := range grpcResp.Data {
		items = append(items, *logic.CategoryProtoToResponse(item))
	}

	resp = &types.GameCategoryListResp{
		List:  items,
		Total: grpcResp.Total,
	}

	l.Infof("[API GameCategoryList] success: total=%d", grpcResp.Total)
	return resp, nil
}
