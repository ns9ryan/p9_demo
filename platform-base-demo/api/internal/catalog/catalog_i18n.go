package catalog

import (
	"oa.98ent.com/p9/core/common/i18n"
	"oa.98ent.com/p9/core/rpc/coreclient"
)

// menuI18n 返回菜单多语言数据
func menuI18n() []*coreclient.I18NItem {
	var out []*coreclient.I18NItem

	add := func(key, zh, hk, en string) {
		out = append(out,
			&coreclient.I18NItem{
				I18NGroup: i18n.GroupMenu,
				TransKey:  key,
				Lang:      i18n.LangZH,
				Value:     zh,
			},
			&coreclient.I18NItem{
				I18NGroup: i18n.GroupMenu,
				TransKey:  key,
				Lang:      i18n.LangHK,
				Value:     hk,
			},
			&coreclient.I18NItem{
				I18NGroup: i18n.GroupMenu,
				TransKey:  key,
				Lang:      i18n.LangEN,
				Value:     en,
			},
		)
	}

	// 基础数据
	add("menu.route.baseData", "基础数据", "基礎數據", "Base data")

	// 语言
	add("menu.route.language", "语言管理", "語言管理", "Languages")
	add("menu.route.languageCreate", "新建语言", "新建語言", "Create language")
	add("menu.route.languageUpdate", "编辑语言", "編輯語言", "Edit language")
	add("menu.route.languageReorder", "调整语言排序", "調整語言排序", "Reorder languages")

	// 时区
	add("menu.route.timezone", "时区管理", "時區管理", "Timezones")
	add("menu.route.timezoneCreate", "新建时区", "新建時區", "Create timezone")
	add("menu.route.timezoneUpdate", "编辑时区", "編輯時區", "Edit timezone")
	add("menu.route.timezoneReorder", "调整时区排序", "調整時區排序", "Reorder timezones")

	// 货币
	add("menu.route.currency", "货币管理", "貨幣管理", "Currencies")
	add("menu.route.currencyCreate", "新建货币", "新建貨幣", "Create currency")
	add("menu.route.currencyUpdate", "编辑货币", "編輯貨幣", "Edit currency")
	add("menu.route.currencyReorder", "调整货币排序", "調整貨幣排序", "Reorder currencies")

	// 国家地区
	add("menu.route.region", "国家地区", "國家地區", "Countries and regions")
	add("menu.route.regionCreate", "新建国家地区", "新建國家地區", "Create country or region")
	add("menu.route.regionUpdate", "编辑国家地区", "編輯國家地區", "Edit country or region")
	add("menu.route.regionReorder", "调整国家地区排序", "調整國家地區排序", "Reorder countries and regions")

	return out
}

// apiI18n 返回API多语言数据
func apiI18n() []*coreclient.I18NItem {
	var out []*coreclient.I18NItem

	add := func(key, zh, hk, en string) {
		out = append(out,
			&coreclient.I18NItem{
				I18NGroup: i18n.GroupAPI,
				TransKey:  key,
				Lang:      i18n.LangZH,
				Value:     zh,
			},
			&coreclient.I18NItem{
				I18NGroup: i18n.GroupAPI,
				TransKey:  key,
				Lang:      i18n.LangHK,
				Value:     hk,
			},
			&coreclient.I18NItem{
				I18NGroup: i18n.GroupAPI,
				TransKey:  key,
				Lang:      i18n.LangEN,
				Value:     en,
			},
		)
	}

	// 语言
	add("api.languageCreate", "创建语言", "創建語言", "Create language")
	add("api.languageUpdate", "更新语言", "更新語言", "Update language")
	add("api.languageGet", "语言详情", "語言詳情", "Language detail")
	add("api.languageList", "语言列表", "語言列表", "Language list")
	add("api.languageListAll", "全部语言", "全部語言", "All languages")
	add("api.languageReorder", "调整语言排序", "調整語言排序", "Reorder languages")

	// 时区
	add("api.timezoneCreate", "创建时区", "創建時區", "Create timezone")
	add("api.timezoneUpdate", "更新时区", "更新時區", "Update timezone")
	add("api.timezoneGet", "时区详情", "時區詳情", "Timezone detail")
	add("api.timezoneList", "时区列表", "時區列表", "Timezone list")
	add("api.timezoneListAll", "全部时区", "全部時區", "All timezones")
	add("api.timezoneReorder", "调整时区排序", "調整時區排序", "Reorder timezones")

	// 货币
	add("api.currencyCreate", "创建货币", "創建貨幣", "Create currency")
	add("api.currencyUpdate", "更新货币", "更新貨幣", "Update currency")
	add("api.currencyGet", "货币详情", "貨幣詳情", "Currency detail")
	add("api.currencyList", "货币列表", "貨幣列表", "Currency list")
	add("api.currencyListAll", "全部货币", "全部貨幣", "All currencies")
	add("api.currencyReorder", "调整货币排序", "調整貨幣排序", "Reorder currencies")

	// 国家地区
	add("api.regionCreate", "创建国家地区", "創建國家地區", "Create country or region")
	add("api.regionUpdate", "更新国家地区", "更新國家地區", "Update country or region")
	add("api.regionGet", "国家地区详情", "國家地區詳情", "Country or region detail")
	add("api.regionList", "国家地区列表", "國家地區列表", "Country and region list")
	add("api.regionListAll", "全部国家地区", "全部國家地區", "All countries and regions")
	add("api.regionReorder", "调整国家地区排序", "調整國家地區排序", "Reorder countries and regions")

	return out
}
