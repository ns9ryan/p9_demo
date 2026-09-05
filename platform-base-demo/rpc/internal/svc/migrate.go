package svc

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"oa.98ent.com/p9/platform-base/rpc/ent"
	"oa.98ent.com/p9/platform-base/rpc/ent/currency"
	"oa.98ent.com/p9/platform-base/rpc/ent/language"
	"oa.98ent.com/p9/platform-base/rpc/ent/region"
	"oa.98ent.com/p9/platform-base/rpc/ent/timezone"
)

// MustMigrate 执行数据库自动迁移
func (s *ServiceContext) MustMigrate() {
	ctx := context.Background()

	// 根据 Ent Schema 自动创建或更新数据库结构
	logx.Must(s.DB.Schema.Create(ctx))

	// 初始化系统默认数据
	s.mustInitLanguage(ctx)
	s.mustInitTimezone(ctx)
	s.mustInitCurrency(ctx)
	s.mustInitRegion(ctx)
}

// mustInitLanguage 初始化系统语言
func (s *ServiceContext) mustInitLanguage(ctx context.Context) {
	data := []struct {
		code     string
		nameI18n map[string]string
		status   int64
		sortNo   int64
	}{
		{
			code: "zh-CN",
			nameI18n: map[string]string{
				"zh-CN": "简体中文",
				"en-US": "Simplified Chinese",
				"vi-VN": "Tiếng Trung giản thể",
			},
			status: 1,
			sortNo: 10,
		},
		{
			code: "zh-TW",
			nameI18n: map[string]string{
				"zh-CN": "繁体中文",
				"en-US": "Traditional Chinese",
				"vi-VN": "Tiếng Trung phồn thể",
			},
			status: 2,
			sortNo: 20,
		},
		{
			code: "en-US",
			nameI18n: map[string]string{
				"zh-CN": "英语（美国）",
				"en-US": "English (United States)",
				"vi-VN": "Tiếng Anh (Hoa Kỳ)",
			},
			status: 2,
			sortNo: 30,
		},
		{
			code: "vi-VN",
			nameI18n: map[string]string{
				"zh-CN": "越南语",
				"en-US": "Vietnamese",
				"vi-VN": "Tiếng Việt",
			},
			status: 2,
			sortNo: 40,
		},
		{
			code: "th-TH",
			nameI18n: map[string]string{
				"zh-CN": "泰语",
				"en-US": "Thai",
				"vi-VN": "Tiếng Thái",
			},
			status: 2,
			sortNo: 50,
		},
		{
			code: "id-ID",
			nameI18n: map[string]string{
				"zh-CN": "印度尼西亚语",
				"en-US": "Indonesian",
				"vi-VN": "Tiếng Indonesia",
			},
			status: 2,
			sortNo: 60,
		},
		{
			code: "ms-MY",
			nameI18n: map[string]string{
				"zh-CN": "马来语",
				"en-US": "Malay",
				"vi-VN": "Tiếng Mã Lai",
			},
			status: 2,
			sortNo: 70,
		},
		{
			code: "ja-JP",
			nameI18n: map[string]string{
				"zh-CN": "日语",
				"en-US": "Japanese",
				"vi-VN": "Tiếng Nhật",
			},
			status: 2,
			sortNo: 80,
		},
		{
			code: "ko-KR",
			nameI18n: map[string]string{
				"zh-CN": "韩语",
				"en-US": "Korean",
				"vi-VN": "Tiếng Hàn",
			},
			status: 2,
			sortNo: 90,
		},
		{
			code: "pt-BR",
			nameI18n: map[string]string{
				"zh-CN": "葡萄牙语（巴西）",
				"en-US": "Portuguese (Brazil)",
				"vi-VN": "Tiếng Bồ Đào Nha (Brazil)",
			},
			status: 2,
			sortNo: 100,
		},
	}

	builders := make([]*ent.LanguageCreate, 0, len(data))

	for _, item := range data {
		builders = append(builders,
			s.DB.Language.
				Create().
				SetCode(item.code).
				SetNameI18n(item.nameI18n).
				SetStatus(item.status).
				SetSortNo(item.sortNo),
		)
	}

	// 只补充缺失的系统语言，已存在的数据保持不变
	logx.Must(
		s.DB.Language.
			CreateBulk(builders...).
			OnConflictColumns(language.FieldCode).
			Ignore().
			Exec(ctx),
	)
}

