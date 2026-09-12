package agentline

import "slices"

const (
	// CashProduction 现金正式线路
	CashProduction = "CASH_PRODUCTION"

	// CashDemo 现金试玩线路
	CashDemo = "CASH_DEMO"

	// CashTest 现金测试线路
	CashTest = "CASH_TEST"

	// CreditDemo 信誉试玩线路
	CreditDemo = "CREDIT_DEMO"

	// CreditTest 信誉测试线路
	CreditTest = "CREDIT_TEST"
)

var all = []string{
	CashProduction,
	CashDemo,
	CashTest,
	CreditDemo,
	CreditTest,
}

// IsValid 判断代理子线路编码是否合法
func IsValid(code string) bool {
	return slices.Contains(all, code)
}

// All 获取全部代理子线路编码
func All() []string {
	return slices.Clone(all)
}
