package svc

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"oa.98ent.com/p9/platform-base/rpc/ent"
	"oa.98ent.com/p9/platform-base/rpc/ent/currency"
	"oa.98ent.com/p9/platform-base/rpc/ent/migrate"
	"oa.98ent.com/p9/platform-base/rpc/ent/region"
	"oa.98ent.com/p9/platform-base/rpc/ent/timezone"
)

// MustMigrate 执行数据库自动迁移
func (s *ServiceContext) MustMigrate() {
	ctx := context.Background()

	// 根据 Ent Schema 自动创建或更新数据库结构
	logx.Must(
		s.DB.Schema.Create(
			ctx,
			migrate.WithForeignKeys(false), // 不创建数据库外键
			migrate.WithDropIndex(true),    // 允许删除废弃索引
		),
	)

	// 初始化系统默认数据
	s.mustInitTimezone(ctx)
	s.mustInitCurrency(ctx)
	s.mustInitRegion(ctx)
}

// mustInitTimezone 初始化系统时区
func (s *ServiceContext) mustInitTimezone(ctx context.Context) {
	data := []struct {
		code    string // IANA时区编码
		nameKey string // 名称翻译 Key
		status  int64  // 状态: 1启用, 2停用
		sortNo  int64  // 排序值, 数值越小越靠前
	}{
		{code: "Etc/UTC", nameKey: "timezone.etc_utc.name", status: 1, sortNo: 10},
		{code: "Asia/Tokyo", nameKey: "timezone.asia_tokyo.name", status: 1, sortNo: 20},
		{code: "Asia/Ho_Chi_Minh", nameKey: "timezone.asia_ho_chi_minh.name", status: 1, sortNo: 30},
		{code: "Asia/Manila", nameKey: "timezone.asia_manila.name", status: 1, sortNo: 40},
		{code: "Asia/Bangkok", nameKey: "timezone.asia_bangkok.name", status: 1, sortNo: 50},
		{code: "Asia/Kuala_Lumpur", nameKey: "timezone.asia_kuala_lumpur.name", status: 1, sortNo: 60},
		{code: "Asia/Singapore", nameKey: "timezone.asia_singapore.name", status: 1, sortNo: 70},
		{code: "Asia/Jakarta", nameKey: "timezone.asia_jakarta.name", status: 1, sortNo: 80},
		{code: "Asia/Seoul", nameKey: "timezone.asia_seoul.name", status: 1, sortNo: 90},
		{code: "America/Sao_Paulo", nameKey: "timezone.america_sao_paulo.name", status: 1, sortNo: 100},
	}

	builders := make([]*ent.TimezoneCreate, 0, len(data))

	for _, item := range data {
		builders = append(builders,
			s.DB.Timezone.
				Create().
				SetCode(item.code).
				SetNameKey(item.nameKey).
				SetStatus(item.status).
				SetSortNo(item.sortNo),
		)
	}

	// 只补充缺失的系统时区, 已存在的数据保持不变
	logx.Must(
		s.DB.Timezone.
			CreateBulk(builders...).
			OnConflictColumns(timezone.FieldCode).
			Ignore().
			Exec(ctx),
	)
}

