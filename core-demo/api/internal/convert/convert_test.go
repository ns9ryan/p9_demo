package convert

import (
	"context"
	"testing"

	"oa.98ent.com/p9/core/api/internal/types"
	"oa.98ent.com/p9/core/common/i18n"
	"oa.98ent.com/p9/core/rpc/coreclient"
	"oa.98ent.com/p9/core/rpc/model"
)

func TestUpdateMenuReqKeepsZeroFlags(t *testing.T) {
	zero := int32(0)
	parent := int64(0)
	in := &types.UpdateMenuReq{
		Id:       1,
		ParentId: &parent,
		HideMenu: &zero,
		Sort:     &zero,
		Disabled: &zero,
	}
	got := UpdateMenuReq(in)
	if got.HideMenu == nil || *got.HideMenu != 0 {
		t.Fatalf("HideMenu=%v", got.HideMenu)
	}
	if got.Disabled == nil || *got.Disabled != 0 {
		t.Fatalf("Disabled=%v", got.Disabled)
	}
	if got.Sort == nil || *got.Sort != 0 {
		t.Fatalf("Sort=%v", got.Sort)
	}
	if got.ParentId == nil || *got.ParentId != 0 {
		t.Fatalf("ParentId=%v", got.ParentId)
	}
}

func TestUpdateMenuReqOmitsUnset(t *testing.T) {
	got := UpdateMenuReq(&types.UpdateMenuReq{Id: 1})
	if got.HideMenu != nil || got.Disabled != nil || got.Sort != nil || got.ParentId != nil || got.MenuType != nil {
		t.Fatalf("expected omitted fields to stay nil, got %+v", got)
	}
}

func TestUpdateI18nLangReqKeepsZeroFlags(t *testing.T) {
	zero := int32(0)
	got := UpdateI18nLangReq(&types.UpdateI18nLangReq{Id: 1, Disabled: &zero, SortNo: &zero})
	if got.Disabled == nil || *got.Disabled != 0 {
		t.Fatalf("Disabled=%v", got.Disabled)
	}
	if got.SortNo == nil || *got.SortNo != 0 {
		t.Fatalf("SortNo=%v", got.SortNo)
	}
}

func TestUpdateI18nLangReqOmitsUnset(t *testing.T) {
	got := UpdateI18nLangReq(&types.UpdateI18nLangReq{Id: 1})
	if got.Lang != nil || got.Name != nil || got.Disabled != nil || got.SortNo != nil {
		t.Fatalf("expected omitted fields to stay nil, got %+v", got)
	}
}

func TestLoginLogListTranslatesResult(t *testing.T) {
	reason := i18n.AuthPasswordIncorrect
	in := &coreclient.LoginLogListResp{
		Total: 2,
		List: []*coreclient.LoginLogInfo{
			{Id: 1, Username: "ok", LoginResult: int32(model.LoginResultSuccess)},
			{Id: 2, Username: "bad", LoginResult: int32(model.LoginResultFail), FailureReason: &reason},
		},
	}
	zh := LoginLogList(i18n.WithLang(context.Background(), i18n.LangZH), in)
	if zh.List[0].LoginResult != "成功" {
		t.Fatalf("zh success=%q", zh.List[0].LoginResult)
	}
	if zh.List[1].LoginResult != "失败" {
		t.Fatalf("zh fail=%q", zh.List[1].LoginResult)
	}
	if zh.List[1].FailureReason != "用户名或密码错误" {
		t.Fatalf("zh reason=%q", zh.List[1].FailureReason)
	}
	en := LoginLogList(i18n.WithLang(context.Background(), i18n.LangEN), in)
	if en.List[0].LoginResult != "Success" || en.List[1].LoginResult != "Failed" {
		t.Fatalf("en results=%q %q", en.List[0].LoginResult, en.List[1].LoginResult)
	}
}

func TestAdminActionLogListTranslatesResult(t *testing.T) {
	in := &coreclient.AdminActionLogListResp{
		Total: 2,
		List: []*coreclient.AdminActionLogInfo{
			{Id: 1, Username: "a", ActionResult: int32(model.ActionResultSuccess)},
			{Id: 2, Username: "b", ActionResult: int32(model.ActionResultFail)},
		},
	}
	zh := AdminActionLogList(i18n.WithLang(context.Background(), i18n.LangZH), in)
	if zh.List[0].ActionResult != "成功" || zh.List[1].ActionResult != "失败" {
		t.Fatalf("zh %q %q", zh.List[0].ActionResult, zh.List[1].ActionResult)
	}
	en := AdminActionLogList(i18n.WithLang(context.Background(), i18n.LangEN), in)
	if en.List[0].ActionResult != "Success" || en.List[1].ActionResult != "Failed" {
		t.Fatalf("en %q %q", en.List[0].ActionResult, en.List[1].ActionResult)
	}
}

