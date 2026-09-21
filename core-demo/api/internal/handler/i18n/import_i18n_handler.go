package i18n

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"oa.98ent.com/p9/core/api/internal/convert"
	"oa.98ent.com/p9/core/api/internal/logic/i18n"
	"oa.98ent.com/p9/core/api/internal/svc"
)

func ImportI18nHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, convert.I18nFileMaxBytes)
		req, err := convert.ReadI18nImportReq(r)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := i18n.NewImportI18nLogic(r.Context(), svcCtx)
		resp, err := l.ImportI18n(req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
