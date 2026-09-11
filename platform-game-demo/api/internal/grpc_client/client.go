package grpc_client

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"oa.98ent.com/p9/platform-game/api/internal/config"
	pb "oa.98ent.com/p9/platform-game/rpc/pb/platform_game"
)

// ClientManager gRPC客户端管理器
type ClientManager struct {
	conn *grpc.ClientConn
	cfg  config.GrpcClientConfig

	// Individual service clients
	pingServiceClient               pb.PingServiceClient
	gameServiceClient               pb.GameServiceClient
	gameCategoryServiceClient       pb.GameCategoryServiceClient
	gameProviderServiceClient       pb.GameProviderServiceClient
	gameChannelServiceClient        pb.GameChannelServiceClient
	gameCurrencyServiceClient       pb.GameCurrencyServiceClient
	gameSyncCheckpointServiceClient pb.GameSyncCheckpointServiceClient
	syncServiceClient               pb.SyncServiceClient
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
		conn:                            conn,
		cfg:                             cfg,
		pingServiceClient:               pb.NewPingServiceClient(conn),
		gameServiceClient:               pb.NewGameServiceClient(conn),
		gameCategoryServiceClient:       pb.NewGameCategoryServiceClient(conn),
		gameProviderServiceClient:       pb.NewGameProviderServiceClient(conn),
		gameChannelServiceClient:        pb.NewGameChannelServiceClient(conn),
		gameCurrencyServiceClient:       pb.NewGameCurrencyServiceClient(conn),
		gameSyncCheckpointServiceClient: pb.NewGameSyncCheckpointServiceClient(conn),
		syncServiceClient:               pb.NewSyncServiceClient(conn),
	}, nil
}

// GetPingServiceClient 获取 Ping 服务客户端
func (cm *ClientManager) GetPingServiceClient() pb.PingServiceClient {
	return cm.pingServiceClient
}

// GetGameServiceClient 获取游戏服务客户端
func (cm *ClientManager) GetGameServiceClient() pb.GameServiceClient {
	return cm.gameServiceClient
}

// GetGameCategoryServiceClient 获取游戏分类服务客户端
func (cm *ClientManager) GetGameCategoryServiceClient() pb.GameCategoryServiceClient {
	return cm.gameCategoryServiceClient
}

// GetGameProviderServiceClient 获取游戏供应商服务客户端
func (cm *ClientManager) GetGameProviderServiceClient() pb.GameProviderServiceClient {
	return cm.gameProviderServiceClient
}

// GetGameChannelServiceClient 获取游戏渠道服务客户端
func (cm *ClientManager) GetGameChannelServiceClient() pb.GameChannelServiceClient {
	return cm.gameChannelServiceClient
}

// GetGameCurrencyServiceClient 获取游戏货币服务客户端
func (cm *ClientManager) GetGameCurrencyServiceClient() pb.GameCurrencyServiceClient {
	return cm.gameCurrencyServiceClient
}

// GetGameSyncCheckpointServiceClient 获取游戏同步检查点服务客户端
func (cm *ClientManager) GetGameSyncCheckpointServiceClient() pb.GameSyncCheckpointServiceClient {
	return cm.gameSyncCheckpointServiceClient
}

// GetSyncServiceClient 获取同步服务客户端
func (cm *ClientManager) GetSyncServiceClient() pb.SyncServiceClient {
	return cm.syncServiceClient
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
