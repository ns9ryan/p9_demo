// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/zeromicro/go-zero/rest/httpx"
	"oa.98ent.com/p9/platform-game/api/internal/logic"
	"oa.98ent.com/p9/platform-game/api/internal/svc"
	"oa.98ent.com/p9/platform-game/api/internal/types"
	"oa.98ent.com/p9/platform-game/common/response"
)

func GameSyncCheckpointListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GameSyncCheckpointListReq
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
		req.SyncScope = query.Get("sync_scope")
		if startTimeStr := query.Get("start_time"); startTimeStr != "" {
			if val, err := strconv.ParseInt(startTimeStr, 10, 64); err == nil {
				req.StartTime = val
			}
		}
		if endTimeStr := query.Get("end_time"); endTimeStr != "" {
			if val, err := strconv.ParseInt(endTimeStr, 10, 64); err == nil {
				req.EndTime = val
			}
		}
		log.Printf("[API GameSyncCheckpointListHandler] parsed query: page=%d, page_size=%d, sync_scope=%q, start_time=%d, end_time=%d", req.Page, req.PageSize, req.SyncScope, req.StartTime, req.EndTime)

		l := logic.NewGameSyncCheckpointListLogic(r.Context(), svcCtx)
		resp, err := l.GameSyncCheckpointList(&req)
		if err != nil {
			response.ErrorWithStatusCode(w, http.StatusInternalServerError, err.Error())
		} else {
			httpx.OkJsonCtx(r.Context(), w, response.Success(response.ListResponseData(resp.Items, resp.Total)))
		}
	}
}
