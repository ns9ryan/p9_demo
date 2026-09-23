package middleware

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"oa.98ent.com/p9/common/ctxdata"
)

type chanRecorder struct {
	ch chan ActionRecord
}

func (c chanRecorder) RecordAction(_ context.Context, rec ActionRecord) {
	c.ch <- rec
}

func TestIsAdminWrite(t *testing.T) {
	if IsAdminWrite(http.MethodPost, "/admin/user/list") {
		t.Fatal("list")
	}
	if !IsAdminWrite(http.MethodPost, "/admin/user/create") {
		t.Fatal("create")
	}
	if !IsAdminWrite(http.MethodPost, "/admin/logout") {
		t.Fatal("logout")
	}
	if IsAdminWrite(http.MethodGet, "/admin/user/detail") {
		t.Fatal("detail")
	}
}

func waitRec(t *testing.T, ch <-chan ActionRecord) ActionRecord {
	t.Helper()
	select {
	case rec := <-ch:
		return rec
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting action log")
		return ActionRecord{}
	}
}

func TestActionLogMiddlewareWrites(t *testing.T) {
	ch := make(chan ActionRecord, 1)
	h := ActionLog(chanRecorder{ch: ch})(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"code":0,"msg":"ok"}`))
	})
	req := httptest.NewRequest(http.MethodPost, "/admin/user/create", strings.NewReader(`{"username":"a","password":"secret"}`))
	req = req.WithContext(ctxdata.WithClaims(req.Context(), &ctxdata.Claims{UserID: 9}))
	rr := httptest.NewRecorder()
	h(rr, req)
	rec := waitRec(t, ch)
	if rec.UserID != 9 || rec.ActionResult != 1 || rec.ResponseStatus != http.StatusOK {
		t.Fatalf("%+v", rec)
	}
	if rec.RequestPath != "/admin/user/create" {
		t.Fatalf("path %s", rec.RequestPath)
	}
	if strings.Contains(rec.RequestBody, "secret") || !strings.Contains(rec.RequestBody, "***") {
		t.Fatalf("body %s", rec.RequestBody)
	}
}

func TestActionLogMiddlewareFailStatus(t *testing.T) {
	ch := make(chan ActionRecord, 1)
	h := ActionLog(chanRecorder{ch: ch})(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte(`{"code":403}`))
	})
	req := httptest.NewRequest(http.MethodPost, "/admin/user/delete", strings.NewReader(`{}`))
	req = req.WithContext(ctxdata.WithClaims(req.Context(), &ctxdata.Claims{UserID: 1}))
	h(httptest.NewRecorder(), req)
	rec := waitRec(t, ch)
	if rec.ActionResult != 2 || rec.ResponseStatus != http.StatusForbidden {
		t.Fatalf("%+v", rec)
	}
}

func TestActionLogMiddlewareSkips(t *testing.T) {
	ch := make(chan ActionRecord, 1)
	h := ActionLog(chanRecorder{ch: ch})(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	list := httptest.NewRequest(http.MethodPost, "/admin/user/list", strings.NewReader(`{}`))
	list = list.WithContext(ctxdata.WithClaims(list.Context(), &ctxdata.Claims{UserID: 1}))
	h(httptest.NewRecorder(), list)

	noClaims := httptest.NewRequest(http.MethodPost, "/admin/user/create", strings.NewReader(`{}`))
	h(httptest.NewRecorder(), noClaims)

	select {
	case rec := <-ch:
		t.Fatalf("unexpected %+v", rec)
	case <-time.After(200 * time.Millisecond):
	}
}

func TestActionLogMiddlewareLeavesMultipartBody(t *testing.T) {
	payload := bytes.Repeat([]byte("x"), maxActionBody+1024)
	var buf bytes.Buffer
	mw := multipart.NewWriter(&buf)
	fw, err := mw.CreateFormFile("file", "i18n.json")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := fw.Write(payload); err != nil {
		t.Fatal(err)
	}
	if err := mw.Close(); err != nil {
		t.Fatal(err)
	}

	ch := make(chan ActionRecord, 1)
	h := ActionLog(chanRecorder{ch: ch})(func(w http.ResponseWriter, r *http.Request) {
		f, _, err := r.FormFile("file")
		if err != nil {
			t.Errorf("form file: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer f.Close()
		got, err := io.ReadAll(f)
		if err != nil {
			t.Errorf("read file: %v", err)
		}
		if !bytes.Equal(got, payload) {
			t.Errorf("truncated file len=%d want=%d", len(got), len(payload))
		}
		w.WriteHeader(http.StatusOK)
	})
	req := httptest.NewRequest(http.MethodPost, "/admin/i18n/import", bytes.NewReader(buf.Bytes()))
	req.Header.Set("Content-Type", mw.FormDataContentType())
	req = req.WithContext(ctxdata.WithClaims(req.Context(), &ctxdata.Claims{UserID: 1}))
	h(httptest.NewRecorder(), req)
	rec := waitRec(t, ch)
	if rec.RequestBody != multipartBodyMark {
		t.Fatalf("body=%q", rec.RequestBody)
	}
}

func TestMaskJSON(t *testing.T) {
	got := MaskJSON(`{"username":"a","password":"secret","nested":{"refresh_token":"abc","keep":1}}`)
	if !strings.Contains(got, `"password":"***"`) {
		t.Fatalf("password not masked: %s", got)
	}
	if strings.Contains(got, `"secret"`) {
		t.Fatalf("secret leaked: %s", got)
	}
	if !strings.Contains(got, `"keep":1`) {
		t.Fatalf("keep lost: %s", got)
	}
	if MaskJSON("not-json") != "not-json" {
		t.Fatal("passthrough")
	}
	if MaskJSON("") != "" {
		t.Fatal("empty")
	}
}