func TestMenuInfoTranslatesTitleFromDict(t *testing.T) {
	i18n.InvalidateAll()
	t.Cleanup(func() {
		i18n.SetDictLoader(nil)
		i18n.InvalidateAll()
	})
	i18n.SetDictLoader(func(_ context.Context, code, group, lang string) (map[string]string, error) {
		if group != i18n.GroupMenu || lang != i18n.LangZH {
			return map[string]string{}, nil
		}
		if code == i18n.CodeOperator {
			return map[string]string{"menu.route.dashboard": "分站工作台"}, nil
		}
		return map[string]string{"menu.route.dashboard": "工作台"}, nil
	})
	got := MenuInfo(i18n.WithLang(context.Background(), i18n.LangZH), i18n.CodePlatform, &coreclient.MenuInfo{
		Id: 1, Title: "menu.route.dashboard", Name: "Dashboard",
	})
	if got.Title != "menu.route.dashboard" {
		t.Fatalf("title=%q", got.Title)
	}
	if got.TransTitle != "工作台" {
		t.Fatalf("transTitle=%q", got.TransTitle)
	}
	op := MenuInfo(i18n.WithLang(context.Background(), i18n.LangZH), i18n.CodeOperator, &coreclient.MenuInfo{
		Id: 1, Title: "menu.route.dashboard", Name: "Dashboard",
	})
	if op.TransTitle != "分站工作台" {
		t.Fatalf("operator transTitle=%q", op.TransTitle)
	}
	miss := MenuInfo(i18n.WithLang(context.Background(), i18n.LangZH), i18n.CodePlatform, &coreclient.MenuInfo{
		Id: 2, Title: "custom", Name: "Custom",
	})
	if miss.TransTitle != "custom" {
		t.Fatalf("passthrough=%q", miss.TransTitle)
	}
}

func TestApiInfoTranslatesDescriptionFromDict(t *testing.T) {
	i18n.InvalidateAll()
	t.Cleanup(func() {
		i18n.SetDictLoader(nil)
		i18n.InvalidateAll()
	})
	i18n.SetDictLoader(func(_ context.Context, code, group, lang string) (map[string]string, error) {
		if group != i18n.GroupAPI || lang != i18n.LangZH {
			return map[string]string{}, nil
		}
		return map[string]string{"api.userCreate": "创建后台用户"}, nil
	})
	got := ApiInfo(i18n.WithLang(context.Background(), i18n.LangZH), i18n.CodePlatform, &coreclient.ApiInfo{
		Id: 1, Description: "api.userCreate", Path: "/admin/user/create", Method: "POST",
	})
	if got.Description != "api.userCreate" {
		t.Fatalf("description=%q", got.Description)
	}
	if got.TransDescription != "创建后台用户" {
		t.Fatalf("transDescription=%q", got.TransDescription)
	}
	miss := ApiInfo(i18n.WithLang(context.Background(), i18n.LangZH), i18n.CodePlatform, &coreclient.ApiInfo{
		Id: 2, Description: "custom", Path: "/admin/custom", Method: "GET",
	})
	if miss.TransDescription != "custom" {
		t.Fatalf("passthrough=%q", miss.TransDescription)
	}
}

func TestUserPublicTranslatesRoleNames(t *testing.T) {
	zh := UserPublic(i18n.WithLang(context.Background(), i18n.LangZH), &coreclient.UserPublic{
		Id: 1, RoleCodes: []string{"super_admin", "editor"},
		RoleNames: []string{i18n.RoleSuperAdmin, "运营"},
	})
	if zh.RoleNames[0] != "超级管理员" {
		t.Fatalf("zh super=%q", zh.RoleNames[0])
	}
	if zh.RoleNames[1] != "运营" {
		t.Fatalf("zh custom=%q", zh.RoleNames[1])
	}
	if !equalStrings(zh.RoleCodes, []string{"super_admin", "editor"}) {
		t.Fatalf("codes=%v", zh.RoleCodes)
	}
	en := UserPublic(i18n.WithLang(context.Background(), i18n.LangEN), &coreclient.UserPublic{
		Id: 1, RoleNames: []string{i18n.RoleSuperAdmin},
	})
	if en.RoleNames[0] != "Super Admin" {
		t.Fatalf("en=%q", en.RoleNames[0])
	}
	empty := UserPublic(context.Background(), &coreclient.UserPublic{Id: 2})
	if empty.RoleCodes == nil || empty.RoleNames == nil {
		t.Fatal("empty slices should not be nil")
	}
}

func equalStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
