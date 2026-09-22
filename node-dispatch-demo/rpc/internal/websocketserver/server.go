package websocketserver

import (
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"net/http"
	"strings"
	"time"

	coderws "github.com/coder/websocket"
	"github.com/zeromicro/go-zero/core/logx"
	"oa.98ent.com/p9/node-dispatch/rpc/ent"
	"oa.98ent.com/p9/node-dispatch/rpc/ent/node"
	"oa.98ent.com/p9/node-dispatch/rpc/internal/connection"
	"oa.98ent.com/p9/node-dispatch/rpc/internal/svc"
)

const (
	nodeWebSocketPath = "/ws/node"    // 节点WebSocket访问路径
	nodeCodeHeader    = "X-Node-Code" // 节点编码请求头

	heartbeatInterval = 30 * time.Second // 心跳间隔
	heartbeatTimeout  = 10 * time.Second // 心跳超时时间
)

// errUnauthorized 节点认证失败
var errUnauthorized = errors.New("unauthorized")

// Server 提供节点WebSocket服务
type Server struct {
	httpServer *http.Server        // HTTP服务
	svcCtx     *svc.ServiceContext // 服务上下文
}

// NewServer 创建节点WebSocket服务
func NewServer(svcCtx *svc.ServiceContext) *Server {
	// 创建WebSocket服务
	server := &Server{
		svcCtx: svcCtx, // 服务上下文
	}

	// 注册节点WebSocket路由
	mux := http.NewServeMux()
	mux.HandleFunc(nodeWebSocketPath, server.handleNode)

	// 创建HTTP服务
	server.httpServer = &http.Server{
		Addr:              svcCtx.Config.WebSocketListenOn, // WebSocket监听地址
		Handler:           mux,                             // HTTP路由
		ReadHeaderTimeout: 5 * time.Second,                 // 请求头读取超时时间
	}

	// 返回WebSocket服务
	return server
}

// Start 启动节点WebSocket服务
func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

// Stop 停止节点WebSocket服务
func (s *Server) Stop() error {
	// 停止HTTP服务, 不再接收新的节点连接
	err := s.httpServer.Close()

	// 主动断开当前全部节点连接
	s.svcCtx.Connections.DisconnectAll()

	return err
}

// handleNode 处理节点WebSocket连接
func (s *Server) handleNode(w http.ResponseWriter, r *http.Request) {
	logger := logx.WithContext(r.Context())

	// 校验请求方法
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	// 验证节点身份
	data, err := s.authenticate(r)
	if err != nil {
		if errors.Is(err, errUnauthorized) {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		logger.Errorw("节点WebSocket认证失败", logx.Field("error", err.Error()))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// 升级为WebSocket连接
	conn, err := coderws.Accept(w, r, nil)
	if err != nil {
		logger.Errorw("节点WebSocket握手失败", logx.Field("node_code", data.Code), logx.Field("error", err.Error()))
		return
	}

	// 创建连接生命周期上下文
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 创建节点连接
	nodeConnection := &connection.Connection{
		NodeCode: data.Code, // 节点业务编码
		Conn:     conn,      // WebSocket连接
	}

	// 注册节点连接
	s.svcCtx.Connections.Register(nodeConnection)
	defer s.svcCtx.Connections.Unregister(nodeConnection)

	logger.Infow("节点WebSocket连接已建立", logx.Field("node_code", data.Code))

	// 更新节点最近活动时间
	s.updateLastSeenAt(ctx, data.ID, data.Code)

	// 启动节点心跳检测
	go s.runHeartbeat(ctx, conn, data.ID, data.Code)

	// 持续读取节点消息
	for {
		messageType, messageData, err := conn.Read(ctx)
		if err != nil {
			logger.Infow("节点WebSocket连接已断开", logx.Field("node_code", data.Code), logx.Field("error", err.Error()))
			return
		}

		// 处理节点上报消息
		if err = s.handleMessage(ctx, data.ID, data.Code, messageType, messageData); err != nil {
			logger.Errorw("处理节点WebSocket消息失败", logx.Field("node_code", data.Code), logx.Field("error", err.Error()))
		}
	}
}

// runHeartbeat 持续检测节点WebSocket连接
func (s *Server) runHeartbeat(ctx context.Context, conn *coderws.Conn, nodeID int64, nodeCode string) {
	logger := logx.WithContext(ctx)

	// 创建心跳定时器
	ticker := time.NewTicker(heartbeatInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return

		case <-ticker.C:
			// 设置本次心跳超时时间
			pingCtx, pingCancel := context.WithTimeout(ctx, heartbeatTimeout)
			err := conn.Ping(pingCtx)
			pingCancel()

			if err != nil {
				if ctx.Err() == nil {
					logger.Errorw("节点WebSocket心跳失败", logx.Field("node_code", nodeCode), logx.Field("error", err.Error()))
				}

				conn.CloseNow()
				return
			}

			// 更新节点最近活动时间
			s.updateLastSeenAt(ctx, nodeID, nodeCode)

			// 临时测试心跳是否正常
			// logger.Infow("节点WebSocket心跳正常", logx.Field("node_code", nodeCode))
		}
	}
}

// updateLastSeenAt 更新节点最近活动时间
func (s *Server) updateLastSeenAt(ctx context.Context, nodeID int64, nodeCode string) {
	// 更新节点最近活动时间
	_, err := s.svcCtx.DB.Node.
		Update().
		Where(node.IDEQ(nodeID)).
		SetLastSeenAt(time.Now()).
		Save(ctx)
	if err != nil {
		logx.WithContext(ctx).Errorw("更新节点最近活动时间失败", logx.Field("node_code", nodeCode), logx.Field("error", err.Error()))
	}
}

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
	data, err := s.svcCtx.DB.Node.
		Query().
		Where(node.CodeEQ(nodeCode)).
		Only(r.Context())
	if err != nil {
		// 节点不存在时统一返回认证失败
		if ent.IsNotFound(err) {
			return nil, errUnauthorized
		}

		return nil, err
	}

	// 计算节点认证密钥哈希
	hash := sha256.Sum256([]byte(authSecret))
	authSecretHash := hex.EncodeToString(hash[:])

	// 比较节点认证密钥哈希
	if subtle.ConstantTimeCompare([]byte(authSecretHash), []byte(data.AuthSecretHash)) != 1 {
		return nil, errUnauthorized
	}

	// 返回认证通过的节点
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

	// 返回认证密钥
	return token, true
}
