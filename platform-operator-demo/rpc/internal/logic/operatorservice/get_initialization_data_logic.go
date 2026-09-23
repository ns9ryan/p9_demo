package operatorservicelogic

import (
	"context"

	"oa.98ent.com/p9/platform-operator/rpc/internal/svc"
	"oa.98ent.com/p9/platform-operator/rpc/pb/platformoperatorrpc/operatorpb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetInitializationDataLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetInitializationDataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetInitializationDataLogic {
	return &GetInitializationDataLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取分站初始化数据
func (l *GetInitializationDataLogic) GetInitializationData(in *operatorpb.GetInitializationDataRequest) (*operatorpb.GetInitializationDataResponse, error) {
	// todo: add your logic here and delete this line

	return &operatorpb.GetInitializationDataResponse{}, nil
}
