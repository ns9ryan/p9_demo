package grpc_client

import (
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/zrpc"
	"oa.98ent.com/p9/platform-base/pkg/api/rpcerror"
	"oa.98ent.com/p9/platform-game/rpc/client/gamecategoryservice"
	"oa.98ent.com/p9/platform-game/rpc/client/gamechannelservice"
	"oa.98ent.com/p9/platform-game/rpc/client/gamecurrencyservice"
	"oa.98ent.com/p9/platform-game/rpc/client/gameproviderservice"
	"oa.98ent.com/p9/platform-game/rpc/client/gameservice"
	"oa.98ent.com/p9/platform-game/rpc/client/gamesynccheckpointservice"
	"oa.98ent.com/p9/platform-game/rpc/client/operatorgameallocationservice"
	"oa.98ent.com/p9/platform-game/rpc/client/operatorgamecategoryservice"
	"oa.98ent.com/p9/platform-game/rpc/client/operatorgamechannelservice"
	"oa.98ent.com/p9/platform-game/rpc/client/operatorgameproviderservice"
	"oa.98ent.com/p9/platform-game/rpc/client/operatorgameservice"
	"oa.98ent.com/p9/platform-game/rpc/client/syncservice"
)

// ClientManager gRPC客户端管理器
type GameClientManager struct {
	gameServiceClient               gameservice.GameService
	gameCategoryServiceClient       gamecategoryservice.GameCategoryService
	gameProviderServiceClient       gameproviderservice.GameProviderService
	gameChannelServiceClient        gamechannelservice.GameChannelService
	gameCurrencyServiceClient       gamecurrencyservice.GameCurrencyService
	gameSyncCheckpointServiceClient gamesynccheckpointservice.GameSyncCheckpointService
	syncServiceClient               syncservice.SyncService

	operatorGameAllocationServiceClient operatorgameallocationservice.OperatorGameAllocationService
	operatorGameCategoryServiceClient   operatorgamecategoryservice.OperatorGameCategoryService
	operatorGameChannelServiceClient    operatorgamechannelservice.OperatorGameChannelService
	operatorGameProviderServiceClient   operatorgameproviderservice.OperatorGameProviderService
	operatorGameServiceClient           operatorgameservice.OperatorGameService
}

// 创建新的游戏gRPC客户端管理器
func NewGameClientManager(cfg zrpc.RpcClientConf) (*GameClientManager, error) {
	logx.Infof("create version 1.0.0 game client manager")
	client := zrpc.MustNewClient(
		cfg,
		zrpc.WithUnaryClientInterceptor(rpcerror.UnaryClientInterceptor),
	)

	return &GameClientManager{
		gameServiceClient:                   gameservice.NewGameService(client),
		gameCategoryServiceClient:           gamecategoryservice.NewGameCategoryService(client),
		gameProviderServiceClient:           gameproviderservice.NewGameProviderService(client),
		gameChannelServiceClient:            gamechannelservice.NewGameChannelService(client),
		gameCurrencyServiceClient:           gamecurrencyservice.NewGameCurrencyService(client),
		gameSyncCheckpointServiceClient:     gamesynccheckpointservice.NewGameSyncCheckpointService(client),
		syncServiceClient:                   syncservice.NewSyncService(client),
		operatorGameAllocationServiceClient: operatorgameallocationservice.NewOperatorGameAllocationService(client),
		operatorGameCategoryServiceClient:   operatorgamecategoryservice.NewOperatorGameCategoryService(client),
		operatorGameChannelServiceClient:    operatorgamechannelservice.NewOperatorGameChannelService(client),
		operatorGameProviderServiceClient:   operatorgameproviderservice.NewOperatorGameProviderService(client),
		operatorGameServiceClient:           operatorgameservice.NewOperatorGameService(client),
	}, nil
}

// GetGameServiceClient 获取游戏服务客户端
func (cm *GameClientManager) GetGameServiceClient() gameservice.GameService {
	return cm.gameServiceClient
}

// GetGameCategoryServiceClient 获取游戏分类服务客户端
func (cm *GameClientManager) GetGameCategoryServiceClient() gamecategoryservice.GameCategoryService {
	return cm.gameCategoryServiceClient
}

// GetGameProviderServiceClient 获取游戏供应商服务客户端
func (cm *GameClientManager) GetGameProviderServiceClient() gameproviderservice.GameProviderService {
	return cm.gameProviderServiceClient
}

// GetGameChannelServiceClient 获取游戏渠道服务客户端
func (cm *GameClientManager) GetGameChannelServiceClient() gamechannelservice.GameChannelService {
	return cm.gameChannelServiceClient
}

// GetGameCurrencyServiceClient 获取游戏货币服务客户端
func (cm *GameClientManager) GetGameCurrencyServiceClient() gamecurrencyservice.GameCurrencyService {
	return cm.gameCurrencyServiceClient
}

// GetGameSyncCheckpointServiceClient 获取游戏同步检查点服务客户端
func (cm *GameClientManager) GetGameSyncCheckpointServiceClient() gamesynccheckpointservice.GameSyncCheckpointService {
	return cm.gameSyncCheckpointServiceClient
}

// GetSyncServiceClient 获取同步服务客户端
func (cm *GameClientManager) GetSyncServiceClient() syncservice.SyncService {
	return cm.syncServiceClient
}

// GetOperatorGameAllocationServiceClient 获取运营商游戏分配服务客户端
func (cm *GameClientManager) GetOperatorGameAllocationServiceClient() operatorgameallocationservice.OperatorGameAllocationService {
	return cm.operatorGameAllocationServiceClient
}

// GetOperatorGameCategoryServiceClient 获取运营商游戏分类服务客户端
func (cm *GameClientManager) GetOperatorGameCategoryServiceClient() operatorgamecategoryservice.OperatorGameCategoryService {
	return cm.operatorGameCategoryServiceClient
}

// GetOperatorGameChannelServiceClient 获取运营商游戏渠道服务客户端
func (cm *GameClientManager) GetOperatorGameChannelServiceClient() operatorgamechannelservice.OperatorGameChannelService {
	return cm.operatorGameChannelServiceClient
}

// GetOperatorGameProviderServiceClient 获取运营商游戏供应商服务客户端
func (cm *GameClientManager) GetOperatorGameProviderServiceClient() operatorgameproviderservice.OperatorGameProviderService {
	return cm.operatorGameProviderServiceClient
}

// GetOperatorGameServiceClient 获取运营商游戏服务客户端
func (cm *GameClientManager) GetOperatorGameServiceClient() operatorgameservice.OperatorGameService {
	return cm.operatorGameServiceClient
}
