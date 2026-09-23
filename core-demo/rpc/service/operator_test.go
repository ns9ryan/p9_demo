package service

import (
	"context"
	"testing"

	"oa.98ent.com/p9/common/ctxdata"
	"oa.98ent.com/p9/common/xerr"
	"oa.98ent.com/p9/core/rpc/model"
)

func TestIssuePreviewTokenModeAndCode(t *testing.T) {
	d := &Deps{Mode: ModeOff, JWTSecret: "secret", JWTExpire: 60}
	_, err := d.IssuePreviewToken(context.Background(), IssuePreviewTokenReq{OperatorCode: "demo"})
	if got := xerr.AsError(err); got.Status != 400 {
		t.Fatalf("mode off %+v", got)
	}
	d.Mode = ModeOn
	_, err = d.IssuePreviewToken(context.Background(), IssuePreviewTokenReq{})
	if got := xerr.AsError(err); got.Status != 400 {
		t.Fatalf("empty code %+v", got)
	}
}

func TestIssuePreviewTokenBySuperAdmin(t *testing.T) {
	d := testDeps(t, ModeOn)
	ctx := ctxdata.WithClientIP(context.Background(), "10.0.0.1")
	code := "demo"
	hash, err := HashPassword("pass")
	if err != nil {
		t.Fatal(err)
	}
	_, err = d.Client.User.Create().
		SetUserCode(NewUserCode()).
		SetUsername("admin").
		SetPasswordHash(hash).
		SetSalt("salt").
		SetDisplayName("admin").
		SetStatus(model.StatusNormal).
		SetIsSuperAdmin(true).
		SetOperatorCode(code).
		Save(ctx)
	if err != nil {
		t.Fatal(err)
	}
	tok, err := d.IssuePreviewToken(ctx, IssuePreviewTokenReq{OperatorCode: code})
	if err != nil {
		t.Fatal(err)
	}
	if tok.AccessToken == "" || tok.OperatorCode != code || tok.HomePath != "/dashboard" {
		t.Fatalf("%+v", tok)
	}
	_, err = d.IssuePreviewToken(ctx, IssuePreviewTokenReq{OperatorCode: "missing"})
	if got := xerr.AsError(err); got.Status != 404 {
		t.Fatalf("missing %+v", got)
	}
}
