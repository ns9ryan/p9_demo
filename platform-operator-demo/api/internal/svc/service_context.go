// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package svc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"oa.98ent.com/p9/platform-operator/api/internal/config"
	"oa.98ent.com/p9/platform-operator/rpc/client/pingservice"
)

type ServiceContext struct {
	Config      config.Config
	PingService pingservice.PingService
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:      c,
		PingService: pingservice.NewPingService(zrpc.MustNewClient(c.PlatformOperatorRpc)),
	}
}
