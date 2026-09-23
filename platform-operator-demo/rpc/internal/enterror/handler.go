package enterror

import (
	"github.com/zeromicro/go-zero/core/logx"
	"oa.98ent.com/p9/common/xerr"
	"oa.98ent.com/p9/platform-operator/pkg/i18nkey"
	"oa.98ent.com/p9/platform-operator/rpc/ent"
)

// Handle 处理 Ent 错误并转换为 gRPC 错误
func Handle(logger logx.Logger, err error) error {
	switch {
	case ent.IsNotFound(err):
		logger.Errorw("数据不存在", logx.Field("error", err.Error()))
		return xerr.RpcErr(xerr.EntNotFound(i18nkey.DataNotFound, err))

	case ent.IsConstraintError(err):
		logger.Errorw("数据约束冲突", logx.Field("error", err.Error()))
		return xerr.RpcErr(xerr.EntConstraintError(i18nkey.ConstraintError, err))

	case ent.IsValidationError(err):
		logger.Errorw("数据校验失败", logx.Field("error", err.Error()))
		return xerr.RpcErr(xerr.EntValidationError(i18nkey.ValidationError, err))

	case ent.IsNotSingular(err):
		logger.Errorw("查询结果不唯一", logx.Field("error", err.Error()))
		return xerr.RpcErr(xerr.EntInternalServerError(i18nkey.DatabaseError, err))

	default:
		logger.Errorw("数据库操作失败", logx.Field("error", err.Error()))
		return xerr.RpcErr(xerr.EntInternalServerError(i18nkey.DatabaseError, err))
	}
}
