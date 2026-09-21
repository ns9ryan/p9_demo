package logic

import (
	"encoding/json"

	"context"

	corei18n "oa.98ent.com/p9/core/common/i18n"
	"oa.98ent.com/p9/platform-game/api/internal/types"
	platformgame "oa.98ent.com/p9/platform-game/rpc/pb/platform_game"
)

func CategoryProtoToResponse(ctx context.Context, category *platformgame.GameCategoryInfo) *types.GameCategoryResp {
	// 调用 TG 进行翻译
	name := corei18n.TG(ctx, corei18n.CodePlatform, "game", category.NameKey)
	return &types.GameCategoryResp{
		ID:                 category.Id,
		SourceID:           category.SourceId,
		CategoryCode:       category.CategoryCode,
		SourceCategoryCode: category.SourceCategoryCode,
		Name:               name,
		SortNo:             int32(category.SortNo),
		Status:             int16(category.Status),
		SourceStatus:       int16(category.SourceStatus),
		UpdatedAt:          category.UpdatedAt,
		CreatedAt:          category.CreatedAt,
		IsDeleted:          int16(category.IsDeleted),
	}
}

func ChannelProtoToResponse(ctx context.Context, channel *platformgame.GameChannelInfo) *types.GameChannelResp {
	name := corei18n.TG(ctx, corei18n.CodePlatform, "game", channel.NameKey)
	return &types.GameChannelResp{
		ID:                channel.Id,
		SourceID:          channel.SourceId,
		ChannelCode:       channel.ChannelCode,
		SourceChannelCode: channel.SourceChannelCode,
		Name:              name,
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

func CurrencyProtoToResponse(ctx context.Context, currency *platformgame.GameCurrencyInfo) *types.GameCurrencyResp {
	currencyName := corei18n.TG(ctx, corei18n.CodePlatform, "base", currency.CurrencyNameKey)
	return &types.GameCurrencyResp{
		ID:           currency.Id,
		GameID:       currency.GameId,
		GameCode:     currency.GameCode,
		GameName:     currency.GameName,
		CurrencyID:   currency.CurrencyId,
		CurrencyCode: currency.CurrencyCode,
		CurrencyName: currencyName,
		Status:       int16(currency.Status),
		SourceStatus: int16(currency.SourceStatus),
		IsDeleted:    int16(currency.IsDeleted),
		UpdatedAt:    currency.UpdatedAt,
		CreatedAt:    currency.CreatedAt,
	}
}

func GameProtoToResponse(ctx context.Context, game *platformgame.GameInfo) *types.GameResp {
	gameCurrencyArray := []map[string]interface{}{}
	gameCurrencyInfo := []types.GameCurrencyInfo{}
	json.Unmarshal([]byte(game.GameCurrencyInfo), &gameCurrencyArray)
	categoryName := corei18n.TG(ctx, corei18n.CodePlatform, "game", game.CategoryNameKey)
	providerName := corei18n.TG(ctx, corei18n.CodePlatform, "game", game.ProviderNameKey)
	channelName := corei18n.TG(ctx, corei18n.CodePlatform, "game", game.ChannelNameKey)
	for _, currencyMap := range gameCurrencyArray {
		currency := types.GameCurrencyInfo{
			CurrencyID:   int64(currencyMap["currency_id"].(float64)),
			CurrencyName: "",
		}
		currencyName := corei18n.TG(ctx, corei18n.CodePlatform, "base", currencyMap["currency_name_key"].(string))
		currency.CurrencyName = currencyName
		gameCurrencyInfo = append(gameCurrencyInfo, currency)
	}
	return &types.GameResp{
		ID:               game.Id,
		SourceID:         game.SourceId,
		GameCode:         game.GameCode,
		SourceGameCode:   game.SourceGameCode,
		Name:             game.Name,
		SortNo:           int32(game.SortNo),
		Status:           int16(game.Status),
		SourceStatus:     int16(game.SourceStatus),
		ImageUrl:         game.ImageUrl,
		SourceImageUrl:   game.SourceImageUrl,
		CategoryID:       game.CatId,
		ProviderID:       game.VenId,
		ChannelID:        game.ChanId,
		CategoryName:     categoryName,
		ProviderName:     providerName,
		ChannelName:      channelName,
		GameCurrencyInfo: gameCurrencyInfo,
		ProviderKey:      game.ProviderKey,
		SupportsEmbed:    game.SupportsEmbed,
		SupportsRedirect: game.SupportsRedirect,
		UpdatedAt:        game.UpdatedAt,
		CreatedAt:        game.CreatedAt,
		IsDeleted:        int16(game.IsDeleted),
	}
}

func ProviderProtoToResponse(ctx context.Context, provider *platformgame.ProviderInfo) *types.GameProviderResp {
	name := corei18n.TG(ctx, corei18n.CodePlatform, "game", provider.NameKey)
	return &types.GameProviderResp{
		ID:                 provider.Id,
		SourceID:           provider.SourceId,
		ProviderCode:       provider.ProviderCode,
		SourceProviderCode: provider.SourceProviderCode,
		Name:               name,
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
