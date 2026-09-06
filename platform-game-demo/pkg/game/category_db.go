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

// GameCategoryGetParams 定义 GameCategoryGet 的参数
type GameCategoryGetParams struct {
	ID int64
}

// GameCategoryGetDB 获取单个分类
func GameCategoryGetDB(ctx context.Context, db *gorm.DB, params *GameCategoryGetParams) (*model.Category, error) {
	if db == nil {
		return nil, fmt.Errorf("database not available")
	}

	var category model.Category
	if err := db.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", params.ID).
		First(&category).Error; err != nil {
		return nil, fmt.Errorf("[GameCategoryGetDB] query failed: %v", err)
	}

	return &category, nil
}

// GameCategoryListParams 定义 GameCategoryList 的参数
type GameCategoryListParams struct {
	Page         int64
	PageSize     int64
	CategoryCode string
	Name         string
	Status       int16
	IsDeleted    int8 // 0: not deleted, 1: deleted, -1: all
	SortBy       string
	SortOrder    string
}

// GameCategoryListResult 定义 GameCategoryList 的返回结果
type GameCategoryListResult struct {
	Categories []*model.Category
	Total      int64
}

// GameCategoryListDB 分页列表查询分类
func GameCategoryListDB(ctx context.Context, db *gorm.DB, params *GameCategoryListParams) (*GameCategoryListResult, error) {
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
	if params.CategoryCode != "" {
		query = query.Where("category_code LIKE ?", "%"+params.CategoryCode+"%")
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
	if err := query.Model(&model.Category{}).Count(&total).Error; err != nil {
		return nil, fmt.Errorf("[GameCategoryListDB] count failed: %v", err)
	}

	var categories []*model.Category
	if err := query.
		Offset(int((page - 1) * pageSize)).
		Limit(int(pageSize)).
		Order(fmt.Sprintf("%s %s", sortBy, sortOrder)).
		Find(&categories).Error; err != nil {
		return nil, fmt.Errorf("[GameCategoryListDB] query failed: %v", err)
	}

	return &GameCategoryListResult{
		Categories: categories,
		Total:      total,
	}, nil
}

// GameCategoryUpdateParams 定义 GameCategoryUpdate 的参数
type GameCategoryUpdateParams struct {
	ID       int64
	NameI18n string
	SortNo   int32
	Status   int16
}

// GameCategoryUpdateDB 更新分类信息
func GameCategoryUpdateDB(ctx context.Context, db *gorm.DB, params *GameCategoryUpdateParams) (*model.Category, error) {
	if db == nil {
		return nil, fmt.Errorf("database not available")
	}

	var category model.Category
	if err := db.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", params.ID).
		First(&category).Error; err != nil {
		return nil, fmt.Errorf("[GameCategoryUpdateDB] query failed: %v", err)
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
		Model(&category).
		Updates(updateData).Error; err != nil {
		return nil, fmt.Errorf("[GameCategoryUpdateDB] update failed: %v", err)
	}

	// 重新查询最新数据
	if err := db.WithContext(ctx).
		Where("id = ?", params.ID).
		First(&category).Error; err != nil {
		return nil, fmt.Errorf("[GameCategoryUpdateDB] query after update failed: %v", err)
	}

	return &category, nil
}
