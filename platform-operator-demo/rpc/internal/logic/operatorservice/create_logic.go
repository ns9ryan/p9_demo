package operatorservicelogic

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"

	"oa.98ent.com/p9/platform-operator/rpc/internal/enterror"
	"oa.98ent.com/p9/platform-operator/rpc/internal/svc"
	"oa.98ent.com/p9/platform-operator/rpc/pb/operator/operator"
)

type CreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
	return &CreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// Create 创建分站
func (l *CreateLogic) Create(in *operator.CreateOperatorRequest) (*operator.CreateOperatorResponse, error) {
	// 生成分站业务编码
	code := "OP_" + strings.ToUpper(strings.ReplaceAll(uuid.NewString(), "-", ""))

	// 创建分站
	data, err := l.svcCtx.DB.Operator.
		Create().
		SetCode(code).                                        // 分站业务编码
		SetName(in.Name).                                     // 分站名称
		SetTimezoneCode(in.TimezoneCode).                     // 时区编码
		SetSettlementCurrencyCode(in.SettlementCurrencyCode). // 结算币种编码
		SetNillableStatus(in.Status).                         // 分站状态
		SetNillableRemark(in.Remark).                         // 内部备注
		Save(l.ctx)
	if err != nil {
		// 处理数据库错误
		return nil, enterror.Handle(l.Logger, err)
	}

	// 返回创建结果
	return &operator.CreateOperatorResponse{
		Id:   data.ID,   // 分站ID
		Code: data.Code, // 分站业务编码
	}, nil
}
