// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package execute

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"oa.98ent.com/p9/platform-game/api/internal/logic/execute"
	"oa.98ent.com/p9/platform-game/api/internal/svc"
)

func ExecuteV2sqlGetHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := execute.NewExecuteV2sqlGetLogic(r.Context(), svcCtx)
		resp, err := l.ExecuteV2sqlGet()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
