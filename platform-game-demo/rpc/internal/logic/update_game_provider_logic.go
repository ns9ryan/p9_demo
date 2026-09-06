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

// 更新游戏供应�?
func (l *UpdateGameProviderLogic) UpdateGameProvider(in *platformgame.UpdateGameProviderRequest) (*platformgame.GetGameProviderListResp, error) {
	l.Infof("[RPC UpdateGameProvider] received request: id=%d", in.Id)

	if l.svcCtx == nil || l.svcCtx.DB == nil {
		l.Errorf("[RPC UpdateGameProvider] Database not available")
		return &platformgame.GetGameProviderListResp{
			Code:    constant.CodeInternalError,
			Message: "Database not available",
			Data:    nil,
		}, nil
	}

	// 查询现有记录
	var provider model.Provider
	if err := l.svcCtx.DB.WithContext(l.ctx).
		Where("id = ? AND deleted_at IS NULL", in.Id).
		First(&provider).Error; err != nil {
		l.Errorf("[RPC UpdateGameProvider] failed to get provider: %v", err)
		return &platformgame.GetGameProviderListResp{
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

	// 执行更新
	if err := l.svcCtx.DB.WithContext(l.ctx).
		Model(&provider).
		Updates(updateData).Error; err != nil {
		l.Errorf("[RPC UpdateGameProvider] failed to update provider: %v", err)
		return &platformgame.GetGameProviderListResp{
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
		return &platformgame.GetGameProviderListResp{
			Code:    constant.CodeInternalError,
			Message: "failed to get updated provider: " + err.Error(),
			Data:    nil,
		}, nil
	}

	// 转换�?proto message
	var providerName string
	if provider.NameI18n != nil {
		if name, ok := provider.NameI18n["default"].(string); ok {
			providerName = name
		} else if name, ok := provider.NameI18n["en"].(string); ok {
			providerName = name
		} else {
			for _, v := range provider.NameI18n {
				if str, ok := v.(string); ok {
					providerName = str
					break
				}
			}
		}
	}

	item := &platformgame.ProviderInfo{
		Id:     provider.ID,
		Code:   provider.ProviderCode,
		Name:   providerName,
		Status: int32(provider.Status),
	}

	l.Infof("[RPC UpdateGameProvider] updated successfully: id=%d", provider.ID)

	// 如果需要强制踢线且状态变更为停用，发送 kafka 事件
	if in.ForceLogout && in.Status == 2 && oldStatus != 2 {
		go func() {
			l.Infof("[RPC UpdateGameProvider] sending force-quit event for provider id=%d", provider.ID)
			if err := kafka.SendForceQuitEvent(l.ctx, provider.ID, "game_provider"); err != nil {
				l.Infof("[RPC UpdateGameProvider] send force-quit event failed: %v", err)
			}
		}()
	}

	return &platformgame.GetGameProviderListResp{
		Code:    constant.CodeSuccess,
		Message: "ok",
		Data:    []*platformgame.ProviderInfo{item},
	}, nil
}
