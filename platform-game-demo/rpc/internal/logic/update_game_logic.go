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

type UpdateGameLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateGameLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateGameLogic {
	return &UpdateGameLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新游戏
func (l *UpdateGameLogic) UpdateGame(in *platformgame.UpdateGameRequest) (*platformgame.GetGameResp, error) {
	l.Infof("[RPC UpdateGame] received request: id=%d", in.Id)

	if l.svcCtx == nil || l.svcCtx.DB == nil {
		return &platformgame.GetGameResp{
			Code:    constant.CodeInternalError,
			Message: "Database not available",
		}, nil
	}

	var game model.Game
	if err := l.svcCtx.DB.WithContext(l.ctx).
		Where("id = ? AND deleted_at IS NULL", in.Id).
		First(&game).Error; err != nil {
		l.Errorf("[RPC UpdateGame] failed: %v", err)
		return &platformgame.GetGameResp{
			Code:    constant.CodeInternalError,
			Message: err.Error(),
		}, nil
	}

	oldStatus := game.Status

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

	if err := l.svcCtx.DB.WithContext(l.ctx).Model(&game).Updates(updateData).Error; err != nil {
		return &platformgame.GetGameResp{
			Code:    constant.CodeInternalError,
			Message: err.Error(),
		}, nil
	}

	l.svcCtx.DB.WithContext(l.ctx).Where("id = ?", in.Id).First(&game)

	var gameName string
	if game.NameI18n != nil {
		if name, ok := game.NameI18n["default"].(string); ok {
			gameName = name
		} else if name, ok := game.NameI18n["en"].(string); ok {
			gameName = name
		}
	}

	// 如果需要强制踢线且状态变更为停用，发送 kafka 事件
	if in.ForceLogout && in.Status == 2 && oldStatus != 2 {
		go func() {
			l.Infof("[RPC UpdateGame] sending force-quit event for game id=%d", game.ID)
			if err := kafka.SendForceQuitEvent(l.ctx, game.ID, "game"); err != nil {
				l.Infof("[RPC UpdateGame] send force-quit event failed: %v", err)
			}
		}()
	}

	return &platformgame.GetGameResp{
		Code:    constant.CodeSuccess,
		Message: "ok",
		Data: &platformgame.GameInfo{
			Id:     game.ID,
			Code:   game.GameCode,
			Name:   gameName,
			Status: int32(game.Status),
		},
	}, nil
}
