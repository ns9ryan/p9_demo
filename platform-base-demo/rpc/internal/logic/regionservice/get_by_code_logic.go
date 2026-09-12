package regionservicelogic

import (
	"context"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"

	"oa.98ent.com/p9/platform-base/pkg/i18nkey"
	"oa.98ent.com/p9/platform-base/pkg/rpc/grpcerror"
	entregion "oa.98ent.com/p9/platform-base/rpc/ent/region"
	"oa.98ent.com/p9/platform-base/rpc/internal/enterror"
	"oa.98ent.com/p9/platform-base/rpc/internal/svc"
	"oa.98ent.com/p9/platform-base/rpc/pb/base/region"
)

type GetByCodeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetByCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetByCodeLogic {
	return &GetByCodeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetByCode 按编码获取国家地区
func (l *GetByCodeLogic) GetByCode(in *region.GetRegionByCodeRequest) (*region.GetRegionByCodeResponse, error) {
	// 整理国家地区编码
	code := strings.ToUpper(strings.TrimSpace(in.Code))
	if code == "" {
		return nil, grpcerror.InvalidArgument(i18nkey.ValidationError)
	}

	// 按编码获取国家地区
	result, err := l.svcCtx.DB.Region.
		Query().
		Where(entregion.CodeEQ(code)).
		Only(l.ctx)
	if err != nil {
		// 转换Ent错误为gRPC错误
		return nil, enterror.Handle(l.Logger, err)
	}

	// 返回国家地区信息
	return &region.GetRegionByCodeResponse{
		Region: toRegionInfo(result), // 国家地区信息
	}, nil
}
