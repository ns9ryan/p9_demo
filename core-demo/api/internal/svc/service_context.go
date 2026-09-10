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
	Core      coreclient.Core // Core RPC客户端
	Authority rest.Middleware // 权限中间件
	Jwt       rest.Middleware // JWT中间件
	ActionLog rest.Middleware // 操作日志中间件
	ErrorLog  rest.Middleware // 错误日志中间件
}

func NewServiceContext(c config.Config) *ServiceContext {
	cli := zrpc.MustNewClient(c.CoreRpc)
	coreCli := coreclient.NewCore(cli)
	// 权限中间件适配器
	auth := coreadapt.Auth(coreCli)
	// 设置多语言字典加载器
	coreadapt.SetDictLoader(coreCli)

	return &ServiceContext{
		Config:    c,
		Core:      coreCli,
		Jwt:       middleware.JWT(auth),
		Authority: middleware.Authority(auth),
		ActionLog: middleware.ActionLog(coreadapt.ActionRecorder(coreCli)),
		ErrorLog:  middleware.ErrorLog(c.Name, coreadapt.ErrorRecorder(coreCli)),
	}
}
