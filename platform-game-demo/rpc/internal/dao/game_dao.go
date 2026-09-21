package dao

import (
	"context"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"oa.98ent.com/p9/platform-game/rpc/ent"
	"oa.98ent.com/p9/platform-game/rpc/ent/currency"
	"oa.98ent.com/p9/platform-game/rpc/ent/game"
	"oa.98ent.com/p9/platform-game/rpc/ent/gamecategory"
	"oa.98ent.com/p9/platform-game/rpc/ent/gamechannel"
	"oa.98ent.com/p9/platform-game/rpc/ent/gamecurrency"
	"oa.98ent.com/p9/platform-game/rpc/ent/gameprovider"
)

type GameDAO struct {
	db *ent.Client
}

func NewGameDAO(db *ent.Client) *GameDAO {
	return &GameDAO{
		db: db,
	}
}

// GetGameByID 根据ID获取游戏
func (d *GameDAO) GetGameByID(ctx context.Context, id int64) (*ent.Game, error) {
	return d.db.Game.Query().
		Where(game.SourceIDEQ(id)).
		Where(game.DeletedAtIsNil()).
		Only(ctx)
}

// GetGameList 获取游戏列表
func (d *GameDAO) GetGameList(ctx context.Context, isDeleted int32, name string, status int32, gameCode string, providerId, categoryId, channelId int64, offset, limit int64) ([]*ent.Game, int, error) {
	logx.Infof("[DAO GameList] params: isDeleted=%d, name=%s, status=%d, gameCode=%s, providerId=%d, categoryId=%d, channelId=%d, offset=%d, limit=%d",
		isDeleted, name, status, gameCode, providerId, categoryId, channelId, offset, limit)

	query := d.db.Game.Query()

	// 处理软删除条件
	if isDeleted == 0 {
		query = query.Where(game.DeletedAtIsNil())
		logx.Infof("[DAO GameList] apply filter: DeletedAtIsNil")
	} else if isDeleted == 1 {
		query = query.Where(game.DeletedAtNotNil())
		logx.Infof("[DAO GameList] apply filter: DeletedAtNotNil")
	} else {
		logx.Infof("[DAO GameList] no soft delete filter (isDeleted=%d)", isDeleted)
	}

	if providerId > 0 {
		query = query.Where(game.ProviderIDEQ(providerId))
		logx.Infof("[DAO GameList] apply filter: providerId=%d", providerId)
	}
	if categoryId > 0 {
		query = query.Where(game.CategoryIDEQ(categoryId))
		logx.Infof("[DAO GameList] apply filter: categoryId=%d", categoryId)
	}
	if channelId > 0 {
		query = query.Where(game.ChannelIDEQ(channelId))
		logx.Infof("[DAO GameList] apply filter: channelId=%d", channelId)
	}
	if name != "" {
		query = query.Where(game.NameContains(name))
		logx.Infof("[DAO GameList] apply filter: name contains %s", name)
	}
	if status > 0 {
		query = query.Where(game.StatusEQ(int64(status)))
		logx.Infof("[DAO GameList] apply filter: status=%d", status)
	}

	if gameCode != "" {
		query = query.Where(game.GameCodeContains(gameCode))
		logx.Infof("[DAO GameList] apply filter: gameCode contains %s", gameCode)
	}

	// 获取总数
	total, err := query.Clone().Count(ctx)
	if err != nil {
		logx.Errorf("[DAO GameList] count failed: %v", err)
		return nil, 0, fmt.Errorf("count failed: %w", err)
	}
	logx.Infof("[DAO GameList] total count result: %d", total)

	// 获取分页数据
	games, err := query.
		Order(
			game.BySortNo(),
			game.ByID(),
		).
		Offset(int(offset)).
		Limit(int(limit)).
		All(ctx)
	if err != nil {
		logx.Errorf("[DAO GameList] query failed: %v", err)
		return nil, 0, fmt.Errorf("query failed: %w", err)
	}

	logx.Infof("[DAO GameList] query result: data_count=%d, total=%d", len(games), total)

	return games, total, nil
}

func (d *GameDAO) GetAllGame(ctx context.Context, providerId, categoryId, channelId int64) ([]*ent.Game, error) {
	query := d.db.Game.Query()

	if providerId > 0 {
		query = query.Where(game.ProviderIDEQ(providerId))
	}
	if categoryId > 0 {
		query = query.Where(game.CategoryIDEQ(categoryId))
	}
	if channelId > 0 {
		query = query.Where(game.ChannelIDEQ(channelId))
	}

	// 获取分页数据
	games, err := query.All(ctx)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}

	return games, nil
}

func (d *GameDAO) GetAllGames(ctx context.Context) ([]*ent.Game, error) {
	games, err := d.db.Game.Query().Where(game.DeletedAtIsNil()).All(ctx)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	return games, nil
}

