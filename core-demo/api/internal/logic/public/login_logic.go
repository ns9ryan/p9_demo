package public

import (
	"context"
	"net/http"

	"oa.98ent.com/p9/core/api/internal/convert"
	"oa.98ent.com/p9/core/api/internal/svc"
	"oa.98ent.com/p9/core/api/internal/types"
	"oa.98ent.com/p9/core/common/ctxdata"
	"oa.98ent.com/p9/core/common/utils"
	"oa.98ent.com/p9/core/rpc/coreclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	r      *http.Request
}

func NewLoginLogic(r *http.Request, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(r.Context()),
		ctx:    r.Context(),
		svcCtx: svcCtx,
		r:      r,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginResp, err error) {
	ctx := l.ctx
	ip := ctxdata.ClientIPFromCtx(ctx)
	if ip == "" {
		ip = utils.ClientIP(l.r)
		ctx = ctxdata.WithClientIP(ctx, ip)
	}
	out, err := l.svcCtx.Core.Login(ctx, &coreclient.LoginReq{
		Username:     req.Username,
		Password:     req.Password,
		OperatorCode: req.OperatorCode,
		ClientIp:     ip,
		UserAgent:    utils.UserAgent(l.r),
	})
	if err != nil {
		return nil, err
	}
	return convert.LoginResp(ctx, out), nil
}
