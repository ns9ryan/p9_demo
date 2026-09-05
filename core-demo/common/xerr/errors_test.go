package xerr

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"oa.98ent.com/p9/core/common/i18n"
)

func TestAsErrorKeepsDomainError(t *testing.T) {
	orig := BadRequest("invalid json")
	got := AsError(fmt.Errorf("wrap: %w", orig))
	if got.Status != http.StatusBadRequest || got.Message != "invalid json" {
		t.Fatalf("got %+v", got)
	}
}

func TestAsErrorHidesInternal(t *testing.T) {
	raw := errors.New("ERROR: duplicate key value violates unique constraint (SQLSTATE 23505)")
	got := AsError(raw)
	if got.Status != http.StatusInternalServerError || got.Message != i18n.InternalError {
		t.Fatalf("got %+v", got)
	}
	if got.Cause == nil || got.Cause.Error() != raw.Error() {
		t.Fatalf("cause %+v", got.Cause)
	}
	if got.Stack == "" {
		t.Fatal("missing stack")
	}
	if Subject(got) != raw.Error() {
		t.Fatalf("subject %q", Subject(got))
	}
	if StackOf(got) == "" {
		t.Fatal("StackOf empty")
	}
}

func TestRpcErrAttachesInternalDetails(t *testing.T) {
	err := RpcErr(errors.New("db down"))
	st, ok := status.FromError(err)
	if !ok || st.Message() != i18n.InternalError {
		t.Fatalf("%v", err)
	}
	if st.Code() != codes.Internal && int(st.Code()) != http.StatusInternalServerError {
		t.Fatalf("code %v", st.Code())
	}
	sub, stack := FromStatus(st)
	if sub != "db down" || stack == "" {
		t.Fatalf("details sub=%q stack empty=%v", sub, stack == "")
	}
}

func TestRpcErrKeeps4xx(t *testing.T) {
	err := RpcErr(BadRequest("bad"))
	st, _ := status.FromError(err)
	if st.Code() != codes.InvalidArgument {
		t.Fatalf("%v", err)
	}
	sub, stack := FromStatus(st)
	if sub != "" || stack != "" {
		t.Fatalf("4xx should not attach details")
	}
}
