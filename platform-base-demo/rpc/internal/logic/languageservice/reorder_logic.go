package languageservicelogic

import (
	"context"

	"oa.98ent.com/p9/platform-base/pkg/i18nkey"
	"oa.98ent.com/p9/platform-base/pkg/rpc/grpcerror"
	entlanguage "oa.98ent.com/p9/platform-base/rpc/ent/language"
	"oa.98ent.com/p9/platform-base/rpc/internal/enterror"
	"oa.98ent.com/p9/platform-base/rpc/internal/svc"
	"oa.98ent.com/p9/platform-base/rpc/pb/base/language"

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

// Reorder 调整语言排序
func (l *ReorderLogic) Reorder(in *language.ReorderLanguageRequest) (*language.ReorderLanguageResponse, error) {
	// 校验语言ID
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

	// 获取当前语言排序
	results, err := tx.Language.
		Query().
		Order(
			entlanguage.BySortNo(), // 按排序值升序
			entlanguage.ByID(),     // 排序值相同时按ID升序
		).
		All(l.ctx)
	if err != nil {
		return nil, enterror.Handle(l.Logger, err)
	}

	// 查找移动语言和目标语言的位置
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

	// 移动语言或目标语言不存在
	if sourceIndex == -1 || targetIndex == -1 {
		return nil, grpcerror.NotFound(i18nkey.DataNotFound)
	}

	// 获取需要移动的语言
	moved := results[sourceIndex]

	if sourceIndex < targetIndex {
		// 向下移动: 放到目标语言后面,中间元素整体向前移动
		copy(
			results[sourceIndex:targetIndex],
			results[sourceIndex+1:targetIndex+1],
		)
		results[targetIndex] = moved
	} else {
		// 向上移动: 放到目标语言前面,中间元素整体向后移动
		copy(
			results[targetIndex+1:sourceIndex+1],
			results[targetIndex:sourceIndex],
		)
		results[targetIndex] = moved
	}

	// 重新整理为连续排序值
	for index, result := range results {
		sortNo := int64(index + 1)

		// 排序值没有变化时不更新
		if result.SortNo == sortNo {
			continue
		}

		err := tx.Language.
			UpdateOneID(result.ID).
			SetSortNo(sortNo). // 排序值
			Exec(l.ctx)
		if err != nil {
			return nil, enterror.Handle(l.Logger, err)
		}
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		return nil, enterror.Handle(l.Logger, err)
	}
	committed = true

	return &language.ReorderLanguageResponse{}, nil
}
