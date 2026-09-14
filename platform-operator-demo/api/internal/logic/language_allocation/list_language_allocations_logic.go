// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package language_allocation

import (
	"context"

	"oa.98ent.com/p9/platform-operator/api/internal/svc"
	"oa.98ent.com/p9/platform-operator/api/internal/types"
	"oa.98ent.com/p9/platform-operator/rpc/pb/platformoperatorrpc/languageallocationpb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListLanguageAllocationsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListLanguageAllocationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListLanguageAllocationsLogic {
	return &ListLanguageAllocationsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// ListLanguageAllocations 获取语言分配列表
func (l *ListLanguageAllocationsLogic) ListLanguageAllocations(req *types.ListLanguageAllocationsRequest) (resp *types.ListLanguageAllocationsResponse, err error) {
	// 获取分站当前语言分配关系
	result, err := l.svcCtx.LanguageAllocationRpc.List(
		l.ctx,
		&languageallocationpb.ListLanguageAllocationsRequest{
			OperatorId: req.OperatorId, // 分站ID
		},
	)
	if err != nil {
		return nil, err
	}

	// 转换语言分配列表
	list := make([]types.LanguageAllocationInfo, 0, len(result.List))
	for _, item := range result.List {
		list = append(list, types.LanguageAllocationInfo{
			LanguageCode: item.LanguageCode, // 语言编码
			AllocatedAt:  item.CreatedAt,    // 分配时间, Unix毫秒时间戳
		})
	}

	// 返回语言分配列表
	return &types.ListLanguageAllocationsResponse{
		List: list, // 已分配语言列表
	}, nil
}
