package dao

import (
	"context"
	"fmt"
	"time"

	"oa.98ent.com/p9/platform-game/rpc/ent"
	"oa.98ent.com/p9/platform-game/rpc/ent/operatorgamechannel"
)

type OperatorGameChannelDAO struct {
	db *ent.Client
}

func NewOperatorGameChannelDAO(db *ent.Client) *OperatorGameChannelDAO {
	return &OperatorGameChannelDAO{
		db: db,
	}
}

// GetOperatorGameChannelList 获取分站游戏渠道列表
func (d *OperatorGameChannelDAO) GetOperatorGameChannelList(ctx context.Context, opCode, channelCode string, status int32, offset, limit int64) ([]*ent.OperatorGameChannel, int, error) {
	query := d.db.OperatorGameChannel.Query()

	if opCode != "" {
		query = query.Where(operatorgamechannel.OpCodeEQ(opCode))
	}

	if channelCode != "" {
		query = query.Where(operatorgamechannel.ChannelCodeEQ(channelCode))
	}

	if status > 0 {
		query = query.Where(operatorgamechannel.StatusEQ(int16(status)))
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

// BatchCreateOperatorGameChannel 批量创建分站游戏渠道
func (d *OperatorGameChannelDAO) BatchCreateOperatorGameChannel(ctx context.Context, items []*ent.OperatorGameChannelCreate) ([]*ent.OperatorGameChannel, error) {
	if len(items) == 0 {
		return []*ent.OperatorGameChannel{}, nil
	}

	records, err := d.db.OperatorGameChannel.CreateBulk(items...).Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("batch create failed: %w", err)
	}
	return records, nil
}

// BatchUpdateOperatorGameChannelStatus 批量修改分站游戏渠道状态
func (d *OperatorGameChannelDAO) BatchUpdateOperatorGameChannelStatus(ctx context.Context, ids []int64, status int16) (int, error) {
	if len(ids) == 0 {
		return 0, nil
	}

	affected, err := d.db.OperatorGameChannel.Update().
		Where(operatorgamechannel.IDIn(ids...)).
		SetStatus(status).
		Save(ctx)
	if err != nil {
		return 0, fmt.Errorf("batch update status failed: %w", err)
	}
	return affected, nil
}

func (d *OperatorGameChannelDAO) ExistByOpCodeAndChannelCode(ctx context.Context, opCode, channelCode string) (bool, error) {
	exists, err := d.db.OperatorGameChannel.Query().
		Where(operatorgamechannel.OpCodeEQ(opCode)).
		Where(operatorgamechannel.ChannelCodeEQ(channelCode)).
		Exist(ctx)
	if err != nil {
		return false, fmt.Errorf("exist check failed: %w", err)
	}
	return exists, nil
}

func (d *OperatorGameChannelDAO) BatchDeleteOperatorGameChannel(ctx context.Context, ids []int64) (int, error) {
	if len(ids) == 0 {
		return 0, nil
	}

	affected, err := d.db.OperatorGameChannel.Delete().
		Where(operatorgamechannel.IDIn(ids...)).
		Exec(ctx)
	if err != nil {
		return 0, fmt.Errorf("batch delete failed: %w", err)
	}
	return affected, nil
}

// GetByOpCodeAndChannelCode 获取分配记录
func (d *OperatorGameChannelDAO) GetByOpCodeAndChannelCode(ctx context.Context, opCode, channelCode string) (*ent.OperatorGameChannel, error) {
	record, err := d.db.OperatorGameChannel.Query().
		Where(
			operatorgamechannel.OpCodeEQ(opCode),
			operatorgamechannel.ChannelCodeEQ(channelCode),
		).
		Only(ctx)
	if err != nil && !ent.IsNotFound(err) {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	return record, nil
}

// CreateAllocation 创建分配记录
func (d *OperatorGameChannelDAO) CreateAllocation(ctx context.Context, opCode, channelCode string) (*ent.OperatorGameChannel, error) {
	record, err := d.db.OperatorGameChannel.Create().
		SetOpCode(opCode).
		SetChannelCode(channelCode).
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
func (d *OperatorGameChannelDAO) DeleteAllocationByID(ctx context.Context, id int64) error {
	err := d.db.OperatorGameChannel.DeleteOneID(id).Exec(ctx)
	if err != nil {
		return fmt.Errorf("delete allocation failed: %w", err)
	}
	return nil
}

// FindAllByOpCodeAndChannelCodes 根据opCode和GameChannel表的channel_code列表查询对应OperatorGameChannel表中的记录
func (d *OperatorGameChannelDAO) FindAllByOpCodeAndChannelCodes(ctx context.Context, opCode string, channelCodes []string) ([]*ent.OperatorGameChannel, error) {
	if len(channelCodes) == 0 {
		return []*ent.OperatorGameChannel{}, nil
	}

	records, err := d.db.OperatorGameChannel.Query().
		Where(operatorgamechannel.OpCodeEQ(opCode)).
		Where(operatorgamechannel.ChannelCodeIn(channelCodes...)).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	return records, nil
}

func (d *OperatorGameChannelDAO) FindAll(ctx context.Context, opCode, channelCode string, offset, limit int64) ([]*ent.OperatorGameChannel, error) {
	query := d.db.OperatorGameChannel.Query().
		Where(operatorgamechannel.OpCodeEQ(opCode)).
		Offset(int(offset)).
		Limit(int(limit))
	if channelCode != "" {
		query = query.Where(operatorgamechannel.ChannelCodeEQ(channelCode))
	}
	records, err := query.All(ctx)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	return records, nil
}
