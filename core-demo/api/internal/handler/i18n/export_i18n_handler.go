package i18n

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"oa.98ent.com/p9/core/api/internal/convert"
	"oa.98ent.com/p9/core/api/internal/logic/i18n"
	"oa.98ent.com/p9/core/api/internal/svc"
	"oa.98ent.com/p9/core/api/internal/types"
)

func ExportI18nHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ExportI18nReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := i18n.NewExportI18nLogic(r.Context(), svcCtx)
		resp, err := l.ExportI18n(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		body, err := convert.MarshalI18nFile(resp)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("Content-Disposition", `attachment; filename="`+convert.I18nExportFilename(resp.Lang)+`"`)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(body)
	}
}
