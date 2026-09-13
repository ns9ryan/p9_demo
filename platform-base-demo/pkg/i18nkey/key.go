package i18nkey

// 请求与数据错误
const (
	InvalidRequest  = "common.invalid_request"  // 请求参数格式错误
	DataNotFound    = "common.data_not_found"   // 数据不存在
	ConstraintError = "common.constraint_error" // 数据约束冲突
	ValidationError = "common.validation_error" // 数据校验失败
	DatabaseError   = "common.database_error"   // 数据库操作失败
)

// 系统与服务错误
const (
	InternalError      = "common.internal_error"      // 系统内部错误
	TooManyRequests    = "common.too_many_requests"   // 请求过于频繁
	ServiceUnavailable = "common.service_unavailable" // 服务暂不可用
	RequestTimeout     = "common.request_timeout"     // 请求超时
)
