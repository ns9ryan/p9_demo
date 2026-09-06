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

type GameChannelListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGameChannelListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GameChannelListLogic {
	return &GameChannelListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GameChannelListLogic) GameChannelList(req *types.GameChannelListReq) (resp *types.GameChannelListResp, err error) {
	l.Infof("[API GameChannelList] received req: page=%d, page_size=%d, channel_code=%s, name=%s, status=%d, is_deleted=%d",
		req.Page, req.PageSize, req.ChannelCode, req.Name, req.Status, req.IsDeleted)

	// 检查 gRPC 客户端是否可用
	if l.svcCtx == nil || l.svcCtx.GrpcClient == nil {
		l.Error("[API GameChannelList] gRPC client not available")
		return nil, fmt.Errorf("gRPC client not available")
	}

	// 构建 gRPC 请求（类型转换）
	grpcReq := &platformgame.GetGameChannelListRequest{
		Page:        int32(req.Page),
		PageSize:    int32(req.PageSize),
		ChannelCode: req.ChannelCode,
		Name:        req.Name,
		Status:      int32(req.Status),
		IsDeleted:   int32(req.IsDeleted),
	}

	// 调用 RPC 服务
	client := l.svcCtx.GrpcClient.GetPlatformGameServiceClient()
	grpcResp, err := client.GetGameChannelList(l.ctx, grpcReq)
	if err != nil {
		l.Errorf("[API GameChannelList] gRPC call failed: %v", err)
		return nil, err
	}

	// 检查 gRPC 响应
	if grpcResp == nil {
		l.Error("[API GameChannelList] gRPC response is nil")
		return nil, fmt.Errorf("gRPC response is nil")
	}

	// 转换响应
	if grpcResp.Code != constant.CodeSuccess {
		l.Errorf("[API GameChannelList] gRPC error: Code=%d, Message=%s", grpcResp.Code, grpcResp.Message)
		return nil, fmt.Errorf("gRPC error: %s", grpcResp.Message)
	}

	items := make([]types.GameChannelResp, 0, len(grpcResp.Items))
	for _, item := range grpcResp.Items {
		if item == nil {
			continue
		}

		catIds := make([]types.GameCategoryInfo, 0, len(item.CatIds))
		for _, catId := range item.CatIds {
			if catId != nil {
				catIds = append(catIds, types.GameCategoryInfo{
					ID:       catId.Id,
					NameI18n: catId.NameI18N, // 注意字段名大小写
				})
			}
		}

		items = append(items, types.GameChannelResp{
			ID:             item.Id,
			SourceID:       item.SourceId,
			ChannelCode:    item.ChannelCode,
			NameI18n:       item.NameI18N,
			SourceNameI18n: item.SourceNameI18N, // 注意字段名大小写
			Status:         int16(item.Status),
			VendorCount:    item.VendorCount,
			CatIds:         catIds,
			GameCount:      item.GameCount,
			CreatedAt:      item.CreatedAt,
			UpdatedAt:      item.UpdatedAt,
		})
	}

	resp = &types.GameChannelListResp{
		Items: items,
		Total: grpcResp.Total,
	}

	l.Infof("[API GameChannelList] success: total=%d, returned=%d", grpcResp.Total, len(items))
	return resp, nil
}
