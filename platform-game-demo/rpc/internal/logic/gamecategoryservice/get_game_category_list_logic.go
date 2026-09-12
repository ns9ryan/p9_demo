package gamecategoryservicelogic

import (
	"context"
	"fmt"
	"strings"

	"oa.98ent.com/p9/platform-game/rpc/ent"
	"oa.98ent.com/p9/platform-game/rpc/internal/constant"
	"oa.98ent.com/p9/platform-game/rpc/internal/logic"
	"oa.98ent.com/p9/platform-game/rpc/internal/svc"
	"oa.98ent.com/p9/platform-game/rpc/internal/utils"
	"oa.98ent.com/p9/platform-game/rpc/pb/platform_game"

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
func (l *GetGameCategoryListLogic) GetGameCategoryList(in *platform_game.GetGameCategoryListRequest) (*platform_game.GetGameCategoryListResp, error) {
	l.Infof("[RPC GetGameCategoryList] received req: page=%d, page_size=%d", in.Page, in.PageSize)
	if l.svcCtx == nil || l.svcCtx.DB == nil {
		l.Errorf("[RPC GetGameCategoryList] Database not available")
		return &platform_game.GetGameCategoryListResp{
			Code:    constant.CodeInternalError,
			Message: "Database not available",
		}, nil
	}

	page, pageSize := utils.HandlePage(int64(in.Page), int64(in.PageSize))

	query := l.svcCtx.DB.WithContext(l.ctx)

	// 处理软删除条件
	if in.GetIsDeleted() == 0 {
		query = query.Where("deleted_at IS NULL")
	} else if in.GetIsDeleted() == 1 {
		query = query.Where("deleted_at IS NOT NULL")
	}

	if in.GetStatus() > 0 {
		query = query.Where("status = ?", in.GetStatus())
	}
	if in.GetCategoryCode() != "" {
		query = query.Where("category_code LIKE ?", "%"+in.GetCategoryCode()+"%")
	}
	if in.GetName() != "" {
		query = query.Where("name_i18n->>'default' LIKE ?", "%"+in.GetName()+"%")
	}

	sortBy := strings.TrimSpace(in.SortBy)
	sortOrder := strings.TrimSpace(in.SortOrder)
	if sortBy == "" {
		sortBy = "sort_no"
	}
	if sortOrder == "" {
		sortOrder = "asc"
	}
	sortBy = strings.ToLower(sortBy)
	sortOrder = strings.ToLower(sortOrder)

	switch sortBy {
	case "id", "sort_no", "created_at":
	default:
		sortBy = "id"
	}

	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "asc"
	}

	var total int64
	if err := query.Model(&ent.GameCategory{}).Count(&total).Error; err != nil {
		l.Errorf("[RPC GetGameCategoryList] count failed: %v", err)
		return &platform_game.GetGameCategoryListResp{
			Code:    constant.CodeInternalError,
			Message: fmt.Sprintf("count failed: %v", err),
		}, nil
	}

	var categories []*ent.GameCategory
	if err := query.
		Offset(int((page - 1) * pageSize)).
		Limit(int(pageSize)).
		Order(fmt.Sprintf("%s %s", sortBy, sortOrder)).
		Find(&categories).Error; err != nil {
		l.Errorf("[RPC GetGameCategoryList] query failed: %v", err)
		return &platform_game.GetGameCategoryListResp{
			Code:    constant.CodeInternalError,
			Message: fmt.Sprintf("query failed: %v", err),
		}, nil
	}

	l.Infof("[RPC GetGameCategoryList] success: total=%d", total)
	return &platform_game.GetGameCategoryListResp{
		Code:     constant.CodeSuccess,
		Message:  "success",
		Data:     logic.CategoryModelToProtoList(categories),
		Total:    total,
		Page:     int32(page),
		PageSize: int32(pageSize),
	}, nil
}
