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

type ExistsByOperatorIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewExistsByOperatorIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ExistsByOperatorIdLogic {
	return &ExistsByOperatorIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ExistsByOperatorId 检测分站管理员是否已存在
func (l *ExistsByOperatorIdLogic) ExistsByOperatorId(in *adminpb.ExistsAdminByOperatorIdRequest) (*adminpb.ExistsAdminByOperatorIdResponse, error) {
	if in.OperatorId <= 0 {
		return nil, xerr.RpcErr(xerr.BadRequest(i18nkey.ValidationError))
	}

	exists, err := l.svcCtx.DB.OperatorAdmin.Query().Where(operatoradmin.OperatorIDEQ(in.OperatorId)).Exist(l.ctx)
	if err != nil {
		return nil, enterror.Handle(l.Logger, err)
	}

	return &adminpb.ExistsAdminByOperatorIdResponse{
		Exists: exists,
	}, nil
}
