package svc

import (
	"context"
	"fmt"
	"strings"

	"oa.98ent.com/p9/core/common/entdb"
	"oa.98ent.com/p9/core/common/entmixin"
	"oa.98ent.com/p9/core/rpc/casbinx"
	"oa.98ent.com/p9/core/rpc/ent"
	"oa.98ent.com/p9/core/rpc/ent/intercept"
	"oa.98ent.com/p9/core/rpc/ent/migrate"
	_ "oa.98ent.com/p9/core/rpc/ent/runtime"
	"oa.98ent.com/p9/core/rpc/internal/config"
	"oa.98ent.com/p9/core/rpc/service"

	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
)

type ServiceContext struct {
	Config config.Config
	Deps   *service.Deps
}

func NewServiceContext(c config.Config) *ServiceContext {
	mode := strings.ToLower(strings.TrimSpace(c.PartnerMode))
	if mode != service.ModeOff && mode != service.ModeOn {
		logx.Must(fmt.Errorf("PartnerMode must be off or on"))
	}
	ctx := context.Background()
	drv, err := entdb.Open(c.DB.Driver, c.DB.DSN)
	logx.Must(err)
	client := ent.NewClient(ent.Driver(drv), ent.Debug())
	attachTenant(client)
	logx.Must(client.Schema.Create(ctx, migrate.WithDropIndex(true), migrate.WithDropColumn(true)))

	var rds redis.UniversalClient
	if c.DataRedis.Host != "" {
		rds = redis.NewClient(&redis.Options{
			Addr:     c.DataRedis.Host,
			Password: c.DataRedis.Pass,
			DB:       c.DataRedis.DB,
		})
	}
	enforcer, err := casbinx.New(client, rds)
	logx.Must(err)

	prefix := c.APIPrefix
	if prefix == "" {
		prefix = "/admin"
	}
	deps := &service.Deps{
		Client:           client,
		Enforcer:         enforcer,
		Redis:            rds,
		Mode:             mode,
		JWTSecret:        c.Jwt.AccessSecret,
		JWTExpire:        c.Jwt.AccessExpire,
		JWTRefreshSecret: c.Jwt.RefreshSecret,
		JWTRefreshExpire: c.Jwt.RefreshExpire,
		APIPrefix:        prefix,
		InitToken:        c.InitToken,
	}
	return &ServiceContext{Config: c, Deps: deps}
}

func attachTenant(c *ent.Client) {
	fOperatorID := intercept.TraverseFunc(func(ctx context.Context, q intercept.Query) error {
		entmixin.FilterOperatorID(ctx, q)
		return nil
	})
	fSoftDelete := intercept.TraverseFunc(func(ctx context.Context, q intercept.Query) error {
		entmixin.FilterSoftDelete(ctx, q)
		return nil
	})
	c.User.Intercept(fOperatorID, fSoftDelete)
	c.Role.Intercept(fOperatorID, fSoftDelete)
	c.LoginLog.Intercept(fOperatorID)
	c.AdminActionLog.Intercept(fOperatorID)
	c.ErrorLog.Intercept(fOperatorID)
}
