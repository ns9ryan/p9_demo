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

// GameChannelGetParams 定义 GameChannelGet 的参数
type GameChannelGetParams struct {
	ID int64
}

// GameChannelGetResult 定义 GameChannelGet 的返回结果
type GameChannelGetResult struct {
	Channel   *model.Channel
	GameCount int64
}

// GameChannelGetDB 获取单个渠道
func GameChannelGetDB(ctx context.Context, db *gorm.DB, params *GameChannelGetParams) (*GameChannelGetResult, error) {
	if db == nil {
		return nil, fmt.Errorf("database not available")
	}

	var channel model.Channel
	if err := db.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", params.ID).
		First(&channel).Error; err != nil {
		return nil, fmt.Errorf("[GameChannelGetDB] query failed: %v", err)
	}

	result := &GameChannelGetResult{
		Channel: &channel,
	}

	// 查询该渠道下的游戏数量
	var gameCount int64
	if err := db.WithContext(ctx).
		Model(&model.Game{}).
		Where("channel_id = ? AND deleted_at IS NULL", channel.ID).
		Count(&gameCount).Error; err != nil {
		gameCount = 0
	}
	result.GameCount = gameCount

	return result, nil
}

// GameChannelListParams 定义 GameChannelList 的参数
type GameChannelListParams struct {
	Page        int64
	PageSize    int64
	ChannelCode string
	Name        string
	Status      int16
	IsDeleted   int8 // 0: not deleted, 1: deleted, -1: all
	SortBy      string
	SortOrder   string
}

// GameChannelListResult 定义 GameChannelList 的返回结果
type GameChannelListResult struct {
	Channels []*GameChannelListItem
	Total    int64
}

type GameChannelCategoryInfo struct {
	ID       int64
	NameI18n string
}

type GameChannelListItem struct {
	Channel     *model.Channel
	VendorCount int64
	GameCount   int64
	Categories  []GameChannelCategoryInfo
}

