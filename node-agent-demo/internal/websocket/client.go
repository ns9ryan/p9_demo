package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"oa.98ent.com/p9/node-agent/internal/config"

	coderws "github.com/coder/websocket"
	"github.com/zeromicro/go-zero/core/logx"
)

const (
	nodeCodeHeader    = "X-Node-Code"   // 节点编码请求头
	reconnectInterval = 5 * time.Second // 重连间隔
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

// handleMessage 处理调度中心WebSocket消息
func (c *Client) handleMessage(ctx context.Context, conn *coderws.Conn, messageType coderws.MessageType, data []byte) error {
	// 只处理文本消息
	if messageType != coderws.MessageText {
		return fmt.Errorf("不支持的WebSocket消息类型: %d", messageType)
	}

	// 解析业务消息
	var message Message
	if err := json.Unmarshal(data, &message); err != nil {
		return fmt.Errorf("解析WebSocket消息失败: %w", err)
	}

	// 根据消息类型处理
	switch message.Type {
	case MessageTypeTaskDispatch:
		return c.handleTaskDispatch(ctx, conn, message.Data)

	default:
		return fmt.Errorf("不支持的WebSocket业务消息类型: %s", message.Type)
	}
}

// handleTaskDispatch 处理任务下发消息
func (c *Client) handleTaskDispatch(ctx context.Context, conn *coderws.Conn, data json.RawMessage) error {
	// 解析任务下发数据
	var task TaskDispatchData
	if err := json.Unmarshal(data, &task); err != nil {
		return fmt.Errorf("解析任务下发数据失败: %w", err)
	}

	// 校验任务基本信息
	if task.TaskNo == "" {
		return fmt.Errorf("任务编号不能为空")
	}
	if task.TaskType == "" {
		return fmt.Errorf("任务类型不能为空")
	}

	logx.WithContext(ctx).Infow("收到调度任务", logx.Field("task_no", task.TaskNo), logx.Field("task_type", task.TaskType))

	// TODO 根据任务类型执行具体任务

	// 编码任务接收确认数据
	ackData, err := json.Marshal(TaskAckData{
		TaskNo: task.TaskNo, // 任务编号
	})
	if err != nil {
		return fmt.Errorf("编码任务接收确认数据失败: %w", err)
	}

	// 发送任务接收确认
	err = c.sendMessage(ctx, conn, Message{
		Type: MessageTypeTaskAck, // 消息类型
		Data: ackData,            // 确认数据
	})
	if err != nil {
		return fmt.Errorf("发送任务接收确认失败: %w", err)
	}

	logx.WithContext(ctx).Infow("调度任务已确认接收", logx.Field("task_no", task.TaskNo))

	return nil
}

// sendMessage 发送WebSocket业务消息
func (c *Client) sendMessage(ctx context.Context, conn *coderws.Conn, message Message) error {
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
