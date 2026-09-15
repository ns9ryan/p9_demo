// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package svc

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
	"oa.98ent.com/p9/core/common/coreadapt"
	"oa.98ent.com/p9/core/common/middleware"
	coremiddleware "oa.98ent.com/p9/core/common/middleware"
	"oa.98ent.com/p9/core/rpc/coreclient"
	"oa.98ent.com/p9/platform-game/api/internal/config"
	"oa.98ent.com/p9/platform-game/api/internal/grpc_client"
	"oa.98ent.com/p9/platform-game/common/auth"
	"oa.98ent.com/p9/platform-game/common/logger"
)

type ServiceContext struct {
	Config     config.Config
	GrpcClient *grpc_client.ClientManager
	Core       coreclient.Core // Core RPC客户端
	CoreAuth   *auth.CoreAuth  // Core 服务鉴权客户端
	// Auth            rest.Middleware  // 认证中间件
	SyncRateLimiter *SyncRateLimiter // 同步操作频率限制器

	Jwt       rest.Middleware // JWT认证中间件
	Authority rest.Middleware // 权限校验中间件

	// 日志
	ActionLog rest.Middleware // 操作日志中间件
	ErrorLog  rest.Middleware // 错误日志中间件
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 初始化 gRPC 客户端
	logger.Info("[ServiceContext] 开始初始化 gRPC 客户端...")
	var grpcClient *grpc_client.ClientManager
	var err error

	if c.GrpcClient.Target == "" {
		logger.Error("[ServiceContext] ✗ GrpcClient.Target 未配置，无法初始化 gRPC 客户端")
		panic("GrpcClient.Target is required for API service")
	}

	grpcClient, err = grpc_client.NewClientManager(c.GrpcClient)
	if err != nil {
		logger.Errorf("[ServiceContext] ✗ gRPC 客户端初始化失败: %v", err)
		panic("failed to initialize gRPC client: " + err.Error())
	}
	logger.Info("[ServiceContext] ✓ gRPC 客户端初始化成功")

	// 初始化 Core 服务鉴权客户端
	logger.Info("[ServiceContext] 开始初始化 Core 鉴权客户端...")
	logger.Infof("[ServiceContext] CoreRpc.Target = '%s'", c.CoreRpc.Target)
	var coreAuth middleware.Client
	var coreCli coreclient.Core

	if c.CoreRpc.Target == "" {
		logger.Warn("[ServiceContext] ⚠ CoreRpc.Target 未配置，鉴权功能不可用")
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

		logger.Info("[ServiceContext] ✓ Core 鉴权客户端初始化成功")
	}

	logger.Info("[ServiceContext] ✓ ServiceContext 初始化完成")

	// 创建认证中间件
	// authMiddleware := createAuthMiddleware(coreAuth)
	jwt := coremiddleware.JWT(coreAuth)
	authority := coremiddleware.Authority(coreAuth)
	logger.Infof("[ServiceContext] isLocal = %v", c.IsLocal)
	if isLocal := c.IsLocal; isLocal {
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
		Config:     c,
		GrpcClient: grpcClient,
		Core:       coreCli,
		// Auth:            authMiddleware,
		SyncRateLimiter: NewSyncRateLimiter(),
		Jwt:             jwt,       // JWT认证中间件
		Authority:       authority, // 权限校验中间件

		// 日志
		ActionLog: coremiddleware.ActionLog(coreadapt.ActionRecorder(coreCli)),       // 操作日志中间件
		ErrorLog:  coremiddleware.ErrorLog(c.Name, coreadapt.ErrorRecorder(coreCli)), // 错误日志中间件
	}
}

// getClientIP 获取客户端 IP，支持代理情况下的 X-Forwarded-For 头
func getClientIP(r *http.Request) string {
	if r == nil {
		return ""
	}
	// 优先从 X-Forwarded-For 获取（由代理设置）
	if x := r.Header.Get("X-Forwarded-For"); x != "" {
		return normalizeIP(strings.TrimSpace(strings.Split(x, ",")[0]))
	}
	// 其次尝试从 X-Real-IP 获取
	if x := r.Header.Get("X-Real-IP"); x != "" {
		return normalizeIP(strings.TrimSpace(x))
	}
	// 最后从 RemoteAddr 获取
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return normalizeIP(r.RemoteAddr)
	}
	return normalizeIP(host)
}

// normalizeIP 规范化 IP：去空白、去方括号，并把 IPv4-mapped IPv6 转成 IPv4
func normalizeIP(raw string) string {
	raw = strings.TrimSpace(raw)
	raw = strings.TrimPrefix(raw, "[")
	raw = strings.TrimSuffix(raw, "]")
	ip := net.ParseIP(raw)
	if ip == nil {
		return raw
	}
	// 如果是 IPv4-mapped IPv6，转成 IPv4
	if v4 := ip.To4(); v4 != nil {
		return v4.String()
	}
	return ip.String()
}

// respondJSON 返回 JSON 格式的响应
func respondJSON(w http.ResponseWriter, statusCode int, code int, msg string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	response := map[string]interface{}{
		"code": code,
		"msg":  msg,
	}
	json.NewEncoder(w).Encode(response)
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
