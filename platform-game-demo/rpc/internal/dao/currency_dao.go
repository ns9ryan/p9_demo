package dao

import (
	"context"
	"fmt"

	"oa.98ent.com/p9/platform-game/rpc/ent"
	"oa.98ent.com/p9/platform-game/rpc/ent/currency"
)

type CurrencyDAO struct {
	db *ent.Client
}

func NewCurrencyDAO(db *ent.Client) *CurrencyDAO {
	return &CurrencyDAO{
		db: db,
	}
}

// 获取货币列表
func (d *CurrencyDAO) GetCurrencyList(ctx context.Context) ([]*ent.Currency, int, error) {
	query := d.db.Currency.Query()

	// 获取总数
	total, err := query.Clone().Count(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("count failed: %w", err)
	}

	// 获取分页数据
	currencies, err := query.
		Order(currency.ByID()).
		All(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("query failed: %w", err)
	}

	return currencies, total, nil
}

func (d *CurrencyDAO) GetCurrencyMap(ctx context.Context) (map[int64]*ent.Currency, error) {
	currencies, _, err := d.GetCurrencyList(ctx)
	if err != nil {
		return nil, err
	}

	currencyMap := make(map[int64]*ent.Currency)
	for _, c := range currencies {
		currencyMap[c.ID] = c
	}
	return currencyMap, nil
}

func (d *CurrencyDAO) GetCurrencyByID(ctx context.Context, id int64) (*ent.Currency, error) {
	return d.db.Currency.Query().Where(currency.IDEQ(id)).Only(ctx)
}
