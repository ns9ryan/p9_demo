package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"oa.98ent.com/p9/node-agent/internal/config"
	"oa.98ent.com/p9/node-agent/internal/protocol"
	"oa.98ent.com/p9/node-agent/internal/task"

	coderws "github.com/coder/websocket"
	"github.com/zeromicro/go-zero/core/logx"
)

const (
	nodeCodeHeader    = "X-Node-Code"   // 节点编码请求头
	reconnectInterval = 5 * time.Second // 重连间隔
)

// Client 节点WebSocket客户端
type Client struct {
	webSocketURL string        // WebSocket连接地址
	nodeCode     string        // 节点业务编码
	authSecret   string        // 节点认证密钥
	task         *task.Service // 调度任务执行服务
}

// NewClient 创建节点WebSocket客户端
func NewClient(c config.Config, taskService *task.Service) *Client {
	return &Client{
		webSocketURL: c.Dispatch.WebSocketURL, // WebSocket连接地址
		nodeCode:     c.NodeCode,              // 节点业务编码
		authSecret:   c.AuthSecret,            // 节点认证密钥
		task:         taskService,             // 调度任务执行服务
	}
}

// Run 连接调度中心并自动重连
func (c *Client) Run(ctx context.Context) error {
	for {
		// 建立连接并持续读取消息
		err := c.runConnection(ctx)
		if ctx.Err() != nil {
			return nil
		}

		// 记录连接断开
		logx.WithContext(ctx).Errorw("节点WebSocket连接已断开, 等待重新连接", logx.Field("error", err.Error()), logx.Field("retry_after", reconnectInterval.String()))

		// 等待重新连接
		select {
		case <-ctx.Done():
			return nil

		case <-time.After(reconnectInterval):
		}
	}
}

// runConnection 建立WebSocket连接并持续读取消息
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

	logger := logx.WithContext(ctx)
	logger.Infow("节点WebSocket连接已建立", logx.Field("node_code", c.nodeCode))

	// 持续读取调度中心消息
	for {
		messageType, data, err := conn.Read(ctx)
		if err != nil {
			return fmt.Errorf("读取调度中心WebSocket消息失败: %w", err)
		}

		// 处理调度中心消息
		if err = c.handleMessage(ctx, conn, messageType, data); err != nil {
			logger.Errorw("处理调度中心WebSocket消息失败", logx.Field("error", err.Error()))
		}
	}
}

// sendMessage 发送WebSocket业务消息
func (c *Client) sendMessage(ctx context.Context, conn *coderws.Conn, message protocol.Message) error {
	// 编码业务消息
	data, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("编码WebSocket消息失败: %w", err)
	}

	// 发送文本消息
	if err = conn.Write(ctx, coderws.MessageText, data); err != nil {
		return fmt.Errorf("发送WebSocket消息失败: %w", err)
	}

	return nil
}
