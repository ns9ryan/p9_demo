package service

import (
	"context"
	"reflect"
	"testing"

	"oa.98ent.com/p9/core/common/i18n"
	"oa.98ent.com/p9/core/rpc/model"
)

func TestRolesOfUsers(t *testing.T) {
	d := testDeps(t, ModeOff)
	ctx := context.Background()

	empty, err := d.RolesOfUsers(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(empty) != 0 {
		t.Fatalf("empty ids=%v", empty)
	}

	on, err := d.Client.Role.Create().
		SetRoleCode("editor").SetRoleName("运营").SetStatus(model.StatusNormal).SetSortNo(2).Save(ctx)
	if err != nil {
		t.Fatal(err)
	}
	sys, err := d.Client.Role.Create().
		SetRoleCode(RoleSuperAdmin).SetRoleName(i18n.RoleSuperAdmin).SetStatus(model.StatusNormal).SetSortNo(1).Save(ctx)
	if err != nil {
		t.Fatal(err)
	}
	off, err := d.Client.Role.Create().
		SetRoleCode("guest").SetRoleName("访客").SetStatus(model.StatusDisabled).SetSortNo(3).Save(ctx)
	if err != nil {
		t.Fatal(err)
	}

	admin, err := d.Client.User.Create().
		SetUserCode(NewUserCode()).SetUsername("admin").SetPasswordHash("x").SetSalt("s").
		SetDisplayName("admin").SetStatus(model.StatusNormal).
		AddRoleIDs(sys.ID, on.ID).Save(ctx)
	if err != nil {
		t.Fatal(err)
	}
	mixed, err := d.Client.User.Create().
		SetUserCode(NewUserCode()).SetUsername("mixed").SetPasswordHash("x").SetSalt("s").
		SetDisplayName("mixed").SetStatus(model.StatusNormal).
		AddRoleIDs(on.ID, off.ID).Save(ctx)
	if err != nil {
		t.Fatal(err)
	}
	none, err := d.Client.User.Create().
		SetUserCode(NewUserCode()).SetUsername("none").SetPasswordHash("x").SetSalt("s").
		SetDisplayName("none").SetStatus(model.StatusNormal).Save(ctx)
	if err != nil {
		t.Fatal(err)
	}

	got, err := d.RolesOfUsers(ctx, []int64{admin.ID, mixed.ID, none.ID})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got[admin.ID].Codes, []string{RoleSuperAdmin, "editor"}) {
		t.Fatalf("admin codes=%v", got[admin.ID].Codes)
	}
	if !reflect.DeepEqual(got[admin.ID].Names, []string{i18n.RoleSuperAdmin, "运营"}) {
		t.Fatalf("admin names=%v", got[admin.ID].Names)
	}
	if !reflect.DeepEqual(got[mixed.ID].Codes, []string{"editor"}) || !reflect.DeepEqual(got[mixed.ID].Names, []string{"运营"}) {
		t.Fatalf("mixed=%+v", got[mixed.ID])
	}
	if len(got[none.ID].Codes) != 0 || len(got[none.ID].Names) != 0 {
		t.Fatalf("none=%+v", got[none.ID])
	}
	if got[none.ID].Codes == nil || got[none.ID].Names == nil {
		t.Fatal("none slices should be empty not nil")
	}
}
