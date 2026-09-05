package regionservicelogic

import (
	"context"

	"oa.98ent.com/p9/platform-base/pkg/i18nkey"
	"oa.98ent.com/p9/platform-base/pkg/rpc/grpcerror"
	entregion "oa.98ent.com/p9/platform-base/rpc/ent/region"
	"oa.98ent.com/p9/platform-base/rpc/internal/enterror"
	"oa.98ent.com/p9/platform-base/rpc/internal/svc"
	"oa.98ent.com/p9/platform-base/rpc/pb/base/region"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReorderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewReorderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReorderLogic {
	return &ReorderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// Reorder 调整国家地区排序
func (l *ReorderLogic) Reorder(in *region.ReorderRegionRequest) (*region.ReorderRegionResponse, error) {
	// 校验国家地区ID
	if in.Id <= 0 || in.TargetId <= 0 || in.Id == in.TargetId {
		return nil, grpcerror.InvalidArgument(i18nkey.ValidationError)
	}

	// 开启事务
	tx, err := l.svcCtx.DB.Tx(l.ctx)
	if err != nil {
		return nil, enterror.Handle(l.Logger, err)
	}

	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback()
		}
	}()

	// 获取当前国家地区排序
	results, err := tx.Region.
		Query().
		Order(
			entregion.BySortNo(), // 按排序值升序
			entregion.ByID(),     // 排序值相同时按ID升序
		).
		All(l.ctx)
	if err != nil {
		return nil, enterror.Handle(l.Logger, err)
	}

	// 查找移动国家地区和目标国家地区的位置
	sourceIndex := -1
	targetIndex := -1

	for index, result := range results {
		switch result.ID {
		case in.Id:
			sourceIndex = index
		case in.TargetId:
			targetIndex = index
		}
	}

	// 移动国家地区或目标国家地区不存在
	if sourceIndex == -1 || targetIndex == -1 {
		return nil, grpcerror.NotFound(i18nkey.DataNotFound)
	}

	// 保存需要移动的国家地区
	moved := results[sourceIndex]

	if sourceIndex < targetIndex {
		// 向下移动时将中间国家地区整体向前移动
		copy(
			results[sourceIndex:targetIndex],
			results[sourceIndex+1:targetIndex+1],
		)

		// 将移动国家地区放到目标国家地区后面
		results[targetIndex] = moved
	} else {
		// 向上移动时将中间国家地区整体向后移动
		copy(
			results[targetIndex+1:sourceIndex+1],
			results[targetIndex:sourceIndex],
		)

		// 将移动国家地区放到目标国家地区前面
		results[targetIndex] = moved
	}

	// 重新整理为连续排序值
	for index, result := range results {
		sortNo := int64(index + 1)

		// 排序值没有变化时不更新
		if result.SortNo == sortNo {
			continue
		}

		if err := tx.Region.
			UpdateOneID(result.ID).
			SetSortNo(sortNo). // 排序值
			Exec(l.ctx); err != nil {
			return nil, enterror.Handle(l.Logger, err)
		}
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		return nil, enterror.Handle(l.Logger, err)
	}
	committed = true

	return &region.ReorderRegionResponse{}, nil
}
