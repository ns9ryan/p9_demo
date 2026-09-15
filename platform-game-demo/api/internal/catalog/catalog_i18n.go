package catalog

import (
	"oa.98ent.com/p9/core/common/i18n"
	"oa.98ent.com/p9/core/rpc/coreclient"
)

const i18nGroupGame = "game" // 游戏管理多语言分组

func langSeeds() []*coreclient.CreateI18NLangReq {
	return []*coreclient.CreateI18NLangReq{
		{Lang: i18n.LangZH, Name: "简体中文", SortNo: 1},
		{Lang: i18n.LangHK, Name: "繁體中文", SortNo: 2},
		{Lang: i18n.LangEN, Name: "English", SortNo: 3},
	}
}

func menuI18n(i18nCode string) []*coreclient.I18NItem {
	var out []*coreclient.I18NItem
	add := func(key, zh, hk, en string) {
		out = append(out,
			&coreclient.I18NItem{I18NCode: i18nCode, I18NGroup: i18n.GroupMenu, TransKey: key, Lang: i18n.LangZH, Value: zh},
			&coreclient.I18NItem{I18NCode: i18nCode, I18NGroup: i18n.GroupMenu, TransKey: key, Lang: i18n.LangHK, Value: hk},
			&coreclient.I18NItem{I18NCode: i18nCode, I18NGroup: i18n.GroupMenu, TransKey: key, Lang: i18n.LangEN, Value: en},
		)
	}
	add("menu.route.game", "游戏管理", "遊戲管理", "Game Management")
	add("menu.route.gameCategory", "游戏分类", "遊戲分類", "Game Category")
	add("menu.route.gameChannel", "游戏渠道", "遊戲渠道", "Game Channel")
	add("menu.route.gameManufacturer", "游戏厂商", "遊戲廠商", "Game Manufacturer")
	add("menu.route.gameIndie", "独立游戏", "獨立遊戲", "Indie Game")
	add("menu.route.gameCurrency", "游戏币种", "遊戲幣種", "Game Currency")
	add("menu.route.gameSyncList", "游戏同步列表", "遊戲同步列表", "Game Sync List")
	return out
}

func apiI18n(i18nCode string) []*coreclient.I18NItem {
	var out []*coreclient.I18NItem
	add := func(key, zh, hk, en string) {
		out = append(out,
			&coreclient.I18NItem{I18NCode: i18nCode, I18NGroup: i18n.GroupAPI, TransKey: key, Lang: i18n.LangZH, Value: zh},
			&coreclient.I18NItem{I18NCode: i18nCode, I18NGroup: i18n.GroupAPI, TransKey: key, Lang: i18n.LangHK, Value: hk},
			&coreclient.I18NItem{I18NCode: i18nCode, I18NGroup: i18n.GroupAPI, TransKey: key, Lang: i18n.LangEN, Value: en},
		)
	}
	add("api.gameCategoryList", "游戏分类列表", "遊戲分類列表", "Game Category list")
	add("api.gameChannelList", "游戏渠道列表", "遊戲渠道列表", "Game Channel list")
	add("api.gameProviderList", "游戏厂商列表", "遊戲廠商列表", "Game Provider list")
	add("api.gameList", "游戏列表", "遊戲列表", "Game list")
	add("api.gameCurrencyList", "游戏币种列表", "遊戲幣種列表", "Game Currency list")
	add("api.gameSyncCheckpointList", "游戏同步检查点列表", "遊戲同步檢查點列表", "Game Sync Checkpoint list")
	add("api.gameCategoryUpdate", "更新游戏分类", "更新遊戲分類", "Update Game Category")
	add("api.gameChannelUpdate", "更新游戏渠道", "更新遊戲渠道", "Update Game Channel")
	add("api.gameProviderUpdate", "更新游戏厂商", "更新遊戲廠商", "Update Game Provider")
	add("api.gameUpdate", "更新游戏", "更新遊戲", "Update Game")
	add("api.gameCurrencyUpdate", "更新游戏币种", "更新遊戲幣種", "Update Game Currency")
	add("api.gameSyncCheckpointUpdate", "更新游戏同步检查点", "更新遊戲同步檢查點", "Update Game Sync Checkpoint")
	add("api.syncCategoriesPreview", "预览同步游戏分类", "預覽同步遊戲分類", "Preview Sync Game Category")
	add("api.syncChannelPreview", "预览同步游戏渠道", "預覽同步遊戲渠道", "Preview Sync Game Channel")
	add("api.syncProviderPreview", "预览同步游戏厂商", "預覽同步遊戲廠商", "Preview Sync Game Provider")
	add("api.syncGamePreview", "预览同步游戏", "預覽同步遊戲", "Preview Sync Game")
	add("api.syncCurrencyPreview", "预览同步游戏币种", "預覽同步遊戲幣種", "Preview Sync Game Currency")
	add("api.syncCategoriesRun", "执行同步游戏分类", "執行同步遊戲分類", "Run Sync Game Category")
	add("api.syncChannelRun", "执行同步游戏渠道", "執行同步遊戲渠道", "Run Sync Game Channel")
	add("api.syncProviderRun", "执行同步游戏厂商", "執行同步遊戲廠商", "Run Sync Game Provider")
	add("api.syncGameRun", "执行同步游戏", "執行同步遊戲", "Run Sync Game")
	add("api.syncCurrencyRun", "执行同步游戏币种", "執行同步遊戲幣種", "Run Sync Game Currency")
	return out
}
