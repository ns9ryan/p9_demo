package websocket

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"oa.98ent.com/p9/node-agent/internal/config"

	coderws "github.com/coder/websocket"
	"github.com/zeromicro/go-zero/core/logx"
)

const (
	nodeCodeHeader    = "X-Node-Code"   // 节点编码请求头
	reconnectInterval = 5 * time.Second // WebSocket重连间隔
)

// Client 节点WebSocket客户端
type Client struct {
	webSocketURL string // WebSocket连接地址
	nodeCode     string // 节点业务编码
	authSecret   string // 节点认证密钥
}

// NewClient 创建节点WebSocket客户端
func NewClient(c config.Config) *Client {
	return &Client{
		webSocketURL: c.Dispatch.WebSocketURL, // WebSocket连接地址
		nodeCode:     c.NodeCode,              // 节点业务编码
		authSecret:   c.AuthSecret,            // 节点认证密钥
	}
}

// Run 连接调度中心并保持WebSocket连接
func (c *Client) Run(ctx context.Context) error {
	logger := logx.WithContext(ctx)

	for {
		// 建立并保持一次WebSocket连接
		err := c.runConnection(ctx)
		if ctx.Err() != nil {
			return nil
		}

		logger.Errorw("节点WebSocket连接已断开, 等待重新连接", logx.Field("error", err.Error()), logx.Field("retry_after", reconnectInterval.String()))

		// 等待重新连接
		select {
		case <-ctx.Done():
			return nil
		case <-time.After(reconnectInterval):
		}
	}
}

// runConnection 建立一次WebSocket连接并持续读取消息
func (c *Client) runConnection(ctx context.Context) error {
	// 构造节点认证请求头
	header := make(http.Header)
	header.Set(nodeCodeHeader, c.nodeCode)
	header.Set("Authorization", "Bearer "+c.authSecret)

	// 连接调度中心
	conn, response, err := coderws.Dial(ctx, c.webSocketURL, &coderws.DialOptions{
		HTTPHeader: header,
	})
	if err != nil {
		if response != nil {
			return fmt.Errorf("连接调度中心WebSocket失败, HTTP状态码: %d: %w", response.StatusCode, err)
		}

		return fmt.Errorf("连接调度中心WebSocket失败: %w", err)
	}
	defer conn.CloseNow()

	logx.WithContext(ctx).Infow("节点WebSocket连接已建立", logx.Field("node_code", c.nodeCode))

	// 持续读取调度中心消息
	for {
		_, _, err = conn.Read(ctx)
		if err != nil {
			return fmt.Errorf("读取调度中心WebSocket消息失败: %w", err)
		}

		// TODO Node Agent消息协议完成后处理调度中心下发消息
	}
}
