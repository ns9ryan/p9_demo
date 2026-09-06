package logic

import (
	"context"

	"oa.98ent.com/p9/platform-game/common/constant"
	"oa.98ent.com/p9/platform-game/common/utils"
	"oa.98ent.com/p9/platform-game/pkg/game"
	"oa.98ent.com/p9/platform-game/rpc/internal/svc"
	platformgame "oa.98ent.com/p9/platform-game/rpc/pb/platform_game"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetGameChannelListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetGameChannelListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGameChannelListLogic {
	return &GetGameChannelListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取游戏渠道列表
func (l *GetGameChannelListLogic) GetGameChannelList(in *platformgame.GetGameChannelListRequest) (*platformgame.GetGameChannelListResp, error) {
	l.Infof("GetGameChannelList called with page=%d, page_size=%d", in.Page, in.PageSize)

	// 检查数据库连接
	if l.svcCtx.DB == nil {
		l.Errorf("Database not available")
		return &platformgame.GetGameChannelListResp{
			Code:    constant.CodeInternalError,
			Message: "Database not available",
		}, nil
	}

	// 构建查询参数
	page := int64(in.Page)
	if page <= 0 {
		page = 1
	}
	pageSize := int64(in.PageSize)
	if pageSize <= 0 {
		pageSize = 15
	}

	params := &game.GameChannelListParams{
		Page:        page,
		PageSize:    pageSize,
		ChannelCode: in.GetChannelCode(),
		Name:        in.GetName(),
		Status:      int16(in.Status),
		IsDeleted:   int8(in.IsDeleted),
		SortBy:      in.GetSortBy(),
		SortOrder:   in.GetSortOrder(),
	}

	// 调用数据库查询函数
	result, err := game.GameChannelListDB(l.ctx, l.svcCtx.DB, params)
	if err != nil {
		l.Errorf("GameChannelListDB failed: %v", err)
		return &platformgame.GetGameChannelListResp{
			Code:    constant.CodeInternalError,
			Message: "Failed to fetch channels: " + err.Error(),
		}, nil
	}

	// 转换结果为 protobuf 消息
	items := make([]*platformgame.GameChannelResp, 0, len(result.Channels))
	for _, item := range result.Channels {
		if item == nil || item.Channel == nil {
			continue
		}
		channel := item.Channel
		catIDs := make([]*platformgame.GameCategoryInfo, 0, len(item.Categories))
		for _, category := range item.Categories {
			catIDs = append(catIDs, &platformgame.GameCategoryInfo{
				Id:       category.ID,
				NameI18N: category.NameI18n,
			})
		}

		items = append(items, &platformgame.GameChannelResp{
			Id:             channel.ID,
			SourceId:       channel.SourceID,
			ChannelCode:    channel.ChannelCode,
			NameI18N:       string(utils.MustMarshalJSON(channel.NameI18n)),
			SourceNameI18N: string(utils.MustMarshalJSON(channel.SourceNameI18n)),
			Status:         int32(channel.Status),
			VendorCount:    item.VendorCount,
			CatIds:         catIDs,
			GameCount:      item.GameCount,
			CreatedAt:      channel.CreatedAt.Unix(),
			UpdatedAt:      channel.UpdatedAt.Unix(),
		})
	}

	return &platformgame.GetGameChannelListResp{
		Code:     constant.CodeSuccess,
		Message:  "success",
		Items:    items,
		Total:    result.Total,
		Page:     int32(page),
		PageSize: int32(pageSize),
	}, nil
}
