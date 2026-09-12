package gamecategoryservicelogic

import (
	"context"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
	"oa.98ent.com/p9/platform-game/rpc/ent"
	"oa.98ent.com/p9/platform-game/rpc/internal/constant"
	"oa.98ent.com/p9/platform-game/rpc/internal/logic"
	"oa.98ent.com/p9/platform-game/rpc/internal/svc"
	"oa.98ent.com/p9/platform-game/rpc/internal/utils"
	"oa.98ent.com/p9/platform-game/rpc/pb/platform_game"
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
func (l *UpdateGameCategoryLogic) UpdateGameCategory(in *platform_game.UpdateGameCategoryRequest) (*platform_game.UpdateGameCategoryResp, error) {
	l.Infof("[RPC UpdateGameCategory] received request: id=%d", in.Id)

	if l.svcCtx == nil || l.svcCtx.DB == nil {
		l.Errorf("[RPC UpdateGameCategory] Database not available")
		return &platform_game.UpdateGameCategoryResp{
			Code:    constant.CodeInternalError,
			Message: "Database not available",
		}, nil
	}

	category := &ent.GameCategory{}
	if err := l.svcCtx.DB.WithContext(l.ctx).
		Where("id = ? AND deleted_at IS NULL", in.Id).
		First(category).Error; err != nil {
		l.Errorf("[RPC UpdateGameCategory] failed: %v", err)
		return &platform_game.UpdateGameCategoryResp{
			Code:    constant.CodeInternalError,
			Message: err.Error(),
		}, nil
	}

	oldStatus := category.Status

	updateData := make(map[string]interface{})
	updateData["updated_at"] = time.Now()

	if in.NameI18N != "" {
		updateData["name_i18n"] = in.NameI18N
	}
	if in.SortNo > 0 {
		updateData["sort_no"] = in.SortNo
	}
	if in.Status > 0 {
		updateData["status"] = in.Status
	}

	if err := l.svcCtx.DB.WithContext(l.ctx).Model(category).Updates(updateData).Error; err != nil {
		return &platform_game.UpdateGameCategoryResp{
			Code:    constant.CodeInternalError,
			Message: err.Error(),
		}, nil
	}

	l.svcCtx.DB.WithContext(l.ctx).Where("id = ?", in.Id).First(category)

	// 如果需要强制踢线且状态变更为停用，发送 kafka 事件
	if in.ForceLogout && in.Status == 2 && oldStatus != 2 {
		go func() {
			l.Infof("[RPC UpdateGameCategory] sending force-quit event for category id=%d", category.SourceId)
			// 查出所有游戏分类下的游戏，并发送强制踢线事件
			gameIds, err := getGamesByCategoryID(context.Background(), l.svcCtx.DB, category.SourceId)
			if err != nil {
				l.Errorf("[RPC UpdateGameCategory] failed to get games by category id=%d: %v", category.SourceId, err)
				return
			}
			forceQuit := utils.ForceQuitEvent{
				Basis:   utils.KafkaPayloadBasis{Stage: 0, Producer: "platform-game-rpc"},
				Scope:   constant.ScopeGameCategory,
				CatID:   category.SourceId,
				GameIDs: gameIds,
				QuitAt:  time.Now(),
			}
			l.Infof("[RPC UpdateGameCategory] sending force-quit event json: %s", utils.JSON(forceQuit))
			err = utils.SendForceQuitEvent(context.Background(), l.svcCtx.Config.Kafka.Brokers, &forceQuit)
			if err != nil {
				l.Errorf("[RPC UpdateGameCategory] failed to send force-quit event for category id=%d: %v", category.SourceId, err)
			}
		}()
	}

	return &platform_game.UpdateGameCategoryResp{
		Code:    constant.CodeSuccess,
		Message: "ok",
		Data:    logic.CategoryModelToProto(category),
	}, nil
}

func getGamesByCategoryID(ctx context.Context, db *gorm.DB, categoryID int64) ([]int64, error) {
	var games []ent.Game
	err := db.WithContext(ctx).
		Where("category_id = ? AND deleted_at IS NULL", categoryID).
		Find(&games).Error
	if err != nil {
		return nil, err
	}
	ids := make([]int64, len(games))
	for i, game := range games {
		ids[i] = game.SourceId.Int64
	}
	return ids, nil
}
