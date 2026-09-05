package timezoneservicelogic

import (
	"context"
	"strings"

	"oa.98ent.com/p9/platform-base/pkg/i18nkey"
	"oa.98ent.com/p9/platform-base/pkg/rpc/grpcerror"
	entlanguage "oa.98ent.com/p9/platform-base/rpc/ent/language"
	"oa.98ent.com/p9/platform-base/rpc/internal/enterror"
	"oa.98ent.com/p9/platform-base/rpc/internal/svc"
	"oa.98ent.com/p9/platform-base/rpc/pb/base/timezone"

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

// Update 修改时区
func (l *UpdateLogic) Update(in *timezone.UpdateTimezoneRequest) (*timezone.UpdateTimezoneResponse, error) {
	// 时区ID必须大于0
	if in.Id <= 0 {
		return nil, grpcerror.InvalidArgument(i18nkey.ValidationError)
	}

	// 至少需要修改一个字段
	if len(in.NameI18N) == 0 && in.Status == nil {
		return nil, grpcerror.InvalidArgument(i18nkey.ValidationError)
	}

	var nameI18N map[string]string

	// 传入多语言名称时进行校验和整理
	if len(in.NameI18N) > 0 {
		nameI18N = make(map[string]string, len(in.NameI18N))
		codes := make([]string, 0, len(in.NameI18N))

		for code, name := range in.NameI18N {
			name = strings.TrimSpace(name)
			if code == "" || name == "" {
				return nil, grpcerror.InvalidArgument(i18nkey.ValidationError)
			}

			nameI18N[code] = name
			codes = append(codes, code)
		}

		// 校验多语言名称中的语言编码是否已经存在
		count, err := l.svcCtx.DB.Language.
			Query().
			Where(entlanguage.CodeIn(codes...)).
			Count(l.ctx)
		if err != nil {
			// 转换Ent错误为gRPC错误
			return nil, enterror.Handle(l.Logger, err)
		}

		if count != len(codes) {
			return nil, grpcerror.InvalidArgument(i18nkey.ValidationError)
		}
	}

	// 修改时区
	update := l.svcCtx.DB.Timezone.
		UpdateOneID(in.Id).
		SetNillableStatus(in.Status) // 状态: 1启用, 2停用

	// 传入多语言名称时才修改
	if nameI18N != nil {
		update.SetNameI18n(nameI18N) // 多语言名称
	}

	if err := update.Exec(l.ctx); err != nil {
		// 转换Ent错误为gRPC错误
		return nil, enterror.Handle(l.Logger, err)
	}

	return &timezone.UpdateTimezoneResponse{}, nil
}
