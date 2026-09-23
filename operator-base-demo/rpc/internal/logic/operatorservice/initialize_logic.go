package operatorservicelogic

import (
	"context"
	"slices"

	operatorservice "oa.98ent.com/p9/operator-base/rpc/internal/operator"
	"oa.98ent.com/p9/operator-base/rpc/internal/svc"
	"oa.98ent.com/p9/operator-base/rpc/pb/operatorbaserpc/operatorpb"

	"github.com/zeromicro/go-zero/core/logx"
)

type InitializeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewInitializeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InitializeLogic {
	return &InitializeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// Initialize 初始化厅 operator
func (l *InitializeLogic) Initialize(in *operatorpb.InitializeOperatorRequest) (*operatorpb.InitializeOperatorResponse, error) {
	// 转换域名初始化数据
	domains := make([]operatorservice.InitializeDomain, 0, len(in.GetDomains()))
	for _, item := range in.GetDomains() {
		domains = append(domains, operatorservice.InitializeDomain{
			DomainName: item.GetDomainName(), // 域名
			DomainType: item.GetDomainType(), // 域名类型
		})
	}

	// 获取 operator 基础信息
	operatorInfo := in.GetOperator()

	// 初始化 operator
	err := l.svcCtx.Operator.Initialize(l.ctx, operatorservice.InitializeRequest{
		Code:                   operatorInfo.GetCode(),                   // operator 全局唯一业务编码
		Name:                   operatorInfo.GetName(),                   // operator 名称
		TimezoneCode:           operatorInfo.GetTimezoneCode(),           // IANA 时区编码
		SettlementCurrencyCode: operatorInfo.GetSettlementCurrencyCode(), // 结算货币编码
		Status:                 operatorInfo.GetStatus(),                 // operator 状态
		Domains:                domains,                                  // 当前有效域名
		LanguageCodes:          slices.Clone(in.GetLanguageCodes()),      // 当前有效语言编码
		RegionCodes:            slices.Clone(in.GetRegionCodes()),        // 当前有效经营地区编码
		AgentLineCodes:         slices.Clone(in.GetAgentLineCodes()),     // 当前有效代理子线路编码
	})
	if err != nil {
		l.Logger.Errorw(
			"初始化厅 operator 失败",
			logx.Field("operator_code", operatorInfo.GetCode()),
			logx.Field("error", err.Error()),
		)
		return nil, err
	}

	// 返回初始化结果
	return &operatorpb.InitializeOperatorResponse{}, nil
}
