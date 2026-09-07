package i18n

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"oa.98ent.com/p9/core/api/internal/logic/i18n"
	"oa.98ent.com/p9/core/api/internal/svc"
)

func GetEnabledI18nLangsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := i18n.NewGetEnabledI18nLangsLogic(r.Context(), svcCtx)
		resp, err := l.GetEnabledI18nLangs()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
