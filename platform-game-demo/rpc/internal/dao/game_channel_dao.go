package dao

import (
	"context"
	"fmt"
	"time"

	"oa.98ent.com/p9/platform-game/rpc/ent"
	"oa.98ent.com/p9/platform-game/rpc/ent/gamechannel"
)

type GameChannelDAO struct {
	db *ent.Client
}

func NewGameChannelDAO(db *ent.Client) *GameChannelDAO {
	return &GameChannelDAO{
		db: db,
	}
}

// GetGameChannelByID 根据ID获取游戏渠道
func (d *GameChannelDAO) GetGameChannelByID(ctx context.Context, id int64) (*ent.GameChannel, error) {
	return d.db.GameChannel.Query().
		Where(gamechannel.IDEQ(id)).
		Where(gamechannel.DeletedAtIsNil()).
		Only(ctx)
}

func (d *GameChannelDAO) GetGameChannelBySourceId(ctx context.Context, sourceID int64) (*ent.GameChannel, error) {
	return d.db.GameChannel.Query().
		Where(gamechannel.SourceIDEQ(sourceID)).
		Where(gamechannel.DeletedAtIsNil()).
		Only(ctx)
}

func (d *GameChannelDAO) GetGameChannelBySourceCode(ctx context.Context, sourceChannelCode string) (*ent.GameChannel, error) {
	return d.db.GameChannel.Query().
		Where(gamechannel.SourceChannelCodeEQ(sourceChannelCode)).
		Where(gamechannel.DeletedAtIsNil()).
		Only(ctx)
}

// GetGameChannelList 获取游戏渠道列表
func (d *GameChannelDAO) GetGameChannelList(ctx context.Context, isDeleted int32, status int32, channelCode string, offset, limit int64) ([]*ent.GameChannel, int, error) {
	query := d.db.GameChannel.Query()

	// 处理软删除条件
	if isDeleted == 0 {
		query = query.Where(gamechannel.DeletedAtIsNil())
	} else if isDeleted == 1 {
		query = query.Where(gamechannel.DeletedAtNotNil())
	}

	if status > 0 {
		query = query.Where(gamechannel.StatusEQ(int64(status)))
	}

	if channelCode != "" {
		query = query.Where(gamechannel.ChannelCodeContains(channelCode))
	}

	// 获取总数
	total, err := query.Clone().Count(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("count failed: %w", err)
	}

	// 获取分页数据
	channels, err := query.
		Order(
			gamechannel.BySortNo(),
			gamechannel.ByID(),
		).
		Offset(int(offset)).
		Limit(int(limit)).
		All(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("query failed: %w", err)
	}

	return channels, total, nil
}

func (d *GameChannelDAO) GetAllGameChannel(ctx context.Context) ([]*ent.GameChannel, error) {
	return d.db.GameChannel.Query().
		Where(gamechannel.DeletedAtIsNil()).
		All(ctx)
}

// UpdateGameChannel 更新游戏渠道
func (d *GameChannelDAO) UpdateGameChannel(ctx context.Context, id int64, updates map[string]interface{}) (*ent.GameChannel, error) {
	update := d.db.GameChannel.UpdateOneID(id)

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

func (d *GameChannelDAO) BatchCreateGameChannel(ctx context.Context, items []*ent.GameChannelCreate) ([]*ent.GameChannel, error) {
	if len(items) == 0 {
		return []*ent.GameChannel{}, nil
	}

	records, err := d.db.GameChannel.CreateBulk(items...).Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("batch create failed: %w", err)
	}
	return records, nil
}

func (d *GameChannelDAO) GetGameChannelCount(ctx context.Context) (int, error) {
	count, err := d.db.GameChannel.Query().Count(ctx)
	if err != nil {
		return 0, fmt.Errorf("count failed: %w", err)
	}
	return count, nil
}

func (d *GameChannelDAO) ExistByCode(ctx context.Context, channelCode string) (bool, error) {
	exists, err := d.db.GameChannel.Query().
		Where(gamechannel.SourceChannelCodeEQ(channelCode)).
		Where(gamechannel.DeletedAtIsNil()).
		Exist(ctx)
	if err != nil {
		return false, fmt.Errorf("exist check failed: %w", err)
	}
	return exists, nil
}
