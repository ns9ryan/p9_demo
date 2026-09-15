package logic

import (
	"fmt"
	"time"

	"oa.98ent.com/p9/platform-game/rpc/ent"
	"oa.98ent.com/p9/platform-game/rpc/internal/constant"
	platformgame "oa.98ent.com/p9/platform-game/rpc/pb/platform_game"
)

func GameCategoryModelToProto(category *ent.GameCategory) *platformgame.GameCategoryInfo {
	return &platformgame.GameCategoryInfo{
		Id:                 category.ID,
		SourceId:           category.SourceID,
		NameKey:            fmt.Sprintf("%s.%s.name", constant.CategoryBiz, category.SourceCategoryCode),
		CategoryCode:       category.CategoryCode,
		SourceCategoryCode: category.SourceCategoryCode,
		SortNo:             int32(category.SortNo),
		Status:             int32(category.Status),
		SourceStatus:       int32(category.SourceStatus),
		IsDeleted:          IsDel(category.DeletedAt),
		CreatedAt:          category.CreatedAt.Unix(),
		UpdatedAt:          category.UpdatedAt.Unix(),
	}
}

func GameCategoryModelToProtoList(categories []*ent.GameCategory) []*platformgame.GameCategoryInfo {
	categoriesProto := make([]*platformgame.GameCategoryInfo, 0, len(categories))
	for _, category := range categories {
		if category == nil {
			continue
		}
		categoriesProto = append(categoriesProto, GameCategoryModelToProto(category))
	}
	return categoriesProto
}

func ChannelModelToProto(channel *ent.GameChannel) *platformgame.GameChannelInfo {
	return &platformgame.GameChannelInfo{
		Id:                channel.ID,
		SourceId:          channel.SourceID,
		NameKey:           fmt.Sprintf("%s.%s.name", constant.ChannelBiz, channel.SourceChannelCode),
		ChannelCode:       channel.ChannelCode,
		SourceChannelCode: channel.SourceChannelCode,
		SortNo:            int32(channel.SortNo),
		SourceSortNo:      int32(channel.SourceSortNo),
		Status:            int32(channel.Status),
		SourceStatus:      int32(channel.SourceStatus),
		IsDeleted:         IsDel(channel.DeletedAt),
		CreatedAt:         channel.CreatedAt.Unix(),
		UpdatedAt:         channel.UpdatedAt.Unix(),
	}
}

func ChannelModelToProtoList(channels []*ent.GameChannel) []*platformgame.GameChannelInfo {
	channelsProto := make([]*platformgame.GameChannelInfo, 0, len(channels))
	for _, channel := range channels {
		if channel == nil {
			continue
		}
		channelsProto = append(channelsProto, ChannelModelToProto(channel))
	}
	return channelsProto
}

func CurrencyModelToProto(currency *ent.GameCurrency, gameRecord *ent.Game, sysCurrencyMap map[int64]*ent.Currency) *platformgame.GameCurrencyInfo {
	currencyCode := ""
	currencyNameKey := ""
	if sysCurrency, ok := sysCurrencyMap[currency.CurrencyID]; ok {
		currencyCode = sysCurrency.Code
		currencyNameKey = sysCurrency.NameKey
	}
	return &platformgame.GameCurrencyInfo{
		Id:              currency.ID,
		GameId:          currency.GameID,
		GameCode:        gameRecord.GameCode,
		GameName:        gameRecord.Name,
		CurrencyId:      currency.CurrencyID,
		CurrencyCode:    currencyCode,
		CurrencyNameKey: currencyNameKey,
		Status:          int32(currency.Status),
		SourceStatus:    int32(currency.SourceStatus),
		IsDeleted:       IsDel(currency.DeletedAt),
		CreatedAt:       currency.CreatedAt.Unix(),
		UpdatedAt:       currency.UpdatedAt.Unix(),
	}
}

func ProviderModelToProto(provider *ent.GameProvider) *platformgame.ProviderInfo {
	return &platformgame.ProviderInfo{
		Id:                 provider.ID,
		SourceId:           provider.SourceID,
		NameKey:            fmt.Sprintf("%s.%s.name", constant.ProviderBiz, provider.SourceProviderCode),
		ProviderCode:       provider.ProviderCode,
		SourceProviderCode: provider.SourceProviderCode,
		LogoUrl:            provider.LogoURL,
		SourceLogoUrl:      provider.SourceLogoURL,
		SortNo:             int32(provider.SortNo),
		Status:             int32(provider.Status),
		SourceStatus:       int32(provider.SourceStatus),
		IsDeleted:          IsDel(provider.DeletedAt),
		CreatedAt:          provider.CreatedAt.Unix(),
		UpdatedAt:          provider.UpdatedAt.Unix(),
	}
}

