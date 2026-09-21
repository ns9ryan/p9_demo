package operatorgamecategoryservicelogic

import (
	"context"

	"oa.98ent.com/p9/platform-game/rpc/internal/svc"
	"oa.98ent.com/p9/platform-game/rpc/pb/platform_game"

	"github.com/zeromicro/go-zero/core/logx"
)

type SaveOperatorGameCategoryAllocationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSaveOperatorGameCategoryAllocationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveOperatorGameCategoryAllocationLogic {
	return &SaveOperatorGameCategoryAllocationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 保存分站游戏分类分配
func (l *SaveOperatorGameCategoryAllocationLogic) SaveOperatorGameCategoryAllocation(in *platform_game.SaveOperatorGameCategoryAllocationRequest) (*platform_game.SaveOperatorGameCategoryAllocationResp, error) {
	var resp platform_game.SaveOperatorGameCategoryAllocationResp
	resp.Total = int64(len(in.GetItems()))

	if l.svcCtx == nil || l.svcCtx.DAOManager == nil {
		l.Errorf("[SaveOperatorGameCategoryAllocation] database not available")
		return &resp, nil
	}

	opCode := in.GetOpCode()
	if opCode == "" {
		l.Errorf("[SaveOperatorGameCategoryAllocation] op_code is empty")
		return &resp, nil
	}
	// 校验opCode是否存在
	isExist, err := l.svcCtx.DAOManager.Operator.ExistByCode(l.ctx, opCode)
	if err != nil || !isExist {
		l.Errorf("[SaveOperatorGameCategoryAllocation] operator not found: op_code=%s, err=%v",
			opCode, err)
		return &resp, nil
	}

	for _, item := range in.GetItems() {
		if item == nil {
			resp.Failed++
			continue
		}

		categoryCode := item.GetCode()
		checkStatus := item.GetCheckStatus()

		// 查询数据库
		existRecord, err := l.svcCtx.DAOManager.OperatorGameCategory.GetByOpCodeAndCategoryCode(l.ctx, opCode, categoryCode)
		if err != nil {
			l.Errorf("[SaveOperatorGameCategoryAllocation] query failed: op_code=%s, category_code=%s, err=%v",
				opCode, categoryCode, err)
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
				err := l.svcCtx.DAOManager.OperatorGameCategory.DeleteAllocationByID(l.ctx, existRecord.ID)
				if err != nil {
					l.Errorf("[SaveOperatorGameCategoryAllocation] delete failed: id=%d, err=%v",
						existRecord.ID, err)
					resp.Failed++
					continue
				}
				resp.Deleted++
			}
		} else {
			// 校验是否存在gameCategory
			isExist, err := l.svcCtx.DAOManager.GameCategory.ExistByCode(l.ctx, categoryCode)
			if err != nil || !isExist {
				l.Errorf("[SaveOperatorGameCategoryAllocation] game category not found: category_code=%s, err=%v",
					categoryCode, err)
				resp.Failed++
				continue
			}
			if checkStatus == 1 {
				// 创建记录
				_, err := l.svcCtx.DAOManager.OperatorGameCategory.CreateAllocation(l.ctx, opCode, categoryCode)
				if err != nil {
					l.Errorf("[SaveOperatorGameCategoryAllocation] create failed: op_code=%s, category_code=%s, err=%v",
						opCode, categoryCode, err)
					resp.Failed++
					continue
				}
				resp.Created++
			}
		}
	}

	return &resp, nil
}
