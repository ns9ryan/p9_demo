package gameproviderservicelogic

import (
	"context"
	"time"

	"gorm.io/gorm"
	"oa.98ent.com/p9/platform-game/rpc/ent"
	"oa.98ent.com/p9/platform-game/rpc/internal/constant"
	"oa.98ent.com/p9/platform-game/rpc/internal/logic"
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

	if l.svcCtx == nil || l.svcCtx.DB == nil {
		l.Errorf("[RPC UpdateGameProvider] Database not available")
		return &platformgame.UpdateGameProviderResp{
			Code:    constant.CodeInternalError,
			Message: "Database not available",
			Data:    nil,
		}, nil
	}

	// 查询现有记录
	var provider *ent.GameProvider
	if err := l.svcCtx.DB.WithContext(l.ctx).
		Where("id = ? AND deleted_at IS NULL", in.Id).
		First(&provider).Error; err != nil {
		l.Errorf("[RPC UpdateGameProvider] failed to get provider: %v", err)
		return &platformgame.UpdateGameProviderResp{
			Code:    constant.CodeInternalError,
			Message: "failed to get provider: " + err.Error(),
			Data:    nil,
		}, nil
	}

	oldStatus := provider.Status

	// 构建更新数据
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
	if in.LogoUrl != "" {
		updateData["logo_url"] = in.LogoUrl
	}

	// 执行更新
	if err := l.svcCtx.DB.WithContext(l.ctx).
		Model(&ent.GameProvider{Id: provider.Id}).
		Updates(updateData).Error; err != nil {
		l.Errorf("[RPC UpdateGameProvider] failed to update provider: %v", err)
		return &platformgame.UpdateGameProviderResp{
			Code:    constant.CodeInternalError,
			Message: "failed to update provider: " + err.Error(),
			Data:    nil,
		}, nil
	}

	// 重新查询
	if err := l.svcCtx.DB.WithContext(l.ctx).
		Where("id = ?", in.Id).
		First(&provider).Error; err != nil {
		l.Errorf("[RPC UpdateGameProvider] failed to get updated provider: %v", err)
		return &platformgame.UpdateGameProviderResp{
			Code:    constant.CodeInternalError,
			Message: "failed to get updated provider: " + err.Error(),
			Data:    nil,
		}, nil
	}

	item := logic.ProviderModelToProto(provider)

	l.Infof("[RPC UpdateGameProvider] updated successfully: id=%d", provider.Id)

	// 如果需要强制踢线且状态变更为停用，发送 kafka 事件
	if in.ForceLogout && in.Status == 2 && oldStatus != 2 {
		go func() {
			l.Infof("[RPC UpdateGameProvider] sending force-quit event for provider id=%d", provider.SourceId)
			// 查出所有游戏渠道下的游戏，并发送强制踢线事件
			gameIds, err := getGamesByVendorID(context.Background(), l.svcCtx.DB, provider.SourceId)
			if err != nil {
				l.Errorf("[RPC UpdateGameProvider] failed to get games by provider id=%d: %v", provider.SourceId, err)
				return
			}
			forceQuit := utils.ForceQuitEvent{
				Basis:   utils.KafkaPayloadBasis{Stage: 0, Producer: "platform-game-rpc"},
				Scope:   constant.ScopeGameVendor,
				VenID:   provider.SourceId,
				GameIDs: gameIds,
				QuitAt:  time.Now(),
			}
			l.Infof("[RPC UpdateGameProvider] sending force-quit event json: %s", utils.JSON(forceQuit))
			err = utils.SendForceQuitEvent(context.Background(), l.svcCtx.Config.Kafka.Brokers, &forceQuit)
			if err != nil {
				l.Errorf("[RPC UpdateGameProvider] failed to send force-quit event for provider id=%d: %v", provider.SourceId, err)
			}
		}()
	}

	return &platformgame.UpdateGameProviderResp{
		Code:    constant.CodeSuccess,
		Message: "ok",
		Data:    item,
	}, nil
}

func getGamesByVendorID(ctx context.Context, db *gorm.DB, vendorID int64) ([]int64, error) {
	var games []ent.Game
	err := db.WithContext(ctx).
		Where("vendor_id = ? AND deleted_at IS NULL", vendorID).
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
