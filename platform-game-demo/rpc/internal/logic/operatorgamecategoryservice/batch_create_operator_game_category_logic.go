package operatorgamecategoryservicelogic

import (
	"context"
	"time"

	"oa.98ent.com/p9/platform-game/rpc/ent"
	"oa.98ent.com/p9/platform-game/rpc/internal/svc"
	"oa.98ent.com/p9/platform-game/rpc/pb/platform_game"

	"github.com/zeromicro/go-zero/core/logx"
)

type BatchCreateOperatorGameCategoryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBatchCreateOperatorGameCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchCreateOperatorGameCategoryLogic {
	return &BatchCreateOperatorGameCategoryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 批量创建分站游戏分类
func (l *BatchCreateOperatorGameCategoryLogic) BatchCreateOperatorGameCategory(in *platform_game.BatchCreateOperatorGameCategoryRequest) (*platform_game.BatchCreateOperatorGameCategoryResp, error) {
	l.Infof("[RPC BatchCreateOperatorGameCategory] received req: %d items", len(in.Items))
	if l.svcCtx == nil || l.svcCtx.DAOManager == nil {
		l.Errorf("[RPC BatchCreateOperatorGameCategory] DAO Manager not available")
		return &platform_game.BatchCreateOperatorGameCategoryResp{}, nil
	}

	if len(in.Items) == 0 {
		return &platform_game.BatchCreateOperatorGameCategoryResp{
			Total:   0,
			Success: 0,
			Failed:  0,
		}, nil
	}

	// 构建创建对象列表
	createList := make([]*ent.OperatorGameCategoryCreate, 0, len(in.Items))
	for _, item := range in.Items {
		// 判断对应code的分站是否存在
		exists, err := l.svcCtx.DAOManager.Operator.ExistByCode(l.ctx, item.OpCode)
		if err != nil {
			l.Errorf("[RPC BatchCreateOperatorGameCategory] 检查分站 %s 是否存在失败: %v", item.OpCode, err)
			continue
		}
		if !exists {
			l.Errorf("[RPC BatchCreateOperatorGameCategory] 分站 %s 不存在", item.OpCode)
			continue
		}

		// 判断对应code的游戏是否存在
		exists, err = l.svcCtx.DAOManager.GameCategory.ExistByCode(l.ctx, item.CategoryCode)
		if err != nil {
			l.Errorf("[RPC BatchCreateOperatorGameCategory] 检查游戏分类 %s 是否存在失败: %v", item.CategoryCode, err)
			continue
		}
		if !exists {
			l.Errorf("[RPC BatchCreateOperatorGameCategory] 游戏分类 %s 不存在", item.CategoryCode)
			continue
		}
		// 判断是否已经存在相同的分站游戏记录
		exists, err = l.svcCtx.DAOManager.OperatorGameCategory.ExistByOpCodeAndCategoryCode(l.ctx, item.OpCode, item.CategoryCode)
		if err != nil {
			l.Errorf("[RPC BatchCreateOperatorGameCategory] 检查分站游戏分类 %s-%s 是否存在失败: %v", item.OpCode, item.CategoryCode, err)
			continue
		}
		if exists {
			l.Errorf("[RPC BatchCreateOperatorGameCategory] 分站游戏分类 %s-%s 已存在", item.OpCode, item.CategoryCode)
			continue
		}
		createList = append(createList, l.svcCtx.DB.OperatorGameCategory.Create().
			SetOpCode(item.OpCode).
			SetCategoryCode(item.CategoryCode).
			SetStatus(int16(item.Status)).
			SetCreatedAt(time.Now()).
			SetUpdatedAt(time.Now()))
	}
	if len(createList) == 0 {
		l.Errorf("[RPC BatchCreateOperatorGameCategory] no valid items to create")
		return &platform_game.BatchCreateOperatorGameCategoryResp{
			Total:   int64(len(in.Items)),
			Success: 0,
			Failed:  int64(len(in.Items)) - int64(len(createList)),
		}, nil
	}
	records, err := l.svcCtx.DAOManager.OperatorGameCategory.BatchCreateOperatorGameCategory(l.ctx, createList)
	if err != nil {
		l.Errorf("[RPC BatchCreateOperatorGameCategory] batch create failed: %v", err)
		return &platform_game.BatchCreateOperatorGameCategoryResp{
			Total:   int64(len(in.Items)),
			Success: 0,
			Failed:  int64(len(in.Items)),
		}, nil
	}

	l.Infof("[RPC BatchCreateOperatorGameCategory] success: created %d records", len(records))
	return &platform_game.BatchCreateOperatorGameCategoryResp{
		Total:   int64(len(in.Items)),
		Success: int64(len(records)),
		Failed:  int64(len(in.Items)) - int64(len(records)),
	}, nil
}
