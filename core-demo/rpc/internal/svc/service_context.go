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
	gozeroservice "github.com/zeromicro/go-zero/core/service"
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

	// 创建Ent客户端配置
	entOpts := []ent.Option{
		// ent.Log(logx.Info), // 使用go-zero日志输出SQL
		ent.Driver(drv), // 设置数据库驱动
	}

	// 开发和测试环境开启Ent调试模式
	if c.Mode == gozeroservice.DevMode || c.Mode == gozeroservice.TestMode {
		entOpts = append(entOpts, ent.Debug())
	}
	// 创建Ent客户端
	client := ent.NewClient(entOpts...)
	// 附加租户过滤器
	attachTenant(client)
	// 创建数据库表
	logx.Must(client.Schema.Create(ctx, migrate.WithDropIndex(true), migrate.WithDropColumn(true)))
	// 初始化Redis客户端
	var rds redis.UniversalClient
	if c.DataRedis.Host != "" {
		rds = redis.NewClient(&redis.Options{
			Addr:     c.DataRedis.Host,
			Password: c.DataRedis.Pass,
			DB:       c.DataRedis.DB,
		})
	}
	// 初始化Casbin
	enforcer, err := casbinx.New(client, rds)
	logx.Must(err)

	prefix := c.APIPrefix
	if prefix == "" {
		prefix = "/admin"
	}
	// 初始化deps，封装数据库操作依赖
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

// attachTenant 附加租户过滤器
func attachTenant(c *ent.Client) {
	fOperatorCode := intercept.TraverseFunc(func(ctx context.Context, q intercept.Query) error {
		entmixin.FilterOperatorCode(ctx, q)
		return nil
	})
	fSoftDelete := intercept.TraverseFunc(func(ctx context.Context, q intercept.Query) error {
		entmixin.FilterSoftDelete(ctx, q)
		return nil
	})
	c.User.Intercept(fOperatorCode, fSoftDelete)
	c.Role.Intercept(fOperatorCode, fSoftDelete)
	c.LoginLog.Intercept(fOperatorCode)
	c.AdminActionLog.Intercept(fOperatorCode)
	c.ErrorLog.Intercept(fOperatorCode)
}
