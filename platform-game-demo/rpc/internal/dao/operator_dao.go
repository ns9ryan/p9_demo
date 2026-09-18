package dao

import (
	"context"
	"fmt"

	"oa.98ent.com/p9/platform-game/rpc/ent"
	"oa.98ent.com/p9/platform-game/rpc/ent/operator"
	"oa.98ent.com/p9/platform-game/rpc/ent/operatorgame"
	"oa.98ent.com/p9/platform-game/rpc/ent/operatorgamecategory"
	"oa.98ent.com/p9/platform-game/rpc/ent/operatorgamechannel"
	"oa.98ent.com/p9/platform-game/rpc/ent/operatorgameprovider"
)

type OperatorDAO struct {
	db *ent.Client
}

func NewOperatorDAO(db *ent.Client) *OperatorDAO {
	return &OperatorDAO{
		db: db,
	}
}

// GetAllOperators 获取所有分站列表
func (d *OperatorDAO) GetAllOperators(ctx context.Context) ([]*ent.Operator, error) {
	return d.db.Operator.Query().Order(operator.ByCode()).All(ctx)
}

// GetAllOperatorsByCode 获取指定分站编码的分站列表
func (d *OperatorDAO) GetAllOperatorsByCode(ctx context.Context, code string) ([]*ent.Operator, error) {
	if code == "" {
		return d.GetAllOperators(ctx)
	}
	return d.db.Operator.Query().
		Where(operator.CodeEQ(code)).
		Order(operator.ByCode()).
		All(ctx)
}

// GetOperatorByCode 根据分站编码获取分站信息
func (d *OperatorDAO) GetOperatorByCode(ctx context.Context, code string) (*ent.Operator, error) {
	return d.db.Operator.Query().
		Where(operator.CodeEQ(code)).
		Only(ctx)
}

// CountOperatorGamesByOpCode 统计分站的游戏数量
func (d *OperatorDAO) CountOperatorGamesByOpCode(ctx context.Context, opCode string) (int, error) {
	count, err := d.db.OperatorGame.Query().
		Where(operatorgame.OpCodeEQ(opCode)).
		Count(ctx)
	if err != nil {
		return 0, fmt.Errorf("count operator games failed: %w", err)
	}
	return count, nil
}

// CountOperatorGameCategoriesByOpCode 统计分站的游戏分类数量
func (d *OperatorDAO) CountOperatorGameCategoriesByOpCode(ctx context.Context, opCode string) (int, error) {
	count, err := d.db.OperatorGameCategory.Query().
		Where(operatorgamecategory.OpCodeEQ(opCode)).
		Count(ctx)
	if err != nil {
		return 0, fmt.Errorf("count operator game categories failed: %w", err)
	}
	return count, nil
}

// CountOperatorGameChannelsByOpCode 统计分站的游戏渠道数量
func (d *OperatorDAO) CountOperatorGameChannelsByOpCode(ctx context.Context, opCode string) (int, error) {
	count, err := d.db.OperatorGameChannel.Query().
		Where(operatorgamechannel.OpCodeEQ(opCode)).
		Count(ctx)
	if err != nil {
		return 0, fmt.Errorf("count operator game channels failed: %w", err)
	}
	return count, nil
}

// CountOperatorGameProvidersByOpCode 统计分站的游戏提供商数量
func (d *OperatorDAO) CountOperatorGameProvidersByOpCode(ctx context.Context, opCode string) (int, error) {
	count, err := d.db.OperatorGameProvider.Query().
		Where(operatorgameprovider.OpCodeEQ(opCode)).
		Count(ctx)
	if err != nil {
		return 0, fmt.Errorf("count operator game providers failed: %w", err)
	}
	return count, nil
}

// ExistByCode 检查分站是否存在
func (d *OperatorDAO) ExistByCode(ctx context.Context, code string) (bool, error) {
	count, err := d.db.Operator.Query().
		Where(operator.CodeEQ(code)).
		Count(ctx)
	if err != nil {
		return false, fmt.Errorf("check operator existence by code failed: %w", err)
	}
	return count > 0, nil
}
