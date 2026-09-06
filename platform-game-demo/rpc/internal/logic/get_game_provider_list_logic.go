package logic

import (
	"context"

	"oa.98ent.com/p9/platform-game/common/constant"
	"oa.98ent.com/p9/platform-game/common/utils"
	"oa.98ent.com/p9/platform-game/pkg/game"
	"oa.98ent.com/p9/platform-game/rpc/internal/svc"
	platformgame "oa.98ent.com/p9/platform-game/rpc/pb/platform_game"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetGameProviderListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetGameProviderListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGameProviderListLogic {
	return &GetGameProviderListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取游戏供应商列�?
func (l *GetGameProviderListLogic) GetGameProviderList(in *platformgame.GetGameProviderListRequest) (*platformgame.GetGameProviderListResp, error) {
	l.Infof("GetGameProviderList called with page=%d, page_size=%d", in.Page, in.PageSize)
	if l.svcCtx.DB == nil {
		l.Errorf("Database not available")
		return &platformgame.GetGameProviderListResp{
			Code:    constant.CodeInternalError,
			Message: "Database not available",
		}, nil
	}

	// 构建查询参数
	page := int64(in.Page)
	if page <= 0 {
		page = 1
	}
	pageSize := int64(in.PageSize)
	if pageSize <= 0 {
		pageSize = 15
	}

	params := &game.GameProviderListParams{
		Page:      page,
		PageSize:  pageSize,
		Name:      in.GetName(),
		Status:    int16(in.Status),
		IsDeleted: int8(in.IsDeleted),
		SortBy:    in.GetSortBy(),
		SortOrder: in.GetSortOrder(),
	}

	// 调用数据库查询函数
	result, err := game.GameProviderListDB(l.ctx, l.svcCtx.DB, params)
	if err != nil {
		l.Errorf("GameProviderListDB failed: %v", err)
		return &platformgame.GetGameProviderListResp{
			Code:    constant.CodeInternalError,
			Message: "Failed to fetch providers: " + err.Error(),
		}, nil
	}

	// 转换结果为protobuf 消息
	providers := make([]*platformgame.ProviderInfo, 0, len(result.Providers))
	for _, provider := range result.Providers {
		if provider == nil {
			continue
		}

		providers = append(providers, &platformgame.ProviderInfo{
			Id:             provider.ID,
			Code:           provider.ProviderCode,
			Name:           provider.ProviderCode,
			Status:         int32(provider.Status),
			NameI18N:       string(utils.MustMarshalJSON(provider.NameI18n)),
			SourceNameI18N: string(utils.MustMarshalJSON(provider.SourceNameI18n)),
		})
	}

	return &platformgame.GetGameProviderListResp{
		Code:    constant.CodeSuccess,
		Message: "success",
		Data:    providers,
		Total:   result.Total,
	}, nil
}
