// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package timezone

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"oa.98ent.com/p9/platform-base/api/internal/logic/timezone"
	"oa.98ent.com/p9/platform-base/api/internal/svc"
	"oa.98ent.com/p9/platform-base/api/internal/types"
)

func ListAllTimezonesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ListAllTimezonesRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := timezone.NewListAllTimezonesLogic(r.Context(), svcCtx)
		resp, err := l.ListAllTimezones(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
