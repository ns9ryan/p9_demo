package user

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"oa.98ent.com/p9/core/api/internal/logic/user"
	"oa.98ent.com/p9/core/api/internal/svc"
)

// swagger:route post /admin/logout/all user LogoutAll
//

//

//

func LogoutAllHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := user.NewLogoutAllLogic(r.Context(), svcCtx)
		resp, err := l.LogoutAll()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
