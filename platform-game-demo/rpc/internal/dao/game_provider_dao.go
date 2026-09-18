package dao

import (
	"context"
	"fmt"

	"oa.98ent.com/p9/platform-game/rpc/ent"
	"oa.98ent.com/p9/platform-game/rpc/ent/gameprovider"
)

type GameProviderDAO struct {
	db *ent.Client
}

func NewGameProviderDAO(db *ent.Client) *GameProviderDAO {
	return &GameProviderDAO{
		db: db,
	}
}

// GetGameProviderByID 根据ID获取游戏提供商
func (d *GameProviderDAO) GetGameProviderByID(ctx context.Context, id int64) (*ent.GameProvider, error) {
	return d.db.GameProvider.Query().
		Where(gameprovider.IDEQ(id)).
		Where(gameprovider.DeletedAtIsNil()).
		Only(ctx)
}

func (d *GameProviderDAO) GetGameProviderBySourceId(ctx context.Context, sourceID int64) (*ent.GameProvider, error) {
	return d.db.GameProvider.Query().
		Where(gameprovider.SourceIDEQ(sourceID)).
		Where(gameprovider.DeletedAtIsNil()).
		Only(ctx)
}

func (d *GameProviderDAO) GetGameProviderBySourceCode(ctx context.Context, sourceProviderCode string) (*ent.GameProvider, error) {
	return d.db.GameProvider.Query().
		Where(gameprovider.SourceProviderCodeEQ(sourceProviderCode)).
		Where(gameprovider.DeletedAtIsNil()).
		Only(ctx)
}

// GetGameProviderList 获取游戏提供商列表
func (d *GameProviderDAO) GetGameProviderList(ctx context.Context, isDeleted int32, status int32, providerCode string, offset, limit int64) ([]*ent.GameProvider, int, error) {
	query := d.db.GameProvider.Query()

	// 处理软删除条件
	if isDeleted == 0 {
		query = query.Where(gameprovider.DeletedAtIsNil())
	} else if isDeleted == 1 {
		query = query.Where(gameprovider.DeletedAtNotNil())
	}

	if status > 0 {
		query = query.Where(gameprovider.StatusEQ(int64(status)))
	}

	if providerCode != "" {
		query = query.Where(gameprovider.ProviderCodeContains(providerCode))
	}

	// 获取总数
	total, err := query.Clone().Count(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("count failed: %w", err)
	}

	// 获取分页数据
	providers, err := query.
		Order(
			gameprovider.BySortNo(),
			gameprovider.ByID(),
		).
		Offset(int(offset)).
		Limit(int(limit)).
		All(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("query failed: %w", err)
	}

	return providers, total, nil
}

// UpdateGameProvider 更新游戏提供商
func (d *GameProviderDAO) UpdateGameProvider(ctx context.Context, id int64, updates map[string]interface{}) (*ent.GameProvider, error) {
	update := d.db.GameProvider.UpdateOneID(id)

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
		case "logo_url":
			if val, ok := value.(string); ok {
				update = update.SetLogoURL(val)
			}
		}
	}

	return update.Save(ctx)
}

func (d *GameProviderDAO) BatchCreateGameProvider(ctx context.Context, items []*ent.GameProviderCreate) ([]*ent.GameProvider, error) {
	if len(items) == 0 {
		return []*ent.GameProvider{}, nil
	}

	records, err := d.db.GameProvider.CreateBulk(items...).Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("batch create failed: %w", err)
	}
	return records, nil
}
func (d *GameProviderDAO) GetGameProviderCount(ctx context.Context) (int, error) {
	count, err := d.db.GameProvider.Query().Count(ctx)
	if err != nil {
		return 0, fmt.Errorf("count failed: %w", err)
	}
	return count, nil
}

func (d *GameProviderDAO) GetAllGameProviders(ctx context.Context) ([]*ent.GameProvider, error) {
	providers, err := d.db.GameProvider.Query().Where(gameprovider.DeletedAtIsNil()).All(ctx)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	return providers, nil
}

func (d *GameProviderDAO) ExistByCode(ctx context.Context, sourceProviderCode string) (bool, error) {
	exists, err := d.db.GameProvider.Query().
		Where(gameprovider.SourceProviderCodeEQ(sourceProviderCode)).
		Where(gameprovider.DeletedAtIsNil()).
		Exist(ctx)
	if err != nil {
		return false, fmt.Errorf("exist check failed: %w", err)
	}
	return exists, nil
}
