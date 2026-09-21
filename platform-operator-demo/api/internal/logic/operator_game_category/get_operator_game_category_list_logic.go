// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package operator_game_category

import (
	"context"
	"fmt"

	"oa.98ent.com/p9/core/common/i18n"
	"oa.98ent.com/p9/platform-game/rpc/pb/platform_game"
	"oa.98ent.com/p9/platform-operator/api/internal/svc"
	"oa.98ent.com/p9/platform-operator/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetOperatorGameCategoryListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取分站游戏分类列表
func NewGetOperatorGameCategoryListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOperatorGameCategoryListLogic {
	return &GetOperatorGameCategoryListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetOperatorGameCategoryListLogic) GetOperatorGameCategoryList(req *types.GetOperatorGameCategoryListRequest) (resp *types.GetOperatorGameCategoryListResponse, err error) {
	// 调用 RPC 获取分站游戏分类列表
	rpcResp, err := l.svcCtx.GameGrpcClient.GetOperatorGameCategoryServiceClient().GetOperatorGameCategoryList(l.ctx, &platform_game.GetOperatorGameCategoryListRequest{
		Page:         int32(req.Page),
		PageSize:     int32(req.PageSize),
		OpCode:       req.OpCode,
		CategoryCode: req.CategoryCode,
		CheckStatus:  req.CheckStatus,
	})
	if err != nil {
		l.Logger.Error("GetOperatorGameCategoryList error:", err)
		return nil, err
	}

	// 构建响应
	resp = &types.GetOperatorGameCategoryListResponse{
		Total:    rpcResp.Total,
		Page:     rpcResp.Page,
		PageSize: rpcResp.PageSize,
	}

	// 转换数据，並使用 corei18n.TG 获取 name
	for _, item := range rpcResp.Items {
		name := i18n.TG(l.ctx, i18n.CodePlatform, "game", fmt.Sprintf("game.category.%s.name", item.CategoryCode))
		resp.Items = append(resp.Items, types.OperatorGameCategoryInfo{
			Id:           item.Id,
			OpCode:       item.OpCode,
			CategoryCode: item.CategoryCode,
			Name:         name,
			Status:       item.Status,
			CheckStatus:  item.CheckStatus,
			CreatedAt:    item.CreatedAt,
			UpdatedAt:    item.UpdatedAt,
		})
	}

	return resp, nil
}
