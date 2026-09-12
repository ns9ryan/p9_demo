package service

import (
	"context"
	"testing"

	"oa.98ent.com/p9/core/common/ctxdata"
	"oa.98ent.com/p9/core/common/i18n"
	"oa.98ent.com/p9/core/common/xerr"
	"oa.98ent.com/p9/core/rpc/ent"
	"oa.98ent.com/p9/core/rpc/ent/loginlog"
)

func TestUpdateUserIpWhitelistAndLogin(t *testing.T) {
	d := testDeps(t, ModeOff)
	ctx := context.Background()
	u := createUserWithRole(t, d, "admin", "pass", nil)
	claims := &ctxdata.Claims{UserID: u.ID}

	if _, err := d.Login(ctx, LoginReq{Username: "admin", Password: "pass", ClientIP: "10.0.0.1"}); err != nil {
		t.Fatalf("default off: %v", err)
	}

	if err := d.UpdateUserIpWhitelist(ctx, claims, UpdateUserIpWhitelistReq{
		ID: u.ID, IPWhitelistEnabled: 1,
	}); err == nil || xerr.AsError(err).Message != i18n.UserInvalidIpWhitelist {
		t.Fatalf("enable empty: %v", err)
	}
	if err := d.UpdateUserIpWhitelist(ctx, claims, UpdateUserIpWhitelistReq{
		ID: u.ID, IPWhitelistEnabled: 1, IPWhitelist: []string{"not-an-ip"},
	}); err == nil || xerr.AsError(err).Message != i18n.UserInvalidIpWhitelist {
		t.Fatalf("bad ip: %v", err)
	}
	if err := d.UpdateUserIpWhitelist(ctx, claims, UpdateUserIpWhitelistReq{
		ID: u.ID, IPWhitelistEnabled: 3, IPWhitelist: []string{"1.1.1.1"},
	}); err == nil || xerr.AsError(err).Message != i18n.InvalidParam {
		t.Fatalf("bad flag: %v", err)
	}

	if err := d.UpdateUserIpWhitelist(ctx, claims, UpdateUserIpWhitelistReq{
		ID: u.ID, IPWhitelistEnabled: 1, IPWhitelist: []string{"::ffff:192.168.0.6", "10.0.0.0/8"},
	}); err != nil {
		t.Fatalf("update: %v", err)
	}

	got, _, err := d.GetUser(ctx, claims, u.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.IPWhitelistEnabled != 1 || len(got.IPWhitelist) != 2 || got.IPWhitelist[0] != "192.168.0.6" || got.IPWhitelist[1] != "10.0.0.0/8" {
		t.Fatalf("stored %+v %v", got.IPWhitelistEnabled, got.IPWhitelist)
	}

	if _, err := d.Login(ctx, LoginReq{Username: "admin", Password: "pass", ClientIP: "192.168.0.6"}); err != nil {
		t.Fatalf("exact: %v", err)
	}
	if _, err := d.Login(ctx, LoginReq{Username: "admin", Password: "pass", ClientIP: "::ffff:10.1.2.3"}); err != nil {
		t.Fatalf("cidr mapped: %v", err)
	}

	_, err = d.Login(ctx, LoginReq{Username: "admin", Password: "pass", ClientIP: "11.0.0.1"})
	if err == nil || xerr.AsError(err).Message != i18n.AuthIPNotAllowed {
		t.Fatalf("deny: %v", err)
	}
	row, err := d.Client.LoginLog.Query().Order(ent.Desc(loginlog.FieldID)).First(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if row.LoginResult == 1 || row.FailureReason == nil || *row.FailureReason != i18n.AuthIPNotAllowed {
		t.Fatalf("fail log %+v", row)
	}

	if err := d.UpdateUserIpWhitelist(ctx, claims, UpdateUserIpWhitelistReq{
		ID: u.ID, IPWhitelistEnabled: 0,
	}); err != nil {
		t.Fatalf("disable: %v", err)
	}
	if _, err := d.Login(ctx, LoginReq{Username: "admin", Password: "pass", ClientIP: "11.0.0.1"}); err != nil {
		t.Fatalf("disabled: %v", err)
	}
}
