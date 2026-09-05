package timezoneservicelogic

import (
	"context"

	"oa.98ent.com/p9/platform-base/pkg/i18nkey"
	"oa.98ent.com/p9/platform-base/pkg/rpc/grpcerror"
	enttimezone "oa.98ent.com/p9/platform-base/rpc/ent/timezone"
	"oa.98ent.com/p9/platform-base/rpc/internal/enterror"
	"oa.98ent.com/p9/platform-base/rpc/internal/svc"
	"oa.98ent.com/p9/platform-base/rpc/pb/base/timezone"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListAllLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListAllLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListAllLogic {
	return &ListAllLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ListAll 获取全部时区
func (l *ListAllLogic) ListAll(in *timezone.ListAllTimezonesRequest) (*timezone.ListAllTimezonesResponse, error) {
	// 校验状态
	if in.Status != nil && (*in.Status < 1 || *in.Status > 2) {
		return nil, grpcerror.InvalidArgument(i18nkey.ValidationError)
	}

	// 创建时区查询
	query := l.svcCtx.DB.Timezone.Query()

	// 按状态筛选
	if in.Status != nil {
		query = query.Where(enttimezone.StatusEQ(*in.Status))
	}

	// 获取全部时区
	results, err := query.
		Order(
			enttimezone.BySortNo(), // 按排序值升序
			enttimezone.ByID(),     // 排序值相同时按ID升序
		).
		All(l.ctx)
	if err != nil {
		// 转换Ent错误为gRPC错误
		return nil, enterror.Handle(l.Logger, err)
	}

	// 转换时区列表
	list := make([]*timezone.TimezoneInfo, 0, len(results))
	for _, result := range results {
		list = append(list, toTimezoneInfo(result))
	}

	return &timezone.ListAllTimezonesResponse{
		List: list,
	}, nil
}
