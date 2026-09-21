package operatorgamechannelservicelogic

import (
	"context"

	"oa.98ent.com/p9/platform-game/rpc/internal/svc"
	"oa.98ent.com/p9/platform-game/rpc/pb/platform_game"

	"github.com/zeromicro/go-zero/core/logx"
)

type SaveOperatorGameChannelAllocationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSaveOperatorGameChannelAllocationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveOperatorGameChannelAllocationLogic {
	return &SaveOperatorGameChannelAllocationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 保存分站游戏渠道分配
func (l *SaveOperatorGameChannelAllocationLogic) SaveOperatorGameChannelAllocation(in *platform_game.SaveOperatorGameChannelAllocationRequest) (*platform_game.SaveOperatorGameChannelAllocationResp, error) {
	var resp platform_game.SaveOperatorGameChannelAllocationResp
	resp.Total = int64(len(in.GetItems()))

	if l.svcCtx == nil || l.svcCtx.DAOManager == nil {
		l.Errorf("[SaveOperatorGameChannelAllocation] database not available")
		return &resp, nil
	}

	opCode := in.GetOpCode()
	if opCode == "" {
		l.Errorf("[SaveOperatorGameChannelAllocation] op_code is empty")
		return &resp, nil
	}
	// 校验opCode是否存在
	isExist, err := l.svcCtx.DAOManager.Operator.ExistByCode(l.ctx, opCode)
	if err != nil || !isExist {
		l.Errorf("[SaveOperatorGameChannelAllocation] operator not found: op_code=%s, err=%v",
			opCode, err)
		return &resp, nil
	}

	for _, item := range in.GetItems() {
		if item == nil {
			resp.Failed++
			continue
		}

		channelCode := item.GetCode()
		checkStatus := item.GetCheckStatus()

		// 查询数据库
		existRecord, err := l.svcCtx.DAOManager.OperatorGameChannel.GetByOpCodeAndChannelCode(l.ctx, opCode, channelCode)
		if err != nil {
			l.Errorf("[SaveOperatorGameChannelAllocation] query failed: op_code=%s, channel_code=%s, err=%v",
				opCode, channelCode, err)
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
				err := l.svcCtx.DAOManager.OperatorGameChannel.DeleteAllocationByID(l.ctx, existRecord.ID)
				if err != nil {
					l.Errorf("[SaveOperatorGameChannelAllocation] delete failed: id=%d, err=%v",
						existRecord.ID, err)
					resp.Failed++
					continue
				}
				resp.Deleted++
			}
		} else {
			// 校验是否存在gameChannel
			isExist, err := l.svcCtx.DAOManager.GameChannel.ExistByCode(l.ctx, channelCode)
			if err != nil || !isExist {
				l.Errorf("[SaveOperatorGameChannelAllocation] game channel not found: channel_code=%s, err=%v",
					channelCode, err)
				resp.Failed++
				continue
			}
			if checkStatus == 1 {
				// 创建记录
				_, err := l.svcCtx.DAOManager.OperatorGameChannel.CreateAllocation(l.ctx, opCode, channelCode)
				if err != nil {
					l.Errorf("[SaveOperatorGameChannelAllocation] create failed: op_code=%s, channel_code=%s, err=%v",
						opCode, channelCode, err)
					resp.Failed++
					continue
				}
				resp.Created++
			}
		}
	}

	return &resp, nil
}
