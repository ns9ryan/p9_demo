// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package svc

import (
	"github.com/zeromicro/go-zero/rest"
	"oa.98ent.com/p9/platform-operator/api/internal/config"
	"oa.98ent.com/p9/platform-operator/api/internal/middleware"
)

type ServiceContext struct {
	Config    config.Config
	Jwt       rest.Middleware
	ActionLog rest.Middleware
	Authority rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:    c,
		Jwt:       middleware.NewJwtMiddleware().Handle,
		ActionLog: middleware.NewActionLogMiddleware().Handle,
		Authority: middleware.NewAuthorityMiddleware().Handle,
	}
}
