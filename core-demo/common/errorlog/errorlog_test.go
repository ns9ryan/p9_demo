package errorlog

import (
	"context"
	"testing"
	"time"

	"oa.98ent.com/p9/core/common/ctxdata"
)

type chanRecorder struct {
	ch chan Record
}

func (c chanRecorder) RecordError(_ context.Context, rec Record) {
	c.ch <- rec
}

func TestNoteAndBag(t *testing.T) {
	ctx, b := WithBag(context.Background())
	Note(ctx, "boom", "stack")
	if b.Subject != "boom" || b.Detail != "stack" {
		t.Fatalf("%+v", b)
	}
	if BagFromCtx(context.Background()) != nil {
		t.Fatal("empty ctx")
	}
}

func TestNoteDoesNotOverwrite(t *testing.T) {
	ctx, b := WithBag(context.Background())
	Note(ctx, "panic", "stack1")
	Note(ctx, "internal", "stack2")
	if b.Subject != "panic" || b.Detail != "stack1" {
		t.Fatalf("%+v", b)
	}
}

func TestShouldCollect(t *testing.T) {
	if ShouldCollect(400) || ShouldCollect(499) {
		t.Fatal("4xx")
	}
	if !ShouldCollect(500) || !ShouldCollect(503) {
		t.Fatal("5xx")
	}
}

func TestReportKeepsClaims(t *testing.T) {
	ch := make(chan Record, 1)
	ctx := ctxdata.WithClaims(context.Background(), &ctxdata.Claims{UserID: 7})
	Report(ctx, chanRecorder{ch: ch}, Record{Subject: "x", ServiceName: "core-api"})
	select {
	case rec := <-ch:
		if rec.Subject != "x" || rec.ServiceName != "core-api" {
			t.Fatalf("%+v", rec)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timeout")
	}
}