// mustInitTimezone 初始化系统时区
func (s *ServiceContext) mustInitTimezone(ctx context.Context) {
	data := []struct {
		code     string
		nameI18n map[string]string
		status   int64
		sortNo   int64
	}{
		{
			code: "Etc/UTC",
			nameI18n: map[string]string{
				"zh-CN": "协调世界时",
				"en-US": "Coordinated Universal Time",
				"vi-VN": "Giờ Phối hợp Quốc tế",
			},
			status: 1,
			sortNo: 10,
		},
		{
			code: "Asia/Tokyo",
			nameI18n: map[string]string{
				"zh-CN": "日本标准时间",
				"en-US": "Japan Standard Time",
				"vi-VN": "Giờ chuẩn Nhật Bản",
			},
			status: 1,
			sortNo: 20,
		},
		{
			code: "Asia/Ho_Chi_Minh",
			nameI18n: map[string]string{
				"zh-CN": "越南时间",
				"en-US": "Vietnam Time",
				"vi-VN": "Giờ Việt Nam",
			},
			status: 1,
			sortNo: 30,
		},
		{
			code: "Asia/Manila",
			nameI18n: map[string]string{
				"zh-CN": "菲律宾时间",
				"en-US": "Philippine Time",
				"vi-VN": "Giờ Philippines",
			},
			status: 1,
			sortNo: 40,
		},
		{
			code: "Asia/Bangkok",
			nameI18n: map[string]string{
				"zh-CN": "泰国时间",
				"en-US": "Thailand Time",
				"vi-VN": "Giờ Thái Lan",
			},
			status: 1,
			sortNo: 50,
		},
		{
			code: "Asia/Kuala_Lumpur",
			nameI18n: map[string]string{
				"zh-CN": "马来西亚时间",
				"en-US": "Malaysia Time",
				"vi-VN": "Giờ Malaysia",
			},
			status: 1,
			sortNo: 60,
		},
		{
			code: "Asia/Singapore",
			nameI18n: map[string]string{
				"zh-CN": "新加坡标准时间",
				"en-US": "Singapore Standard Time",
				"vi-VN": "Giờ chuẩn Singapore",
			},
			status: 1,
			sortNo: 70,
		},
		{
			code: "Asia/Jakarta",
			nameI18n: map[string]string{
				"zh-CN": "印度尼西亚西部时间",
				"en-US": "Western Indonesia Time",
				"vi-VN": "Giờ Tây Indonesia",
			},
			status: 1,
			sortNo: 80,
		},
		{
			code: "Asia/Seoul",
			nameI18n: map[string]string{
				"zh-CN": "韩国标准时间",
				"en-US": "Korea Standard Time",
				"vi-VN": "Giờ chuẩn Hàn Quốc",
			},
			status: 1,
			sortNo: 90,
		},
		{
			code: "America/Sao_Paulo",
			nameI18n: map[string]string{
				"zh-CN": "巴西利亚时间",
				"en-US": "Brasilia Time",
				"vi-VN": "Giờ Brasilia",
			},
			status: 1,
			sortNo: 100,
		},
	}

	builders := make([]*ent.TimezoneCreate, 0, len(data))

	for _, item := range data {
		builders = append(builders,
			s.DB.Timezone.
				Create().
				SetCode(item.code).
				SetNameI18n(item.nameI18n).
				SetStatus(item.status).
				SetSortNo(item.sortNo),
		)
	}

	// 只补充缺失的系统时区，已存在的数据保持不变
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
		code         string
		nameI18n     map[string]string
		currencyType int64
		symbol       string
		amountFactor int64
		status       int64
		sortNo       int64
	}{
		{
			code: "USD",
			nameI18n: map[string]string{
				"zh-CN": "美元",
				"en-US": "US Dollar",
				"vi-VN": "Đô la Mỹ",
			},
			currencyType: 1,
			symbol:       "$",
			amountFactor: 100,
			status:       1,
			sortNo:       10,
		},
		{
			code: "PHP",
			nameI18n: map[string]string{
				"zh-CN": "菲律宾比索",
				"en-US": "Philippine Peso",
				"vi-VN": "Peso Philippines",
			},
			currencyType: 1,
			symbol:       "₱",
			amountFactor: 100,
			status:       1,
			sortNo:       20,
		},
		{
			code: "VND",
			nameI18n: map[string]string{
				"zh-CN": "越南盾",
				"en-US": "Vietnamese Dong",
				"vi-VN": "Đồng Việt Nam",
			},
			currencyType: 1,
			symbol:       "₫",
			amountFactor: 1,
			status:       1,
			sortNo:       30,
		},
		{
			code: "JPY",
			nameI18n: map[string]string{
				"zh-CN": "日元",
				"en-US": "Japanese Yen",
				"vi-VN": "Yên Nhật",
			},
			currencyType: 1,
			symbol:       "¥",
			amountFactor: 1,
			status:       1,
			sortNo:       40,
		},
		{
			code: "MYR",
			nameI18n: map[string]string{
				"zh-CN": "马来西亚林吉特",
				"en-US": "Malaysian Ringgit",
				"vi-VN": "Ringgit Malaysia",
			},
			currencyType: 1,
			symbol:       "RM",
			amountFactor: 100,
			status:       1,
			sortNo:       50,
		},
		{
			code: "THB",
			nameI18n: map[string]string{
				"zh-CN": "泰铢",
				"en-US": "Thai Baht",
				"vi-VN": "Baht Thái",
			},
			currencyType: 1,
			symbol:       "฿",
			amountFactor: 100,
			status:       1,
			sortNo:       60,
		},
		{
			code: "IDR",
			nameI18n: map[string]string{
				"zh-CN": "印度尼西亚卢比",
				"en-US": "Indonesian Rupiah",
				"vi-VN": "Rupiah Indonesia",
			},
			currencyType: 1,
			symbol:       "Rp",
			amountFactor: 1,
			status:       1,
			sortNo:       70,
		},
		{
			code: "SGD",
			nameI18n: map[string]string{
				"zh-CN": "新加坡元",
				"en-US": "Singapore Dollar",
				"vi-VN": "Đô la Singapore",
			},
			currencyType: 1,
			symbol:       "S$",
			amountFactor: 100,
			status:       1,
			sortNo:       80,
		},
		{
			code: "KRW",
			nameI18n: map[string]string{
				"zh-CN": "韩元",
				"en-US": "South Korean Won",
				"vi-VN": "Won Hàn Quốc",
			},
			currencyType: 1,
			symbol:       "₩",
			amountFactor: 1,
			status:       1,
			sortNo:       90,
		},
		{
			code: "BRL",
			nameI18n: map[string]string{
				"zh-CN": "巴西雷亚尔",
				"en-US": "Brazilian Real",
				"vi-VN": "Real Brazil",
			},
			currencyType: 1,
			symbol:       "R$",
			amountFactor: 100,
			status:       1,
			sortNo:       100,
		},
		{
			code: "USDT",
			nameI18n: map[string]string{
				"zh-CN": "泰达币",
				"en-US": "Tether",
				"vi-VN": "Tether",
			},
			currencyType: 2,
			symbol:       "USDT",
			amountFactor: 100,
			status:       2,
			sortNo:       110,
		},
	}

	builders := make([]*ent.CurrencyCreate, 0, len(data))

	for _, item := range data {
		builders = append(builders,
			s.DB.Currency.
				Create().
				SetCode(item.code).
				SetNameI18n(item.nameI18n).
				SetCurrencyType(item.currencyType).
				SetSymbol(item.symbol).
				SetAmountFactor(item.amountFactor).
				SetStatus(item.status).
				SetSortNo(item.sortNo),
		)
	}

	// 只补充缺失的系统货币，已存在的数据保持不变
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
		code        string
		callingCode string
		nameI18n    map[string]string
		status      int64
		sortNo      int64
	}{
		{
			code:        "JP",
			callingCode: "81",
			nameI18n: map[string]string{
				"zh-CN": "日本",
				"en-US": "Japan",
				"vi-VN": "Nhật Bản",
			},
			status: 1,
			sortNo: 10,
		},
		{
			code:        "VN",
			callingCode: "84",
			nameI18n: map[string]string{
				"zh-CN": "越南",
				"en-US": "Vietnam",
				"vi-VN": "Việt Nam",
			},
			status: 1,
			sortNo: 20,
		},
		{
			code:        "PH",
			callingCode: "63",
			nameI18n: map[string]string{
				"zh-CN": "菲律宾",
				"en-US": "Philippines",
				"vi-VN": "Philippines",
			},
			status: 1,
			sortNo: 30,
		},
		{
			code:        "MY",
			callingCode: "60",
			nameI18n: map[string]string{
				"zh-CN": "马来西亚",
				"en-US": "Malaysia",
				"vi-VN": "Malaysia",
			},
			status: 1,
			sortNo: 40,
		},
		{
			code:        "TH",
			callingCode: "66",
			nameI18n: map[string]string{
				"zh-CN": "泰国",
				"en-US": "Thailand",
				"vi-VN": "Thái Lan",
			},
			status: 1,
			sortNo: 50,
		},
		{
			code:        "ID",
			callingCode: "62",
			nameI18n: map[string]string{
				"zh-CN": "印度尼西亚",
				"en-US": "Indonesia",
				"vi-VN": "Indonesia",
			},
			status: 1,
			sortNo: 60,
		},
		{
			code:        "SG",
			callingCode: "65",
			nameI18n: map[string]string{
				"zh-CN": "新加坡",
				"en-US": "Singapore",
				"vi-VN": "Singapore",
			},
			status: 1,
			sortNo: 70,
		},
		{
			code:        "KR",
			callingCode: "82",
			nameI18n: map[string]string{
				"zh-CN": "韩国",
				"en-US": "South Korea",
				"vi-VN": "Hàn Quốc",
			},
			status: 1,
			sortNo: 80,
		},
		{
			code:        "BR",
			callingCode: "55",
			nameI18n: map[string]string{
				"zh-CN": "巴西",
				"en-US": "Brazil",
				"vi-VN": "Brazil",
			},
			status: 1,
			sortNo: 90,
		},
		{
			code:        "US",
			callingCode: "1",
			nameI18n: map[string]string{
				"zh-CN": "美国",
				"en-US": "United States",
				"vi-VN": "Hoa Kỳ",
			},
			status: 1,
			sortNo: 100,
		},
	}

	builders := make([]*ent.RegionCreate, 0, len(data))

	for _, item := range data {
		builders = append(builders,
			s.DB.Region.
				Create().
				SetCode(item.code).
				SetCallingCode(item.callingCode).
				SetNameI18n(item.nameI18n).
				SetStatus(item.status).
				SetSortNo(item.sortNo),
		)
	}

	// 只补充缺失的系统国家地区，已存在的数据保持不变
	logx.Must(
		s.DB.Region.
			CreateBulk(builders...).
			OnConflictColumns(region.FieldCode).
			Ignore().
			Exec(ctx),
	)
}
