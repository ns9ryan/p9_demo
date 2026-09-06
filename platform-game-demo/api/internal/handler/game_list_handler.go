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

func GameListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GameListReq
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
		req.GameCode = query.Get("game_code")
		req.Name = query.Get("name")
		if categoryIDStr := query.Get("category_id"); categoryIDStr != "" {
			if val, err := strconv.ParseInt(categoryIDStr, 10, 64); err == nil {
				req.CategoryID = val
			}
		}
		if providerIDStr := query.Get("provider_id"); providerIDStr != "" {
			if val, err := strconv.ParseInt(providerIDStr, 10, 64); err == nil {
				req.ProviderID = val
			}
		}
		if channelIDStr := query.Get("channel_id"); channelIDStr != "" {
			if val, err := strconv.ParseInt(channelIDStr, 10, 64); err == nil {
				req.ChannelID = val
			}
		}
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

		l := logic.NewGameListLogic(r.Context(), svcCtx)
		resp, err := l.GameList(&req)
		if err != nil {
			response.ErrorWithStatusCode(w, http.StatusInternalServerError, err.Error())
		} else {
			httpx.OkJsonCtx(r.Context(), w, response.Success(response.ListResponseData(resp.Items, resp.Total)))
		}
	}
}
