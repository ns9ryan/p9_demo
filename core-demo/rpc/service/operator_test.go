package service

import (
	"context"
	"testing"

	"oa.98ent.com/p9/core/common/xerr"
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
