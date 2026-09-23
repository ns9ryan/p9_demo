// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package svc

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
	"oa.98ent.com/p9/common/i18n"
	"oa.98ent.com/p9/core/common/coreadapt"
	coremiddleware "oa.98ent.com/p9/core/common/middleware"
	"oa.98ent.com/p9/core/rpc/coreclient"
	"oa.98ent.com/p9/platform-game/api/internal/config"
	"oa.98ent.com/p9/platform-game/api/internal/locales"
	"oa.98ent.com/p9/platform-game/pkg/grpc_client"
)

type ServiceContext struct {
	Config          config.Config
	GameGrpcClient  *grpc_client.GameClientManager
	Core            coreclient.Core  // Core RPC客户端
	SyncRateLimiter *SyncRateLimiter // 同步操作频率限制器

	// 多语言
	Trans    *i18n.Translator // API翻译器
	Language rest.Middleware  // API语言中间件

	Jwt       rest.Middleware // JWT认证中间件
	Authority rest.Middleware // 权限校验中间件

	// 日志
	ActionLog rest.Middleware // 操作日志中间件
	ErrorLog  rest.Middleware // 错误日志中间件
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 创建翻译器
	trans, err := i18n.New(c.I18n, locales.FS)
	logx.Must(err)

	// 初始化 gRPC 客户端
	logx.Info("[ServiceContext] 开始初始化 gRPC 客户端...")
	gameGrpcClient, err := grpc_client.NewGameClientManager(c.PlatformGameRpc)

	if err != nil {
		logx.Errorf("[ServiceContext] ✗ 游戏 gRPC 客户端初始化失败: %v", err)
		panic("failed to initialize 游戏 gRPC 客户端: " + err.Error())
	}
	var coreAuth coremiddleware.Client
	var coreCli coreclient.Core

	if c.CoreRpc.Target == "" {
		logx.Error("[ServiceContext] ⚠ CoreRpc.Target 未配置，鉴权功能不可用")
	} else {
		ctx := context.Background()
		if c.CoreRpc.Timeout > 0 {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, time.Duration(c.CoreRpc.Timeout)*time.Second)
			defer cancel()
		}

		// 创建 Core 客户端并初始化鉴权
		// coreCli = coreclient.NewCore(zrpcClient)
		coreClient := zrpc.MustNewClient(zrpc.RpcClientConf{
			Target:  c.CoreRpc.Target,
			Timeout: int64(c.CoreRpc.Timeout),
		})
		coreCli = coreclient.NewCore(coreClient)

		// 设置Core多语言词典加载器
		coreadapt.SetDictLoader(coreCli)

		// 创建Core认证适配器
		coreAuth = coreadapt.Auth(coreCli)

		logx.Info("[ServiceContext] ✓ Core 鉴权客户端初始化成功")
	}

	logx.Info("[ServiceContext] ✓ ServiceContext 初始化完成")

	// 创建认证中间件
	jwt := coremiddleware.JWT(coreAuth)
	authority := coremiddleware.Authority(coreAuth)
	logx.Infof("[ServiceContext] mode = %s", c.Mode)
	if c.IsLocal() {
		// 如果是本地环境，跳过权限校验
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

	return &ServiceContext{
		Config:          c,
		GameGrpcClient:  gameGrpcClient,
		Core:            coreCli,
		SyncRateLimiter: NewSyncRateLimiter(),
		Trans:           trans,
		Jwt:             jwt,       // JWT认证中间件
		Authority:       authority, // 权限校验中间件

		// 日志
		ActionLog: coremiddleware.ActionLog(coreadapt.ActionRecorder(coreCli)),       // 操作日志中间件
		ErrorLog:  coremiddleware.ErrorLog(c.Name, coreadapt.ErrorRecorder(coreCli)), // 错误日志中间件
	}
}

// ===== 同步操作频率限制 =====

// SyncRateLimiter 同步操作频率限制器（5秒限制）
type SyncRateLimiter struct {
	mu           sync.Mutex
	lastSyncTime map[string]time.Time // 记录每种对象类型的最后一次同步时间
	minInterval  time.Duration        // 最小间隔：5秒
}

// NewSyncRateLimiter 创建频率限制器
func NewSyncRateLimiter() *SyncRateLimiter {
	return &SyncRateLimiter{
		lastSyncTime: make(map[string]time.Time),
		minInterval:  5 * time.Second,
	}
}

// CheckRateLimit 检查频率限制，如果超过限制则返回 true，否则返回 false 并记录时间
func (r *SyncRateLimiter) CheckRateLimit(objectType string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	lastTime, exists := r.lastSyncTime[objectType]
	now := time.Now()

	// 如果从未调用过，或距离上次调用已超过5秒
	if !exists || now.Sub(lastTime) >= r.minInterval {
		r.lastSyncTime[objectType] = now
		return true
	}

	// 未超过5秒，拒绝请求
	return false
}
