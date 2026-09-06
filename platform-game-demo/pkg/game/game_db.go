package game

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"oa.98ent.com/p9/platform-game/common/model"
	"oa.98ent.com/p9/platform-game/common/utils"

	"gorm.io/gorm"
)

// GameGetParams 定义 GameGet 的参数
type GameGetParams struct {
	ID int64
}

// GameGetResult 定义 GameGet 的返回结果
type GameGetResult struct {
	Game              *model.Game
	CategoryName      string
	ProviderName      string
	ChannelName       string
	GameCurrencyInfos []GameCurrencyInfo
}

// GameCurrencyInfo 游戏货币信息
type GameCurrencyInfo struct {
	CurrencyID       int64
	CurrencyNameI18n string
}

// GameGetDB 获取单个游戏信息
func GameGetDB(ctx context.Context, db *gorm.DB, params *GameGetParams) (*GameGetResult, error) {
	if db == nil {
		return nil, fmt.Errorf("database not available")
	}

	var game model.Game
	if err := db.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", params.ID).
		First(&game).Error; err != nil {
		return nil, fmt.Errorf("[GameGetDB] query failed: %v", err)
	}

	result := &GameGetResult{
		Game: &game,
	}

	// 获取分类信息
	var category model.Category
	if err := db.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", game.CategoryID).
		First(&category).Error; err == nil {
		result.CategoryName = string(utils.MustMarshalJSON(category.NameI18n))
	}

	// 获取供应商信息
	var provider model.Provider
	if err := db.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", game.ProviderID).
		First(&provider).Error; err == nil {
		result.ProviderName = string(utils.MustMarshalJSON(provider.NameI18n))
	}

	// 获取渠道信息
	if game.ChannelID != 0 {
		var channel model.Channel
		if err := db.WithContext(ctx).
			Where("id = ? AND deleted_at IS NULL", game.ChannelID).
			First(&channel).Error; err == nil {
			result.ChannelName = string(utils.MustMarshalJSON(channel.NameI18n))
		}
	}

	// 获取游戏货币信息
	type currencyInfoRow struct {
		CurrencyID       int64           `gorm:"column:currency_id"`
		CurrencyNameI18n json.RawMessage `gorm:"column:name_i18n"`
	}
	var currencyRows []currencyInfoRow
	if err := db.WithContext(ctx).
		Table("game_currency gc").
		Select("gc.currency_id, sc.name_i18n").
		Joins("JOIN sys_currency sc ON gc.currency_id = sc.id").
		Where("gc.game_id = ? AND gc.deleted_at IS NULL", game.ID).
		Find(&currencyRows).Error; err == nil {
		for _, row := range currencyRows {
			result.GameCurrencyInfos = append(result.GameCurrencyInfos, GameCurrencyInfo{
				CurrencyID:       row.CurrencyID,
				CurrencyNameI18n: string(row.CurrencyNameI18n),
			})
		}
	}

	return result, nil
}

// GameListParams 定义 GameList 的参数
type GameListParams struct {
	Page       int64
	PageSize   int64
	GameCode   string
	Name       string
	CategoryID int64
	ProviderID int64
	ChannelID  int64
	Status     int16
	IsDeleted  int8 // 0: not deleted, 1: deleted, -1: all
	SortBy     string
	SortOrder  string
}

// GameListResult 定义 GameList 的返回结果
type GameListResult struct {
	Games []*GameListItem
	Total int64
}

// GameListItem 游戏列表项
type GameListItem struct {
	Game              *model.Game
	CategoryNameI18n  string
	ProviderNameI18n  string
	ChannelNameI18n   string
	GameCurrencyInfos []GameCurrencyInfo
}

