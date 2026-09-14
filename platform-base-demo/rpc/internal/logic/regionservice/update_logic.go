package regionservicelogic

import (
	"context"
	"strings"

	"oa.98ent.com/p9/platform-base/pkg/i18nkey"
	"oa.98ent.com/p9/platform-base/pkg/rpc/grpcerror"
	"oa.98ent.com/p9/platform-base/rpc/internal/enterror"
	"oa.98ent.com/p9/platform-base/rpc/internal/svc"
	"oa.98ent.com/p9/platform-base/rpc/pb/platformbaserpc/regionpb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLogic {
	return &UpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// Update 修改国家地区
func (l *UpdateLogic) Update(in *regionpb.UpdateRegionRequest) (*regionpb.UpdateRegionResponse, error) {
	// 国家地区ID必须大于0
	if in.Id <= 0 {
		return nil, grpcerror.InvalidArgument(i18nkey.ValidationError)
	}

	// 至少需要修改一个字段
	if in.CallingCode == nil && in.Status == nil {
		return nil, grpcerror.InvalidArgument(i18nkey.ValidationError)
	}

	var callingCode *string

	// 传入国际电话区号时进行整理
	if in.CallingCode != nil {
		value := strings.TrimSpace(*in.CallingCode)
		callingCode = &value
	}

	// 修改国家地区
	update := l.svcCtx.DB.Region.
		UpdateOneID(in.Id).
		SetNillableStatus(in.Status) // 状态: 1启用, 2停用

	// 传入国际电话区号时进行修改
	if callingCode != nil {
		update.SetCallingCode(*callingCode) // 国际电话区号
	}

	if err := update.Exec(l.ctx); err != nil {
		// 转换Ent错误为gRPC错误
		return nil, enterror.Handle(l.Logger, err)
	}

	// 返回修改结果
	return &regionpb.UpdateRegionResponse{}, nil
}
