package dao

import (
	"context"
	"fmt"

	"oa.98ent.com/p9/platform-game/rpc/ent"
	"oa.98ent.com/p9/platform-game/rpc/ent/operatorgameprovider"
)

type OperatorGameProviderDAO struct {
	db *ent.Client
}

func NewOperatorGameProviderDAO(db *ent.Client) *OperatorGameProviderDAO {
	return &OperatorGameProviderDAO{
		db: db,
	}
}

// GetOperatorGameProviderList 获取分站游戏提供商列表
func (d *OperatorGameProviderDAO) GetOperatorGameProviderList(ctx context.Context, opCode, providerCode string, status int32, offset, limit int64) ([]*ent.OperatorGameProvider, int, error) {
	query := d.db.OperatorGameProvider.Query()

	if opCode != "" {
		query = query.Where(operatorgameprovider.OpCodeEQ(opCode))
	}

	if providerCode != "" {
		query = query.Where(operatorgameprovider.ProviderCodeEQ(providerCode))
	}

	if status > 0 {
		query = query.Where(operatorgameprovider.StatusEQ(int16(status)))
	}

	// 获取总数
	total, err := query.Clone().Count(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("count failed: %w", err)
	}

	// 获取分页数据
	records, err := query.
		Offset(int(offset)).
		Limit(int(limit)).
		All(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("query failed: %w", err)
	}

	return records, total, nil
}

// BatchCreateOperatorGameProvider 批量创建分站游戏提供商
func (d *OperatorGameProviderDAO) BatchCreateOperatorGameProvider(ctx context.Context, items []*ent.OperatorGameProviderCreate) ([]*ent.OperatorGameProvider, error) {
	if len(items) == 0 {
		return []*ent.OperatorGameProvider{}, nil
	}

	records, err := d.db.OperatorGameProvider.CreateBulk(items...).Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("batch create failed: %w", err)
	}
	return records, nil
}

// BatchUpdateOperatorGameProviderStatus 批量修改分站游戏提供商状态
func (d *OperatorGameProviderDAO) BatchUpdateOperatorGameProviderStatus(ctx context.Context, ids []int64, status int16) (int, error) {
	if len(ids) == 0 {
		return 0, nil
	}

	affected, err := d.db.OperatorGameProvider.Update().
		Where(operatorgameprovider.IDIn(ids...)).
		SetStatus(status).
		Save(ctx)
	if err != nil {
		return 0, fmt.Errorf("batch update status failed: %w", err)
	}
	return affected, nil
}

func (d *OperatorGameProviderDAO) ExistByOpCodeAndProviderCode(ctx context.Context, opCode, providerCode string) (bool, error) {
	exists, err := d.db.OperatorGameProvider.Query().
		Where(operatorgameprovider.OpCodeEQ(opCode)).
		Where(operatorgameprovider.ProviderCodeEQ(providerCode)).
		Exist(ctx)
	if err != nil {
		return false, fmt.Errorf("exist check failed: %w", err)
	}
	return exists, nil
}

func (d *OperatorGameProviderDAO) BatchDeleteOperatorGameProvider(ctx context.Context, ids []int64) (int, error) {
	if len(ids) == 0 {
		return 0, nil
	}

	affected, err := d.db.OperatorGameProvider.Delete().
		Where(operatorgameprovider.IDIn(ids...)).
		Exec(ctx)
	if err != nil {
		return 0, fmt.Errorf("batch delete failed: %w", err)
	}
	return affected, nil
}
