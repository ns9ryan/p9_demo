package gameproviderservicelogic

import (
	"context"

	"fmt"
	"strings"

	"gorm.io/gorm"
	"oa.98ent.com/p9/platform-game/rpc/ent"
	"oa.98ent.com/p9/platform-game/rpc/internal/constant"
	"oa.98ent.com/p9/platform-game/rpc/internal/logic"
	"oa.98ent.com/p9/platform-game/rpc/internal/svc"
	"oa.98ent.com/p9/platform-game/rpc/internal/utils"
	platformgame "oa.98ent.com/p9/platform-game/rpc/pb/platform_game"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetGameProviderListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetGameProviderListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGameProviderListLogic {
	return &GetGameProviderListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取游戏供应商列表
func (l *GetGameProviderListLogic) GetGameProviderList(in *platformgame.GetGameProviderListRequest) (*platformgame.GetGameProviderListResp, error) {
	l.Infof("GetGameProviderList called with page=%d, page_size=%d", in.Page, in.PageSize)
	if l.svcCtx.DB == nil {
		l.Errorf("Database not available")
		return &platformgame.GetGameProviderListResp{
			Code:    constant.CodeInternalError,
			Message: "Database not available",
		}, nil
	}

	// 调用数据库查询函数
	providers, total, err := gameProviderListDB(l.ctx, l.svcCtx.DB, in)
	if err != nil {
		l.Errorf("GameProviderListDB failed: %v", err)
		return &platformgame.GetGameProviderListResp{
			Code:    constant.CodeInternalError,
			Message: "Failed to fetch providers: " + err.Error(),
		}, nil
	}

	return &platformgame.GetGameProviderListResp{
		Code:    constant.CodeSuccess,
		Message: "success",
		Data:    logic.ProviderModelToProtoList(providers),
		Total:   total,
	}, nil
}

func gameProviderListDB(ctx context.Context, db *gorm.DB, params *platformgame.GetGameProviderListRequest) ([]*ent.GameProvider, int64, error) {
	if db == nil {
		return nil, 0, fmt.Errorf("database not available")
	}

	page, pageSize := utils.HandlePage(int64(params.Page), int64(params.PageSize))

	query := db.WithContext(ctx)

	// 处理软删除条件
	if params.GetIsDeleted() == 0 {
		query = query.Where("deleted_at IS NULL")
	} else if params.GetIsDeleted() == 1 {
		query = query.Where("deleted_at IS NOT NULL")
	}

	if params.GetStatus() > 0 {
		query = query.Where("status = ?", params.GetStatus())
	}
	if params.GetName() != "" {
		query = query.Where("name_i18n->>'default' LIKE ?", "%"+params.GetName()+"%")
	}
	if params.GetProviderCode() != "" {
		query = query.Where("provider_code LIKE ?", "%"+params.GetProviderCode()+"%")
	}

	sortBy := strings.TrimSpace(params.SortBy)
	sortOrder := strings.TrimSpace(params.SortOrder)
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
	if err := query.Model(&ent.GameProvider{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("[GameProviderListDB] count failed: %v", err)
	}

	var providers []*ent.GameProvider
	if err := query.
		Offset(int((page - 1) * pageSize)).
		Limit(int(pageSize)).
		Order(fmt.Sprintf("%s %s", sortBy, sortOrder)).
		Find(&providers).Error; err != nil {
		return nil, 0, fmt.Errorf("[GameProviderListDB] query failed: %v", err)
	}

	return providers, total, nil
}
