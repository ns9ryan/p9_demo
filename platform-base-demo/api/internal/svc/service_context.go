// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package svc

import (
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"

	"oa.98ent.com/p9/core/common/coreadapt"
	coremiddleware "oa.98ent.com/p9/core/common/middleware"
	"oa.98ent.com/p9/core/rpc/coreclient"
	"oa.98ent.com/p9/platform-base/api/internal/config"
	"oa.98ent.com/p9/platform-base/api/internal/locales"
	apimiddleware "oa.98ent.com/p9/platform-base/pkg/api/middleware"
	"oa.98ent.com/p9/platform-base/pkg/api/rpcerror"
	"oa.98ent.com/p9/platform-base/pkg/i18n"
	"oa.98ent.com/p9/platform-base/rpc/client/currencyservice"
	"oa.98ent.com/p9/platform-base/rpc/client/pingservice"
	"oa.98ent.com/p9/platform-base/rpc/client/regionservice"
	"oa.98ent.com/p9/platform-base/rpc/client/timezoneservice"
)

// ServiceContext 服务上下文
type ServiceContext struct {
	// 服务配置
	Config config.Config // 服务配置

	// Core
	Core coreclient.Core // Core RPC客户端

	// Platform Base RPC
	PingRpc     pingservice.PingService         // Ping RPC
	TimezoneRpc timezoneservice.TimezoneService // 时区RPC
	CurrencyRpc currencyservice.CurrencyService // 货币RPC
	RegionRpc   regionservice.RegionService     // 国家地区RPC

	// 多语言
	Trans    *i18n.Translator // API翻译器
	Language rest.Middleware  // API语言中间件
	CoreI18n rest.Middleware  // Core多语言中间件

	// 认证权限
	Jwt       rest.Middleware // JWT认证中间件
	Authority rest.Middleware // 权限校验中间件

	// 日志
	ActionLog rest.Middleware // 操作日志中间件
	ErrorLog  rest.Middleware // 错误日志中间件
}

// NewServiceContext 创建服务上下文
func NewServiceContext(c config.Config) *ServiceContext {
	// 创建API翻译器
	trans, err := i18n.New(c.I18n, locales.FS)
	logx.Must(err)

	// ============================== Platform Base RPC ==============================

	// 创建Platform Base RPC客户端，并注册RPC错误拦截器
	platformBaseClient := zrpc.MustNewClient(
		c.PlatformBaseRpc,
		zrpc.WithUnaryClientInterceptor(rpcerror.UnaryClientInterceptor),
	)

	// ============================== Core RPC ==============================

	// 创建Core RPC客户端
	coreClient := zrpc.MustNewClient(c.CoreRpc)
	coreCli := coreclient.NewCore(coreClient)

	// 设置Core多语言词典加载器
	coreadapt.SetDictLoader(coreCli)

	// 创建Core认证适配器
	auth := coreadapt.Auth(coreCli)

	// ============================== Service Context ==============================

	return &ServiceContext{
		// 服务配置
		Config: c, // 服务配置

		// Core
		Core: coreCli, // Core RPC客户端

		// Platform Base RPC
		PingRpc:     pingservice.NewPingService(platformBaseClient),         // Ping RPC
		TimezoneRpc: timezoneservice.NewTimezoneService(platformBaseClient), // 时区RPC
		CurrencyRpc: currencyservice.NewCurrencyService(platformBaseClient), // 货币RPC
		RegionRpc:   regionservice.NewRegionService(platformBaseClient),     // 国家地区RPC

		// 多语言
		Trans:    trans,                                        // API翻译器
		Language: apimiddleware.NewLanguageMiddleware().Handle, // API语言中间件
		CoreI18n: coremiddleware.I18n,                          // Core多语言中间件

		// 认证权限
		Jwt:       coremiddleware.JWT(auth),       // JWT认证中间件
		Authority: coremiddleware.Authority(auth), // 权限校验中间件

		// 日志
		ActionLog: coremiddleware.ActionLog(coreadapt.ActionRecorder(coreCli)),       // 操作日志中间件
		ErrorLog:  coremiddleware.ErrorLog(c.Name, coreadapt.ErrorRecorder(coreCli)), // 错误日志中间件
	}
}
