package constant

// 错误码定义 - 遵循Simple Admin规范

// 标准HTTP状态码 - 遵循 core-api 规范
const (
	CodeSuccess       = 0   // 成功
	CodeBadRequest    = 400 // 参数错误
	CodeUnauthorized  = 401 // 未登录/无效token
	CodeForbidden     = 403 // 无权限或资源被禁用
	CodeNotFound      = 404 // 不存在
	CodeAccessExpired = 498 // access token 过期
	CodeInternalError = 500 // 内部错误
)

// 业务错误码 (10000-19999)
const (
	// 游戏相关
	CodeGameNotFound      = 10001
	CodeGameAlreadyExists = 10002
	CodeInvalidGameStatus = 10003
	CodeGameSyncFailed    = 10004

	// 分类相关
	CodeCategoryNotFound      = 10011
	CodeCategoryAlreadyExists = 10012

	// 厂商相关
	CodeProviderNotFound      = 10021
	CodeProviderAlreadyExists = 10022

	// 渠道相关
	CodeChannelNotFound      = 10031
	CodeChannelAlreadyExists = 10032

	// 币种相关
	CodeCurrencyNotFound = 10041

	// 同步相关
	CodeSyncConflict      = 10051
	CodeSyncPreviewFailed = 10052
	CodeSyncExecuteFailed = 10053
)

// 系统错误码 (20000-29999)
const (
	CodeDatabaseError  = 20001
	CodeInvalidRequest = 20002
	CodeNetworkError   = 20006
)

// 验证错误码 (30000-39999)
const (
	CodeValidationFailed = 30001
	CodeFieldInvalid     = 30002
	CodeFieldRequired    = 30003
)
