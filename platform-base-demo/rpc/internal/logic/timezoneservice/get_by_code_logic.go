package timezoneservicelogic

import (
	"context"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"

	"oa.98ent.com/p9/platform-base/pkg/i18nkey"
	"oa.98ent.com/p9/platform-base/pkg/rpc/grpcerror"
	enttimezone "oa.98ent.com/p9/platform-base/rpc/ent/timezone"
	"oa.98ent.com/p9/platform-base/rpc/internal/enterror"
	"oa.98ent.com/p9/platform-base/rpc/internal/svc"
	"oa.98ent.com/p9/platform-base/rpc/pb/base/timezone"
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

// GetByCode 按编码获取时区
func (l *GetByCodeLogic) GetByCode(in *timezone.GetTimezoneByCodeRequest) (*timezone.GetTimezoneByCodeResponse, error) {
	// 整理时区编码
	code := strings.TrimSpace(in.Code)
	if code == "" {
		return nil, grpcerror.InvalidArgument(i18nkey.ValidationError)
	}

	// 按编码获取时区
	result, err := l.svcCtx.DB.Timezone.
		Query().
		Where(enttimezone.CodeEQ(code)).
		Only(l.ctx)
	if err != nil {
		// 转换Ent错误为gRPC错误
		return nil, enterror.Handle(l.Logger, err)
	}

	// 返回时区信息
	return &timezone.GetTimezoneByCodeResponse{
		Timezone: toTimezoneInfo(result), // 时区信息
	}, nil
}
