package websocket

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"net/http"
	"strings"

	"oa.98ent.com/p9/node-dispatch/rpc/ent"
	"oa.98ent.com/p9/node-dispatch/rpc/ent/node"
)

const nodeCodeHeader = "X-Node-Code" // 节点编码请求头

// errUnauthorized 节点认证失败
var errUnauthorized = errors.New("unauthorized")

// authenticate 验证节点身份
func (s *Server) authenticate(r *http.Request) (*ent.Node, error) {
	// 获取节点编码
	nodeCode := strings.TrimSpace(r.Header.Get(nodeCodeHeader))
	if nodeCode == "" {
		return nil, errUnauthorized
	}

	// 获取节点认证密钥
	authSecret, ok := parseBearerToken(r.Header.Get("Authorization"))
	if !ok {
		return nil, errUnauthorized
	}

	// 查询节点
	data, err := s.svcCtx.DB.Node.Query().Where(node.CodeEQ(nodeCode)).Only(r.Context())
	if err != nil {
		// 节点不存在时统一返回认证失败
		if ent.IsNotFound(err) {
			return nil, errUnauthorized
		}

		return nil, err
	}

	// 停用节点不允许建立连接
	if data.Status != 1 {
		return nil, errUnauthorized
	}

	// 计算节点认证密钥哈希
	hash := sha256.Sum256([]byte(authSecret))
	authSecretHash := hex.EncodeToString(hash[:])

	// 比较节点认证密钥哈希
	if subtle.ConstantTimeCompare([]byte(authSecretHash), []byte(data.AuthSecretHash)) != 1 {
		return nil, errUnauthorized
	}

	return data, nil
}

// parseBearerToken 解析Bearer认证密钥
func parseBearerToken(value string) (string, bool) {
	// 拆分认证类型和认证密钥
	parts := strings.SplitN(strings.TrimSpace(value), " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return "", false
	}

	// 整理认证密钥
	token := strings.TrimSpace(parts[1])
	if token == "" {
		return "", false
	}

	return token, true
}
