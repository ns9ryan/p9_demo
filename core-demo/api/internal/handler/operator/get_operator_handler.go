package operator

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"oa.98ent.com/p9/core/api/internal/logic/operator"
	"oa.98ent.com/p9/core/api/internal/svc"
)

// swagger:route get /admin/operator/self operator GetOperator
//

//

//
// Responses:
//  200: OperatorInfo

func GetOperatorHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := operator.NewGetOperatorLogic(r.Context(), svcCtx)
		resp, err := l.GetOperator()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
