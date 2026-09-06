package game

import (
	"context"
	"fmt"
	"strings"
	"time"

	"oa.98ent.com/p9/platform-game/common/model"
	"oa.98ent.com/p9/platform-game/common/utils"

	"gorm.io/gorm"
)

// GameProviderGetParams 定义 GameProviderGet 的参数
type GameProviderGetParams struct {
	ID int64
}

// GameProviderGetDB 获取单个供应商
func GameProviderGetDB(ctx context.Context, db *gorm.DB, params *GameProviderGetParams) (*model.Provider, error) {
	if db == nil {
		return nil, fmt.Errorf("database not available")
	}

	var provider model.Provider
	if err := db.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", params.ID).
		First(&provider).Error; err != nil {
		return nil, fmt.Errorf("[GameProviderGetDB] query failed: %v", err)
	}

	return &provider, nil
}

// GameProviderListParams 定义 GameProviderList 的参数
type GameProviderListParams struct {
	Page      int64
	PageSize  int64
	Name      string
	Status    int16
	IsDeleted int8 // 0: not deleted, 1: deleted, -1: all
	SortBy    string
	SortOrder string
}

// GameProviderListResult 定义 GameProviderList 的返回结果
type GameProviderListResult struct {
	Providers []*model.Provider
	Total     int64
}

// GameProviderListDB 分页列表查询供应商
func GameProviderListDB(ctx context.Context, db *gorm.DB, params *GameProviderListParams) (*GameProviderListResult, error) {
	if db == nil {
		return nil, fmt.Errorf("database not available")
	}

	page := params.Page
	if page <= 0 {
		page = 1
	}
	pageSize := params.PageSize
	if pageSize <= 0 {
		pageSize = 15
	}
	if pageSize > 100 {
		pageSize = 100
	}

	query := db.WithContext(ctx)

	// 处理软删除条件
	if params.IsDeleted == 0 {
		query = query.Where("deleted_at IS NULL")
	} else if params.IsDeleted == 1 {
		query = query.Where("deleted_at IS NOT NULL")
	}

	if params.Status > 0 {
		query = query.Where("status = ?", params.Status)
	}
	if params.Name != "" {
		query = query.Where("name_i18n->>'default' LIKE ?", "%"+params.Name+"%")
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
	if err := query.Model(&model.Provider{}).Count(&total).Error; err != nil {
		return nil, fmt.Errorf("[GameProviderListDB] count failed: %v", err)
	}

	var providers []*model.Provider
	if err := query.
		Offset(int((page - 1) * pageSize)).
		Limit(int(pageSize)).
		Order(fmt.Sprintf("%s %s", sortBy, sortOrder)).
		Find(&providers).Error; err != nil {
		return nil, fmt.Errorf("[GameProviderListDB] query failed: %v", err)
	}

	return &GameProviderListResult{
		Providers: providers,
		Total:     total,
	}, nil
}

// GameProviderUpdateParams 定义 GameProviderUpdate 的参数
type GameProviderUpdateParams struct {
	ID       int64
	NameI18n string
	SortNo   int32
	Status   int16
}

// GameProviderUpdateDB 更新供应商信息
func GameProviderUpdateDB(ctx context.Context, db *gorm.DB, params *GameProviderUpdateParams) (*model.Provider, error) {
	if db == nil {
		return nil, fmt.Errorf("database not available")
	}

	var provider model.Provider
	if err := db.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", params.ID).
		First(&provider).Error; err != nil {
		return nil, fmt.Errorf("[GameProviderUpdateDB] query failed: %v", err)
	}

	updateData := make(map[string]interface{})
	updateData["updated_at"] = time.Now()

	if params.NameI18n != "" {
		updateData["name_i18n"] = utils.MustParseJSON([]byte(params.NameI18n))
	}
	if params.SortNo > 0 {
		updateData["sort_no"] = params.SortNo
	}
	if params.Status > 0 {
		updateData["status"] = params.Status
	}

	if err := db.WithContext(ctx).
		Model(&provider).
		Updates(updateData).Error; err != nil {
		return nil, fmt.Errorf("[GameProviderUpdateDB] update failed: %v", err)
	}

	// 重新查询最新数据
	if err := db.WithContext(ctx).
		Where("id = ?", params.ID).
		First(&provider).Error; err != nil {
		return nil, fmt.Errorf("[GameProviderUpdateDB] query after update failed: %v", err)
	}

	return &provider, nil
}
