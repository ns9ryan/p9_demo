package operatoradminservicelogic

import (
	"context"

	"oa.98ent.com/p9/common/xerr"
	"oa.98ent.com/p9/platform-operator/pkg/i18nkey"
	"oa.98ent.com/p9/platform-operator/rpc/ent/operatoradmin"
	"oa.98ent.com/p9/platform-operator/rpc/internal/enterror"
	"oa.98ent.com/p9/platform-operator/rpc/internal/svc"
	"oa.98ent.com/p9/platform-operator/rpc/pb/platformoperatorrpc/adminpb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteByOperatorIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteByOperatorIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteByOperatorIdLogic {
	return &DeleteByOperatorIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// DeleteByOperatorId 按分站删除管理员
func (l *DeleteByOperatorIdLogic) DeleteByOperatorId(in *adminpb.DeleteAdminsByOperatorIdRequest) (*adminpb.DeleteAdminsByOperatorIdResponse, error) {
	if in.OperatorId <= 0 {
		return nil, xerr.RpcErr(xerr.BadRequest(i18nkey.ValidationError))
	}

	_, err := l.svcCtx.DB.OperatorAdmin.
		Delete().
		Where(operatoradmin.OperatorIDEQ(in.OperatorId)).
		Exec(l.ctx)
	if err != nil {
		return nil, enterror.Handle(l.Logger, err)
	}

	return &adminpb.DeleteAdminsByOperatorIdResponse{}, nil
}
