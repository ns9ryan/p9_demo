// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package operator_game_channel

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"oa.98ent.com/p9/platform-operator/api/internal/logic/operator_game_channel"
	"oa.98ent.com/p9/platform-operator/api/internal/svc"
	"oa.98ent.com/p9/platform-operator/api/internal/types"
)

// 批量修改分站游戏渠道状态
func BatchUpdateOperatorGameChannelStatusHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.BatchUpdateOperatorGameChannelStatusRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := operator_game_channel.NewBatchUpdateOperatorGameChannelStatusLogic(r.Context(), svcCtx)
		resp, err := l.BatchUpdateOperatorGameChannelStatus(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
