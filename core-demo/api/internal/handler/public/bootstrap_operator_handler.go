package public

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"oa.98ent.com/p9/core/api/internal/logic/public"
	"oa.98ent.com/p9/core/api/internal/svc"
	"oa.98ent.com/p9/core/api/internal/types"
)

// swagger:route post /admin/bootstrap/operator public BootstrapOperator
//

//

//
// Parameters:
//  + name: body
//    require: true
//    in: body
//    type: BootstrapOperatorReq
//
// Responses:
//  200: UserPublic

func BootstrapOperatorHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.BootstrapOperatorReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := public.NewBootstrapOperatorLogic(r.Context(), svcCtx)
		resp, err := l.BootstrapOperator(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
