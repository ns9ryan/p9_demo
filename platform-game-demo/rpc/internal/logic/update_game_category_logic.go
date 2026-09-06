package logic

import (
	"context"
	"time"

	"oa.98ent.com/p9/platform-game/common/constant"
	"oa.98ent.com/p9/platform-game/common/kafka"
	"oa.98ent.com/p9/platform-game/common/model"
	"oa.98ent.com/p9/platform-game/common/utils"
	"oa.98ent.com/p9/platform-game/rpc/internal/svc"
	platformgame "oa.98ent.com/p9/platform-game/rpc/pb/platform_game"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateGameCategoryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateGameCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateGameCategoryLogic {
	return &UpdateGameCategoryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新游戏分类
func (l *UpdateGameCategoryLogic) UpdateGameCategory(in *platformgame.UpdateGameCategoryRequest) (*platformgame.GetGameCategoryListResp, error) {
	l.Infof("[RPC UpdateGameCategory] received request: id=%d", in.Id)

	if l.svcCtx == nil || l.svcCtx.DB == nil {
		l.Errorf("[RPC UpdateGameCategory] Database not available")
		return &platformgame.GetGameCategoryListResp{
			Code:    constant.CodeInternalError,
			Message: "Database not available",
			Data:    nil,
		}, nil
	}

	// 查询现有记录
	var category model.Category
	if err := l.svcCtx.DB.WithContext(l.ctx).
		Where("id = ? AND deleted_at IS NULL", in.Id).
		First(&category).Error; err != nil {
		l.Errorf("[RPC UpdateGameCategory] failed to get category: %v", err)
		return &platformgame.GetGameCategoryListResp{
			Code:    constant.CodeInternalError,
			Message: "failed to get category: " + err.Error(),
			Data:    nil,
		}, nil
	}

	// 记录状态变更前的状态
	oldStatus := category.Status

	// 构建更新数据（只允许更新特定字段)
	updateData := make(map[string]interface{})
	updateData["updated_at"] = time.Now()

	if in.NameI18N != "" {
		updateData["name_i18n"] = utils.MustParseJSON([]byte(in.NameI18N))
	}
	if in.SortNo > 0 {
		updateData["sort_no"] = in.SortNo
	}
	if in.Status > 0 {
		updateData["status"] = in.Status
	}

	// 执行更新
	if err := l.svcCtx.DB.WithContext(l.ctx).
		Model(&category).
		Updates(updateData).Error; err != nil {
		l.Errorf("[RPC UpdateGameCategory] failed to update category: %v", err)
		return &platformgame.GetGameCategoryListResp{
			Code:    constant.CodeInternalError,
			Message: "failed to update category: " + err.Error(),
			Data:    nil,
		}, nil
	}

	// 重新查询最新数�?
	if err := l.svcCtx.DB.WithContext(l.ctx).
		Where("id = ?", in.Id).
		First(&category).Error; err != nil {
		l.Errorf("[RPC UpdateGameCategory] failed to get updated category: %v", err)
		return &platformgame.GetGameCategoryListResp{
			Code:    constant.CodeInternalError,
			Message: "failed to get updated category: " + err.Error(),
			Data:    nil,
		}, nil
	}

	// 转换�?proto message
	var categoryName string
	if category.NameI18n != nil {
		if name, ok := category.NameI18n["default"].(string); ok {
			categoryName = name
		} else if name, ok := category.NameI18n["en"].(string); ok {
			categoryName = name
		} else {
			for _, v := range category.NameI18n {
				if str, ok := v.(string); ok {
					categoryName = str
					break
				}
			}
		}
	}

	item := &platformgame.CategoryInfo{
		Id:     category.ID,
		Code:   category.CategoryCode,
		Name:   categoryName,
		Status: int32(category.Status),
	}

	l.Infof("[RPC UpdateGameCategory] updated successfully: id=%d", category.ID)

	// 如果需要强制踢线且状态变更为停用，发�?kafka 事件
	if in.ForceLogout && in.Status == 2 && oldStatus != 2 {
		go func() {
			l.Infof("[RPC UpdateGameCategory] sending force-quit event for category id=%d", category.ID)
			if err := kafka.SendForceQuitEvent(l.ctx, category.ID, "game_category"); err != nil {
				l.Infof("[RPC UpdateGameCategory] send force-quit event failed: %v", err)
			}
		}()
	}

	return &platformgame.GetGameCategoryListResp{
		Code:    constant.CodeSuccess,
		Message: "ok",
		Data:    []*platformgame.CategoryInfo{item},
	}, nil
}
