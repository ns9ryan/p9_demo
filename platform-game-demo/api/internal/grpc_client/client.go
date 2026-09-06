package grpc_client

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"oa.98ent.com/p9/platform-game/api/internal/config"
	platformgame "oa.98ent.com/p9/platform-game/rpc/pb/platform_game"
)

// ClientManager gRPC客户端管理器
type ClientManager struct {
	conn   *grpc.ClientConn
	client platformgame.PlatformGameServiceClient
	cfg    config.GrpcClientConfig
}

// NewClientManager 创建新的gRPC客户端管理器
func NewClientManager(cfg config.GrpcClientConfig) (*ClientManager, error) {
	if cfg.Target == "" {
		return nil, fmt.Errorf("grpc client target is empty")
	}

	// 设置拨号选项
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	// 如果配置了超时，使用带超时的上下文
	ctx := context.Background()
	if cfg.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, time.Duration(cfg.Timeout)*time.Second)
		defer cancel()
	}

	// 建立连接
	conn, err := grpc.DialContext(ctx, cfg.Target, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to gRPC server: %w", err)
	}

	return &ClientManager{
		conn:   conn,
		client: platformgame.NewPlatformGameServiceClient(conn),
		cfg:    cfg,
	}, nil
}

// GetPlatformGameServiceClient 获取平台游戏服务客户端
func (cm *ClientManager) GetPlatformGameServiceClient() platformgame.PlatformGameServiceClient {
	return cm.client
}

// Close 关闭gRPC连接
func (cm *ClientManager) Close() error {
	if cm.conn != nil {
		return cm.conn.Close()
	}
	return nil
}

// ZrpcConnWrapper 包装 *grpc.ClientConn 以实现 zrpc.Client 接口
type ZrpcConnWrapper struct {
	*grpc.ClientConn
}

// Conn 返回 gRPC 连接，实现 zrpc.Client 接口
func (w *ZrpcConnWrapper) Conn() *grpc.ClientConn {
	return w.ClientConn
}
