package validate

// Error 参数校验错误
type Error struct {
	message string // 错误信息
}

// NewError 创建参数校验错误
func NewError(message string) *Error {
	return &Error{
		message: message,
	}
}

// Error 返回错误信息
func (e *Error) Error() string {
	return e.message
}
