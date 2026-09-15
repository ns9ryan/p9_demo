package gameproviderservicelogic

import (
	"context"
	"time"

	"oa.98ent.com/p9/platform-game/rpc/internal/constant"
	"oa.98ent.com/p9/platform-game/rpc/internal/svc"
	"oa.98ent.com/p9/platform-game/rpc/internal/utils"
	platformgame "oa.98ent.com/p9/platform-game/rpc/pb/platform_game"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateGameProviderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateGameProviderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateGameProviderLogic {
	return &UpdateGameProviderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新游戏供应商
func (l *UpdateGameProviderLogic) UpdateGameProvider(in *platformgame.UpdateGameProviderRequest) (*platformgame.UpdateGameProviderResp, error) {
	l.Infof("[RPC UpdateGameProvider] received request: id=%d", in.Id)

	if l.svcCtx == nil || l.svcCtx.DAOManager == nil {
		l.Errorf("[RPC UpdateGameProvider] DAO Manager not available")
		return &platformgame.UpdateGameProviderResp{
			Code:    constant.CodeInternalError,
			Message: "DAO Manager not available",
			Data:    nil,
		}, nil
	}

	// 查询当前供应商数据
	providerRecord, err := l.svcCtx.DAOManager.GameProvider.GetGameProviderByID(l.ctx, in.Id)
	if err != nil {
		l.Errorf("[RPC UpdateGameProvider] failed to get provider: %v", err)
		return &platformgame.UpdateGameProviderResp{
			Code:    constant.CodeInternalError,
			Message: err.Error(),
		}, nil
	}

	oldStatus := providerRecord.Status

	// 准备更新数据
	updates := make(map[string]interface{})
	if in.SortNo > 0 {
		updates["sort_no"] = int64(in.SortNo)
	}
	if in.Status > 0 {
		updates["status"] = int64(in.Status)
	}
	if in.LogoUrl != "" {
		updates["logo_url"] = in.LogoUrl
	}
	updates["updated_at"] = time.Now()

	// 执行更新
	_, err = l.svcCtx.DAOManager.GameProvider.UpdateGameProvider(l.ctx, in.Id, updates)
	if err != nil {
		return &platformgame.UpdateGameProviderResp{
			Code:    constant.CodeInternalError,
			Message: err.Error(),
		}, nil
	}

	// 如果需要强制踢线且状态变更为停用，发送 kafka 事件
	if in.ForceLogout && in.Status == 2 && oldStatus != 2 {
		go func() {
			l.Infof("[RPC UpdateGameProvider] sending force-quit event for provider id=%d", providerRecord.SourceID)
			// 查出所有游戏渠道下的游戏，并发送强制踢线事件
			gameIds, err := getGamesByVendorID(context.Background(), l.svcCtx, providerRecord.SourceID)
			if err != nil {
				l.Errorf("[RPC UpdateGameProvider] failed to get games by provider id=%d: %v", providerRecord.SourceID, err)
				return
			}
			forceQuit := utils.ForceQuitEvent{
				Basis:   utils.KafkaPayloadBasis{Stage: 0, Producer: "platform-game-rpc"},
				Scope:   constant.ScopeGameVendor,
				VenID:   providerRecord.SourceID,
				GameIDs: gameIds,
				QuitAt:  time.Now(),
			}
			l.Infof("[RPC UpdateGameProvider] sending force-quit event json: %s", utils.JSON(forceQuit))
			err = utils.SendForceQuitEvent(context.Background(), l.svcCtx.Config.Kafka.Brokers, &forceQuit)
			if err != nil {
				l.Errorf("[RPC UpdateGameProvider] failed to send force-quit event for provider id=%d: %v", providerRecord.SourceID, err)
			}
		}()
	}

	return &platformgame.UpdateGameProviderResp{
		Code:    constant.CodeSuccess,
		Message: "ok",
	}, nil
}

func getGamesByVendorID(ctx context.Context, svcCtx *svc.ServiceContext, vendorID int64) ([]int64, error) {
	if svcCtx == nil || svcCtx.DAOManager == nil {
		return nil, nil
	}

	games, err := svcCtx.DAOManager.Game.GetAllGame(
		ctx,
		vendorID, // providerId
		0,        // categoryId
		0,        // channelId
	)
	if err != nil {
		return nil, err
	}

	ids := make([]int64, 0, len(games))
	for _, game := range games {
		if game != nil {
			ids = append(ids, game.SourceID)
		}
	}
	return ids, nil
}
