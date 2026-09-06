package logic

import (
	"context"
	"time"

	"oa.98ent.com/p9/platform-game/common/constant"
	"oa.98ent.com/p9/platform-game/common/kafka"
	"oa.98ent.com/p9/platform-game/common/model"
	"oa.98ent.com/p9/platform-game/rpc/internal/svc"
	platformgame "oa.98ent.com/p9/platform-game/rpc/pb/platform_game"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateGameCurrencyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateGameCurrencyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateGameCurrencyLogic {
	return &UpdateGameCurrencyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新游戏货币
func (l *UpdateGameCurrencyLogic) UpdateGameCurrency(in *platformgame.UpdateGameCurrencyRequest) (*platformgame.GetGameCurrencyListResp, error) {
	l.Infof("[RPC UpdateGameCurrency] received request: id=%d", in.Id)

	if l.svcCtx == nil || l.svcCtx.DB == nil {
		return &platformgame.GetGameCurrencyListResp{
			Code:    constant.CodeInternalError,
			Message: "Database not available",
		}, nil
	}

	var currency model.GameCurrency
	if err := l.svcCtx.DB.WithContext(l.ctx).
		Where("id = ? AND deleted_at IS NULL", in.Id).
		First(&currency).Error; err != nil {
		l.Errorf("[RPC UpdateGameCurrency] failed: %v", err)
		return &platformgame.GetGameCurrencyListResp{
			Code:    constant.CodeInternalError,
			Message: err.Error(),
		}, nil
	}

	oldStatus := currency.Status

	updateData := make(map[string]interface{})
	updateData["updated_at"] = time.Now()

	if in.Status > 0 {
		updateData["status"] = in.Status
	}

	if err := l.svcCtx.DB.WithContext(l.ctx).Model(&currency).Updates(updateData).Error; err != nil {
		return &platformgame.GetGameCurrencyListResp{
			Code:    constant.CodeInternalError,
			Message: err.Error(),
		}, nil
	}

	// 如果需要强制踢线且状态变更为停用，发送 kafka 事件
	if in.ForceLogout && in.Status == 2 && oldStatus != 2 {
		go func() {
			l.Infof("[RPC UpdateGameCurrency] sending force-quit event for game id=%d", currency.GameID)
			if err := kafka.SendForceQuitEvent(l.ctx, currency.GameID, "game"); err != nil {
				l.Infof("[RPC UpdateGameCurrency] send force-quit event failed: %v", err)
			}
		}()
	}

	return &platformgame.GetGameCurrencyListResp{
		Code:    constant.CodeSuccess,
		Message: "ok",
	}, nil
}
