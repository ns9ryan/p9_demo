package ctxdata

import (
	"context"
	"testing"

	"google.golang.org/grpc/metadata"
)

func TestWithClaimsWritesOutgoingAndValue(t *testing.T) {
	in := &Claims{
		UserID: 7, UserCode: "u7", Username: "alice",
		OperatorID: 3, OperatorCode: "A",
		RoleCodes: []string{"admin", "ops"},
		Salt:      "s", ExpiresAt: 99,
		IsPlatform: true, TokenType: "preview",
	}
	ctx := WithClaims(context.Background(), in)
	got := ClaimsFromCtx(ctx)
	if got == nil || got.Username != "alice" || got.OperatorCode != "A" || !got.IsPlatform || got.TokenType != "preview" {
		t.Fatalf("value claims: %+v", got)
	}
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok || first(md, headerUsername) != "alice" || first(md, headerOperatorCode) != "A" || first(md, headerIsPlatform) != "1" || first(md, headerTokenType) != "preview" {
		t.Fatalf("outgoing md: %v", md)
	}
}

func TestClaimsFromIncomingMetadata(t *testing.T) {
	md := metadata.Pairs(
		headerUserID, "7",
		headerUserCode, "u7",
		headerUsername, "alice",
		headerOperatorID, "3",
		headerOperatorCode, "A",
		headerRoleCodes, "admin,ops",
		headerSalt, "s",
		headerExpiresAt, "99",
	)
	ctx := metadata.NewIncomingContext(context.Background(), md)
	got := ClaimsFromCtx(ctx)
	if got == nil {
		t.Fatal("nil claims")
	}
	if got.UserID != 7 || got.Username != "alice" || got.OperatorCode != "A" {
		t.Fatalf("%+v", got)
	}
	if len(got.RoleCodes) != 2 || got.RoleCodes[0] != "admin" {
		t.Fatalf("roles %+v", got.RoleCodes)
	}
}

func TestRawTokenRoundTrip(t *testing.T) {
	ctx := WithRawToken(context.Background(), "tok")
	if RawTokenFromCtx(ctx) != "tok" {
		t.Fatal("value")
	}
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok || first(md, headerRawToken) != "tok" {
		t.Fatalf("outgoing %v", md)
	}
	in := metadata.NewIncomingContext(context.Background(), metadata.Pairs(headerRawToken, "tok2"))
	if RawTokenFromCtx(in) != "tok2" {
		t.Fatal("incoming")
	}
}

func TestClaimsHolderSeesChildWithClaims(t *testing.T) {
	parent := WithClaimsHolder(context.Background())
	child := WithClaims(parent, &Claims{UserID: 9, OperatorID: 3, Username: "alice"})
	got := ClaimsFromCtx(parent)
	if got == nil || got.UserID != 9 || got.OperatorID != 3 || got.Username != "alice" {
		t.Fatalf("parent %+v", got)
	}
	if ClaimsFromCtx(child) != got {
		t.Fatal("child should see same claims")
	}
}

func TestWithClaimsWithoutHolderDoesNotLeak(t *testing.T) {
	parent := context.Background()
	_ = WithClaims(parent, &Claims{UserID: 9, OperatorID: 3})
	if ClaimsFromCtx(parent) != nil {
		t.Fatal("parent should stay empty")
	}
}

func TestWithClaimsNil(t *testing.T) {
	ctx := WithClaims(context.Background(), nil)
	if ClaimsFromCtx(ctx) != nil {
		t.Fatal("expected nil")
	}
}

func TestOperatorFromCtx(t *testing.T) {
	if OperatorIDFromCtx(context.Background()) != 0 || OperatorCodeFromCtx(context.Background()) != "" {
		t.Fatal("empty ctx")
	}
	ctx := WithClaims(context.Background(), &Claims{UserID: 1, OperatorID: 3, OperatorCode: "A"})
	if OperatorIDFromCtx(ctx) != 3 || OperatorCodeFromCtx(ctx) != "A" {
		t.Fatal("value")
	}
	in := metadata.NewIncomingContext(context.Background(), metadata.Pairs(
		headerUserID, "1", headerOperatorID, "9", headerOperatorCode, "B",
	))
	if OperatorIDFromCtx(in) != 9 || OperatorCodeFromCtx(in) != "B" {
		t.Fatal("incoming")
	}
}

func TestClaimsFromIncomingPreviewPlatform(t *testing.T) {
	md := metadata.Pairs(
		headerOperatorID, "12",
		headerOperatorCode, "demo",
		headerRoleCodes, "super_admin",
		headerIsPlatform, "1",
		headerTokenType, "preview",
	)
	got := ClaimsFromCtx(metadata.NewIncomingContext(context.Background(), md))
	if got == nil {
		t.Fatal("nil claims")
	}
	if got.UserID != 0 || got.OperatorID != 12 || !got.IsPlatform || got.TokenType != "preview" {
		t.Fatalf("%+v", got)
	}
}
