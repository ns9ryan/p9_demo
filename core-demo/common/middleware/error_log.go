package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"oa.98ent.com/p9/core/common/ctxdata"
	"oa.98ent.com/p9/core/common/errorlog"
	"oa.98ent.com/p9/core/common/i18n"
	"oa.98ent.com/p9/core/common/response"
	"oa.98ent.com/p9/core/common/utils"
	"oa.98ent.com/p9/core/common/xerr"

	"github.com/zeromicro/go-zero/rest"
)

// ErrorLog 错误日志中间件
func ErrorLog(serviceName string, rec errorlog.Recorder) rest.Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if rec == nil {
				next(w, r)
				return
			}
			ctx, bag := errorlog.WithBag(r.Context())
			r = r.WithContext(ctxdata.WithClaimsHolder(ctx))
			reqBody := readBody(r)
			cw := &captureWriter{ResponseWriter: w, status: http.StatusOK}
			start := time.Now()
			defer func() {
				if v := recover(); v != nil {
					bag.Subject = panicSubject(v)
					bag.Detail = xerr.ClipStack(string(debug.Stack()))
					response.FailCtx(r.Context(), cw, xerr.InternalServerError(i18n.InternalError))
				}
				// 如果状态码小于500，则不收集错误日志
				if !errorlog.ShouldCollect(cw.status) {
					return
				}
				claims := ctxdata.ClaimsFromCtx(r.Context())
				record := errorlog.Record{
					RequestMethod:  r.Method,
					RequestPath:    r.URL.Path,
					RequestQuery:   r.URL.RawQuery,
					RequestBody:    MaskJSON(clipBytes(reqBody, maxActionBody)),
					ServiceName:    serviceName,
					ResponseStatus: cw.status,
					ResponseBody:   MaskJSON(clipBytes(cw.body.Bytes(), maxActionBody)),
					Subject:        bag.Subject,
					Detail:         bag.Detail,
					DurationMS:     int(time.Since(start).Milliseconds()),
					ClientIP:       utils.ClientIP(r),
					UserAgent:      utils.UserAgent(r),
				}
				if claims != nil {
					record.UserID = claims.UserID
				}
				errorlog.Report(r.Context(), rec, record)
			}()
			next(cw, r)
		}
	}
}

func panicSubject(v any) string {
	if err, ok := v.(error); ok {
		return err.Error()
	}
	return fmt.Sprint(v)
}
