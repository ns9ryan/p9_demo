// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package logic

import (
	"context"
	"fmt"

	"oa.98ent.com/p9/platform-game/api/internal/svc"
	"oa.98ent.com/p9/platform-game/api/internal/types"
	"oa.98ent.com/p9/platform-game/common/constant"
	platformgame "oa.98ent.com/p9/platform-game/rpc/pb/platform_game"

	"github.com/zeromicro/go-zero/core/logx"
)

type PreviewCurrencyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPreviewCurrencyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PreviewCurrencyLogic {
	return &PreviewCurrencyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PreviewCurrencyLogic) PreviewCurrency(req *types.SyncPreviewReq) (resp *types.SyncPreviewResp, err error) {
	l.Infof("[API PreviewCurrency] received sync preview request")

	if l.svcCtx == nil || l.svcCtx.GrpcClient == nil {
		l.Error("[API PreviewCurrency] gRPC client not available")
		return nil, fmt.Errorf("gRPC client not available")
	}

	grpcReq := &platformgame.SyncPreviewRequest{
		ObjectType: "currency",
	}

	client := l.svcCtx.GrpcClient.GetPlatformGameServiceClient()
	grpcResp, err := client.SyncPreview(l.ctx, grpcReq)
	if err != nil {
		l.Errorf("[API PreviewCurrency] gRPC call failed: %v", err)
		return nil, err
	}

	if grpcResp == nil {
		l.Error("[API PreviewCurrency] gRPC response is nil")
		return nil, fmt.Errorf("gRPC response is nil")
	}

	if grpcResp.Code != constant.CodeSuccess {
		l.Errorf("[API PreviewCurrency] gRPC error: Code=%d, Message=%s", grpcResp.Code, grpcResp.Message)
		return nil, fmt.Errorf("gRPC error: %s", grpcResp.Message)
	}

	resp = &types.SyncPreviewResp{
		Stats: types.SyncStats{
			RemoteTotal: int64(grpcResp.Stats.Total),
			LocalTotal:  int64(grpcResp.Stats.Total),
			CreateTotal: int64(grpcResp.Stats.Added),
			UpdateTotal: int64(grpcResp.Stats.Updated),
		},
	}

	// 转换 diffs
	if grpcResp.Diffs != nil && len(grpcResp.Diffs) > 0 {
		resp.Diffs = make([]types.SyncDiff, 0, len(grpcResp.Diffs))
		for _, diff := range grpcResp.Diffs {
			resp.Diffs = append(resp.Diffs, types.SyncDiff{
				ObjectType:   diff.Type,
				ObjectID:     diff.ObjectNumericId,
				ObjectCode:   diff.ObjectCode,
				Action:       diff.Action,
				RemoteID:     diff.RemoteId,
				RemoteCode:   diff.RemoteCode,
				Reason:       diff.Reason,
				ConflictType: diff.ConflictType,
				Details:      diff.Details,
			})
		}
	} else {
		// 确保返回空数组而不是 nil
		resp.Diffs = make([]types.SyncDiff, 0)
	}

	l.Infof("[API PreviewCurrency] success: total=%d, added=%d, updated=%d, deleted=%d, diffs_count=%d",
		grpcResp.Stats.Total, grpcResp.Stats.Added, grpcResp.Stats.Updated, grpcResp.Stats.Deleted, len(resp.Diffs))
	return resp, nil
}
