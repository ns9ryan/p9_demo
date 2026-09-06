// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"oa.98ent.com/p9/platform-game/api/internal/logic"
	"oa.98ent.com/p9/platform-game/api/internal/svc"
	"oa.98ent.com/p9/platform-game/api/internal/types"
	"oa.98ent.com/p9/platform-game/common/response"
)

func GameUpdateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GameUpdateReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewGameUpdateLogic(r.Context(), svcCtx)
		resp, err := l.GameUpdate(&req)
		if err != nil {
			response.ErrorWithStatusCode(w, http.StatusInternalServerError, err.Error())
		} else {
			httpx.OkJsonCtx(r.Context(), w, response.Success(resp))
		}
	}
}
