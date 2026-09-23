package dao

import "oa.98ent.com/p9/platform-game/rpc/ent"

// Manager 数据库访问对象管理器
type Manager struct {
	DB                      *ent.Client
	GameCategory            *GameCategoryDAO
	GameChannel             *GameChannelDAO
	GameCurrency            *GameCurrencyDAO
	Currency                *CurrencyDAO
	GameProvider            *GameProviderDAO
	Game                    *GameDAO
	GameSyncCheckpoint      *GameSyncCheckpointDAO
	OperatorGame            *OperatorGameDAO
	OperatorGameCategory    *OperatorGameCategoryDAO
	OperatorGameChannel     *OperatorGameChannelDAO
	OperatorGameProvider    *OperatorGameProviderDAO
	Operator                *OperatorDAO
}

// NewManager 创建 DAO 管理器
func NewManager(db *ent.Client) *Manager {
	return &Manager{
		DB:                      db,
		GameCategory:            NewGameCategoryDAO(db),
		GameChannel:             NewGameChannelDAO(db),
		GameCurrency:            NewGameCurrencyDAO(db),
		Currency:                NewCurrencyDAO(db),
		GameProvider:            NewGameProviderDAO(db),
		Game:                    NewGameDAO(db),
		GameSyncCheckpoint:      NewGameSyncCheckpointDAO(db),
		OperatorGame:            NewOperatorGameDAO(db),
		OperatorGameCategory:    NewOperatorGameCategoryDAO(db),
		OperatorGameChannel:     NewOperatorGameChannelDAO(db),
		OperatorGameProvider:    NewOperatorGameProviderDAO(db),
		Operator:                NewOperatorDAO(db),
	}
}
