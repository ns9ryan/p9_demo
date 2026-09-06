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

func GameProviderListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GameProviderListReq
		query := r.URL.Query()
		if pageStr := query.Get("page"); pageStr != "" {
			if val, err := strconv.ParseInt(pageStr, 10, 64); err == nil {
				req.Page = val
			}
		}
		if pageSizeStr := query.Get("page_size"); pageSizeStr != "" {
			if val, err := strconv.ParseInt(pageSizeStr, 10, 64); err == nil {
				req.PageSize = val
			}
		}
		req.ProviderCode = query.Get("provider_code")
		req.Name = query.Get("name")
		if statusStr := query.Get("status"); statusStr != "" {
			if val, err := strconv.ParseInt(statusStr, 10, 16); err == nil {
				req.Status = int16(val)
			}
		}
		if isDeletedStr := query.Get("is_deleted"); isDeletedStr != "" {
			if val, err := strconv.ParseInt(isDeletedStr, 10, 16); err == nil {
				req.IsDeleted = int16(val)
			}
		}
		req.SortBy = query.Get("sort_by")
		req.SortOrder = query.Get("sort_order")

		l := logic.NewGameProviderListLogic(r.Context(), svcCtx)
		resp, err := l.GameProviderList(&req)
		if err != nil {
			response.ErrorWithStatusCode(w, http.StatusInternalServerError, err.Error())
		} else {
			httpx.OkJsonCtx(r.Context(), w, response.Success(response.ListResponseData(resp.Items, resp.Total)))
		}
	}
}
