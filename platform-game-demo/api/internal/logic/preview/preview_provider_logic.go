// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package preview

import (
	"context"
	"fmt"

	"oa.98ent.com/p9/platform-game/api/internal/svc"
	"oa.98ent.com/p9/platform-game/api/internal/types"
	"oa.98ent.com/p9/platform-game/common/constant"
	platformgame "oa.98ent.com/p9/platform-game/rpc/pb/platform_game"

	"github.com/zeromicro/go-zero/core/logx"
)

type PreviewProviderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPreviewProviderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PreviewProviderLogic {
	return &PreviewProviderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PreviewProviderLogic) PreviewProvider(req *types.SyncPreviewReq) (resp *types.SyncPreviewResp, err error) {
	l.Infof("[API PreviewProvider] received sync preview request")

	if l.svcCtx == nil || l.svcCtx.GrpcClient == nil {
		l.Error("[API PreviewProvider] gRPC client not available")
		return nil, fmt.Errorf("gRPC client not available")
	}

	grpcReq := &platformgame.SyncPreviewRequest{
		ObjectType: "provider",
		Page:       req.Page,
		PageSize:   req.PageSize,
		IsSkip:     req.IsSkip,
	}

	client := l.svcCtx.GrpcClient.GetSyncServiceClient()
	grpcResp, err := client.SyncPreview(l.ctx, grpcReq)
	if err != nil {
		l.Errorf("[API PreviewProvider] gRPC call failed: %v", err)
		return nil, err
	}

	if grpcResp == nil {
		l.Error("[API PreviewProvider] gRPC response is nil")
		return nil, fmt.Errorf("gRPC response is nil")
	}

	if grpcResp.Code != constant.CodeSuccess {
		l.Errorf("[API PreviewProvider] gRPC error: Code=%d, Message=%s", grpcResp.Code, grpcResp.Message)
		return nil, fmt.Errorf("gRPC error: %s", grpcResp.Message)
	}

	resp = &types.SyncPreviewResp{
		Stats: types.SyncStats{
			RemoteTotal: int64(grpcResp.Stats.RemoteTotal),
			LocalTotal:  int64(grpcResp.Stats.LocalTotal),
			CreateTotal: int64(grpcResp.Stats.CreateTotal),
			UpdateTotal: int64(grpcResp.Stats.UpdateTotal),
		},
	}

	// 转换 diffs
	var allDiffs []types.SyncDiff
	if grpcResp.Diffs != nil && len(grpcResp.Diffs) > 0 {
		allDiffs = make([]types.SyncDiff, 0, len(grpcResp.Diffs))
		for _, diff := range grpcResp.Diffs {
			allDiffs = append(allDiffs, types.SyncDiff{
				ObjectType:   diff.ObjectType,
				ObjectID:     diff.ObjectNumericId,
				ObjectCode:   diff.ObjectCode,
				Action:       diff.Action,
				RemoteID:     diff.RemoteId,
				RemoteCode:   diff.RemoteCode,
				Reason:       diff.Reason,
				ConflictType: diff.ConflictType,
			})
		}
	} else {
		allDiffs = make([]types.SyncDiff, 0)
	}

	// 处理分页和过滤
	page := req.Page
	pageSize := req.PageSize
	// isSkip := req.IsSkip

	// 设置默认值
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	// 执行分页和过滤
	// paginatedDiffs, totalCount, totalPage := PaginateAndFilterDiffs(allDiffs, page, pageSize, isSkip)

	resp.Diffs = allDiffs
	resp.Page = page
	resp.PageSize = pageSize
	resp.TotalCount = int64(grpcResp.Stats.DiffTotal)
	resp.TotalPage = (int64(grpcResp.Stats.DiffTotal) + pageSize - 1) / pageSize

	l.Infof("[API PreviewProvider] success: remoteTotal=%d, localTotal=%d, createTotal=%d, updateTotal=%d, page=%d, pageSize=%d, filtered=%d, totalPage=%d",
		grpcResp.Stats.RemoteTotal, grpcResp.Stats.LocalTotal, grpcResp.Stats.CreateTotal, grpcResp.Stats.UpdateTotal, page, pageSize, grpcResp.Stats.DiffTotal, (int64(grpcResp.Stats.DiffTotal)+pageSize-1)/pageSize)
	return resp, nil
}
