package user

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"oa.98ent.com/p9/core/api/internal/logic/user"
	"oa.98ent.com/p9/core/api/internal/svc"
	"oa.98ent.com/p9/core/api/internal/types"
)

// swagger:route post /admin/user/ipWhitelist user UpdateUserIpWhitelist
//
//
//
// Parameters:
//  + name: body
//    require: true
//    in: body
//    type: UpdateUserIpWhitelistReq
//

func UpdateUserIpWhitelistHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UpdateUserIpWhitelistReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := user.NewUpdateUserIpWhitelistLogic(r.Context(), svcCtx)
		resp, err := l.UpdateUserIpWhitelist(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
