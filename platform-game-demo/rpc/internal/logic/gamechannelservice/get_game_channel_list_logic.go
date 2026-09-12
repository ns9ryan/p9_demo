package gamechannelservicelogic

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
	l.Infof("[RPC GetGameChannelList] received req: page=%d, page_size=%d, sortBy=%s, sortOrder=%s", in.Page, in.PageSize, in.SortBy, in.SortOrder)

	if l.svcCtx == nil || l.svcCtx.DB == nil {
		l.Error("[RPC GetGameChannelList] database not available")
		return &platformgame.GetGameChannelListResp{
			Code:    constant.CodeInternalError,
			Message: "database not available",
		}, nil
	}

	page, pageSize := utils.HandlePage(int64(in.Page), int64(in.PageSize))

	query := l.svcCtx.DB.WithContext(l.ctx)

	// 处理软删除条件
	if in.GetIsDeleted() == 0 {
		query = query.Where("deleted_at IS NULL")
		l.Infof("[RPC GetGameChannelList] filter: is_deleted=0 (not deleted)")
	} else if in.GetIsDeleted() == 1 {
		query = query.Where("deleted_at IS NOT NULL")
		l.Infof("[RPC GetGameChannelList] filter: is_deleted=1 (deleted)")
	}

	if in.GetStatus() > 0 {
		query = query.Where("status = ?", in.GetStatus())
		l.Infof("[RPC GetGameChannelList] filter: status=%d", in.GetStatus())
	}
	if in.GetChannelCode() != "" {
		query = query.Where("channel_code LIKE ?", "%"+in.GetChannelCode()+"%")
		l.Infof("[RPC GetGameChannelList] filter: channel_code=%s", in.GetChannelCode())
	}
	if in.GetName() != "" {
		query = query.Where("name_i18n->>'default' LIKE ?", "%"+in.GetName()+"%")
		l.Infof("[RPC GetGameChannelList] filter: name=%s", in.GetName())
	}

	sortBy := strings.TrimSpace(in.SortBy)
	sortOrder := strings.TrimSpace(in.SortOrder)
	l.Infof("[RPC GetGameChannelList] before default: sortBy=%s, sortOrder=%s", sortBy, sortOrder)

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
	l.Infof("[RPC GetGameChannelList] after default: sortBy=%s, sortOrder=%s", sortBy, sortOrder)

	var total int64
	if err := query.Model(&ent.GameChannel{}).Count(&total).Error; err != nil {
		l.Errorf("[RPC GetGameChannelList] count failed: %v", err)
		return &platformgame.GetGameChannelListResp{
			Code:    constant.CodeInternalError,
			Message: fmt.Sprintf("count failed: %v", err),
		}, nil
	}

	var channels []*ent.GameChannel
	orderSQL := fmt.Sprintf("%s %s", sortBy, sortOrder)
	l.Infof("[RPC GetGameChannelList] executing query: offset=%d, limit=%d, order=%s", (page-1)*pageSize, pageSize, orderSQL)

	if err := query.
		Offset(int((page - 1) * pageSize)).
		Limit(int(pageSize)).
		Order(orderSQL).
		Find(&channels).Error; err != nil {
		l.Errorf("[RPC GetGameChannelList] query failed: %v", err)
		return &platformgame.GetGameChannelListResp{
			Code:    constant.CodeInternalError,
			Message: fmt.Sprintf("query failed: %v", err),
		}, nil
	}
	items := logic.ChannelModelToProtoList(channels)
	l.Infof("[RPC GetGameChannelList] query result: returned=%d items, total=%d", len(items), total)
	if len(items) > 0 {
		l.Infof("[RPC GetGameChannelList] first item: id=%d, sortNo=%d, name=%v", items[0].Id, items[0].SortNo, items[0].NameI18N)
		if len(items) > 1 {
			l.Infof("[RPC GetGameChannelList] second item: id=%d, sortNo=%d, name=%v", items[1].Id, items[1].SortNo, items[1].NameI18N)
		}
	}
	return &platformgame.GetGameChannelListResp{
		Code:     constant.CodeSuccess,
		Message:  "success",
		Items:    items,
		Total:    total,
		Page:     int32(page),
		PageSize: int32(pageSize),
	}, nil
}
