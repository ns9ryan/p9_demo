package gameservicelogic

import (
	"context"
	"fmt"
	"strings"

	"oa.98ent.com/p9/platform-game/rpc/ent"
	"oa.98ent.com/p9/platform-game/rpc/internal/constant"
	"oa.98ent.com/p9/platform-game/rpc/internal/logic"
	"oa.98ent.com/p9/platform-game/rpc/internal/svc"
	"oa.98ent.com/p9/platform-game/rpc/internal/utils"
	platformgame "oa.98ent.com/p9/platform-game/rpc/pb/platform_game"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetGameListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetGameListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGameListLogic {
	return &GetGameListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取游戏列表
func (l *GetGameListLogic) GetGameList(in *platformgame.GetGameListRequest) (*platformgame.GetGameListResp, error) {
	l.Infof("[RPC GetGameList] received req: page=%d, page_size=%d", in.Page, in.PageSize)
	if l.svcCtx == nil || l.svcCtx.DB == nil {
		l.Errorf("[RPC GetGameList] Database not available")
		return &platformgame.GetGameListResp{
			Code:    constant.CodeInternalError,
			Message: "Database not available",
		}, nil
	}

	page, pageSize := utils.HandlePage(int64(in.Page), int64(in.PageSize))

	query := l.svcCtx.DB.WithContext(l.ctx)

	// 处理软删除条件
	if in.GetIsDeleted() == 0 {
		query = query.Where("deleted_at IS NULL")
	} else if in.GetIsDeleted() == 1 {
		query = query.Where("deleted_at IS NOT NULL")
	}

	if in.GetStatus() > 0 {
		query = query.Where("status = ?", in.GetStatus())
	}
	if in.GetGameCode() != "" {
		query = query.Where("game_code LIKE ?", "%"+in.GetGameCode()+"%")
	}
	if in.GetProviderId() > 0 {
		query = query.Where("provider_id = ?", in.GetProviderId())
	}
	if in.GetCategoryId() > 0 {
		query = query.Where("category_id = ?", in.GetCategoryId())
	}
	if in.GetChannelId() > 0 {
		query = query.Where("channel_id = ?", in.GetChannelId())
	}
	if in.GetName() != "" {
		query = query.Where("name_i18n->>'default' LIKE ?", "%"+in.GetName()+"%")
	}

	sortBy := strings.TrimSpace(in.SortBy)
	sortOrder := strings.TrimSpace(in.SortOrder)
	if sortBy == "" {
		sortBy = "sort_no"
	}
	if sortOrder == "" {
		sortOrder = "asc"
	}
	sortBy = strings.ToLower(sortBy)
	sortOrder = strings.ToLower(sortOrder)

	switch sortBy {
	case "id", "sort_no", "created_at":
	default:
		sortBy = "id"
	}

	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "asc"
	}

	var total int64
	if err := query.Model(&ent.Game{}).Count(&total).Error; err != nil {
		l.Errorf("[RPC GetGameList] count failed: %v", err)
		return &platformgame.GetGameListResp{
			Code:    constant.CodeInternalError,
			Message: fmt.Sprintf("count failed: %v", err),
		}, nil
	}

	var games []*ent.Game
	if err := query.
		Offset(int((page - 1) * pageSize)).
		Limit(int(pageSize)).
		Order(fmt.Sprintf("%s %s", sortBy, sortOrder)).
		Find(&games).Error; err != nil {
		l.Errorf("[RPC GetGameList] query failed: %v", err)
		return &platformgame.GetGameListResp{
			Code:    constant.CodeInternalError,
			Message: fmt.Sprintf("query failed: %v", err),
		}, nil
	}
	var gameProtoList []*platformgame.GameInfo
	for _, game := range games {
		ext, err := GetGameExtInfo(l.ctx, l.svcCtx, game)
		if err != nil {
			l.Errorf("[RPC GetGameList] query extended game info failed: %v", err)
			return &platformgame.GetGameListResp{
				Code:    constant.CodeInternalError,
				Message: fmt.Sprintf("query extended game info failed: %v", err),
			}, nil
		}
		gameProtoList = append(gameProtoList, logic.GameModelToProto(game, ext))
	}

	l.Infof("[RPC GetGameList] success: total=%d", total)
	return &platformgame.GetGameListResp{
		Code:     constant.CodeSuccess,
		Message:  "success",
		Data:     gameProtoList,
		Total:    total,
		Page:     int32(page),
		PageSize: int32(pageSize),
	}, nil
}