func ProviderModelToProtoList(providers []*ent.GameProvider) []*platformgame.ProviderInfo {
	providersProto := make([]*platformgame.ProviderInfo, 0, len(providers))
	for _, provider := range providers {
		if provider == nil {
			continue
		}
		providersProto = append(providersProto, ProviderModelToProto(provider))
	}
	return providersProto
}

type GameInfoExt struct {
	CategoryCode     string
	ProviderCode     string
	ChannelCode      string
	GameCurrencyInfo string
}

func GameModelToProto(gameRecord *ent.Game, ext *GameInfoExt) *platformgame.GameInfo {
	return &platformgame.GameInfo{
		Id:               gameRecord.ID,
		SourceId:         gameRecord.SourceID,
		GameCode:         gameRecord.GameCode,
		SourceGameCode:   gameRecord.SourceGameCode,
		Name:             gameRecord.Name,
		Status:           int32(gameRecord.Status),
		SourceStatus:     int32(gameRecord.SourceStatus),
		CatId:            gameRecord.CategoryID,
		VenId:            gameRecord.ProviderID,
		ChanId:           gameRecord.ChannelID,
		CategoryNameKey:  fmt.Sprintf("%s.%s.name", constant.CategoryBiz, ext.CategoryCode),
		ProviderNameKey:  fmt.Sprintf("%s.%s.name", constant.ProviderBiz, ext.ProviderCode),
		ChannelNameKey:   fmt.Sprintf("%s.%s.name", constant.ChannelBiz, ext.ChannelCode),
		GameCurrencyInfo: ext.GameCurrencyInfo,
		ProviderKey:      gameRecord.ProviderKey,
		ImageUrl:         gameRecord.ImageURL,
		SourceImageUrl:   gameRecord.SourceImageURL,
		SortNo:           gameRecord.SortNo,
		SupportsEmbed:    gameRecord.SupportsEmbed,
		SupportsRedirect: gameRecord.SupportsRedirect,
		IsDeleted:        IsDel(gameRecord.DeletedAt),
		CreatedAt:        gameRecord.CreatedAt.Unix(),
		UpdatedAt:        gameRecord.UpdatedAt.Unix(),
	}
}

func CheckpointModelToProto(checkpoint *ent.GameSyncCheckpoint) *platformgame.GameSyncCheckpointInfo {
	return &platformgame.GameSyncCheckpointInfo{
		Id:              checkpoint.ID,
		SyncScope:       checkpoint.SyncScope,
		CheckpointValue: checkpoint.CheckpointValue,
		RemoteTotal:     checkpoint.RemoteTotal,
		LocalTotal:      checkpoint.LocalTotal,
		CreatedCount:    checkpoint.CreatedCount,
		UpdatedCount:    checkpoint.UpdatedCount,
		DeletedCount:    checkpoint.DeletedCount,
		FailedCount:     checkpoint.FailedCount,
		Progress:        int32(checkpoint.Progress),
		SyncStatus:      int32(checkpoint.SyncStatus),
		LastSyncAt:      checkpoint.LastSyncAt.Unix(),
		LastSuccessAt:   checkpoint.LastSuccessAt.Unix(),
		LastError:       checkpoint.LastErrorMessage,
		CreatedAt:       checkpoint.CreatedAt.Unix(),
		UpdatedAt:       checkpoint.UpdatedAt.Unix(),
	}
}

func CheckpointModelToProtoList(checkpoints []*ent.GameSyncCheckpoint) []*platformgame.GameSyncCheckpointInfo {
	checkpointsProto := make([]*platformgame.GameSyncCheckpointInfo, 0, len(checkpoints))
	for _, checkpoint := range checkpoints {
		if checkpoint == nil {
			continue
		}
		checkpointsProto = append(checkpointsProto, CheckpointModelToProto(checkpoint))
	}
	return checkpointsProto
}

// IsDel 检查是否已被软删除（DeletedAt 不为零值）
// 返回 1 表示已删除，0 表示未删除
func IsDel(deletedAt time.Time) int32 {
	if !deletedAt.IsZero() {
		return 1
	}
	return 0
}
