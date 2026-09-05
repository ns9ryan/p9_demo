package currencyservicelogic

import (
	"oa.98ent.com/p9/platform-base/rpc/ent"
	"oa.98ent.com/p9/platform-base/rpc/pb/base/currency"
)

// toCurrencyInfo 转换货币信息
func toCurrencyInfo(data *ent.Currency) *currency.CurrencyInfo {
	return &currency.CurrencyInfo{
		Id:           data.ID,           // 货币ID
		Code:         data.Code,         // 货币编码
		NameI18N:     data.NameI18n,     // 多语言名称
		CurrencyType: data.CurrencyType, // 货币类型: 1法定货币, 2虚拟货币
		Symbol:       data.Symbol,       // 货币符号
		AmountFactor: data.AmountFactor, // 金额换算倍率
		Status:       data.Status,       // 状态: 1启用, 2停用
		SortNo:       data.SortNo,       // 排序值
	}
}
