package constant

// 错误消息定义

var ErrorMessages = map[int]string{
	// 标准HTTP状态码 - 遵循 core-api 规范
	CodeSuccess:       "ok",
	CodeBadRequest:    "参数错误",
	CodeUnauthorized:  "未登录/无效token",
	CodeForbidden:     "无权限或资源被禁用",
	CodeNotFound:      "不存在",
	CodeAccessExpired: "access token 过期",
	CodeInternalError: "内部错误",

	// 业务错误
	CodeGameNotFound:      "游戏不存在",
	CodeGameAlreadyExists: "游戏已存在",
	CodeInvalidGameStatus: "无效的游戏状态",
	CodeGameSyncFailed:    "游戏同步失败",

	CodeCategoryNotFound:      "分类不存在",
	CodeCategoryAlreadyExists: "分类已存在",

	CodeProviderNotFound:      "厂商不存在",
	CodeProviderAlreadyExists: "厂商已存在",

	CodeChannelNotFound:      "渠道不存在",
	CodeChannelAlreadyExists: "渠道已存在",

	CodeCurrencyNotFound: "币种不存在",

	CodeSyncConflict:      "数据冲突",
	CodeSyncPreviewFailed: "同步预检查失败",
	CodeSyncExecuteFailed: "同步执行失败",

	// 系统错误
	CodeDatabaseError:  "数据库错误",
	CodeInvalidRequest: "无效请求",
	CodeNetworkError:   "网络错误",

	// 验证错误
	CodeValidationFailed: "数据验证失败",
	CodeFieldInvalid:     "字段无效",
	CodeFieldRequired:    "字段必填",
}

func GetMessage(code int) string {
	if msg, ok := ErrorMessages[code]; ok {
		return msg
	}
	return "未知错误"
}
