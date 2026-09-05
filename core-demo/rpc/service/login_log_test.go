package service

import (
	"context"
	"fmt"
	"testing"
	"time"

	"oa.98ent.com/p9/core/common/ctxdata"
	"oa.98ent.com/p9/core/common/entmixin"
	"oa.98ent.com/p9/core/common/i18n"
	"oa.98ent.com/p9/core/rpc/ent"
	"oa.98ent.com/p9/core/rpc/ent/enttest"
	"oa.98ent.com/p9/core/rpc/ent/intercept"
	"oa.98ent.com/p9/core/rpc/ent/loginlog"
	"oa.98ent.com/p9/core/rpc/model"

	"entgo.io/ent/dialect"
	_ "github.com/mattn/go-sqlite3"
)

func testClient(t *testing.T) *ent.Client {
	t.Helper()
	dsn := fmt.Sprintf("file:%s?mode=memory&cache=shared&_fk=1", t.Name())
	client := enttest.Open(t, dialect.SQLite, dsn)
	fOperatorID := intercept.TraverseFunc(func(ctx context.Context, q intercept.Query) error {
		entmixin.FilterOperatorID(ctx, q)
		return nil
	})
	client.LoginLog.Intercept(fOperatorID)
	client.AdminActionLog.Intercept(fOperatorID)
	client.ErrorLog.Intercept(fOperatorID)
	t.Cleanup(func() { _ = client.Close() })
	return client
}

func testDeps(t *testing.T, mode string) *Deps {
	t.Helper()
	return &Deps{
		Client:           testClient(t),
		Mode:             mode,
		JWTSecret:        "secret",
		JWTExpire:        60,
		JWTRefreshSecret: "refresh",
		JWTRefreshExpire: 120,
	}
}

func createUserWithRole(t *testing.T, d *Deps, username, password string, operatorID *int64) *ent.User {
	t.Helper()
	ctx := context.Background()
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatal(err)
	}
	rb := d.Client.Role.Create().
		SetRoleCode("admin_" + username).
		SetRoleName("Admin").
		SetStatus(model.StatusNormal)
	if operatorID != nil {
		rb.SetOperatorID(*operatorID)
	}
	role, err := rb.Save(ctx)
	if err != nil {
		t.Fatal(err)
	}
	ub := d.Client.User.Create().
		SetUserCode(NewUserCode()).
		SetUsername(username).
		SetPasswordHash(hash).
		SetSalt("salt").
		SetDisplayName(username).
		SetStatus(model.StatusNormal).
		AddRoleIDs(role.ID)
	if operatorID != nil {
		ub.SetOperatorID(*operatorID)
	}
	u, err := ub.Save(ctx)
	if err != nil {
		t.Fatal(err)
	}
	return u
}

