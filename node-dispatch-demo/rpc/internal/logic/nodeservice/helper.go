package nodeservicelogic

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"strings"

	"oa.98ent.com/p9/node-dispatch/rpc/ent"
	"oa.98ent.com/p9/node-dispatch/rpc/pb/nodedispatchrpc/nodepb"
)

// generateAuthSecret 生成节点认证密钥及其哈希
func generateAuthSecret() (string, string, error) {
	// 生成32字节安全随机数据
	secretBytes := make([]byte, 32)
	if _, err := rand.Read(secretBytes); err != nil {
		return "", "", err
	}

	// 转换为64位十六进制认证密钥
	authSecret := hex.EncodeToString(secretBytes)

	// 生成认证密钥哈希
	hash := sha256.Sum256([]byte(authSecret))
	authSecretHash := hex.EncodeToString(hash[:])

	return authSecret, authSecretHash, nil
}

// trimOptionalString 整理可选字符串
func trimOptionalString(value *string) *string {
	if value == nil {
		return nil
	}

	return new(strings.TrimSpace(*value))
}

// toNodeInfo 转换节点信息
func toNodeInfo(data *ent.Node, online bool) *nodepb.NodeInfo {
	// 转换最近一次活动时间
	var lastSeenAt *int64
	if data.LastSeenAt != nil {
		value := data.LastSeenAt.UnixMilli()
		lastSeenAt = &value
	}

	return &nodepb.NodeInfo{
		Id:         data.ID,                    // 节点ID
		Code:       data.Code,                  // 节点业务编码
		Name:       data.Name,                  // 节点名称
		Status:     data.Status,                // 节点状态: 1启用, 2停用
		Online:     online,                     // 是否在线
		LastSeenAt: lastSeenAt,                 // 最近一次活动时间, Unix毫秒时间戳
		Remark:     data.Remark,                // 运维备注
		CreatedAt:  data.CreatedAt.UnixMilli(), // 创建时间, Unix毫秒时间戳
		UpdatedAt:  data.UpdatedAt.UnixMilli(), // 更新时间, Unix毫秒时间戳
	}
}
