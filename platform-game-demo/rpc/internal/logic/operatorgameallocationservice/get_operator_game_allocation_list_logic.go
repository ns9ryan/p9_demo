package operatorgameallocationservicelogic

import (
	"context"
	"fmt"
	"sync"

	"oa.98ent.com/p9/platform-game/rpc/internal/constant"
	"oa.98ent.com/p9/platform-game/rpc/internal/svc"
	"oa.98ent.com/p9/platform-game/rpc/internal/utils"
	"oa.98ent.com/p9/platform-game/rpc/pb/platform_game"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetOperatorGameAllocationListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetOperatorGameAllocationListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOperatorGameAllocationListLogic {
	return &GetOperatorGameAllocationListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取分站游戏分配列表
func (l *GetOperatorGameAllocationListLogic) GetOperatorGameAllocationList(in *platform_game.GetOperatorGameAllocationListRequest) (*platform_game.GetOperatorGameAllocationListResp, error) {
	l.Infof("[RPC GetOperatorGameAllocationList] received req: page=%d, page_size=%d", in.Page, in.PageSize)
	if l.svcCtx == nil || l.svcCtx.DAOManager == nil {
		l.Errorf("[RPC GetOperatorGameAllocationList] DAO Manager not available")
		return &platform_game.GetOperatorGameAllocationListResp{
			Code:    constant.CodeInternalError,
			Message: "DAO Manager not available",
		}, nil
	}

	page, pageSize := utils.HandlePage(int64(in.Page), int64(in.PageSize))

	// 获取所有分站
	operators, err := l.svcCtx.DAOManager.Operator.GetAllOperatorsByCode(l.ctx, in.OpCode)
	if err != nil {
		l.Errorf("[RPC GetOperatorGameAllocationList] query operators failed: %v", err)
		return &platform_game.GetOperatorGameAllocationListResp{
			Code:    constant.CodeInternalError,
			Message: fmt.Sprintf("query failed: %v", err),
		}, nil
	}

	total := int64(len(operators))

	// 计算分页范围
	offset := (page - 1) * pageSize
	endOffset := offset + pageSize
	if endOffset > total {
		endOffset = total
	}
	if offset > total {
		offset = total
	}

	// 并发查询每个分站的相关数据
	items := make([]*platform_game.OperatorGameAllocationInfo, 0, endOffset-offset)
	var wg sync.WaitGroup
	var mu sync.Mutex

	for i := int64(offset); i < endOffset; i++ {
		wg.Add(1)
		idx := i
		go func() {
			defer wg.Done()
			operator := operators[idx]
			cacheKey := "operator_game_allocation:" + operator.Code

			// 先查询缓存
			if cachedData, found := l.svcCtx.CacheManager.Get(cacheKey); found {
				if info, ok := cachedData.(*platform_game.OperatorGameAllocationInfo); ok {
					mu.Lock()
					items = append(items, info)
					mu.Unlock()
					return
				}
			}

			// 获取四种关联数据的数量
			gameCount, _ := l.svcCtx.DAOManager.Operator.CountOperatorGamesByOpCode(l.ctx, operator.Code)
			categoryCount, _ := l.svcCtx.DAOManager.Operator.CountOperatorGameCategoriesByOpCode(l.ctx, operator.Code)
			channelCount, _ := l.svcCtx.DAOManager.Operator.CountOperatorGameChannelsByOpCode(l.ctx, operator.Code)
			providerCount, _ := l.svcCtx.DAOManager.Operator.CountOperatorGameProvidersByOpCode(l.ctx, operator.Code)

			info := &platform_game.OperatorGameAllocationInfo{
				Id:                operator.ID,
				OpCode:            operator.Code,
				OpName:            operator.Name,
				GameCount:         int32(gameCount),
				GameCategoryCount: int32(categoryCount),
				GameChannelCount:  int32(channelCount),
				GameProviderCount: int32(providerCount),
				UpdatedAt:         operator.UpdatedAt.Unix(),
			}

			// 缓存2分钟
			l.svcCtx.CacheManager.Set(cacheKey, info, 120)

			mu.Lock()
			items = append(items, info)
			mu.Unlock()
		}()
	}

	wg.Wait()

	// 按分站编码排序
	for i := 0; i < len(items)-1; i++ {
		for j := i + 1; j < len(items); j++ {
			if items[i].OpCode > items[j].OpCode {
				items[i], items[j] = items[j], items[i]
			}
		}
	}

	l.Infof("[RPC GetOperatorGameAllocationList] success: total=%d, returned=%d", total, len(items))
	return &platform_game.GetOperatorGameAllocationListResp{
		Code:     constant.CodeSuccess,
		Message:  "success",
		Items:    items,
		Total:    total,
		Page:     int32(page),
		PageSize: int32(pageSize),
	}, nil
}
