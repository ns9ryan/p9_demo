package service

import (
	"context"
	"testing"
	"time"

	"oa.98ent.com/p9/core/common/ctxdata"
	"oa.98ent.com/p9/core/rpc/ent"
	"oa.98ent.com/p9/core/rpc/ent/adminactionlog"
	"oa.98ent.com/p9/core/rpc/model"
)

func TestCreateAdminActionLogNilClient(t *testing.T) {
	d := &Deps{}
	d.CreateAdminActionLog(context.Background(), CreateAdminActionLogReq{UserID: 1, RequestMethod: "POST", RequestPath: "/x"})
}

func TestListAdminActionLogsFilterAndTenant(t *testing.T) {
	d := testDeps(t, ModeOn)
	ctx := context.Background()
	code1, code2 := "op1", "op2"
	u1 := createUserWithRole(t, d, "alice", "pass", &code1)
	u2 := createUserWithRole(t, d, "bob", "pass", &code2)
	now := time.Now()
	if err := d.Client.AdminActionLog.Create().
		SetUserID(u1.ID).SetRequestMethod("POST").SetRequestPath("/admin/user/create").
		SetActionResult(model.ActionResultSuccess).SetResponseStatus(200).SetClientIP("1.1.1.1").
		SetOperatorCode(code1).SetCreatedAt(now).Exec(ctx); err != nil {
		t.Fatal(err)
	}
	if err := d.Client.AdminActionLog.Create().
		SetUserID(u2.ID).SetRequestMethod("POST").SetRequestPath("/admin/role/update").
		SetActionResult(model.ActionResultFail).SetResponseStatus(400).SetClientIP("2.2.2.2").
		SetOperatorCode(code2).SetCreatedAt(now).Exec(ctx); err != nil {
		t.Fatal(err)
	}

	claims := &ctxdata.Claims{OperatorCode: code1}
	list, total, err := d.ListAdminActionLogs(ctx, claims, AdminActionLogListReq{})
	if err != nil {
		t.Fatal(err)
	}
	if total != 1 || len(list) != 1 || list[0].Username != "alice" {
		t.Fatalf("tenant list total=%d list=%+v", total, list)
	}
	if list[0].OperatorCode == nil || *list[0].OperatorCode != code1 {
		t.Fatalf("tenant operator_code %+v", list[0].OperatorCode)
	}

	d.Mode = ModeOff
	plat := createUserWithRole(t, d, "plat", "pass", nil)
	if err := d.Client.AdminActionLog.Create().
		SetUserID(plat.ID).SetRequestMethod("POST").SetRequestPath("/admin/role/update").
		SetActionResult(model.ActionResultFail).SetResponseStatus(400).SetClientIP("3.3.3.3").
		SetCreatedAt(now).Exec(ctx); err != nil {
		t.Fatal(err)
	}
	if err := d.Client.AdminActionLog.Create().
		SetUserID(plat.ID).SetRequestMethod("POST").SetRequestPath("/admin/user/create").
		SetActionResult(model.ActionResultSuccess).SetResponseStatus(200).SetClientIP("4.4.4.4").
		SetCreatedAt(now).Exec(ctx); err != nil {
		t.Fatal(err)
	}
	fails, total, err := d.ListAdminActionLogs(ctx, &ctxdata.Claims{}, AdminActionLogListReq{ActionResult: model.ActionResultFail})
	if err != nil {
		t.Fatal(err)
	}
	if total != 1 || fails[0].Username != "plat" || fails[0].OperatorCode != nil {
		t.Fatalf("fail filter total=%d list=%+v", total, fails)
	}
	byPath, total, err := d.ListAdminActionLogs(ctx, &ctxdata.Claims{}, AdminActionLogListReq{RequestPath: "/admin/user"})
	if err != nil {
		t.Fatal(err)
	}
	if total != 1 || byPath[0].RequestPath != "/admin/user/create" || byPath[0].OperatorCode != nil {
		t.Fatalf("path filter total=%d list=%+v", total, byPath)
	}
}

func TestCreateAdminActionLogWrites(t *testing.T) {
	d := testDeps(t, ModeOff)
	ctx := context.Background()
	u := createUserWithRole(t, d, "admin", "pass", nil)
	d.CreateAdminActionLog(ctx, CreateAdminActionLogReq{
		UserID:         u.ID,
		RequestMethod:  "post",
		RequestPath:    "/admin/user/update",
		RequestBody:    `{"display_name":"x"}`,
		ActionResult:   model.ActionResultSuccess,
		ResponseStatus: 200,
		ClientIP:       "127.0.0.1",
		UserAgent:      "ua",
	})
	rows, err := d.Client.AdminActionLog.Query().WithUser().Order(ent.Asc(adminactionlog.FieldID)).All(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 {
		t.Fatalf("got %d", len(rows))
	}
	if rows[0].UserID != u.ID || rows[0].RequestMethod != "POST" || rows[0].ActionResult != model.ActionResultSuccess {
		t.Fatalf("%+v", rows[0])
	}
	if rows[0].UserAgent == nil || *rows[0].UserAgent != "ua" {
		t.Fatalf("ua %+v", rows[0].UserAgent)
	}
	if rows[0].OperatorCode != nil {
		t.Fatalf("mode off operator_code=%v", rows[0].OperatorCode)
	}
}

func TestCreateAdminActionLogWritesOperatorCode(t *testing.T) {
	d := testDeps(t, ModeOn)
	ctx := ctxdata.WithClaims(context.Background(), &ctxdata.Claims{OperatorCode: "op1"})
	code := "op1"
	u := createUserWithRole(t, d, "admin", "pass", &code)
	d.CreateAdminActionLog(ctx, CreateAdminActionLogReq{
		UserID:         u.ID,
		RequestMethod:  "POST",
		RequestPath:    "/admin/user/update",
		ActionResult:   model.ActionResultSuccess,
		ResponseStatus: 200,
		ClientIP:       "127.0.0.1",
	})
	rows, err := d.Client.AdminActionLog.Query().All(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 || rows[0].OperatorCode == nil || *rows[0].OperatorCode != "op1" {
		t.Fatalf("operator_code %+v", rows)
	}
}