func TestLoginWritesSuccessAndFailLogs(t *testing.T) {
	d := testDeps(t, ModeOff)
	ctx := context.Background()
	u := createUserWithRole(t, d, "admin", "pass", nil)

	if _, err := d.Login(ctx, LoginReq{Username: "admin", Password: "pass", ClientIP: "127.0.0.1", UserAgent: "ua"}); err != nil {
		t.Fatalf("login: %v", err)
	}
	_, err := d.Login(ctx, LoginReq{Username: "admin", Password: "wrong", ClientIP: "10.0.0.1"})
	if err == nil {
		t.Fatal("expected fail")
	}
	_, err = d.Login(ctx, LoginReq{Username: "nobody", Password: "x", ClientIP: "1.1.1.1"})
	if err == nil {
		t.Fatal("expected unknown user fail")
	}

	rows, err := d.Client.LoginLog.Query().Order(ent.Asc(loginlog.FieldID)).All(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 3 {
		t.Fatalf("got %d logs", len(rows))
	}
	if rows[0].LoginResult != model.LoginResultSuccess || rows[0].UserID == nil || *rows[0].UserID != u.ID {
		t.Fatalf("success log %+v", rows[0])
	}
	if rows[0].UserAgent == nil || *rows[0].UserAgent != "ua" {
		t.Fatalf("user agent %+v", rows[0].UserAgent)
	}
	if rows[1].LoginResult != model.LoginResultFail || rows[1].UserID == nil || *rows[1].UserID != u.ID {
		t.Fatalf("wrong password log %+v", rows[1])
	}
	if rows[1].FailureReason == nil || *rows[1].FailureReason != i18n.AuthPasswordIncorrect {
		t.Fatalf("reason %+v", rows[1].FailureReason)
	}
	if rows[2].LoginResult != model.LoginResultFail || rows[2].UserID != nil {
		t.Fatalf("unknown user log %+v", rows[2])
	}
	for i, row := range rows {
		if row.OperatorID != nil {
			t.Fatalf("mode off log %d operator_id=%v", i, row.OperatorID)
		}
	}
}

func TestLoginWritesOperatorIDModeOn(t *testing.T) {
	d := testDeps(t, ModeOn)
	ctx := context.Background()
	op, err := d.Client.Operator.Create().
		SetOperatorCode("op1").
		SetTimezoneCode("UTC").
		SetSettlementCurrencyCode("USD").
		SetStatus(model.StatusNormal).
		Save(ctx)
	if err != nil {
		t.Fatal(err)
	}
	u := createUserWithRole(t, d, "admin", "pass", &op.ID)

	if _, err := d.Login(ctx, LoginReq{Username: "admin", Password: "pass", OperatorCode: "op1", ClientIP: "127.0.0.1"}); err != nil {
		t.Fatalf("login: %v", err)
	}
	if _, err := d.Login(ctx, LoginReq{Username: "admin", Password: "wrong", OperatorCode: "op1", ClientIP: "10.0.0.1"}); err == nil {
		t.Fatal("expected fail")
	}
	if _, err := d.Login(ctx, LoginReq{Username: "nobody", Password: "x", OperatorCode: "op1", ClientIP: "1.1.1.1"}); err == nil {
		t.Fatal("expected unknown user fail")
	}

	rows, err := d.Client.LoginLog.Query().Order(ent.Asc(loginlog.FieldID)).All(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 3 {
		t.Fatalf("got %d logs", len(rows))
	}
	if rows[0].LoginResult != model.LoginResultSuccess || rows[0].UserID == nil || *rows[0].UserID != u.ID {
		t.Fatalf("success log %+v", rows[0])
	}
	if rows[1].LoginResult != model.LoginResultFail || rows[1].UserID == nil || *rows[1].UserID != u.ID {
		t.Fatalf("wrong password log %+v", rows[1])
	}
	if rows[2].LoginResult != model.LoginResultFail || rows[2].UserID != nil {
		t.Fatalf("unknown user log %+v", rows[2])
	}
	for i, row := range rows {
		if row.OperatorID == nil || *row.OperatorID != op.ID {
			t.Fatalf("log %d operator_id=%v want %d", i, row.OperatorID, op.ID)
		}
	}
}

func TestWriteLoginLogNilClient(t *testing.T) {
	d := &Deps{}
	d.writeLoginLog(context.Background(), LoginReq{Username: "a"}, nil, false, i18n.AuthPasswordIncorrect)
}

func TestListLoginLogsFilterAndTenant(t *testing.T) {
	d := testDeps(t, ModeOn)
	ctx := context.Background()
	op1, err := d.Client.Operator.Create().
		SetOperatorCode("op1").
		SetTimezoneCode("UTC").
		SetSettlementCurrencyCode("USD").
		SetStatus(model.StatusNormal).
		Save(ctx)
	if err != nil {
		t.Fatal(err)
	}
	op2, err := d.Client.Operator.Create().
		SetOperatorCode("op2").
		SetTimezoneCode("UTC").
		SetSettlementCurrencyCode("USD").
		SetStatus(model.StatusNormal).
		Save(ctx)
	if err != nil {
		t.Fatal(err)
	}
	u1 := createUserWithRole(t, d, "alice", "pass", &op1.ID)
	u2 := createUserWithRole(t, d, "bob", "pass", &op2.ID)
	now := time.Now()
	if err := d.Client.LoginLog.Create().SetUsername("alice").SetLoginResult(model.LoginResultSuccess).SetLoginIP("1.1.1.1").SetUserID(u1.ID).SetOperatorID(op1.ID).SetLoginAt(now).Exec(ctx); err != nil {
		t.Fatal(err)
	}
	if err := d.Client.LoginLog.Create().SetUsername("bob").SetLoginResult(model.LoginResultFail).SetLoginIP("2.2.2.2").SetUserID(u2.ID).SetOperatorID(op2.ID).SetFailureReason(i18n.AuthPasswordIncorrect).SetLoginAt(now).Exec(ctx); err != nil {
		t.Fatal(err)
	}
	if err := d.Client.LoginLog.Create().SetUsername("ghost").SetLoginResult(model.LoginResultFail).SetLoginIP("3.3.3.3").SetOperatorID(op1.ID).SetLoginAt(now).Exec(ctx); err != nil {
		t.Fatal(err)
	}

	claims := &ctxdata.Claims{OperatorID: op1.ID}
	list, total, err := d.ListLoginLogs(ctx, claims, LoginLogListReq{})
	if err != nil {
		t.Fatal(err)
	}
	if total != 2 || len(list) != 2 {
		t.Fatalf("tenant list total=%d list=%+v", total, list)
	}
	got := map[string]bool{}
	for _, row := range list {
		got[row.Username] = true
		if row.OperatorID == nil || *row.OperatorID != op1.ID {
			t.Fatalf("tenant row operator_id=%v want %d", row.OperatorID, op1.ID)
		}
	}
	if !got["alice"] || !got["ghost"] || got["bob"] {
		t.Fatalf("tenant names %+v", got)
	}

	d.Mode = ModeOff
	if err := d.Client.LoginLog.Create().SetUsername("bob").SetLoginResult(model.LoginResultFail).SetLoginIP("4.4.4.4").SetLoginAt(now).Exec(ctx); err != nil {
		t.Fatal(err)
	}
	if err := d.Client.LoginLog.Create().SetUsername("plat").SetLoginResult(model.LoginResultFail).SetLoginIP("5.5.5.5").SetLoginAt(now).Exec(ctx); err != nil {
		t.Fatal(err)
	}
	all, total, err := d.ListLoginLogs(ctx, &ctxdata.Claims{}, LoginLogListReq{Username: "bob"})
	if err != nil {
		t.Fatal(err)
	}
	if total != 1 || all[0].Username != "bob" || all[0].OperatorID != nil {
		t.Fatalf("filter username total=%d list=%+v", total, all)
	}
	fails, total, err := d.ListLoginLogs(ctx, &ctxdata.Claims{}, LoginLogListReq{LoginResult: model.LoginResultFail})
	if err != nil {
		t.Fatal(err)
	}
	if total != 2 {
		t.Fatalf("fail filter total=%d list=%+v", total, fails)
	}
}

func TestClip(t *testing.T) {
	if got := clip("abc", 2); got != "ab" {
		t.Fatalf("got %q", got)
	}
	if got := clip("ab", 2); got != "ab" {
		t.Fatalf("got %q", got)
	}
}
