package dao

import (
	"context"
	"fmt"
	"time"

	"oa.98ent.com/p9/platform-game/rpc/ent"
	"oa.98ent.com/p9/platform-game/rpc/ent/operatorgame"
)

type OperatorGameDAO struct {
	db *ent.Client
}

func NewOperatorGameDAO(db *ent.Client) *OperatorGameDAO {
	return &OperatorGameDAO{
		db: db,
	}
}

// GetOperatorGameList 获取分站游戏列表
func (d *OperatorGameDAO) GetOperatorGameList(ctx context.Context, opCode, gameCode string, status int32, offset, limit int64) ([]*ent.OperatorGame, int, error) {
	query := d.db.OperatorGame.Query()

	if opCode != "" {
		query = query.Where(operatorgame.OpCodeEQ(opCode))
	}

	if gameCode != "" {
		query = query.Where(operatorgame.GameCodeEQ(gameCode))
	}

	if status > 0 {
		query = query.Where(operatorgame.StatusEQ(int16(status)))
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

// BatchCreateOperatorGame 批量创建分站游戏
func (d *OperatorGameDAO) BatchCreateOperatorGame(ctx context.Context, items []*ent.OperatorGameCreate) ([]*ent.OperatorGame, error) {
	if len(items) == 0 {
		return []*ent.OperatorGame{}, nil
	}

	records, err := d.db.OperatorGame.CreateBulk(items...).Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("batch create failed: %w", err)
	}
	return records, nil
}

// BatchUpdateOperatorGameStatus 批量修改分站游戏状态
func (d *OperatorGameDAO) BatchUpdateOperatorGameStatus(ctx context.Context, ids []int64, status int16) (int, error) {
	if len(ids) == 0 {
		return 0, nil
	}

	affected, err := d.db.OperatorGame.Update().
		Where(operatorgame.IDIn(ids...)).
		SetStatus(status).
		Save(ctx)
	if err != nil {
		return 0, fmt.Errorf("batch update status failed: %w", err)
	}
	return affected, nil
}

func (d *OperatorGameDAO) ExistByOpCodeAndGameCode(ctx context.Context, opCode, gameCode string) (bool, error) {
	count, err := d.db.OperatorGame.Query().
		Where(operatorgame.OpCodeEQ(opCode)).
		Where(operatorgame.GameCodeEQ(gameCode)).
		Count(ctx)
	if err != nil {
		return false, fmt.Errorf("check operator game existence by opCode and gameCode failed: %w", err)
	}
	return count > 0, nil
}

func (d *OperatorGameDAO) BatchDeleteOperatorGame(ctx context.Context, ids []int64) (int, error) {
	if len(ids) == 0 {
		return 0, nil
	}

	affected, err := d.db.OperatorGame.Delete().
		Where(operatorgame.IDIn(ids...)).
		Exec(ctx)
	if err != nil {
		return 0, fmt.Errorf("batch delete failed: %w", err)
	}
	return affected, nil
}

// GetByOpCodeAndGameCode 获取分配记录
func (d *OperatorGameDAO) GetByOpCodeAndGameCode(ctx context.Context, opCode, gameCode string) (*ent.OperatorGame, error) {
	record, err := d.db.OperatorGame.Query().
		Where(
			operatorgame.OpCodeEQ(opCode),
			operatorgame.GameCodeEQ(gameCode),
		).
		Only(ctx)
	if err != nil && !ent.IsNotFound(err) {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	return record, nil
}

// CreateAllocation 创建分配记录
func (d *OperatorGameDAO) CreateAllocation(ctx context.Context, name, opCode, gameCode string) (*ent.OperatorGame, error) {
	record, err := d.db.OperatorGame.Create().
		SetName(name).
		SetOpCode(opCode).
		SetGameCode(gameCode).
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
func (d *OperatorGameDAO) DeleteAllocationByID(ctx context.Context, id int64) error {
	err := d.db.OperatorGame.DeleteOneID(id).Exec(ctx)
	if err != nil {
		return fmt.Errorf("delete allocation failed: %w", err)
	}
	return nil
}

// FindAllByOpCodeAndGameCodes 根据opCode和Game表的game_code列表查询对应OperatorGame表中的记录
func (d *OperatorGameDAO) FindAllByOpCodeAndGameCodes(ctx context.Context, opCode string, gameCodes []string) ([]*ent.OperatorGame, error) {
	if len(gameCodes) == 0 {
		return []*ent.OperatorGame{}, nil
	}

	records, err := d.db.OperatorGame.Query().
		Where(operatorgame.OpCodeEQ(opCode)).
		Where(operatorgame.GameCodeIn(gameCodes...)).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	return records, nil
}

func (d *OperatorGameDAO) FindAll(ctx context.Context, opCode, gameCode string, offset, limit int64) ([]*ent.OperatorGame, error) {
	query := d.db.OperatorGame.Query().
		Where(operatorgame.OpCodeEQ(opCode)).
		Offset(int(offset)).
		Limit(int(limit))
	if gameCode != "" {
		query = query.Where(operatorgame.GameCodeEQ(gameCode))
	}
	records, err := query.All(ctx)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	return records, nil
}
