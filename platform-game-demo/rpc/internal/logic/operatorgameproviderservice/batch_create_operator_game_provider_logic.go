package operatorgameproviderservicelogic

import (
	"context"
	"time"

	"oa.98ent.com/p9/platform-game/rpc/ent"
	"oa.98ent.com/p9/platform-game/rpc/internal/svc"
	"oa.98ent.com/p9/platform-game/rpc/pb/platform_game"

	"github.com/zeromicro/go-zero/core/logx"
)

type BatchCreateOperatorGameProviderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBatchCreateOperatorGameProviderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchCreateOperatorGameProviderLogic {
	return &BatchCreateOperatorGameProviderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 批量创建分站游戏提供商
func (l *BatchCreateOperatorGameProviderLogic) BatchCreateOperatorGameProvider(in *platform_game.BatchCreateOperatorGameProviderRequest) (*platform_game.BatchCreateOperatorGameProviderResp, error) {
	l.Infof("[RPC BatchCreateOperatorGameProvider] received req: %d items", len(in.Items))
	if l.svcCtx == nil || l.svcCtx.DAOManager == nil {
		l.Errorf("[RPC BatchCreateOperatorGameProvider] DAO Manager not available")
		return &platform_game.BatchCreateOperatorGameProviderResp{}, nil
	}

	if len(in.Items) == 0 {
		return &platform_game.BatchCreateOperatorGameProviderResp{
			Total:   0,
			Success: 0,
			Failed:  0,
		}, nil
	}

	// 构建创建对象列表
	createList := make([]*ent.OperatorGameProviderCreate, 0, len(in.Items))
	for _, item := range in.Items {
		// 判断对应code的分站是否存在
		exists, err := l.svcCtx.DAOManager.Operator.ExistByCode(l.ctx, item.OpCode)
		if err != nil {
			l.Errorf("[RPC BatchCreateOperatorGameProvider] 检查分站 %s 是否存在失败: %v", item.OpCode, err)
			continue
		}
		if !exists {
			l.Errorf("[RPC BatchCreateOperatorGameProvider] 分站 %s 不存在", item.OpCode)
			continue
		}

		// 判断对应code的游戏是否存在
		exists, err = l.svcCtx.DAOManager.GameProvider.ExistByCode(l.ctx, item.ProviderCode)
		if err != nil {
			l.Errorf("[RPC BatchCreateOperatorGameProvider] 检查游戏提供商 %s 是否存在失败: %v", item.ProviderCode, err)
			continue
		}
		if !exists {
			l.Errorf("[RPC BatchCreateOperatorGameProvider] 游戏提供商 %s 不存在", item.ProviderCode)
			continue
		}
		// 判断是否已经存在相同的分站游戏记录
		exists, err = l.svcCtx.DAOManager.OperatorGameProvider.ExistByOpCodeAndProviderCode(l.ctx, item.OpCode, item.ProviderCode)
		if err != nil {
			l.Errorf("[RPC BatchCreateOperatorGameProvider] 检查分站游戏提供商 %s-%s 是否存在失败: %v", item.OpCode, item.ProviderCode, err)
			continue
		}
		if exists {
			l.Errorf("[RPC BatchCreateOperatorGameProvider] 分站游戏提供商 %s-%s 已存在", item.OpCode, item.ProviderCode)
			continue
		}
		createList = append(createList, l.svcCtx.DB.OperatorGameProvider.Create().
			SetOpCode(item.OpCode).
			SetProviderCode(item.ProviderCode).
			SetStatus(int16(item.Status)).
			SetCreatedAt(time.Now()).
			SetUpdatedAt(time.Now()))
	}
	if len(createList) == 0 {
		l.Errorf("[RPC BatchCreateOperatorGameProvider] no valid items to create")
		return &platform_game.BatchCreateOperatorGameProviderResp{
			Total:   int64(len(in.Items)),
			Success: 0,
			Failed:  int64(len(in.Items)) - int64(len(createList)),
		}, nil
	}
	records, err := l.svcCtx.DAOManager.OperatorGameProvider.BatchCreateOperatorGameProvider(l.ctx, createList)
	if err != nil {
		l.Errorf("[RPC BatchCreateOperatorGameProvider] batch create failed: %v", err)
		return &platform_game.BatchCreateOperatorGameProviderResp{
			Total:   int64(len(in.Items)),
			Success: 0,
			Failed:  int64(len(in.Items)),
		}, nil
	}

	l.Infof("[RPC BatchCreateOperatorGameProvider] success: created %d records", len(records))
	return &platform_game.BatchCreateOperatorGameProviderResp{
		Total:   int64(len(in.Items)),
		Success: int64(len(records)),
		Failed:  int64(len(in.Items)) - int64(len(records)),
	}, nil
}
