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

type GetGameCategoryListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetGameCategoryListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGameCategoryListLogic {
	return &GetGameCategoryListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取游戏分类列表
func (l *GetGameCategoryListLogic) GetGameCategoryList(in *platformgame.GetGameCategoryListRequest) (*platformgame.GetGameCategoryListResp, error) {
	l.Infof("GetGameCategoryList called with page=%d, page_size=%d", in.Page, in.PageSize)
	if l.svcCtx.DB == nil {
		l.Errorf("Database not available")
		return &platformgame.GetGameCategoryListResp{
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

	params := &game.GameCategoryListParams{
		Page:         page,
		PageSize:     pageSize,
		CategoryCode: in.GetCategoryCode(),
		Name:         in.GetName(),
		Status:       int16(in.Status),
		IsDeleted:    int8(in.IsDeleted),
		SortBy:       in.GetSortBy(),
		SortOrder:    in.GetSortOrder(),
	}

	// 调用数据库查询函数
	result, err := game.GameCategoryListDB(l.ctx, l.svcCtx.DB, params)
	if err != nil {
		l.Errorf("GameCategoryListDB failed: %v", err)
		return &platformgame.GetGameCategoryListResp{
			Code:    constant.CodeInternalError,
			Message: "Failed to fetch categories: " + err.Error(),
		}, nil
	}

	// 转换结果为 protobuf 消息
	categories := make([]*platformgame.CategoryInfo, 0, len(result.Categories))
	for _, cat := range result.Categories {
		if cat == nil {
			continue
		}

		categories = append(categories, &platformgame.CategoryInfo{
			Id:             cat.ID,
			Code:           cat.CategoryCode,
			Name:           cat.CategoryCode,
			Status:         int32(cat.Status),
			NameI18N:       string(utils.MustMarshalJSON(cat.NameI18n)),
			SourceNameI18N: string(utils.MustMarshalJSON(cat.SourceNameI18n)),
		})
	}

	return &platformgame.GetGameCategoryListResp{
		Code:    constant.CodeSuccess,
		Message: "success",
		Data:    categories,
		Total:   result.Total,
	}, nil
}
