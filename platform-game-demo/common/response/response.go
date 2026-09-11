package response

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"oa.98ent.com/p9/platform-game/common/constant"
)

// Response 统一响应结构 - 遵循 core-api 规范
// 成功: { "code": 0, "msg": "ok", "data": {} }
// 失败: { "code": 40x/50x, "msg": "error message" }
type Response struct {
	Code int         `json:"code" comment:"错误码：0表示成功，其他为HTTP状态码"`
	Msg  string      `json:"msg" comment:"提示消息"`
	Data interface{} `json:"data" comment:"响应数据；无返回体的接口为null"`
}

// Success 返回成功响应
func Success(data interface{}) *Response {
	return &Response{
		Code: constant.CodeSuccess,
		Msg:  "ok",
		Data: data,
	}
}

// Error 返回错误响应
func Error(code int, message string) *Response {
	return &Response{
		Code: code,
		Msg:  message,
		Data: nil,
	}
}

// ErrorCode 根据错误码返回错误响应
func ErrorCode(code int) *Response {
	return Error(code, constant.GetMessage(code))
}

// ErrorWithStatusCode 返回错误响应并设置HTTP状态码
func ErrorWithStatusCode(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	httpx.OkJson(w, Error(code, message))
}

// ListResponse 列表响应结构（用于 data 字段）
type ListResponse struct {
	List  interface{} `json:"list" comment:"列表数据"`
	Total int64       `json:"total" comment:"总条数"`
}

// ListResponseData 构建列表格式的响应数据
// 分页规范：用户列表缺省20，角色列表缺省50，上限100
func ListResponseData(list interface{}, total int64) interface{} {
	return ListResponse{
		List:  list,
		Total: total,
	}
}

// SetupHTTPX 设置全局的 HTTP 响应处理
// 配置 go-zero 的 httpx 包来使用统一的响应格式
func SetupHTTPX() {
	// 设置成功响应处理
	httpx.SetOkHandler(func(ctx context.Context, data any) any {
		return Success(data)
	})

	// 设置错误响应处理
	httpx.SetErrorHandlerCtx(func(ctx context.Context, err error) (int, any) {
		// 默认返回 500 Internal Server Error
		statusCode := http.StatusInternalServerError
		errMsg := err.Error()

		// 可以根据错误类型进行更精细的处理
		// 例如：验证错误返回 400，权限错误返回 403 等
		// 但为了简化，这里直接返回 500 和错误信息

		return statusCode, Error(statusCode, errMsg)
	})
}
