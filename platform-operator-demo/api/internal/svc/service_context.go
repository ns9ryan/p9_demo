// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package svc

import (
	"net/http"

	"oa.98ent.com/p9/core/common/coreadapt"
	coremiddleware "oa.98ent.com/p9/core/common/middleware"
	"oa.98ent.com/p9/core/rpc/coreclient"
	"oa.98ent.com/p9/platform-base/rpc/client/currencyservice"
	"oa.98ent.com/p9/platform-base/rpc/client/regionservice"
	"oa.98ent.com/p9/platform-base/rpc/client/timezoneservice"

	"oa.98ent.com/p9/common/i18n"
	game_grpc_client "oa.98ent.com/p9/platform-game/pkg/grpc_client"
	"oa.98ent.com/p9/platform-operator/api/internal/config"
	"oa.98ent.com/p9/platform-operator/api/internal/locales"
	"oa.98ent.com/p9/platform-operator/rpc/client/agentlineallocationservice"
	"oa.98ent.com/p9/platform-operator/rpc/client/basicresourceallocationservice"
	"oa.98ent.com/p9/platform-operator/rpc/client/languageallocationservice"
	"oa.98ent.com/p9/platform-operator/rpc/client/operatoradminservice"
	"oa.98ent.com/p9/platform-operator/rpc/client/operatordomainservice"
	"oa.98ent.com/p9/platform-operator/rpc/client/operatorprofileservice"
	"oa.98ent.com/p9/platform-operator/rpc/client/operatorservice"
	"oa.98ent.com/p9/platform-operator/rpc/client/pingservice"
	"oa.98ent.com/p9/platform-operator/rpc/client/regionallocationservice"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

// ServiceContext 服务上下文
type ServiceContext struct {
	// 服务配置
	Config config.Config // 服务配置

	// Core
	Core coreclient.Core // Core RPC客户端

	// Platform Operator RPC
	PingRpc                    pingservice.PingService                                       // Ping RPC
	OperatorRpc                operatorservice.OperatorService                               // 分站RPC
	OperatorProfileRpc         operatorprofileservice.OperatorProfileService                 // 分站档案RPC
	OperatorDomainRpc          operatordomainservice.OperatorDomainService                   // 分站域名RPC
	OperatorAdminRpc           operatoradminservice.OperatorAdminService                     // 分站管理员RPC
	BasicResourceAllocationRpc basicresourceallocationservice.BasicResourceAllocationService // 基础资源分配RPC
	LanguageAllocationRpc      languageallocationservice.LanguageAllocationService           // 语言分配RPC
	RegionAllocationRpc        regionallocationservice.RegionAllocationService               // 经营地区分配RPC
	AgentLineAllocationRpc     agentlineallocationservice.AgentLineAllocationService         // 代理子线路分配RPC

	// Platform Base RPC
	TimezoneRpc timezoneservice.TimezoneService // 时区RPC
	CurrencyRpc currencyservice.CurrencyService // 货币RPC
	RegionRpc   regionservice.RegionService     // 国家地区RPC

	// Platform Game RPC
	GameGrpcClient *game_grpc_client.GameClientManager // 分站游戏GRPC客户端

	// 多语言
	Trans    *i18n.Translator // API翻译器
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

	// ============================== Platform Operator RPC ==============================

	// 创建Platform Operator RPC客户端，并注册RPC错误拦截器
	platformOperatorClient := zrpc.MustNewClient(
		c.PlatformOperatorRpc,
	)

	// ============================== Platform Base RPC ==============================

	// 创建Platform Base RPC客户端，并注册RPC错误拦截器
	platformBaseClient := zrpc.MustNewClient(
		c.PlatformBaseRpc,
	)

	// ============================== Platform Game RPC ==============================

	// 创建Platform Game RPC连接
	gameGrpcClient, err := game_grpc_client.NewGameClientManager(c.PlatformGameRpc)
	if err != nil {
		logx.Errorf("创建游戏GRPC客户端失败: %v", err)
		gameGrpcClient = nil
	}

	// ============================== Core RPC ==============================

	// 创建Core RPC客户端
	coreClient := zrpc.MustNewClient(c.CoreRpc)
	coreCli := coreclient.NewCore(coreClient)

	// 设置Core多语言词典加载器
	coreadapt.SetDictLoader(coreCli)

	// 创建Core认证适配器
	auth := coreadapt.Auth(coreCli)

	jwt := coremiddleware.JWT(auth)
	authority := coremiddleware.Authority(auth)
	if c.IsDev() {
		// 如果是开发环境，跳过权限校验
		jwt = func(next http.HandlerFunc) http.HandlerFunc {
			return func(w http.ResponseWriter, r *http.Request) {
				next(w, r)
			}
		}
		authority = func(next http.HandlerFunc) http.HandlerFunc {
			return func(w http.ResponseWriter, r *http.Request) {
				next(w, r)
			}
		}
	}

	// ============================== Service Context ==============================

	return &ServiceContext{
		// 服务配置
		Config: c, // 服务配置

		// Core
		Core: coreCli, // Core RPC客户端

		// Platform Operator RPC
		PingRpc:                    pingservice.NewPingService(platformOperatorClient),                                       // Ping RPC
		OperatorRpc:                operatorservice.NewOperatorService(platformOperatorClient),                               // 分站RPC
		OperatorProfileRpc:         operatorprofileservice.NewOperatorProfileService(platformOperatorClient),                 // 分站档案RPC
		OperatorDomainRpc:          operatordomainservice.NewOperatorDomainService(platformOperatorClient),                   // 分站域名RPC
		OperatorAdminRpc:           operatoradminservice.NewOperatorAdminService(platformOperatorClient),                     // 分站管理员RPC
		BasicResourceAllocationRpc: basicresourceallocationservice.NewBasicResourceAllocationService(platformOperatorClient), // 基础资源分配RPC
		LanguageAllocationRpc:      languageallocationservice.NewLanguageAllocationService(platformOperatorClient),           // 语言分配RPC
		RegionAllocationRpc:        regionallocationservice.NewRegionAllocationService(platformOperatorClient),               // 经营地区分配RPC
		AgentLineAllocationRpc:     agentlineallocationservice.NewAgentLineAllocationService(platformOperatorClient),         // 代理子线路分配RPC

		// Platform Base RPC
		TimezoneRpc: timezoneservice.NewTimezoneService(platformBaseClient), // 时区RPC
		CurrencyRpc: currencyservice.NewCurrencyService(platformBaseClient), // 货币RPC
		RegionRpc:   regionservice.NewRegionService(platformBaseClient),     // 国家地区RPC

		// Platform Game RPC
		GameGrpcClient: gameGrpcClient, // 分站游戏GRPC客户端

		// 多语言
		Trans:    trans,               // API翻译器
		CoreI18n: coremiddleware.I18n, // Core多语言中间件

		// 认证权限
		Jwt:       jwt,       // JWT认证中间件
		Authority: authority, // 权限校验中间件

		// 日志
		ActionLog: coremiddleware.ActionLog(coreadapt.ActionRecorder(coreCli)),       // 操作日志中间件
		ErrorLog:  coremiddleware.ErrorLog(c.Name, coreadapt.ErrorRecorder(coreCli)), // 错误日志中间件
	}
}