// UpdateGame 更新游戏
func (d *GameDAO) UpdateGame(ctx context.Context, id int64, updates map[string]interface{}) (*ent.Game, error) {
	update := d.db.Game.UpdateOneID(id)

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

func (d *GameDAO) BatchCreateGame(ctx context.Context, items []*ent.GameCreate) ([]*ent.Game, error) {
	if len(items) == 0 {
		return []*ent.Game{}, nil
	}

	records, err := d.db.Game.CreateBulk(items...).Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("batch create failed: %w", err)
	}
	return records, nil
}

// GameExtInfo 游戏扩展信息
type GameExtInfo struct {
	CategoryCode string  `json:"category_code"`
	ProviderCode string  `json:"provider_code"`
	ChannelCode  string  `json:"channel_code"`
	Currencies   []int64 `json:"currencies"`
}

// GetGameExtInfo 获取游戏的扩展信息
func (d *GameDAO) GetGameExtInfo(ctx context.Context, gameRecord *ent.Game) (*GameExtInfo, error) {
	ext := &GameExtInfo{}

	// 查询分类信息
	if gameRecord.CategoryID > 0 {
		categoryRecord, err := d.db.GameCategory.Query().
			Select(gamecategory.FieldCategoryCode).
			Where(gamecategory.SourceIDEQ(gameRecord.CategoryID)).
			Where(gamecategory.DeletedAtIsNil()).
			Only(ctx)
		if err != nil {
			return nil, fmt.Errorf("query category failed: %w", err)
		}
		ext.CategoryCode = categoryRecord.CategoryCode
	}

	// 查询提供商信息
	if gameRecord.ProviderID > 0 {
		providerRecord, err := d.db.GameProvider.Query().
			Select(gameprovider.FieldProviderCode).
			Where(gameprovider.SourceIDEQ(gameRecord.ProviderID)).
			Where(gameprovider.DeletedAtIsNil()).
			Only(ctx)
		if err != nil {
			return nil, fmt.Errorf("query provider failed: %w", err)
		}
		ext.ProviderCode = providerRecord.ProviderCode
	}

	// 查询渠道信息（如果存在）
	if gameRecord.ChannelID > 0 {
		channelRecord, err := d.db.GameChannel.Query().
			Select(gamechannel.FieldChannelCode).
			Where(gamechannel.SourceIDEQ(gameRecord.ChannelID)).
			Where(gamechannel.DeletedAtIsNil()).
			Only(ctx)
		if err != nil {
			return nil, fmt.Errorf("query channel failed: %w", err)
		}
		ext.ChannelCode = channelRecord.ChannelCode
	}

	// 查询游戏货币信息
	gameCurrencyRecords, err := d.db.GameCurrency.Query().
		Select(gamecurrency.FieldCurrencyID).
		Where(gamecurrency.GameIDEQ(gameRecord.SourceID)).
		Where(gamecurrency.DeletedAtIsNil()).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("query game currency failed: %w", err)
	}

	currencies := make([]int64, 0, len(gameCurrencyRecords))
	for _, gameCurrencyRecord := range gameCurrencyRecords {
		_, err := d.db.Currency.Query().
			Select(currency.FieldNameKey).
			Where(currency.IDEQ(gameCurrencyRecord.CurrencyID)).
			Only(ctx)
		if err != nil {
			return nil, fmt.Errorf("query currency failed: %w", err)
		}
		currencies = append(currencies, gameCurrencyRecord.CurrencyID)
	}

	ext.Currencies = currencies
	return ext, nil
}

func (d *GameDAO) CreateGameCurrency(ctx context.Context, gameID int64, currencyID int64) (*ent.GameCurrency, error) {
	return d.db.GameCurrency.Create().
		SetGameID(gameID).
		SetCurrencyID(currencyID).
		SetSourceStatus(1).
		SetStatus(1).
		SetCreatedAt(time.Now()).
		SetUpdatedAt(time.Now()).
		Save(ctx)
}

func (d *GameDAO) ExistByCode(ctx context.Context, code string) (bool, error) {
	count, err := d.db.Game.Query().
		Where(game.SourceGameCodeEQ(code)).
		Where(game.DeletedAtIsNil()).
		Count(ctx)
	if err != nil {
		return false, fmt.Errorf("check game existence by code failed: %w", err)
	}
	return count > 0, nil
}

func (d *GameDAO) GetGameBySourceId(ctx context.Context, SourceId int64) (*ent.Game, error) {
	return d.db.Game.Query().
		Where(game.SourceIDEQ(SourceId)).
		Where(game.DeletedAtIsNil()).
		Only(ctx)
}

func (d *GameDAO) GetGameByCode(ctx context.Context, code string) (*ent.Game, error) {
	return d.db.Game.Query().
		Where(game.SourceGameCodeEQ(code)).
		Where(game.DeletedAtIsNil()).
		Only(ctx)
}
