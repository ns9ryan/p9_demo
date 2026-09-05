package svc

import (
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	"oa.98ent.com/p9/platform-base/rpc/ent"
	_ "oa.98ent.com/p9/platform-base/rpc/ent/runtime"
	"oa.98ent.com/p9/platform-base/rpc/internal/config"
)

type ServiceContext struct {
	Config config.Config
	DB     *ent.Client // Ent数据库客户端
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 创建数据库驱动
	driver, err := c.DatabaseConf.NewDriver()
	logx.Must(err)

	// 创建Ent客户端配置
	entOpts := []ent.Option{
		ent.Log(logx.Info), // 使用go-zero日志输出SQL
		ent.Driver(driver), // 设置数据库驱动
	}

	// 开发和测试环境开启Ent调试模式
	if c.Mode == service.DevMode || c.Mode == service.TestMode {
		entOpts = append(entOpts, ent.Debug())
	}

	// 创建Ent数据库客户端
	db := ent.NewClient(entOpts...)

	return &ServiceContext{
		Config: c,
		DB:     db,
	}
}
