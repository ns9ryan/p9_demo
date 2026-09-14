package service

import (
	"context"
	"testing"

	"oa.98ent.com/p9/core/common/ctxdata"
	"oa.98ent.com/p9/core/common/i18n"
	"oa.98ent.com/p9/core/common/jwt"
	"oa.98ent.com/p9/core/common/xerr"
)

func TestLoginRefreshTokenBindsClientIP(t *testing.T) {
	d := testDeps(t, ModeOff)
	createUserWithRole(t, d, "admin", "pass", nil)
	ctx := ctxdata.WithClientIP(context.Background(), "10.0.0.1")
	res, err := d.Login(ctx, LoginReq{Username: "admin", Password: "pass", ClientIP: "10.0.0.1"})
	if err != nil {
		t.Fatal(err)
	}
	access, err := jwt.Parse(d.JWTSecret, res.Token.AccessToken)
	if err != nil || access.ClientIP != "10.0.0.1" {
		t.Fatalf("access %+v %v", access, err)
	}
	refresh, err := jwt.ParseTyped(d.JWTRefreshSecret, res.Token.RefreshToken, jwt.TokenRefresh)
	if err != nil || refresh.ClientIP != "10.0.0.1" {
		t.Fatalf("refresh %+v %v", refresh, err)
	}

	got, err := d.CheckToken(ctx, res.Token.AccessToken)
	if err != nil {
		t.Fatal(err)
	}
	if got.ClientIP != "10.0.0.1" {
		t.Fatalf("claims ip %q", got.ClientIP)
	}
	if _, err := d.CheckToken(ctxdata.WithClientIP(context.Background(), "::ffff:10.0.0.1"), res.Token.AccessToken); err != nil {
		t.Fatalf("mapped: %v", err)
	}
	_, err = d.CheckToken(ctxdata.WithClientIP(context.Background(), "11.0.0.1"), res.Token.AccessToken)
	if got := xerr.AsError(err); got.Status != 401 || got.Message != i18n.AuthIPMismatch {
		t.Fatalf("check mismatch %+v", got)
	}

	same, err := d.Refresh(ctx, RefreshReq{RefreshToken: res.Token.RefreshToken})
	if err != nil {
		t.Fatal(err)
	}
	newAccess, err := jwt.Parse(d.JWTSecret, same.Token.AccessToken)
	if err != nil || newAccess.ClientIP != "10.0.0.1" {
		t.Fatalf("new access %+v %v", newAccess, err)
	}

	_, err = d.Refresh(ctxdata.WithClientIP(context.Background(), "11.0.0.1"), RefreshReq{RefreshToken: same.Token.RefreshToken})
	if got := xerr.AsError(err); got.Status != 401 || got.Message != i18n.AuthIPMismatch {
		t.Fatalf("refresh mismatch %+v", got)
	}
}

func TestLoginReqClientIPFallback(t *testing.T) {
	d := testDeps(t, ModeOff)
	createUserWithRole(t, d, "admin", "pass", nil)
	res, err := d.Login(context.Background(), LoginReq{Username: "admin", Password: "pass", ClientIP: "::ffff:192.168.0.6"})
	if err != nil {
		t.Fatal(err)
	}
	access, err := jwt.Parse(d.JWTSecret, res.Token.AccessToken)
	if err != nil || access.ClientIP != "192.168.0.6" {
		t.Fatalf("fallback %+v %v", access, err)
	}
}
