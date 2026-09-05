package service

import (
	"context"
	"testing"
	"time"

	"oa.98ent.com/p9/core/common/ctxdata"
	"oa.98ent.com/p9/core/rpc/model"
)

func TestCreateErrorLogNilClient(t *testing.T) {
	d := &Deps{}
	d.CreateErrorLog(context.Background(), CreateErrorLogReq{RequestMethod: "POST", RequestPath: "/x", ServiceName: "core-api"})
}

func TestCreateErrorLogWritesWithoutUser(t *testing.T) {
	d := testDeps(t, ModeOff)
	d.CreateErrorLog(context.Background(), CreateErrorLogReq{
		RequestMethod:  "post",
		RequestPath:    "/admin/login",
		ServiceName:    "core-api",
		ResponseStatus: 500,
		Subject:        "db down",
		Detail:         "goroutine 1",
		ClientIP:       "127.0.0.1",
	})
	rows, err := d.Client.ErrorLog.Query().All(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 {
		t.Fatalf("got %d", len(rows))
	}
	if rows[0].UserID != nil || rows[0].RequestMethod != "POST" || rows[0].ServiceName != "core-api" {
		t.Fatalf("%+v", rows[0])
	}
	if rows[0].Subject == nil || *rows[0].Subject != "db down" {
		t.Fatalf("subject %+v", rows[0].Subject)
	}
}

func TestCreateErrorLogSkipsMissingUser(t *testing.T) {
	d := testDeps(t, ModeOff)
	d.CreateErrorLog(context.Background(), CreateErrorLogReq{
		UserID:         99999,
		RequestMethod:  "POST",
		RequestPath:    "/admin/user/list",
		ServiceName:    "core-api",
		ResponseStatus: 500,
		ClientIP:       "10.0.0.1",
	})
	rows, err := d.Client.ErrorLog.Query().All(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 || rows[0].UserID != nil {
		t.Fatalf("user_id %+v", rows)
	}
}

func TestListErrorLogsFilterAndTenant(t *testing.T) {
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
	now := time.Now()
	if err := d.Client.ErrorLog.Create().
		SetUserID(u1.ID).SetRequestMethod("POST").SetRequestPath("/admin/user/create").
		SetServiceName("core-api").SetResponseStatus(500).SetClientIP("1.1.1.1").
		SetSubject("boom").SetOperatorID(op1.ID).SetCreatedAt(now).Exec(ctx); err != nil {
		t.Fatal(err)
	}
	if err := d.Client.ErrorLog.Create().
		SetRequestMethod("GRPC").SetRequestPath("/core.Core/login").
		SetServiceName("core-rpc").SetResponseStatus(500).SetClientIP("2.2.2.2").
		SetOperatorID(op2.ID).SetCreatedAt(now).Exec(ctx); err != nil {
		t.Fatal(err)
	}

	list, total, err := d.ListErrorLogs(ctx, &ctxdata.Claims{OperatorID: op1.ID}, ErrorLogListReq{})
	if err != nil {
		t.Fatal(err)
	}
	if total != 1 || len(list) != 1 || list[0].Username != "alice" {
		t.Fatalf("tenant list total=%d list=%+v", total, list)
	}

	d.Mode = ModeOff
	if err := d.Client.ErrorLog.Create().
		SetRequestMethod("POST").SetRequestPath("/admin/role/update").
		SetServiceName("demo-api").SetResponseStatus(503).SetClientIP("3.3.3.3").
		SetCreatedAt(now).Exec(ctx); err != nil {
		t.Fatal(err)
	}
	bySvc, total, err := d.ListErrorLogs(ctx, &ctxdata.Claims{}, ErrorLogListReq{ServiceName: "demo-api"})
	if err != nil {
		t.Fatal(err)
	}
	if total != 1 || bySvc[0].ServiceName != "demo-api" {
		t.Fatalf("service filter total=%d list=%+v", total, bySvc)
	}
}
