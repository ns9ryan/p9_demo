// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package operator_game_provider

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"oa.98ent.com/p9/platform-operator/api/internal/logic/operator_game_provider"
	"oa.98ent.com/p9/platform-operator/api/internal/svc"
	"oa.98ent.com/p9/platform-operator/api/internal/types"
)

// 批量删除分站游戏提供商
func BatchDeleteOperatorGameProviderHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.BatchDeleteOperatorGameProviderRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := operator_game_provider.NewBatchDeleteOperatorGameProviderLogic(r.Context(), svcCtx)
		resp, err := l.BatchDeleteOperatorGameProvider(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