// GameChannelListDB 分页列表查询渠道
func GameChannelListDB(ctx context.Context, db *gorm.DB, params *GameChannelListParams) (*GameChannelListResult, error) {
	// 验证数据库连接是否可用
	if db == nil {
		return nil, fmt.Errorf("database not available")
	}

	// 处理分页参数：页码最小值为1，页大小默认15，最大100
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

	// 构建基础查询对象，附加上下文和超时控制
	query := db.WithContext(ctx)

	// 处理软删除过滤：0=未删除，1=已删除，其他值不过滤
	if params.IsDeleted == 0 {
		query = query.Where("deleted_at IS NULL")
	} else if params.IsDeleted == 1 {
		query = query.Where("deleted_at IS NOT NULL")
	}

	// 按状态过滤（可选）
	if params.Status > 0 {
		query = query.Where("status = ?", params.Status)
	}
	// 按渠道编码模糊查询（可选）
	if params.ChannelCode != "" {
		query = query.Where("channel_code LIKE ?", "%"+params.ChannelCode+"%")
	}
	// 按名称的 i18n.default 字段模糊查询（可选）
	if params.Name != "" {
		query = query.Where("name_i18n->>'default' LIKE ?", "%"+params.Name+"%")
	}

	// 验证并规范化排序字段和顺序
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

	// 排序字段白名单验证，防止 SQL 注入
	switch sortBy {
	case "id", "sort_no", "created_at":
	default:
		sortBy = "id"
	}

	// 排序顺序验证
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "asc"
	}

	// 执行 COUNT 查询获取总数
	var total int64
	if err := query.Model(&model.Channel{}).Count(&total).Error; err != nil {
		return nil, fmt.Errorf("[GameChannelListDB] count failed: %v", err)
	}

	// 执行分页查询，获取当前页的渠道数据
	var channels []*model.Channel
	if err := query.
		Offset(int((page - 1) * pageSize)).
		Limit(int(pageSize)).
		Order(fmt.Sprintf("%s %s", sortBy, sortOrder)).
		Find(&channels).Error; err != nil {
		return nil, fmt.Errorf("[GameChannelListDB] query failed: %v", err)
	}

	// 如果没有渠道数据，直接返回空结果
	// 提取所有渠道的ID，用于批量聚合查询
	channelIDs := make([]int64, 0, len(channels))
	for _, channel := range channels {
		if channel != nil {
			channelIDs = append(channelIDs, channel.ID)
		}
	}

	// 初始化映射表：存储聚合结果
	vendorCountMap := make(map[int64]int64)         // 渠道ID -> 厂商数量
	gameCountMap := make(map[int64]int64)           // 渠道ID -> 游戏数量
	categoryIDsByChannel := make(map[int64][]int64) // 渠道ID -> 分类ID列表
	categoryNameMap := make(map[int64]string)       // 分类ID -> 分类名称i18n

	// 如果存在渠道数据，执行关联数据聚合
	if len(channelIDs) > 0 {
		// 第一步：聚合查询游戏数和厂商数
		// 按 channel_id 分组，计算每个渠道的游戏总数和不重复厂商数
		type channelGameAggRow struct {
			ChannelID   int64 `gorm:"column:channel_id"`
			GameCount   int64 `gorm:"column:game_count"`
			VendorCount int64 `gorm:"column:vendor_count"`
		}
		var aggRows []channelGameAggRow
		if err := db.WithContext(ctx).
			Table("game").
			Select("channel_id, COUNT(*) AS game_count, COUNT(DISTINCT provider_id) AS vendor_count").
			Where("channel_id IN ? AND deleted_at IS NULL", channelIDs).
			Group("channel_id").
			Find(&aggRows).Error; err == nil {
			for _, row := range aggRows {
				gameCountMap[row.ChannelID] = row.GameCount
				vendorCountMap[row.ChannelID] = row.VendorCount
			}
		}

		// 第二步：查询每个渠道关联的分类列表及分类名称
		// 先查询每个渠道下的不重复分类ID
		type channelCategoryRow struct {
			ChannelID  int64 `gorm:"column:channel_id"`
			CategoryID int64 `gorm:"column:category_id"`
		}
		var channelCategoryRows []channelCategoryRow
		if err := db.WithContext(ctx).
			Table("game").
			Select("DISTINCT channel_id, category_id").
			Where("channel_id IN ? AND deleted_at IS NULL AND category_id > 0", channelIDs).
			Find(&channelCategoryRows).Error; err == nil {
			// 收集所有需要查询的分类ID（去重）
			categoryIDs := make([]int64, 0, len(channelCategoryRows))
			categorySeen := make(map[int64]struct{})
			for _, row := range channelCategoryRows {
				categoryIDsByChannel[row.ChannelID] = append(categoryIDsByChannel[row.ChannelID], row.CategoryID)
				if _, ok := categorySeen[row.CategoryID]; !ok {
					categorySeen[row.CategoryID] = struct{}{}
					categoryIDs = append(categoryIDs, row.CategoryID)
				}
			}

			// 批量查询分类信息，获取分类名称
			if len(categoryIDs) > 0 {
				var categories []model.Category
				if err := db.WithContext(ctx).
					Where("id IN ? AND deleted_at IS NULL", categoryIDs).
					Find(&categories).Error; err == nil {
					for _, category := range categories {
						categoryNameMap[category.ID] = string(utils.MustMarshalJSON(category.NameI18n))
					}
				}
			}
		}
	}

	// 第三步：组装返回结果
	// 将渠道和聚合数据关联，转换成 GameChannelListItem 结构体
	items := make([]*GameChannelListItem, 0, len(channels))
	for _, channel := range channels {
		if channel == nil {
			continue
		}
		// 获取该渠道的分类ID列表
		categoryIDs := categoryIDsByChannel[channel.ID]
		// 将分类ID转换为完整的分类信息（ID + 名称）
		categories := make([]GameChannelCategoryInfo, 0, len(categoryIDs))
		for _, categoryID := range categoryIDs {
			categories = append(categories, GameChannelCategoryInfo{
				ID:       categoryID,
				NameI18n: categoryNameMap[categoryID],
			})
		}
		// 组装单个渠道的完整信息
		items = append(items, &GameChannelListItem{
			Channel:     channel,
			VendorCount: vendorCountMap[channel.ID],
			GameCount:   gameCountMap[channel.ID],
			Categories:  categories,
		})
	}

	// 返回分页后的结果和总数
	return &GameChannelListResult{
		Channels: items,
		Total:    total,
	}, nil
}

// GameChannelUpdateParams 定义 GameChannelUpdate 的参数
type GameChannelUpdateParams struct {
	ID       int64
	NameI18n string
	SortNo   int32
	Status   int16
}

// GameChannelUpdateDB 更新渠道信息
func GameChannelUpdateDB(ctx context.Context, db *gorm.DB, params *GameChannelUpdateParams) (*model.Channel, error) {
	if db == nil {
		return nil, fmt.Errorf("database not available")
	}

	var channel model.Channel
	if err := db.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", params.ID).
		First(&channel).Error; err != nil {
		return nil, fmt.Errorf("[GameChannelUpdateDB] query failed: %v", err)
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
		Model(&channel).
		Updates(updateData).Error; err != nil {
		return nil, fmt.Errorf("[GameChannelUpdateDB] update failed: %v", err)
	}

	// 重新查询最新数据
	if err := db.WithContext(ctx).
		Where("id = ?", params.ID).
		First(&channel).Error; err != nil {
		return nil, fmt.Errorf("[GameChannelUpdateDB] query after update failed: %v", err)
	}

	return &channel, nil
}
