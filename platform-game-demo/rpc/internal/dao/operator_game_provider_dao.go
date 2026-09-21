package dao

import (
	"context"
	"fmt"
	"time"

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

// GetByOpCodeAndProviderCode 获取分配记录
func (d *OperatorGameProviderDAO) GetByOpCodeAndProviderCode(ctx context.Context, opCode, providerCode string) (*ent.OperatorGameProvider, error) {
	record, err := d.db.OperatorGameProvider.Query().
		Where(
			operatorgameprovider.OpCodeEQ(opCode),
			operatorgameprovider.ProviderCodeEQ(providerCode),
		).
		Only(ctx)
	if err != nil && !ent.IsNotFound(err) {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	return record, nil
}

// CreateAllocation 创建分配记录
func (d *OperatorGameProviderDAO) CreateAllocation(ctx context.Context, opCode, providerCode string) (*ent.OperatorGameProvider, error) {
	record, err := d.db.OperatorGameProvider.Create().
		SetOpCode(opCode).
		SetProviderCode(providerCode).
		SetStatus(1).
		SetCreatedAt(time.Now()).
		SetUpdatedAt(time.Now()).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("create allocation failed: %w", err)
	}
	return record, nil
}

// DeleteAllocationByID 删除分配记录
func (d *OperatorGameProviderDAO) DeleteAllocationByID(ctx context.Context, id int64) error {
	err := d.db.OperatorGameProvider.DeleteOneID(id).Exec(ctx)
	if err != nil {
		return fmt.Errorf("delete allocation failed: %w", err)
	}
	return nil
}

// FindAllByOpCodeAndProviderCodes 根据opCode和GameProvider表的provider_code列表查询对应OperatorGameProvider表中的记录
func (d *OperatorGameProviderDAO) FindAllByOpCodeAndProviderCodes(ctx context.Context, opCode string, providerCodes []string) ([]*ent.OperatorGameProvider, error) {
	if len(providerCodes) == 0 {
		return []*ent.OperatorGameProvider{}, nil
	}

	records, err := d.db.OperatorGameProvider.Query().
		Where(operatorgameprovider.OpCodeEQ(opCode)).
		Where(operatorgameprovider.ProviderCodeIn(providerCodes...)).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	return records, nil
}

func (d *OperatorGameProviderDAO) FindAll(ctx context.Context, opCode, providerCode string, offset, limit int64) ([]*ent.OperatorGameProvider, error) {
	query := d.db.OperatorGameProvider.Query().
		Where(operatorgameprovider.OpCodeEQ(opCode)).
		Offset(int(offset)).
		Limit(int(limit))
	if providerCode != "" {
		query = query.Where(operatorgameprovider.ProviderCodeEQ(providerCode))
	}
	records, err := query.All(ctx)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	return records, nil
}
