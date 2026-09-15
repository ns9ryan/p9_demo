package catalog

import (
	"oa.98ent.com/p9/core/common/i18n"
	"oa.98ent.com/p9/core/rpc/coreclient"
)

// platformOperatorI18n 返回Platform Operator全部多语言数据
func platformOperatorI18n() []*coreclient.I18NItem {
	var out []*coreclient.I18NItem

	out = append(out, menuI18n()...)
	out = append(out, apiI18n()...)

	return out
}

// menuI18n 返回菜单多语言数据
func menuI18n() []*coreclient.I18NItem {
	var out []*coreclient.I18NItem

	// 分站管理
	addI18n(&out, i18n.GroupMenu, "menu.route.operatorManagement", "分站管理", "分站管理", "Operator management")

	// 分站列表
	addI18n(&out, i18n.GroupMenu, "menu.route.operatorList", "分站列表", "分站列表", "Operators")
	addI18n(&out, i18n.GroupMenu, "menu.route.operatorCreate", "创建分站", "建立分站", "Create operator")
	addI18n(&out, i18n.GroupMenu, "menu.route.operatorUpdate", "编辑分站", "編輯分站", "Edit operator")
	addI18n(&out, i18n.GroupMenu, "menu.route.operatorPublish", "发布分站", "發布分站", "Publish operator")
	addI18n(&out, i18n.GroupMenu, "menu.route.operatorDelete", "删除分站", "刪除分站", "Delete operator")

	// 域名管理
	addI18n(&out, i18n.GroupMenu, "menu.route.operatorDomain", "域名管理", "網域管理", "Domain management")
	addI18n(&out, i18n.GroupMenu, "menu.route.operatorDomainCreate", "创建域名", "建立網域", "Create domain")
	addI18n(&out, i18n.GroupMenu, "menu.route.operatorDomainUpdate", "编辑域名", "編輯網域", "Edit domain")
	addI18n(&out, i18n.GroupMenu, "menu.route.operatorDomainDelete", "删除域名", "刪除網域", "Delete domain")

	// 管理员账号
	addI18n(&out, i18n.GroupMenu, "menu.route.operatorAdmin", "管理员账号", "管理員帳號", "Administrator accounts")

	// 基础资源分配
	addI18n(&out, i18n.GroupMenu, "menu.route.basicResourceAllocation", "基础资源分配", "基礎資源分配", "Basic resource allocation")
	addI18n(&out, i18n.GroupMenu, "menu.route.basicResourceAllocationSave", "保存基础资源配置", "儲存基礎資源設定", "Save basic resource configuration")

	// 游戏资源分配
	addI18n(&out, i18n.GroupMenu, "menu.route.gameResourceAllocation", "游戏资源分配", "遊戲資源分配", "Game resource allocation")

	return out
}

// apiI18n 返回API多语言数据
func apiI18n() []*coreclient.I18NItem {
	var out []*coreclient.I18NItem

	// 分站
	addI18n(&out, i18n.GroupAPI, "api.operatorCreate", "创建分站", "建立分站", "Create operator")
	addI18n(&out, i18n.GroupAPI, "api.operatorUpdate", "修改分站", "修改分站", "Update operator")
	addI18n(&out, i18n.GroupAPI, "api.operatorGet", "分站详情", "分站詳情", "Operator detail")
	addI18n(&out, i18n.GroupAPI, "api.operatorList", "分站列表", "分站列表", "Operator list")
	addI18n(&out, i18n.GroupAPI, "api.operatorComplete", "完成分站创建", "完成分站建立", "Complete operator creation")
	addI18n(&out, i18n.GroupAPI, "api.operatorPublish", "发布分站", "發布分站", "Publish operator")
	addI18n(&out, i18n.GroupAPI, "api.operatorDelete", "删除分站", "刪除分站", "Delete operator")

	// 分站档案
	addI18n(&out, i18n.GroupAPI, "api.operatorProfileCreate", "创建分站档案", "建立分站檔案", "Create operator profile")
	addI18n(&out, i18n.GroupAPI, "api.operatorProfileUpdate", "修改分站档案", "修改分站檔案", "Update operator profile")
	addI18n(&out, i18n.GroupAPI, "api.operatorProfileGet", "分站档案详情", "分站檔案詳情", "Operator profile detail")

	// 分站域名
	addI18n(&out, i18n.GroupAPI, "api.operatorDomainCreate", "创建分站域名", "建立分站網域", "Create operator domain")
	addI18n(&out, i18n.GroupAPI, "api.operatorDomainUpdate", "修改分站域名", "修改分站網域", "Update operator domain")
	addI18n(&out, i18n.GroupAPI, "api.operatorDomainGet", "分站域名详情", "分站網域詳情", "Operator domain detail")
	addI18n(&out, i18n.GroupAPI, "api.operatorDomainList", "分站域名列表", "分站網域列表", "Operator domain list")
	addI18n(&out, i18n.GroupAPI, "api.operatorDomainDelete", "删除分站域名", "刪除分站網域", "Delete operator domain")

	// 基础资源分配
	addI18n(&out, i18n.GroupAPI, "api.basicResourceAllocationList", "基础资源分配列表", "基礎資源分配列表", "Basic resource allocation list")

	// 语言分配
	addI18n(&out, i18n.GroupAPI, "api.languageAllocationList", "语言分配列表", "語言分配列表", "Language allocation list")
	addI18n(&out, i18n.GroupAPI, "api.languageAllocationSave", "保存语言分配", "儲存語言分配", "Save language allocation")

	// 经营地区分配
	addI18n(&out, i18n.GroupAPI, "api.regionAllocationList", "经营地区分配列表", "經營地區分配列表", "Region allocation list")
	addI18n(&out, i18n.GroupAPI, "api.regionAllocationSave", "保存经营地区分配", "儲存經營地區分配", "Save region allocation")

	// 代理子线路分配
	addI18n(&out, i18n.GroupAPI, "api.agentLineAllocationList", "代理子线路分配列表", "代理子線路分配列表", "Agent line allocation list")
	addI18n(&out, i18n.GroupAPI, "api.agentLineAllocationSave", "保存代理子线路分配", "儲存代理子線路分配", "Save agent line allocation")

	return out
}

// addI18n 添加Platform站点的简体中文、繁体中文和英文翻译
func addI18n(out *[]*coreclient.I18NItem, group, key, zh, hk, en string) {
	*out = append(*out,
		&coreclient.I18NItem{I18NCode: i18n.CodePlatform, I18NGroup: group, TransKey: key, Lang: i18n.LangZH, Value: zh},
		&coreclient.I18NItem{I18NCode: i18n.CodePlatform, I18NGroup: group, TransKey: key, Lang: i18n.LangHK, Value: hk},
		&coreclient.I18NItem{I18NCode: i18n.CodePlatform, I18NGroup: group, TransKey: key, Lang: i18n.LangEN, Value: en},
	)
}
