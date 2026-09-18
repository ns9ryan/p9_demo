package dao

import (
	"context"
	"fmt"
	"time"

	"oa.98ent.com/p9/platform-game/rpc/ent"
	"oa.98ent.com/p9/platform-game/rpc/ent/game"
	"oa.98ent.com/p9/platform-game/rpc/ent/gamecategory"
)

type GameCategoryDAO struct {
	db *ent.Client
}

func NewGameCategoryDAO(db *ent.Client) *GameCategoryDAO {
	return &GameCategoryDAO{
		db: db,
	}
}

// GetGameCategoryByID 根据ID获取游戏分类
func (d *GameCategoryDAO) GetGameCategoryByID(ctx context.Context, id int64) (*ent.GameCategory, error) {
	return d.db.GameCategory.Query().
		Where(gamecategory.IDEQ(id)).
		Where(gamecategory.DeletedAtIsNil()).
		Only(ctx)
}

func (d *GameCategoryDAO) GetGameCategoryBySourceId(ctx context.Context, sourceID int64) (*ent.GameCategory, error) {
	return d.db.GameCategory.Query().
		Where(gamecategory.SourceIDEQ(sourceID)).
		Where(gamecategory.DeletedAtIsNil()).
		Only(ctx)
}

// GetGameCategoryList 获取游戏分类列表
func (d *GameCategoryDAO) GetGameCategoryList(ctx context.Context, isDeleted int32, status int32, categoryCode string, offset, limit int64) ([]*ent.GameCategory, int, error) {
	query := d.db.GameCategory.Query()

	// 处理软删除条件
	if isDeleted == 0 {
		query = query.Where(gamecategory.DeletedAtIsNil())
	} else if isDeleted == 1 {
		query = query.Where(gamecategory.DeletedAtNotNil())
	}

	if status > 0 {
		query = query.Where(gamecategory.StatusEQ(int64(status)))
	}

	if categoryCode != "" {
		query = query.Where(gamecategory.CategoryCodeContains(categoryCode))
	}

	// 获取总数
	total, err := query.Clone().Count(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("count failed: %w", err)
	}

	// 获取分页数据
	categories, err := query.
		Order(
			gamecategory.BySortNo(),
			gamecategory.ByID(),
		).
		Offset(int(offset)).
		Limit(int(limit)).
		All(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("query failed: %w", err)
	}

	return categories, total, nil
}

func (d *GameCategoryDAO) GetAllGameCategories(ctx context.Context) ([]*ent.GameCategory, error) {
	return d.db.GameCategory.Query().
		Where(gamecategory.DeletedAtIsNil()).
		All(ctx)
}

// 更新游戏分类
func (d *GameCategoryDAO) UpdateGameCategory(ctx context.Context, id int64, updates map[string]interface{}) (*ent.GameCategory, error) {
	update := d.db.GameCategory.UpdateOneID(id)

	// 动态应用更新
	for key, value := range updates {
		switch key {
		case "sort_no":
			if val, ok := value.(int64); ok {
				update = update.SetSortNo(val)
			}
		case "status":
			if val, ok := value.(int64); ok {
				update = update.SetStatus(val)
			}
		case "deleted_at":
			if val, ok := value.(time.Time); ok {
				update = update.SetDeletedAt(val)
			}
		}
	}

	return update.Save(ctx)
}

func (d *GameCategoryDAO) BatchCreateGameCategory(ctx context.Context, items []*ent.GameCategoryCreate) ([]*ent.GameCategory, error) {
	if len(items) == 0 {
		return []*ent.GameCategory{}, nil
	}

	records, err := d.db.GameCategory.CreateBulk(items...).Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("batch create failed: %w", err)
	}
	return records, nil
}

func (d *GameCategoryDAO) GetGameCategoryBySourceCode(ctx context.Context, sourceCategoryCode string) (*ent.GameCategory, error) {
	return d.db.GameCategory.Query().
		Where(gamecategory.SourceCategoryCodeEQ(sourceCategoryCode)).
		Where(gamecategory.DeletedAtIsNil()).
		Only(ctx)
}

// GetGamesByCategoryID 获取指定分类下的所有游戏ID
func (d *GameCategoryDAO) GetGamesByCategoryID(ctx context.Context, categoryID int64) ([]int64, error) {
	games, err := d.db.Game.Query().
		Select(game.FieldSourceID).
		Where(game.CategoryIDEQ(categoryID)).
		Where(game.DeletedAtIsNil()).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("query games failed: %w", err)
	}

	gameIds := make([]int64, 0, len(games))
	for _, g := range games {
		gameIds = append(gameIds, g.SourceID)
	}
	return gameIds, nil
}

func (d *GameCategoryDAO) GetGameCategoryCount(ctx context.Context) (int, error) {
	count, err := d.db.GameCategory.Query().Count(ctx)
	if err != nil {
		return 0, fmt.Errorf("count failed: %w", err)
	}
	return count, nil
}

func (d *GameCategoryDAO) ExistByCode(ctx context.Context, categoryCode string) (bool, error) {
	exists, err := d.db.GameCategory.Query().
		Where(gamecategory.SourceCategoryCodeEQ(categoryCode)).
		Where(gamecategory.DeletedAtIsNil()).
		Exist(ctx)
	if err != nil {
		return false, fmt.Errorf("exist check failed: %w", err)
	}
	return exists, nil
}
