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
