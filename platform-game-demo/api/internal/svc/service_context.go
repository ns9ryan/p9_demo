// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package svc

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/rest"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"oa.98ent.com/p9/core/common/ctxdata"
	"oa.98ent.com/p9/core/rpc/coreclient"
	"oa.98ent.com/p9/platform-game/api/internal/config"
	"oa.98ent.com/p9/platform-game/api/internal/grpc_client"
	"oa.98ent.com/p9/platform-game/common/auth"
	"oa.98ent.com/p9/platform-game/common/logger"
)

type ServiceContext struct {
	Config     config.Config
	GrpcClient *grpc_client.ClientManager
	CoreAuth   *auth.CoreAuth  // Core 服务鉴权客户端
	Auth       rest.Middleware // 认证中间件
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
	var coreAuth *auth.CoreAuth

	if c.CoreRpc.Target == "" {
		logger.Warn("[ServiceContext] ⚠ CoreRpc.Target 未配置，鉴权功能不可用")
	} else {
		// 创建 Core 服务的 gRPC 连接
		opts := []grpc.DialOption{
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		}

		ctx := context.Background()
		if c.CoreRpc.Timeout > 0 {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, time.Duration(c.CoreRpc.Timeout)*time.Second)
			defer cancel()
		}

		coreConn, err := grpc.DialContext(ctx, c.CoreRpc.Target, opts...)
		if err != nil {
			logger.Errorf("[ServiceContext] ✗ Core gRPC 连接初始化失败: %v", err)
			panic("failed to connect to Core gRPC service: " + err.Error())
		}

		// 包装连接以实现 zrpc.Client 接口
		zrpcClient := &grpc_client.ZrpcConnWrapper{coreConn}

		// 创建 Core 客户端并初始化鉴权
		coreCli := coreclient.NewCore(zrpcClient)
		coreAuth = auth.NewCoreAuth(coreCli)
		logger.Info("[ServiceContext] ✓ Core 鉴权客户端初始化成功")
	}

	logger.Info("[ServiceContext] ✓ ServiceContext 初始化完成")

	// 创建认证中间件
	authMiddleware := createAuthMiddleware(coreAuth)

	return &ServiceContext{
		Config:     c,
		GrpcClient: grpcClient,
		CoreAuth:   coreAuth,
		Auth:       authMiddleware,
	}
}

// createAuthMiddleware 创建认证中间件
func createAuthMiddleware(coreAuth *auth.CoreAuth) rest.Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// 如果 CoreAuth 未配置，返回错误
			if coreAuth == nil {
				logger.Warn("[Auth] ⚠ CoreAuth 未配置，无法进行鉴权")
				respondJSON(w, http.StatusUnauthorized, 401, "未登录")
				return
			}

			// 获取 Authorization 请求头
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				logger.Warn("[Auth] ⚠ 缺少 Authorization 请求头")
				respondJSON(w, http.StatusUnauthorized, 401, "未登录")
				return
			}

			// 解析 Bearer token
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				logger.Warn("[Auth] ⚠ Authorization 请求头格式不正确")
				respondJSON(w, http.StatusUnauthorized, 401, "未登录")
				return
			}

			token := parts[1]

			// 调用 Core 服务验证 token
			claims, err := coreAuth.CheckToken(r.Context(), token)
			if err != nil {
				logger.Warnf("[Auth] ⚠ Token 验证失败: %v", err)
				respondJSON(w, http.StatusUnauthorized, 401, "未登录")
				return
			}

			// 使用 ctxdata.WithClaims 将 claims 存储到 context 中
			// 这样 handler 就可以通过 ctxdata.ClaimsFromCtx() 获取用户信息
			ctx := ctxdata.WithClaims(r.Context(), claims)
			logger.Infof("[Auth] ✓ Token 验证成功，用户: %s", claims.Username)

			next(w, r.WithContext(ctx))
		}
	}
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
