package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"oa.98ent.com/p9/core/common/ctxdata"
	"oa.98ent.com/p9/core/common/errorlog"
	"oa.98ent.com/p9/core/common/response"
)

type errChanRecorder struct {
	ch chan errorlog.Record
}

func (c errChanRecorder) RecordError(_ context.Context, rec errorlog.Record) {
	c.ch <- rec
}

func waitErr(t *testing.T, ch <-chan errorlog.Record) errorlog.Record {
	t.Helper()
	select {
	case rec := <-ch:
		return rec
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting error log")
		return errorlog.Record{}
	}
}

func TestErrorLogMiddlewareWrites5xx(t *testing.T) {
	ch := make(chan errorlog.Record, 1)
	h := ErrorLog("core-api", errChanRecorder{ch: ch})(func(w http.ResponseWriter, r *http.Request) {
		response.FailCtx(r.Context(), w, http.ErrAbortHandler)
	})
	req := httptest.NewRequest(http.MethodPost, "/admin/user/create?x=1", strings.NewReader(`{"password":"secret"}`))
	req = req.WithContext(ctxdata.WithClaims(req.Context(), &ctxdata.Claims{UserID: 9, OperatorID: 3}))
	rr := httptest.NewRecorder()
	h(rr, req)
	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("status %d", rr.Code)
	}
	if strings.Contains(rr.Body.String(), "net/http") || strings.Contains(rr.Body.String(), "goroutine") {
		t.Fatalf("stack leaked: %s", rr.Body.String())
	}
	rec := waitErr(t, ch)
	if rec.ServiceName != "core-api" || rec.ResponseStatus != http.StatusInternalServerError || rec.UserID != 9 {
		t.Fatalf("%+v", rec)
	}
	if rec.Subject == "" || rec.Detail == "" {
		t.Fatalf("subject/detail %+v", rec)
	}
	if strings.Contains(rec.RequestBody, "secret") {
		t.Fatalf("body %s", rec.RequestBody)
	}
}

func TestErrorLogMiddlewareReadsInnerJWTClaims(t *testing.T) {
	ch := make(chan errorlog.Record, 1)
	jwtLike := func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			ctx := ctxdata.WithClaims(r.Context(), &ctxdata.Claims{UserID: 9, OperatorID: 3})
			next(w, r.WithContext(ctx))
		}
	}
	h := ErrorLog("core-api", errChanRecorder{ch: ch})(jwtLike(func(w http.ResponseWriter, r *http.Request) {
		response.FailCtx(r.Context(), w, http.ErrAbortHandler)
	}))
	h(httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/admin/user/create", strings.NewReader(`{}`)))
	rec := waitErr(t, ch)
	if rec.UserID != 9 {
		t.Fatalf("%+v", rec)
	}
}

func TestErrorLogMiddlewareSkips4xx(t *testing.T) {
	ch := make(chan errorlog.Record, 1)
	h := ErrorLog("core-api", errChanRecorder{ch: ch})(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"code":400}`))
	})
	h(httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/admin/user/create", strings.NewReader(`{}`)))
	select {
	case rec := <-ch:
		t.Fatalf("unexpected %+v", rec)
	case <-time.After(200 * time.Millisecond):
	}
}

func TestErrorLogMiddlewarePanic(t *testing.T) {
	ch := make(chan errorlog.Record, 1)
	h := ErrorLog("core-api", errChanRecorder{ch: ch})(func(w http.ResponseWriter, r *http.Request) {
		panic("kaboom")
	})
	rr := httptest.NewRecorder()
	h(rr, httptest.NewRequest(http.MethodGet, "/admin/login", nil))
	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("status %d body %s", rr.Code, rr.Body.String())
	}
	if strings.Contains(rr.Body.String(), "kaboom") || strings.Contains(rr.Body.String(), "goroutine") {
		t.Fatalf("leaked: %s", rr.Body.String())
	}
	rec := waitErr(t, ch)
	if rec.Subject != "kaboom" || rec.Detail == "" || rec.UserID != 0 {
		t.Fatalf("%+v", rec)
	}
}
