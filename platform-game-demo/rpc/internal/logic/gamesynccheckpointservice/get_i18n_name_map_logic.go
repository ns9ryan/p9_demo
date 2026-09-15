package gamesynccheckpointservicelogic

import (
	"context"

	"oa.98ent.com/p9/platform-game/rpc/internal/svc"
	game_sync "oa.98ent.com/p9/platform-game/rpc/internal/synchro"
	"oa.98ent.com/p9/platform-game/rpc/pb/platform_game"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetI18nNameMapLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetI18nNameMapLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetI18nNameMapLogic {
	return &GetI18nNameMapLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取多语言映射文件name_map
func (l *GetI18nNameMapLogic) GetI18NNameMap(in *platform_game.GetI18NNameMapRequest) (*platform_game.GetI18NNameMapResp, error) {
	// todo: add your logic here and delete this line
	// 读取多语言映射文件name_map
	data := make(map[string]string)
	// 假设这里有逻辑去加载name_map文件到data中
	categoryData, _ := game_sync.LoadLocalI18nNameMap("game_category_name_map.json")
	providerData, _ := game_sync.LoadLocalI18nNameMap("game_provider_name_map.json")
	channelData, _ := game_sync.LoadLocalI18nNameMap("game_channel_name_map.json")
	// 合并各个数据到总的data中
	for k, v := range categoryData {
		data[k] = v
	}
	for k, v := range providerData {
		data[k] = v
	}
	for k, v := range channelData {
		data[k] = v
	}

	return &platform_game.GetI18NNameMapResp{
		Code:    0,
		Message: "success",
		Data:    data,
	}, nil
}
