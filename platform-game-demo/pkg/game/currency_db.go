package game

import (
	"context"
	"fmt"
	"strings"
	"time"

	"oa.98ent.com/p9/platform-game/common/model"

	"gorm.io/gorm"
)

// GameCurrencyWithMeta 带有关联信息的游戏货币
type GameCurrencyWithMeta struct {
	model.GameCurrency
	GameCode     string `gorm:"column:game_code"`
	CurrencyCode string `gorm:"column:currency_code"`
}

// GameCurrencyGetParams 定义 GameCurrencyGet 的参数
type GameCurrencyGetParams struct {
	ID int64
}

// GameCurrencyGetDB 获取单个游戏货币关系
func GameCurrencyGetDB(ctx context.Context, db *gorm.DB, params *GameCurrencyGetParams) (*GameCurrencyWithMeta, error) {
	if db == nil {
		return nil, fmt.Errorf("database not available")
	}

	var currency GameCurrencyWithMeta
	if err := db.WithContext(ctx).
		Model(&model.GameCurrency{}).
		Select("game_currency.*, g.game_code AS game_code, c.currency_code AS currency_code").
		Joins("JOIN game g ON game_currency.game_id = g.id AND g.deleted_at IS NULL").
		Joins("JOIN sys_currency c ON game_currency.currency_id = c.id").
		Where("game_currency.id = ? AND game_currency.deleted_at IS NULL", params.ID).
		First(&currency).Error; err != nil {
		return nil, fmt.Errorf("[GameCurrencyGetDB] query failed: %v", err)
	}

	return &currency, nil
}

// GameCurrencyListParams 定义 GameCurrencyList 的参数
type GameCurrencyListParams struct {
	Page       int64
	PageSize   int64
	GameID     int64
	CurrencyID int64
	Status     int16
	IsDeleted  int8 // 0: not deleted, 1: deleted, -1: all
	SortBy     string
	SortOrder  string
}

// GameCurrencyListResult 定义 GameCurrencyList 的返回结果
type GameCurrencyListResult struct {
	Currencies []*GameCurrencyWithMeta
	Total      int64
}

// GameCurrencyListDB 分页列表查询游戏货币关系
func GameCurrencyListDB(ctx context.Context, db *gorm.DB, params *GameCurrencyListParams) (*GameCurrencyListResult, error) {
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

	query := db.WithContext(ctx).
		Model(&model.GameCurrency{}).
		Joins("JOIN game g ON game_currency.game_id = g.id AND g.deleted_at IS NULL").
		Joins("JOIN sys_currency c ON game_currency.currency_id = c.id")

	// 处理软删除条件
	if params.IsDeleted == 0 {
		query = query.Where("game_currency.deleted_at IS NULL")
	} else if params.IsDeleted == 1 {
		query = query.Where("game_currency.deleted_at IS NOT NULL")
	}

	if params.GameID > 0 {
		query = query.Where("game_currency.game_id = ?", params.GameID)
	}
	if params.CurrencyID > 0 {
		query = query.Where("game_currency.currency_id = ?", params.CurrencyID)
	}
	if params.Status > 0 {
		query = query.Where("game_currency.status = ?", params.Status)
	}

	orderByColumn := "game_currency.id"
	switch strings.ToLower(strings.TrimSpace(params.SortBy)) {
	case "id":
		orderByColumn = "game_currency.id"
	case "created_at":
		orderByColumn = "game_currency.created_at"
	}

	orderDirection := "DESC"
	switch strings.ToLower(strings.TrimSpace(params.SortOrder)) {
	case "asc":
		orderDirection = "ASC"
	case "desc":
		orderDirection = "DESC"
	}
	orderClause := fmt.Sprintf("%s %s", orderByColumn, orderDirection)

	var total int64
	if err := query.Model(&model.GameCurrency{}).Count(&total).Error; err != nil {
		return nil, fmt.Errorf("[GameCurrencyListDB] count failed: %v", err)
	}

	var currencies []*GameCurrencyWithMeta
	if err := query.
		Select("game_currency.*, g.game_code AS game_code, c.currency_code AS currency_code").
		Offset(int((page - 1) * pageSize)).
		Limit(int(pageSize)).
		Order(orderClause).
		Find(&currencies).Error; err != nil {
		return nil, fmt.Errorf("[GameCurrencyListDB] query failed: %v", err)
	}

	return &GameCurrencyListResult{
		Currencies: currencies,
		Total:      total,
	}, nil
}

// GameCurrencyUpdateParams 定义 GameCurrencyUpdate 的参数
type GameCurrencyUpdateParams struct {
	ID     int64
	Status int16
}

// GameCurrencyUpdateDB 更新游戏货币状态
func GameCurrencyUpdateDB(ctx context.Context, db *gorm.DB, params *GameCurrencyUpdateParams) (*model.GameCurrency, error) {
	if db == nil {
		return nil, fmt.Errorf("database not available")
	}

	var currency model.GameCurrency
	if err := db.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", params.ID).
		First(&currency).Error; err != nil {
		return nil, fmt.Errorf("[GameCurrencyUpdateDB] query failed: %v", err)
	}

	updateData := make(map[string]interface{})
	updateData["updated_at"] = time.Now()

	if params.Status > 0 {
		updateData["status"] = params.Status
	}

	if err := db.WithContext(ctx).
		Model(&currency).
		Updates(updateData).Error; err != nil {
		return nil, fmt.Errorf("[GameCurrencyUpdateDB] update failed: %v", err)
	}

	// 重新查询最新数据
	if err := db.WithContext(ctx).
		Where("id = ?", params.ID).
		First(&currency).Error; err != nil {
		return nil, fmt.Errorf("[GameCurrencyUpdateDB] query after update failed: %v", err)
	}

	return &currency, nil
}
