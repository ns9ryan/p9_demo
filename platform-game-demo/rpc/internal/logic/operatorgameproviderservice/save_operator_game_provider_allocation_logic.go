package operatorgameproviderservicelogic

import (
	"context"

	"oa.98ent.com/p9/platform-game/rpc/internal/svc"
	"oa.98ent.com/p9/platform-game/rpc/pb/platform_game"

	"github.com/zeromicro/go-zero/core/logx"
)

type SaveOperatorGameProviderAllocationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSaveOperatorGameProviderAllocationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveOperatorGameProviderAllocationLogic {
	return &SaveOperatorGameProviderAllocationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 保存分站游戏提供商分配
func (l *SaveOperatorGameProviderAllocationLogic) SaveOperatorGameProviderAllocation(in *platform_game.SaveOperatorGameProviderAllocationRequest) (*platform_game.SaveOperatorGameProviderAllocationResp, error) {
	var resp platform_game.SaveOperatorGameProviderAllocationResp
	resp.Total = int64(len(in.GetItems()))

	if l.svcCtx == nil || l.svcCtx.DAOManager == nil {
		l.Errorf("[SaveOperatorGameProviderAllocation] database not available")
		return &resp, nil
	}

	opCode := in.GetOpCode()
	if opCode == "" {
		l.Errorf("[SaveOperatorGameProviderAllocation] op_code is empty")
		return &resp, nil
	}
	// 校验opCode是否存在
	isExist, err := l.svcCtx.DAOManager.Operator.ExistByCode(l.ctx, opCode)
	if err != nil || !isExist {
		l.Errorf("[SaveOperatorGameProviderAllocation] operator not found: op_code=%s, err=%v",
			opCode, err)
		return &resp, nil
	}

	for _, item := range in.GetItems() {
		if item == nil {
			resp.Failed++
			continue
		}

		providerCode := item.GetCode()
		checkStatus := item.GetCheckStatus()

		// 查询数据库
		existRecord, err := l.svcCtx.DAOManager.OperatorGameProvider.GetByOpCodeAndProviderCode(l.ctx, opCode, providerCode)
		if err != nil {
			l.Errorf("[SaveOperatorGameProviderAllocation] query failed: op_code=%s, provider_code=%s, err=%v",
				opCode, providerCode, err)
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
				err := l.svcCtx.DAOManager.OperatorGameProvider.DeleteAllocationByID(l.ctx, existRecord.ID)
				if err != nil {
					l.Errorf("[SaveOperatorGameProviderAllocation] delete failed: id=%d, err=%v",
						existRecord.ID, err)
					resp.Failed++
					continue
				}
				resp.Deleted++
			}
		} else {
			// 校验是否存在gameProvider
			isExist, err := l.svcCtx.DAOManager.GameProvider.ExistByCode(l.ctx, providerCode)
			if err != nil || !isExist {
				l.Errorf("[SaveOperatorGameProviderAllocation] game provider not found: provider_code=%s, err=%v",
					providerCode, err)
				resp.Failed++
				continue
			}
			if checkStatus == 1 {
				// 创建记录
				_, err := l.svcCtx.DAOManager.OperatorGameProvider.CreateAllocation(l.ctx, opCode, providerCode)
				if err != nil {
					l.Errorf("[SaveOperatorGameProviderAllocation] create failed: op_code=%s, provider_code=%s, err=%v",
						opCode, providerCode, err)
					resp.Failed++
					continue
				}
				resp.Created++
			}
		}
	}

	return &resp, nil
}
