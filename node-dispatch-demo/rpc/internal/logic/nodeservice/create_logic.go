package nodeservicelogic

import (
	"context"
	"strings"

	"oa.98ent.com/p9/node-dispatch/rpc/internal/svc"
	"oa.98ent.com/p9/node-dispatch/rpc/pb/nodedispatchrpc/nodepb"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
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

// Create 创建节点
func (l *CreateLogic) Create(in *nodepb.CreateNodeRequest) (*nodepb.CreateNodeResponse, error) {
	// 整理创建参数
	name := strings.TrimSpace(in.Name)
	remark := trimOptionalString(in.Remark)

	// 生成节点业务编码
	code := "NODE_" + strings.ToUpper(strings.ReplaceAll(uuid.NewString(), "-", ""))

	// 生成节点认证密钥
	authSecret, authSecretHash, err := generateAuthSecret()
	if err != nil {
		l.Logger.Errorw("生成节点认证密钥失败", logx.Field("error", err.Error()))
		return nil, err
	}

	// 创建节点
	data, err := l.svcCtx.DB.Node.
		Create().
		SetCode(code).                     // 节点业务编码
		SetName(name).                     // 节点名称
		SetAuthSecretHash(authSecretHash). // 节点认证密钥哈希
		SetNillableStatus(in.Status).      // 节点状态: 1启用, 2停用
		SetNillableRemark(remark).         // 运维备注
		Save(l.ctx)
	if err != nil {
		l.Logger.Errorw("创建节点失败", logx.Field("error", err.Error()))
		return nil, err
	}

	// 返回创建结果
	return &nodepb.CreateNodeResponse{
		Id:         data.ID,    // 节点ID
		Code:       data.Code,  // 节点业务编码
		AuthSecret: authSecret, // 节点认证密钥
	}, nil
}
