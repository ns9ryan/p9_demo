package i18n

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"oa.98ent.com/p9/core/api/internal/logic/i18n"
	"oa.98ent.com/p9/core/api/internal/svc"
	"oa.98ent.com/p9/core/api/internal/types"
)

func GetI18nDictHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetI18nDictReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := i18n.NewGetI18nDictLogic(r.Context(), svcCtx)
		resp, err := l.GetI18nDict(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		items := map[string]string{}
		if resp != nil && resp.Items != nil {
			items = resp.Items
		}
		httpx.OkJsonCtx(r.Context(), w, items)
	}
}
