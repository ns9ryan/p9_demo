package service

import (
	"context"
	"testing"
	"time"

	"oa.98ent.com/p9/core/common/ctxdata"
	"oa.98ent.com/p9/core/common/i18n"
	"oa.98ent.com/p9/core/common/jwt"
	"oa.98ent.com/p9/core/common/xerr"

	jwtv5 "github.com/golang-jwt/jwt/v5"
)

func TestCheckTokenExpired(t *testing.T) {
	d := &Deps{JWTSecret: "secret"}
	tok := mustSignExpired(t, "secret", jwt.Claims{UserID: 1, Salt: "s1", TokenType: jwt.TokenAccess})
	_, err := d.CheckToken(context.Background(), tok)
	got := xerr.AsError(err)
	if got.Status != xerr.StatusTokenExpired || got.Message != i18n.TokenExpired {
		t.Fatalf("got %+v", got)
	}
}

func TestCheckTokenInvalidAndRefreshType(t *testing.T) {
	d := &Deps{JWTSecret: "secret"}
	_, err := d.CheckToken(context.Background(), "not-a-token")
	if got := xerr.AsError(err); got.Status != 401 {
		t.Fatalf("invalid token %+v", got)
	}
	refresh, _, err := jwt.Sign("secret", 60, jwt.Claims{UserID: 1, Salt: "s1", TokenType: jwt.TokenRefresh})
	if err != nil {
		t.Fatal(err)
	}
	_, err = d.CheckToken(context.Background(), refresh)
	if got := xerr.AsError(err); got.Status != 401 {
		t.Fatalf("refresh type %+v", got)
	}
}

func TestCheckTokenPreviewMissingOperator(t *testing.T) {
	d := &Deps{JWTSecret: "secret"}
	tok, _, err := jwt.Sign("secret", 60, jwt.Claims{
		TokenType: jwt.TokenPreview,
		RoleCodes: []string{"super_admin"},
	})
	if err != nil {
		t.Fatal(err)
	}
	_, err = d.CheckToken(context.Background(), tok)
	if got := xerr.AsError(err); got.Status != 401 {
		t.Fatalf("preview without operator %+v", got)
	}
}

func TestCheckTokenClientIP(t *testing.T) {
	d := &Deps{JWTSecret: "secret"}
	tok, _, err := jwt.Sign("secret", 60, jwt.Claims{
		UserID: 1, Salt: "s1", TokenType: jwt.TokenAccess, ClientIP: "10.0.0.1",
	})
	if err != nil {
		t.Fatal(err)
	}
	_, err = d.CheckToken(context.Background(), tok)
	if got := xerr.AsError(err); got.Status != 401 || got.Message != i18n.AuthIPMismatch {
		t.Fatalf("empty ctx %+v", got)
	}
	_, err = d.CheckToken(ctxdata.WithClientIP(context.Background(), "11.0.0.1"), tok)
	if got := xerr.AsError(err); got.Status != 401 || got.Message != i18n.AuthIPMismatch {
		t.Fatalf("mismatch %+v", got)
	}
	old, _, err := jwt.Sign("secret", 60, jwt.Claims{UserID: 1, Salt: "s1", TokenType: jwt.TokenAccess})
	if err != nil {
		t.Fatal(err)
	}
	_, err = d.CheckToken(ctxdata.WithClientIP(context.Background(), "10.0.0.1"), old)
	if got := xerr.AsError(err); got.Status != 401 || got.Message != i18n.AuthIPMismatch {
		t.Fatalf("old token %+v", got)
	}
}

func mustSignExpired(t *testing.T, secret string, c jwt.Claims) string {
	t.Helper()
	c.RegisteredClaims = jwtv5.RegisteredClaims{
		IssuedAt:  jwtv5.NewNumericDate(time.Now().Add(-2 * time.Hour)),
		ExpiresAt: jwtv5.NewNumericDate(time.Now().Add(-time.Hour)),
	}
	s, err := jwtv5.NewWithClaims(jwtv5.SigningMethodHS256, c).SignedString([]byte(secret))
	if err != nil {
		t.Fatal(err)
	}
	return s
}
