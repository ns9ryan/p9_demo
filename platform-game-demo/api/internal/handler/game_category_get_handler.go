// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package handler

import (
	"net/http"
	"strconv"

	"github.com/zeromicro/go-zero/rest/httpx"
	"oa.98ent.com/p9/platform-game/api/internal/logic"
	"oa.98ent.com/p9/platform-game/api/internal/svc"
	"oa.98ent.com/p9/platform-game/api/internal/types"
	"oa.98ent.com/p9/platform-game/common/response"
)

func GameCategoryGetHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GameCategoryGetReq
		query := r.URL.Query()
		idStr := query.Get("id")
		if idStr == "" {
			httpx.OkJsonCtx(r.Context(), w, response.Error(400, "id parameter is required"))
			return
		}
		val, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			httpx.OkJsonCtx(r.Context(), w, response.Error(400, "id parameter is invalid"))
			return
		}
		req.ID = val
		l := logic.NewGameCategoryGetLogic(r.Context(), svcCtx)
		resp, err := l.GameCategoryGet(&req)
		if err != nil {
			response.ErrorWithStatusCode(w, http.StatusInternalServerError, err.Error())
		} else {
			httpx.OkJsonCtx(r.Context(), w, response.Success(resp))
		}
	}
}
