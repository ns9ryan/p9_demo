package dao

import (
	"context"
	"fmt"
	"time"

	"oa.98ent.com/p9/platform-game/rpc/ent"
	"oa.98ent.com/p9/platform-game/rpc/ent/operatorgamecategory"
)

type OperatorGameCategoryDAO struct {
	db *ent.Client
}

func NewOperatorGameCategoryDAO(db *ent.Client) *OperatorGameCategoryDAO {
	return &OperatorGameCategoryDAO{
		db: db,
	}
}

// GetOperatorGameCategoryList 获取分站游戏分类列表
func (d *OperatorGameCategoryDAO) GetOperatorGameCategoryList(ctx context.Context, opCode, categoryCode string, status int32, offset, limit int64) ([]*ent.OperatorGameCategory, int, error) {
	query := d.db.OperatorGameCategory.Query()

	if opCode != "" {
		query = query.Where(operatorgamecategory.OpCodeEQ(opCode))
	}

	if categoryCode != "" {
		query = query.Where(operatorgamecategory.CategoryCodeEQ(categoryCode))
	}

	if status > 0 {
		query = query.Where(operatorgamecategory.StatusEQ(int16(status)))
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

// BatchCreateOperatorGameCategory 批量创建分站游戏分类
func (d *OperatorGameCategoryDAO) BatchCreateOperatorGameCategory(ctx context.Context, items []*ent.OperatorGameCategoryCreate) ([]*ent.OperatorGameCategory, error) {
	if len(items) == 0 {
		return []*ent.OperatorGameCategory{}, nil
	}

	records, err := d.db.OperatorGameCategory.CreateBulk(items...).Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("batch create failed: %w", err)
	}
	return records, nil
}

// BatchUpdateOperatorGameCategoryStatus 批量修改分站游戏分类状态
func (d *OperatorGameCategoryDAO) BatchUpdateOperatorGameCategoryStatus(ctx context.Context, ids []int64, status int16) (int, error) {
	if len(ids) == 0 {
		return 0, nil
	}

	affected, err := d.db.OperatorGameCategory.Update().
		Where(operatorgamecategory.IDIn(ids...)).
		SetStatus(status).
		Save(ctx)
	if err != nil {
		return 0, fmt.Errorf("batch update status failed: %w", err)
	}
	return affected, nil
}

func (d *OperatorGameCategoryDAO) ExistByOpCodeAndCategoryCode(ctx context.Context, opCode, categoryCode string) (bool, error) {
	exists, err := d.db.OperatorGameCategory.Query().
		Where(operatorgamecategory.OpCodeEQ(opCode)).
		Where(operatorgamecategory.CategoryCodeEQ(categoryCode)).
		Exist(ctx)
	if err != nil {
		return false, fmt.Errorf("exist check failed: %w", err)
	}
	return exists, nil
}

func (d *OperatorGameCategoryDAO) BatchDeleteOperatorGameCategory(ctx context.Context, ids []int64) (int, error) {
	if len(ids) == 0 {
		return 0, nil
	}

	affected, err := d.db.OperatorGameCategory.Delete().
		Where(operatorgamecategory.IDIn(ids...)).
		Exec(ctx)
	if err != nil {
		return 0, fmt.Errorf("batch delete failed: %w", err)
	}
	return affected, nil
}

// GetByOpCodeAndCategoryCode 获取分配记录
func (d *OperatorGameCategoryDAO) GetByOpCodeAndCategoryCode(ctx context.Context, opCode, categoryCode string) (*ent.OperatorGameCategory, error) {
	record, err := d.db.OperatorGameCategory.Query().
		Where(
			operatorgamecategory.OpCodeEQ(opCode),
			operatorgamecategory.CategoryCodeEQ(categoryCode),
		).
		Only(ctx)
	if err != nil && !ent.IsNotFound(err) {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	return record, nil
}

// CreateAllocation 创建分配记录
func (d *OperatorGameCategoryDAO) CreateAllocation(ctx context.Context, opCode, categoryCode string) (*ent.OperatorGameCategory, error) {
	record, err := d.db.OperatorGameCategory.Create().
		SetOpCode(opCode).
		SetCategoryCode(categoryCode).
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
func (d *OperatorGameCategoryDAO) DeleteAllocationByID(ctx context.Context, id int64) error {
	err := d.db.OperatorGameCategory.DeleteOneID(id).Exec(ctx)
	if err != nil {
		return fmt.Errorf("delete allocation failed: %w", err)
	}
	return nil
}

// FindAllByOpCodeAndCategoryCodes 根据opCode和GameCategory表的category_code列表查询对应OperatorGameCategory表中的记录
func (d *OperatorGameCategoryDAO) FindAllByOpCodeAndCategoryCodes(ctx context.Context, opCode string, categoryCodes []string) ([]*ent.OperatorGameCategory, error) {
	if len(categoryCodes) == 0 {
		return []*ent.OperatorGameCategory{}, nil
	}

	records, err := d.db.OperatorGameCategory.Query().
		Where(operatorgamecategory.OpCodeEQ(opCode)).
		Where(operatorgamecategory.CategoryCodeIn(categoryCodes...)).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	return records, nil
}

func (d *OperatorGameCategoryDAO) FindAll(ctx context.Context, opCode, categoryCode string, offset, limit int64) ([]*ent.OperatorGameCategory, error) {
	query := d.db.OperatorGameCategory.Query().
		Where(operatorgamecategory.OpCodeEQ(opCode)).
		Offset(int(offset)).
		Limit(int(limit))
	if categoryCode != "" {
		query = query.Where(operatorgamecategory.CategoryCodeEQ(categoryCode))
	}
	records, err := query.All(ctx)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	return records, nil
}
