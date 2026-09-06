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
	"oa.98ent.com/p9/platform-base/rpc/client/languageservice"
	"oa.98ent.com/p9/platform-base/rpc/client/pingservice"
	"oa.98ent.com/p9/platform-base/rpc/client/regionservice"
	"oa.98ent.com/p9/platform-base/rpc/client/timezoneservice"
)

// ServiceContext 服务上下文
type ServiceContext struct {
	Config config.Config

	Core coreclient.Core // Core RPC

	PingRpc     pingservice.PingService         // Ping RPC
	LanguageRpc languageservice.LanguageService // 语言 RPC
	TimezoneRpc timezoneservice.TimezoneService // 时区 RPC
	CurrencyRpc currencyservice.CurrencyService // 货币 RPC
	RegionRpc   regionservice.RegionService     // 国家地区 RPC

	Trans     *i18n.Translator // 翻译器
	Language  rest.Middleware  // API语言中间件
	CoreI18n  rest.Middleware  // Core语言中间件
	Jwt       rest.Middleware  // JWT认证中间件
	Authority rest.Middleware  // 权限校验中间件
	ActionLog rest.Middleware  // 操作日志中间件
	ErrorLog  rest.Middleware  // 错误日志中间件
}

// NewServiceContext 创建服务上下文
func NewServiceContext(c config.Config) *ServiceContext {
	// 创建翻译器
	trans, err := i18n.New(c.I18n, locales.FS)
	logx.Must(err)

	// 创建Platform Base RPC客户端，并注册RPC错误拦截器
	platformBaseClient := zrpc.MustNewClient(
		c.PlatformBaseRpc,
		zrpc.WithUnaryClientInterceptor(rpcerror.UnaryClientInterceptor),
	)

	// 创建Core RPC客户端
	coreClient := zrpc.MustNewClient(c.CoreRpc)
	coreCli := coreclient.NewCore(coreClient)

	// 创建Core认证适配器
	auth := coreadapt.Auth(coreCli)

	return &ServiceContext{
		Config: c,

		Core: coreCli,

		PingRpc:     pingservice.NewPingService(platformBaseClient),
		LanguageRpc: languageservice.NewLanguageService(platformBaseClient),
		TimezoneRpc: timezoneservice.NewTimezoneService(platformBaseClient),
		CurrencyRpc: currencyservice.NewCurrencyService(platformBaseClient),
		RegionRpc:   regionservice.NewRegionService(platformBaseClient),

		Trans:    trans,
		Language: apimiddleware.NewLanguageMiddleware().Handle,

		CoreI18n:  coremiddleware.I18n,
		Jwt:       coremiddleware.JWT(auth),
		Authority: coremiddleware.Authority(auth),
		ActionLog: coremiddleware.ActionLog(coreadapt.ActionRecorder(coreCli)),
		ErrorLog:  coremiddleware.ErrorLog(c.Name, coreadapt.ErrorRecorder(coreCli)),
	}
}
