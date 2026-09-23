package operatoradminservicelogic

import (
	"context"

	"oa.98ent.com/p9/common/xerr"
	"oa.98ent.com/p9/platform-operator/pkg/i18nkey"
	"oa.98ent.com/p9/platform-operator/rpc/ent"
	"oa.98ent.com/p9/platform-operator/rpc/internal/enterror"
	"oa.98ent.com/p9/platform-operator/rpc/internal/svc"
	"oa.98ent.com/p9/platform-operator/rpc/pb/platformoperatorrpc/adminpb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLogic {
	return &GetLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取分站管理员
func (l *GetLogic) Get(in *adminpb.GetAdminRequest) (*adminpb.GetAdminResponse, error) {
	adminInfo, err := l.svcCtx.DB.OperatorAdmin.Get(l.ctx, in.Id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, xerr.RpcErr(xerr.NotFound(i18nkey.DataNotFound))
		}
		return nil, enterror.Handle(l.Logger, err)
	}

	return &adminpb.GetAdminResponse{
		Admin: &adminpb.AdminInfo{
			Id:          adminInfo.ID,
			Username:    adminInfo.Username,
			DisplayName: adminInfo.DisplayName,
			Status:      adminInfo.Status,
			OperatorId:  adminInfo.OperatorID,
			CreatedAt:   adminInfo.CreatedAt.Unix(),
			UpdatedAt:   adminInfo.UpdatedAt.Unix(),
		},
	}, nil
}
