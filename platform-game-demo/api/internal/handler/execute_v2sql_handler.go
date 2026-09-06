package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"oa.98ent.com/p9/platform-game/api/internal/logic"
	"oa.98ent.com/p9/platform-game/api/internal/svc"
	"oa.98ent.com/p9/platform-game/api/internal/types"
	"oa.98ent.com/p9/platform-game/common/response"
)

func ExecuteV2SqlHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ExecuteV2SqlReq
		// if err := httpx.Parse(r, &req); err != nil {
		// 	httpx.OkJsonCtx(r.Context(), w, response.Error(400, "invalid request"))
		// 	return
		// }

		l := logic.NewExecuteV2SqlLogic(r.Context(), svcCtx)
		resp, err := l.ExecuteV2Sql(&req)
		if err != nil {
			response.ErrorWithStatusCode(w, http.StatusServiceUnavailable, err.Error())
			return
		}
		httpx.OkJsonCtx(r.Context(), w, response.Success(resp))
	}
}
