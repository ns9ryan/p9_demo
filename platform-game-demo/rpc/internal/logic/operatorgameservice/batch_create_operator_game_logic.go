package operatorgameservicelogic

import (
	"context"
	"time"

	"oa.98ent.com/p9/platform-game/rpc/ent"
	"oa.98ent.com/p9/platform-game/rpc/internal/svc"
	"oa.98ent.com/p9/platform-game/rpc/pb/platform_game"

	"github.com/zeromicro/go-zero/core/logx"
)

type BatchCreateOperatorGameLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBatchCreateOperatorGameLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchCreateOperatorGameLogic {
	return &BatchCreateOperatorGameLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 批量创建分站游戏
func (l *BatchCreateOperatorGameLogic) BatchCreateOperatorGame(in *platform_game.BatchCreateOperatorGameRequest) (*platform_game.BatchCreateOperatorGameResp, error) {
	l.Infof("[RPC BatchCreateOperatorGame] received req: %d items", len(in.Items))
	if l.svcCtx == nil || l.svcCtx.DAOManager == nil {
		l.Errorf("[RPC BatchCreateOperatorGame] DAO Manager not available")
		return &platform_game.BatchCreateOperatorGameResp{}, nil
	}

	if len(in.Items) == 0 {
		return &platform_game.BatchCreateOperatorGameResp{
			Total:   0,
			Success: 0,
			Failed:  0,
		}, nil
	}

	// 构建创建对象列表
	createList := make([]*ent.OperatorGameCreate, 0, len(in.Items))
	for _, item := range in.Items {
		// 判断对应code的分站是否存在
		exists, err := l.svcCtx.DAOManager.Operator.ExistByCode(l.ctx, item.OpCode)
		if err != nil {
			l.Errorf("[RPC BatchCreateOperatorGame] 检查分站 %s 是否存在失败: %v", item.OpCode, err)
			continue
		}
		if !exists {
			l.Errorf("[RPC BatchCreateOperatorGame] 分站 %s 不存在", item.OpCode)
			continue
		}

		// 判断对应code的游戏是否存在
		exists, err = l.svcCtx.DAOManager.Game.ExistByCode(l.ctx, item.GameCode)
		if err != nil {
			l.Errorf("[RPC BatchCreateOperatorGame] 检查游戏 %s 是否存在失败: %v", item.GameCode, err)
			continue
		}
		if !exists {
			l.Errorf("[RPC BatchCreateOperatorGame] 游戏 %s 不存在", item.GameCode)
			continue
		}
		// 判断是否已经存在相同的分站游戏记录
		exists, err = l.svcCtx.DAOManager.OperatorGame.ExistByOpCodeAndGameCode(l.ctx, item.OpCode, item.GameCode)
		if err != nil {
			l.Errorf("[RPC BatchCreateOperatorGame] 检查分站游戏 %s-%s 是否存在失败: %v", item.OpCode, item.GameCode, err)
			continue
		}
		if exists {
			l.Errorf("[RPC BatchCreateOperatorGame] 分站游戏 %s-%s 已存在", item.OpCode, item.GameCode)
			continue
		}
		createList = append(createList, l.svcCtx.DB.OperatorGame.Create().
			SetOpCode(item.OpCode).
			SetGameCode(item.GameCode).
			SetName(item.Name).
			SetStatus(int16(item.Status)).
			SetCreatedAt(time.Now()).
			SetUpdatedAt(time.Now()))
	}
	if len(createList) == 0 {
		l.Errorf("[RPC BatchCreateOperatorGame] no valid items to create")
		return &platform_game.BatchCreateOperatorGameResp{
			Total:   int64(len(in.Items)),
			Success: 0,
			Failed:  int64(len(in.Items)) - int64(len(createList)),
		}, nil
	}
	records, err := l.svcCtx.DAOManager.OperatorGame.BatchCreateOperatorGame(l.ctx, createList)
	if err != nil {
		l.Errorf("[RPC BatchCreateOperatorGame] batch create failed: %v", err)
		return &platform_game.BatchCreateOperatorGameResp{
			Total:   int64(len(in.Items)),
			Success: 0,
			Failed:  int64(len(in.Items)),
		}, nil
	}

	l.Infof("[RPC BatchCreateOperatorGame] success: created %d records", len(records))
	return &platform_game.BatchCreateOperatorGameResp{
		Total:   int64(len(in.Items)),
		Success: int64(len(records)),
		Failed:  int64(len(in.Items)) - int64(len(records)),
	}, nil
}