// mustInitCurrency 初始化系统货币
func (s *ServiceContext) mustInitCurrency(ctx context.Context) {
	data := []struct {
		code         string // 货币编码
		nameKey      string // 名称翻译 Key
		currencyType int64  // 货币类型: 1法币, 2数字货币
		symbol       string // 货币符号
		amountFactor int64  // 金额换算因子, USD=100, VND=1
		status       int64  // 状态: 1启用, 2停用
		sortNo       int64  // 排序值, 数值越小越靠前
	}{
		{code: "USD", nameKey: "currency.usd.name", currencyType: 1, symbol: "$", amountFactor: 100, status: 1, sortNo: 10},
		{code: "PHP", nameKey: "currency.php.name", currencyType: 1, symbol: "₱", amountFactor: 100, status: 1, sortNo: 20},
		{code: "VND", nameKey: "currency.vnd.name", currencyType: 1, symbol: "₫", amountFactor: 1, status: 1, sortNo: 30},
		{code: "JPY", nameKey: "currency.jpy.name", currencyType: 1, symbol: "¥", amountFactor: 1, status: 1, sortNo: 40},
		{code: "MYR", nameKey: "currency.myr.name", currencyType: 1, symbol: "RM", amountFactor: 100, status: 1, sortNo: 50},
		{code: "THB", nameKey: "currency.thb.name", currencyType: 1, symbol: "฿", amountFactor: 100, status: 1, sortNo: 60},
		{code: "IDR", nameKey: "currency.idr.name", currencyType: 1, symbol: "Rp", amountFactor: 1, status: 1, sortNo: 70},
		{code: "SGD", nameKey: "currency.sgd.name", currencyType: 1, symbol: "S$", amountFactor: 100, status: 1, sortNo: 80},
		{code: "KRW", nameKey: "currency.krw.name", currencyType: 1, symbol: "₩", amountFactor: 1, status: 1, sortNo: 90},
		{code: "BRL", nameKey: "currency.brl.name", currencyType: 1, symbol: "R$", amountFactor: 100, status: 1, sortNo: 100},
		{code: "USDT", nameKey: "currency.usdt.name", currencyType: 2, symbol: "USDT", amountFactor: 100, status: 2, sortNo: 110},
	}

	builders := make([]*ent.CurrencyCreate, 0, len(data))

	for _, item := range data {
		builders = append(builders,
			s.DB.Currency.
				Create().
				SetCode(item.code).
				SetNameKey(item.nameKey).
				SetCurrencyType(item.currencyType).
				SetSymbol(item.symbol).
				SetAmountFactor(item.amountFactor).
				SetStatus(item.status).
				SetSortNo(item.sortNo),
		)
	}

	// 只补充缺失的系统货币, 已存在的数据保持不变
	logx.Must(
		s.DB.Currency.
			CreateBulk(builders...).
			OnConflictColumns(currency.FieldCode).
			Ignore().
			Exec(ctx),
	)
}

// mustInitRegion 初始化系统国家地区
func (s *ServiceContext) mustInitRegion(ctx context.Context) {
	data := []struct {
		code        string // 国家地区编码
		callingCode string // 国际电话区号
		nameKey     string // 名称翻译 Key
		status      int64  // 状态: 1启用, 2停用
		sortNo      int64  // 排序值, 数值越小越靠前
	}{
		{code: "JP", callingCode: "81", nameKey: "region.jp.name", status: 1, sortNo: 10},
		{code: "VN", callingCode: "84", nameKey: "region.vn.name", status: 1, sortNo: 20},
		{code: "PH", callingCode: "63", nameKey: "region.ph.name", status: 1, sortNo: 30},
		{code: "MY", callingCode: "60", nameKey: "region.my.name", status: 1, sortNo: 40},
		{code: "TH", callingCode: "66", nameKey: "region.th.name", status: 1, sortNo: 50},
		{code: "ID", callingCode: "62", nameKey: "region.id.name", status: 1, sortNo: 60},
		{code: "SG", callingCode: "65", nameKey: "region.sg.name", status: 1, sortNo: 70},
		{code: "KR", callingCode: "82", nameKey: "region.kr.name", status: 1, sortNo: 80},
		{code: "BR", callingCode: "55", nameKey: "region.br.name", status: 1, sortNo: 90},
		{code: "US", callingCode: "1", nameKey: "region.us.name", status: 1, sortNo: 100},
	}

	builders := make([]*ent.RegionCreate, 0, len(data))

	for _, item := range data {
		builders = append(builders,
			s.DB.Region.
				Create().
				SetCode(item.code).
				SetCallingCode(item.callingCode).
				SetNameKey(item.nameKey).
				SetStatus(item.status).
				SetSortNo(item.sortNo),
		)
	}

	// 只补充缺失的系统国家地区, 已存在的数据保持不变
	logx.Must(
		s.DB.Region.
			CreateBulk(builders...).
			OnConflictColumns(region.FieldCode).
			Ignore().
			Exec(ctx),
	)
}
