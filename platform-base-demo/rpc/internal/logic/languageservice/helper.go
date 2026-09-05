package languageservicelogic

import (
	"oa.98ent.com/p9/platform-base/rpc/ent"
	"oa.98ent.com/p9/platform-base/rpc/pb/base/language"
)

// toLanguageInfo 转换语言信息
func toLanguageInfo(data *ent.Language) *language.LanguageInfo {
	return &language.LanguageInfo{
		Id:       data.ID,       // 语言ID
		Code:     data.Code,     // 语言编码
		NameI18N: data.NameI18n, // 多语言名称
		Status:   data.Status,   // 状态: 1启用, 2停用
		SortNo:   data.SortNo,   // 排序值
	}
}
