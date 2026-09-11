package gameservicelogic

import (
	"context"
	"time"

	"oa.98ent.com/p9/platform-game/rpc/ent"
	"oa.98ent.com/p9/platform-game/rpc/internal/constant"
	"oa.98ent.com/p9/platform-game/rpc/internal/logic"
	"oa.98ent.com/p9/platform-game/rpc/internal/svc"
	"oa.98ent.com/p9/platform-game/rpc/internal/utils"
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
func (l *UpdateGameLogic) UpdateGame(in *platformgame.UpdateGameRequest) (*platformgame.UpdateGameResp, error) {
	l.Infof("[RPC UpdateGame] received request: id=%d", in.Id)

	if l.svcCtx == nil || l.svcCtx.DB == nil {
		return &platformgame.UpdateGameResp{
			Code:    constant.CodeInternalError,
			Message: "Database not available",
		}, nil
	}

	gameRecord := &ent.Game{}
	if err := l.svcCtx.DB.WithContext(l.ctx).
		Where("id = ? AND deleted_at IS NULL", in.Id).
		First(gameRecord).Error; err != nil {
		l.Errorf("[RPC UpdateGame] failed: %v", err)
		return &platformgame.UpdateGameResp{
			Code:    constant.CodeInternalError,
			Message: err.Error(),
		}, nil
	}

	oldStatus := gameRecord.Status

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
	if in.ImageUrl != "" {
		updateData["image_url"] = in.ImageUrl
	}
	if in.SupportsEmbed != 0 {
		if in.SupportsEmbed == 1 {
			updateData["supports_embed"] = true
		}
		if in.SupportsEmbed == 2 {
			updateData["supports_embed"] = false
		}
	}
	if in.SupportsRedirect != 0 {
		if in.SupportsRedirect == 1 {
			updateData["supports_redirect"] = true
		}
		if in.SupportsRedirect == 2 {
			updateData["supports_redirect"] = false
		}
	}

	if err := l.svcCtx.DB.WithContext(l.ctx).Model(gameRecord).Updates(updateData).Error; err != nil {
		return &platformgame.UpdateGameResp{
			Code:    constant.CodeInternalError,
			Message: err.Error(),
		}, nil
	}

	l.svcCtx.DB.WithContext(l.ctx).Where("id = ?", in.Id).First(gameRecord)

	// 如果需要强制踢线且状态变更为停用，发送 kafka 事件
	if in.ForceLogout && in.Status == 2 && oldStatus != 2 {
		go func() {
			l.Infof("[RPC UpdateGame] sending force-quit event for game id=%d", gameRecord.SourceId.Int64)
			forceQuit := utils.ForceQuitEvent{
				Basis:   utils.KafkaPayloadBasis{Stage: 0, Producer: "platform-game-rpc"},
				Scope:   constant.ScopeGame,
				GameIDs: []int64{gameRecord.SourceId.Int64},
				QuitAt:  time.Now(),
			}
			l.Infof("[RPC UpdateGame] sending force-quit event json: %s", utils.JSON(forceQuit))
			err := utils.SendForceQuitEvent(context.Background(), l.svcCtx.Config.Kafka.Brokers, &forceQuit)
			if err != nil {
				l.Errorf("[RPC UpdateGame] failed to send force-quit event for game id=%d: %v", gameRecord.SourceId.Int64, err)
			}
		}()
	}

	// 查询关联的分类、供应商和渠道信息
	ext, err := GetGameExtInfo(l.ctx, l.svcCtx, gameRecord)
	if err != nil {
		l.Errorf("[RPC UpdateGame] query extended game info failed: %v", err)
		return &platformgame.UpdateGameResp{
			Code:    constant.CodeInternalError,
			Message: "failed to get extended game info: " + err.Error(),
		}, nil
	}

	return &platformgame.UpdateGameResp{
		Code:    constant.CodeSuccess,
		Message: "ok",
		Data:    logic.GameModelToProto(gameRecord, ext),
	}, nil
}