// GameListDB 分页列表查询游戏
func GameListDB(ctx context.Context, db *gorm.DB, params *GameListParams) (*GameListResult, error) {
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

	if params.CategoryID > 0 {
		query = query.Where("category_id = ?", params.CategoryID)
	}
	if params.ProviderID > 0 {
		query = query.Where("provider_id = ?", params.ProviderID)
	}
	if params.ChannelID > 0 {
		query = query.Where("channel_id = ?", params.ChannelID)
	}
	if params.GameCode != "" {
		query = query.Where("game_code LIKE ?", "%"+params.GameCode+"%")
	}
	if params.Name != "" {
		query = query.Where("name_i18n->>'default' LIKE ?", "%"+params.Name+"%")
	}
	if params.Status > 0 {
		query = query.Where("status = ?", params.Status)
	}

	sortBy := strings.ToLower(strings.TrimSpace(params.SortBy))
	if sortBy != "id" && sortBy != "sort_no" && sortBy != "created_at" {
		sortBy = "sort_no"
	}
	sortOrder := strings.ToLower(strings.TrimSpace(params.SortOrder))
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "asc"
	}
	orderExpr := sortBy + " " + sortOrder + ", id DESC"

	var total int64
	if err := query.Model(&model.Game{}).Count(&total).Error; err != nil {
		return nil, fmt.Errorf("[GameListDB] count failed: %v", err)
	}

	var games []*model.Game
	if err := query.
		Offset(int((page - 1) * pageSize)).
		Limit(int(pageSize)).
		Order(orderExpr).
		Find(&games).Error; err != nil {
		return nil, fmt.Errorf("[GameListDB] query failed: %v", err)
	}
	if len(games) > 0 {
		sample := games[0]
		log.Printf("[GameListDB] sample row: game_id=%d, category_id=%d, provider_id=%d, channel_id=%d", sample.ID, sample.CategoryID, sample.ProviderID, sample.ChannelID)
	} else {
		log.Printf("[GameListDB] no rows returned for current filter")
	}

	// 批量获取关联信息
	categoryIDs := make([]int64, 0, len(games))
	providerIDs := make([]int64, 0, len(games))
	channelIDs := make([]int64, 0, len(games))
	gameIDs := make([]int64, 0, len(games))

	for _, game := range games {
		categoryIDs = append(categoryIDs, game.CategoryID)
		providerIDs = append(providerIDs, game.ProviderID)
		if game.ChannelID != 0 {
			channelIDs = append(channelIDs, game.ChannelID)
		}
		gameIDs = append(gameIDs, game.ID)
	}

	// 查询分类信息
	categoryNames := make(map[int64]string)
	if len(categoryIDs) > 0 {
		var categories []model.Category
		if err := db.WithContext(ctx).
			Where("id IN ? AND deleted_at IS NULL", categoryIDs).
			Find(&categories).Error; err == nil {
			for _, c := range categories {
				categoryNames[c.ID] = string(utils.MustMarshalJSON(c.NameI18n))
			}
		}
	}

	// 查询供应商信息
	providerNames := make(map[int64]string)
	if len(providerIDs) > 0 {
		var providers []model.Provider
		if err := db.WithContext(ctx).
			Where("id IN ? AND deleted_at IS NULL", providerIDs).
			Find(&providers).Error; err == nil {
			for _, p := range providers {
				providerNames[p.ID] = string(utils.MustMarshalJSON(p.NameI18n))
			}
		}
	}

	// 查询渠道信息
	channelNames := make(map[int64]string)
	if len(channelIDs) > 0 {
		var channels []model.Channel
		if err := db.WithContext(ctx).
			Where("id IN ? AND deleted_at IS NULL", channelIDs).
			Find(&channels).Error; err == nil {
			for _, ch := range channels {
				channelNames[ch.ID] = string(utils.MustMarshalJSON(ch.NameI18n))
			}
		}
	}

	// 查询货币信息
	currencyInfoMap := make(map[int64][]GameCurrencyInfo)
	if len(gameIDs) > 0 {
		type currencyInfoRow struct {
			GameID           int64           `gorm:"column:game_id"`
			CurrencyID       int64           `gorm:"column:currency_id"`
			CurrencyNameI18n json.RawMessage `gorm:"column:name_i18n"`
		}
		var rows []currencyInfoRow

		if err := db.WithContext(ctx).
			Table("game_currency gc").
			Select("gc.game_id, gc.currency_id, sc.name_i18n").
			Joins("JOIN sys_currency sc ON gc.currency_id = sc.id").
			Where("gc.game_id IN ? AND gc.deleted_at IS NULL", gameIDs).
			Find(&rows).Error; err == nil {
			for _, row := range rows {
				currencyInfoMap[row.GameID] = append(currencyInfoMap[row.GameID], GameCurrencyInfo{
					CurrencyID:       row.CurrencyID,
					CurrencyNameI18n: string(row.CurrencyNameI18n),
				})
			}
		}
	}

	// 组装结果
	items := make([]*GameListItem, len(games))
	for i, game := range games {
		currencyInfos := currencyInfoMap[game.ID]
		if currencyInfos == nil {
			currencyInfos = make([]GameCurrencyInfo, 0)
		}

		items[i] = &GameListItem{
			Game:              game,
			CategoryNameI18n:  categoryNames[game.CategoryID],
			ProviderNameI18n:  providerNames[game.ProviderID],
			ChannelNameI18n:   channelNames[game.ChannelID],
			GameCurrencyInfos: currencyInfos,
		}
	}
	if len(items) > 0 && items[0] != nil && items[0].Game != nil {
		g := items[0].Game
		log.Printf("[GameListDB] sample assembled item: game_id=%d, category_id=%d, provider_id=%d, channel_id=%d", g.ID, g.CategoryID, g.ProviderID, g.ChannelID)
	}

	return &GameListResult{
		Games: items,
		Total: total,
	}, nil
}

// GameUpdateParams 定义 GameUpdate 的参数
type GameUpdateParams struct {
	ID               int64
	NameI18n         string
	ImageURL         string
	SortNo           int32
	ProviderKey      string
	SupportsEmbed    bool
	SupportsRedirect bool
	Status           int16
}

// GameUpdateDB 更新游戏信息
func GameUpdateDB(ctx context.Context, db *gorm.DB, params *GameUpdateParams) (*model.Game, error) {
	if db == nil {
		return nil, fmt.Errorf("database not available")
	}

	var game model.Game
	if err := db.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", params.ID).
		First(&game).Error; err != nil {
		return nil, fmt.Errorf("[GameUpdateDB] query failed: %v", err)
	}

	updateData := make(map[string]interface{})

	if params.NameI18n != "" {
		updateData["name_i18n"] = utils.MustParseJSON([]byte(params.NameI18n))
	}
	if params.ImageURL != "" {
		updateData["image_url"] = params.ImageURL
	}
	if params.SortNo > 0 {
		updateData["sort_no"] = params.SortNo
	}
	if params.ProviderKey != "" {
		updateData["provider_key"] = params.ProviderKey
	}
	updateData["supports_embed"] = params.SupportsEmbed
	updateData["supports_redirect"] = params.SupportsRedirect
	if params.Status > 0 {
		updateData["status"] = params.Status
	}

	if err := db.WithContext(ctx).
		Model(&game).
		Updates(updateData).Error; err != nil {
		return nil, fmt.Errorf("[GameUpdateDB] update failed: %v", err)
	}

	// 重新查询最新数据
	if err := db.WithContext(ctx).
		Where("id = ?", params.ID).
		First(&game).Error; err != nil {
		return nil, fmt.Errorf("[GameUpdateDB] query after update failed: %v", err)
	}

	return &game, nil
}
