package logic

import (
	"encoding/json"
	"fmt"

	"oa.98ent.com/p9/platform-game/api/internal/types"
	platformgame "oa.98ent.com/p9/platform-game/rpc/pb/platform_game"
)

func CategoryProtoToResponse(category *platformgame.CategoryInfo) *types.GameCategoryResp {
	return &types.GameCategoryResp{
		ID:                 category.Id,
		SourceID:           category.SourceId,
		CategoryCode:       category.CategoryCode,
		SourceCategoryCode: category.SourceCategoryCode,
		NameI18n:           category.NameI18N,
		SourceNameI18n:     category.SourceNameI18N,
		SortNo:             int32(category.SortNo),
		Status:             int16(category.Status),
		SourceStatus:       int16(category.SourceStatus),
		UpdatedAt:          category.UpdatedAt,
		CreatedAt:          category.CreatedAt,
		IsDeleted:          int16(category.IsDeleted),
	}
}

func ChannelProtoToResponse(channel *platformgame.GameChannelResp) *types.GameChannelResp {
	return &types.GameChannelResp{
		ID:                channel.Id,
		SourceID:          channel.SourceId,
		ChannelCode:       channel.ChannelCode,
		SourceChannelCode: channel.SourceChannelCode,
		NameI18n:          channel.NameI18N,
		SourceNameI18n:    channel.SourceNameI18N,
		SortNo:            int32(channel.SortNo),
		SourceSortNo:      int32(channel.SourceSortNo),
		Status:            int16(channel.Status),
		SourceStatus:      int16(channel.SourceStatus),
		UpdatedAt:         channel.UpdatedAt,
		CreatedAt:         channel.CreatedAt,
		IsDeleted:         int16(channel.IsDeleted),
	}
}

func CheckpointProtoToResponse(checkpoint *platformgame.GameSyncCheckpointInfo) *types.GameSyncCheckpointResp {
	return &types.GameSyncCheckpointResp{
		ID:               checkpoint.Id,
		SyncScope:        checkpoint.SyncScope,
		CheckpointValue:  checkpoint.CheckpointValue,
		RemoteTotal:      checkpoint.RemoteTotal,
		LocalTotal:       checkpoint.LocalTotal,
		CreatedCount:     checkpoint.CreatedCount,
		UpdatedCount:     checkpoint.UpdatedCount,
		DeletedCount:     checkpoint.DeletedCount,
		FailedCount:      checkpoint.FailedCount,
		Progress:         int16(checkpoint.Progress),
		SyncStatus:       int16(checkpoint.SyncStatus),
		LastSyncAt:       checkpoint.LastSyncAt,
		LastSuccessAt:    checkpoint.LastSuccessAt,
		LastErrorMessage: checkpoint.LastError,
		CreatedAt:        checkpoint.CreatedAt,
		UpdatedAt:        checkpoint.UpdatedAt,
	}
}

func CurrencyProtoToResponse(currency *platformgame.GameCurrencyInfo) *types.GameCurrencyResp {
	return &types.GameCurrencyResp{
		ID:               currency.Id,
		GameID:           currency.GameId,
		GameCode:         currency.GameCode,
		GameNameI18n:     currency.GameNameI18N,
		CurrencyID:       currency.CurrencyId,
		CurrencyCode:     currency.CurrencyCode,
		CurrencyNameI18n: currency.CurrencyNameI18N,
		Status:           int16(currency.Status),
		SourceStatus:     int16(currency.SourceStatus),
		IsDeleted:        int16(currency.IsDeleted),
		UpdatedAt:        currency.UpdatedAt,
		CreatedAt:        currency.CreatedAt,
	}
}

func GameProtoToResponse(game *platformgame.GameInfo) *types.GameResp {
	gameCurrencyInfo := []types.GameCurrencyInfo{}
	json.Unmarshal([]byte(game.GameCurrencyInfo), &gameCurrencyInfo)
	fmt.Println("gameId:", game.Id, "game.GameCurrencyInfo", game.GameCurrencyInfo, "gameCurrencyInfo:", gameCurrencyInfo)

	return &types.GameResp{
		ID:               game.Id,
		SourceID:         game.SourceId,
		GameCode:         game.GameCode,
		SourceGameCode:   game.SourceGameCode,
		NameI18n:         game.NameI18N,
		SourceNameI18n:   game.SourceNameI18N,
		SortNo:           int32(game.SortNo),
		Status:           int16(game.Status),
		SourceStatus:     int16(game.SourceStatus),
		ImageUrl:         game.ImageUrl,
		SourceImageUrl:   game.SourceImageUrl,
		CategoryID:       game.CatId,
		ProviderID:       game.VenId,
		ChannelID:        game.ChanId,
		CategoryNameI18n: game.CategoryNameI18N,
		ProviderNameI18n: game.ProviderNameI18N,
		ChannelNameI18n:  game.ChannelNameI18N,
		GameCurrencyInfo: gameCurrencyInfo,
		ProviderKey:      game.ProviderKey,
		SupportsEmbed:    game.SupportsEmbed,
		SupportsRedirect: game.SupportsRedirect,
		UpdatedAt:        game.UpdatedAt,
		CreatedAt:        game.CreatedAt,
		IsDeleted:        int16(game.IsDeleted),
	}
}

func ProviderProtoToResponse(provider *platformgame.ProviderInfo) *types.GameProviderResp {
	return &types.GameProviderResp{
		ID:                 provider.Id,
		SourceID:           provider.SourceId,
		ProviderCode:       provider.ProviderCode,
		SourceProviderCode: provider.SourceProviderCode,
		NameI18n:           provider.NameI18N,
		SourceNameI18n:     provider.SourceNameI18N,
		SortNo:             int32(provider.SortNo),
		Status:             int16(provider.Status),
		SourceStatus:       int16(provider.SourceStatus),
		LogoUrl:            provider.LogoUrl,
		SourceLogoUrl:      provider.SourceLogoUrl,
		UpdatedAt:          provider.UpdatedAt,
		CreatedAt:          provider.CreatedAt,
		IsDeleted:          int16(provider.IsDeleted),
	}
}
