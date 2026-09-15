package dao

import (
	"context"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"oa.98ent.com/p9/platform-game/rpc/ent"
	"oa.98ent.com/p9/platform-game/rpc/ent/gamecurrency"
)

type GameCurrencyDAO struct {
	db *ent.Client
}

func NewGameCurrencyDAO(db *ent.Client) *GameCurrencyDAO {
	return &GameCurrencyDAO{
		db: db,
	}
}

// GetGameCurrencyByID 根据ID获取游戏货币
func (d *GameCurrencyDAO) GetGameCurrencyByID(ctx context.Context, id int64) (*ent.GameCurrency, error) {
	return d.db.GameCurrency.Query().
		Where(gamecurrency.IDEQ(id)).
		Where(gamecurrency.DeletedAtIsNil()).
		Only(ctx)
}

// GetGameCurrencyList 获取游戏货币列表
func (d *GameCurrencyDAO) GetGameCurrencyList(ctx context.Context, isDeleted int32, status int32, gameID, currencyID int64, offset, limit int64) ([]*ent.GameCurrency, int, error) {
	query := d.db.GameCurrency.Query()
	if status > 0 {
		query = query.Where(gamecurrency.StatusEQ(int64(status)))
		logx.Infof("[DAO GameCurrencyList] apply filter: status=%d", status)
	} else if status == 0 {
		logx.Infof("[DAO GameCurrencyList] skip status filter (status=0 means all)")
	}
	if gameID > 0 {
		query = query.Where(gamecurrency.GameIDEQ(gameID))
		logx.Infof("[DAO GameCurrencyList] apply filter: gameID=%d", gameID)
	}
	if currencyID > 0 {
		query = query.Where(gamecurrency.CurrencyIDEQ(currencyID))
		logx.Infof("[DAO GameCurrencyList] apply filter: currencyID=%d", currencyID)
	}

	// 处理软删除条件
	if isDeleted == 0 {
		query = query.Where(gamecurrency.DeletedAtIsNil())
		logx.Infof("[DAO GameCurrencyList] apply filter: DeletedAtIsNil")
	} else if isDeleted == 1 {
		query = query.Where(gamecurrency.DeletedAtNotNil())
		logx.Infof("[DAO GameCurrencyList] apply filter: DeletedAtNotNil")
	} else {
		logx.Infof("[DAO GameCurrencyList] no soft delete filter (isDeleted=%d)", isDeleted)
	}

	// 获取总数
	total, err := query.Clone().Count(ctx)
	if err != nil {
		logx.Errorf("[DAO GameCurrencyList] count failed: %v", err)
		return nil, 0, fmt.Errorf("count failed: %w", err)
	}
	logx.Infof("[DAO GameCurrencyList] total count result: %d", total)

	// 获取分页数据
	currencies, err := query.
		Order(gamecurrency.ByID()).
		Offset(int(offset)).
		Limit(int(limit)).
		All(ctx)
	if err != nil {
		logx.Errorf("[DAO GameCurrencyList] query failed: %v", err)
		return nil, 0, fmt.Errorf("query failed: %w", err)
	}

	logx.Infof("[DAO GameCurrencyList] query result: data_count=%d, total=%d", len(currencies), total)

	return currencies, total, nil
}

func (d *GameCurrencyDAO) GetAllGameCurrency(ctx context.Context, gameID, currencyID int64) ([]*ent.GameCurrency, error) {
	query := d.db.GameCurrency.Query()
	if gameID > 0 {
		query = query.Where(gamecurrency.GameIDEQ(gameID))
	}
	if currencyID > 0 {
		query = query.Where(gamecurrency.CurrencyIDEQ(currencyID))
	}
	query = query.Where(gamecurrency.DeletedAtIsNil())

	// 获取分页数据
	currencies, err := query.All(ctx)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}

	return currencies, nil
}

// UpdateGameCurrency 更新游戏货币
func (d *GameCurrencyDAO) UpdateGameCurrency(ctx context.Context, id int64, updates map[string]interface{}) (*ent.GameCurrency, error) {
	update := d.db.GameCurrency.UpdateOneID(id)

	// 动态应用更新
	for key, value := range updates {
		switch key {
		case "game_id":
			if val, ok := value.(int64); ok {
				update = update.SetGameID(val)
			}

		case "currency_id":
			if val, ok := value.(int64); ok {
				update = update.SetCurrencyID(val)
			}
		case "status":
			if val, ok := value.(int32); ok {
				update = update.SetStatus(int64(val))
			}
		case "source_status":
			if val, ok := value.(int32); ok {
				update = update.SetSourceStatus(int64(val))
			}
		case "deleted_at":
			if val, ok := value.(time.Time); ok {
				update = update.SetDeletedAt(val)
			}
		}
	}

	return update.Save(ctx)
}
