package logic

import (
	"database/sql"
	"time"

	"oa.98ent.com/p9/platform-game/rpc/ent"
	platformgame "oa.98ent.com/p9/platform-game/rpc/pb/platform_game"
)

func CategoryModelToProto(category *ent.GameCategory) *platformgame.CategoryInfo {
	return &platformgame.CategoryInfo{
		Id:                 category.Id,
		SourceId:           category.SourceId,
		CategoryCode:       category.CategoryCode,
		SourceCategoryCode: category.SourceCategoryCode,
		NameI18N:           category.NameI18n,
		SourceNameI18N:     category.SourceNameI18n,
		SortNo:             int32(category.SortNo),
		Status:             int32(category.Status),
		SourceStatus:       int32(category.SourceStatus),
		IsDeleted:          IsDel(category.DeletedAt.Time),
		CreatedAt:          category.CreatedAt.Unix(),
		UpdatedAt:          category.UpdatedAt.Unix(),
	}
}

func CategoryModelToProtoList(categories []*ent.GameCategory) []*platformgame.CategoryInfo {
	categoriesProto := make([]*platformgame.CategoryInfo, 0, len(categories))
	for _, category := range categories {
		if category == nil {
			continue
		}
		categoriesProto = append(categoriesProto, CategoryModelToProto(category))
	}
	return categoriesProto
}

func ChannelModelToProto(channel *ent.GameChannel) *platformgame.GameChannelResp {
	return &platformgame.GameChannelResp{
		Id:                channel.Id,
		SourceId:          channel.SourceId,
		ChannelCode:       channel.ChannelCode,
		SourceChannelCode: channel.SourceChannelCode,
		NameI18N:          channel.NameI18n,
		SourceNameI18N:    channel.SourceNameI18n,
		SortNo:            int32(channel.SortNo),
		SourceSortNo:      int32(channel.SourceSortNo),
		Status:            int32(channel.Status),
		SourceStatus:      int32(channel.SourceStatus),
		IsDeleted:         IsDel(channel.DeletedAt.Time),
		CreatedAt:         channel.CreatedAt.Unix(),
		UpdatedAt:         channel.UpdatedAt.Unix(),
	}
}

func ChannelModelToProtoList(channels []*ent.GameChannel) []*platformgame.GameChannelResp {
	channelsProto := make([]*platformgame.GameChannelResp, 0, len(channels))
	for _, channel := range channels {
		if channel == nil {
			continue
		}
		channelsProto = append(channelsProto, ChannelModelToProto(channel))
	}
	return channelsProto
}

func CurrencyModelToProto(currency *ent.GameCurrency, gameRecord *ent.Game, sysCurrencyMap map[int64]interface{}) *platformgame.GameCurrencyInfo {
	currencyCode := ""
	currencyNameI18n := ""
	if sysCurrency, ok := sysCurrencyMap[currency.CurrencyId]; ok {
		currencyCode = sysCurrency.(map[string]string)["currency_code"]
		currencyNameI18n = sysCurrency.(map[string]string)["name_i18n"]
	}
	return &platformgame.GameCurrencyInfo{
		Id:               currency.Id,
		GameId:           currency.GameId,
		GameCode:         gameRecord.GameCode,
		GameNameI18N:     gameRecord.NameI18n,
		CurrencyId:       currency.CurrencyId,
		CurrencyCode:     currencyCode,
		CurrencyNameI18N: currencyNameI18n,
		Status:           int32(currency.Status),
		SourceStatus:     int32(currency.SourceStatus),
		IsDeleted:        IsDel(currency.DeletedAt.Time),
		CreatedAt:        currency.CreatedAt.Unix(),
		UpdatedAt:        currency.UpdatedAt.Unix(),
	}
}

func ProviderModelToProto(provider *ent.GameProvider) *platformgame.ProviderInfo {
	return &platformgame.ProviderInfo{
		Id:                 provider.Id,
		SourceId:           provider.SourceId,
		ProviderCode:       provider.ProviderCode,
		SourceProviderCode: provider.SourceProviderCode,
		NameI18N:           provider.NameI18n,
		SourceNameI18N:     provider.SourceNameI18n,
		LogoUrl:            provider.LogoUrl.String,
		SourceLogoUrl:      provider.SourceLogoUrl.String,
		SortNo:             int32(provider.SortNo),
		Status:             int32(provider.Status),
		SourceStatus:       int32(provider.SourceStatus),
		IsDeleted:          IsDel(provider.DeletedAt.Time),
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
	CategoryNameI18N string
	ProviderNameI18N string
	ChannelNameI18N  string
	GameCurrencyInfo string
}

func GameModelToProto(gameRecord *ent.Game, ext *GameInfoExt) *platformgame.GameInfo {
	return &platformgame.GameInfo{
		Id:               gameRecord.Id,
		SourceId:         derefNullInt64(gameRecord.SourceId),
		GameCode:         gameRecord.GameCode,
		SourceGameCode:   gameRecord.SourceGameCode,
		NameI18N:         gameRecord.NameI18n,
		SourceNameI18N:   gameRecord.SourceNameI18n,
		Status:           int32(gameRecord.Status),
		SourceStatus:     int32(gameRecord.SourceStatus),
		CatId:            gameRecord.CategoryId,
		VenId:            gameRecord.ProviderId,
		ChanId:           derefNullInt64(gameRecord.ChannelId),
		CategoryNameI18N: ext.CategoryNameI18N,
		ProviderNameI18N: ext.ProviderNameI18N,
		ChannelNameI18N:  ext.ChannelNameI18N,
		GameCurrencyInfo: ext.GameCurrencyInfo,
		ProviderKey:      gameRecord.ProviderKey.String,
		ImageUrl:         derefNullString(gameRecord.ImageUrl),
		SourceImageUrl:   derefNullString(gameRecord.SourceImageUrl),
		SortNo:           gameRecord.SortNo,
		SupportsEmbed:    gameRecord.SupportsEmbed,
		SupportsRedirect: gameRecord.SupportsRedirect,
		IsDeleted:        IsDel(gameRecord.DeletedAt.Time),
		CreatedAt:        gameRecord.CreatedAt.Unix(),
		UpdatedAt:        gameRecord.UpdatedAt.Unix(),
	}
}

func CheckpointModelToProto(checkpoint *ent.GameSyncCheckpoint) *platformgame.GameSyncCheckpointInfo {
	return &platformgame.GameSyncCheckpointInfo{
		Id:              checkpoint.Id,
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
		LastError:       derefNullString(checkpoint.LastErrorMessage),
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

func derefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func derefNullString(s sql.NullString) string {
	if !s.Valid {
		return ""
	}
	return s.String
}

func derefNullInt64(n sql.NullInt64) int64 {
	if !n.Valid {
		return 0
	}
	return n.Int64
}

// IsDel 检查是否已被软删除（DeletedAt 不为零值）
// 返回 1 表示已删除，0 表示未删除
func IsDel(deletedAt time.Time) int32 {
	if !deletedAt.IsZero() {
		return 1
	}
	return 0
}
