package svc

import (
	"oa.98ent.com/p9/core/api/internal/config"
	"oa.98ent.com/p9/core/common/coreadapt"
	"oa.98ent.com/p9/core/common/middleware"
	"oa.98ent.com/p9/core/rpc/coreclient"

	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config    config.Config
	Core      coreclient.Core
	Authority rest.Middleware
	Jwt       rest.Middleware
	ActionLog rest.Middleware
	ErrorLog  rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {
	cli := zrpc.MustNewClient(c.CoreRpc)
	coreCli := coreclient.NewCore(cli)
	auth := coreadapt.Auth(coreCli)
	return &ServiceContext{
		Config:    c,
		Core:      coreCli,
		Jwt:       middleware.JWT(auth),
		Authority: middleware.Authority(auth),
		ActionLog: middleware.ActionLog(coreadapt.ActionRecorder(coreCli)),
		ErrorLog:  middleware.ErrorLog(c.Name, coreadapt.ErrorRecorder(coreCli)),
	}
}
