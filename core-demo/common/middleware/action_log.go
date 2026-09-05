package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"oa.98ent.com/p9/core/common/ctxdata"
	"oa.98ent.com/p9/core/common/jwt"
	"oa.98ent.com/p9/core/common/utils"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
)

const maxActionBody = 16 * 1024

type ActionRecord struct {
	UserID         int64
	RequestMethod  string
	RequestPath    string
	RequestQuery   string
	RequestBody    string
	ActionResult   int16
	ResponseStatus int
	ResponseBody   string
	DurationMS     int
	ClientIP       string
	UserAgent      string
}

// ActionRecorder 动作记录器
type ActionRecorder interface {
	RecordAction(ctx context.Context, rec ActionRecord)
}

// ActionLog 动作日志中间件
func ActionLog(rec ActionRecorder) rest.Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if rec == nil || !IsAdminWrite(r.Method, r.URL.Path) {
				next(w, r)
				return
			}
			reqBody := readBody(r)
			cw := &captureWriter{ResponseWriter: w, status: http.StatusOK}
			start := time.Now()
			next(cw, r)
			claims := ctxdata.ClaimsFromCtx(r.Context())
			if claims == nil || claims.UserID == 0 || claims.TokenType == jwt.TokenPreview {
				return
			}
			status := cw.status
			result := int16(2)
			if status >= 200 && status < 300 {
				result = 1
			}
			record := ActionRecord{
				UserID:         claims.UserID,
				RequestMethod:  r.Method,
				RequestPath:    r.URL.Path,
				RequestQuery:   r.URL.RawQuery,
				RequestBody:    MaskJSON(clipBytes(reqBody, maxActionBody)),
				ActionResult:   result,
				ResponseStatus: status,
				ResponseBody:   MaskJSON(clipBytes(cw.body.Bytes(), maxActionBody)),
				DurationMS:     int(time.Since(start).Milliseconds()),
				ClientIP:       utils.ClientIP(r),
				UserAgent:      utils.UserAgent(r),
			}
			go func() {
				ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
				defer cancel()
				if claims != nil {
					ctx = ctxdata.WithClaims(ctx, claims)
				}
				rec.RecordAction(ctx, record)
			}()
		}
	}
}

func readBody(r *http.Request) []byte {
	if r.Body == nil {
		return nil
	}
	body, err := io.ReadAll(io.LimitReader(r.Body, maxActionBody+1))
	_ = r.Body.Close()
	if err != nil {
		logx.Errorf("action log read body: %v", err)
		r.Body = io.NopCloser(bytes.NewReader(nil))
		return nil
	}
	r.Body = io.NopCloser(bytes.NewReader(body))
	return body
}

func clipBytes(b []byte, n int) string {
	if n <= 0 || len(b) <= n {
		return string(b)
	}
	return string(b[:n])
}

type captureWriter struct {
	http.ResponseWriter
	status int
	body   bytes.Buffer
}

func (w *captureWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *captureWriter) Write(p []byte) (int, error) {
	if w.status == 0 {
		w.status = http.StatusOK
	}
	if w.body.Len() < maxActionBody {
		remain := maxActionBody - w.body.Len()
		if len(p) > remain {
			_, _ = w.body.Write(p[:remain])
		} else {
			_, _ = w.body.Write(p)
		}
	}
	return w.ResponseWriter.Write(p)
}

const maskedValue = "***"

var sensitiveKeys = map[string]struct{}{
	"password":      {},
	"old_password":  {},
	"new_password":  {},
	"refresh_token": {},
	"access_token":  {},
	"token":         {},
	"init_token":    {},
}

// MaskJSON 屏蔽敏感信息
func MaskJSON(raw string) string {
	s := strings.TrimSpace(raw)
	if s == "" {
		return ""
	}
	var v any
	if err := json.Unmarshal([]byte(s), &v); err != nil {
		return raw
	}
	maskValue(v)
	out, err := json.Marshal(v)
	if err != nil {
		return raw
	}
	if bytes.Equal(out, []byte("null")) {
		return raw
	}
	return string(out)
}

func maskValue(v any) {
	switch t := v.(type) {
	case map[string]any:
		for k, child := range t {
			if _, ok := sensitiveKeys[strings.ToLower(k)]; ok {
				t[k] = maskedValue
				continue
			}
			maskValue(child)
		}
	case []any:
		for _, child := range t {
			maskValue(child)
		}
	}
}
