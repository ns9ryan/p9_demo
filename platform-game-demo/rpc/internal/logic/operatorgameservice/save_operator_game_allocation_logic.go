package operatorgameservicelogic

import (
	"context"

	"oa.98ent.com/p9/platform-game/rpc/internal/svc"
	"oa.98ent.com/p9/platform-game/rpc/pb/platform_game"

	"github.com/zeromicro/go-zero/core/logx"
)

type SaveOperatorGameAllocationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSaveOperatorGameAllocationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveOperatorGameAllocationLogic {
	return &SaveOperatorGameAllocationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 保存分站游戏分配
func (l *SaveOperatorGameAllocationLogic) SaveOperatorGameAllocation(in *platform_game.SaveOperatorGameAllocationRequest) (*platform_game.SaveOperatorGameAllocationResp, error) {
	var resp platform_game.SaveOperatorGameAllocationResp
	resp.Total = int64(len(in.GetItems()))

	if l.svcCtx == nil || l.svcCtx.DAOManager == nil {
		l.Errorf("[SaveOperatorGameAllocation] database not available")
		return &resp, nil
	}

	opCode := in.GetOpCode()
	if opCode == "" {
		l.Errorf("[SaveOperatorGameAllocation] op_code is empty")
		return &resp, nil
	}
	// 校验opCode是否存在
	isExist, err := l.svcCtx.DAOManager.Operator.ExistByCode(l.ctx, opCode)
	if err != nil || !isExist {
		l.Errorf("[SaveOperatorGameAllocation] operator not found: op_code=%s, err=%v",
			opCode, err)
		return &resp, nil
	}

	for _, item := range in.GetItems() {
		if item == nil {
			resp.Failed++
			continue
		}

		gameCode := item.GetCode()
		checkStatus := item.GetCheckStatus()

		// 查询数据库
		existRecord, err := l.svcCtx.DAOManager.OperatorGame.GetByOpCodeAndGameCode(l.ctx, opCode, gameCode)
		if err != nil {
			l.Errorf("[SaveOperatorGameAllocation] query failed: op_code=%s, game_code=%s, err=%v",
				opCode, gameCode, err)
			resp.Failed++
			continue
		}

		// 数据库存在
		if existRecord != nil {
			if checkStatus == 1 {
				// exist++
				resp.Exist++
			} else if checkStatus == 2 {
				// 删除记录
				err := l.svcCtx.DAOManager.OperatorGame.DeleteAllocationByID(l.ctx, existRecord.ID)
				if err != nil {
					l.Errorf("[SaveOperatorGameAllocation] delete failed: id=%d, err=%v",
						existRecord.ID, err)
					resp.Failed++
					continue
				}
				resp.Deleted++
			}
		} else {
			// 校验是否存在game
			gameRecord, err := l.svcCtx.DAOManager.Game.GetGameByCode(l.ctx, gameCode)
			if err != nil || gameRecord == nil {
				l.Errorf("[SaveOperatorGameAllocation] game not found: game_code=%s, err=%v",
					gameCode, err)
				resp.Failed++
				continue
			}

			if checkStatus == 1 {
				// 创建记录
				_, err := l.svcCtx.DAOManager.OperatorGame.CreateAllocation(l.ctx, gameRecord.Name, opCode, gameCode)
				if err != nil {
					l.Errorf("[SaveOperatorGameAllocation] create failed: op_code=%s, game_code=%s, err=%v",
						opCode, gameCode, err)
					resp.Failed++
					continue
				}
				resp.Created++
			}
		}
	}

	return &resp, nil
}
