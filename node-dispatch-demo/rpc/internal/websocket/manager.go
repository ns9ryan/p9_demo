package websocket

import (
	"sync"

	coderws "github.com/coder/websocket"
)

// Connection 表示一个已认证的节点WebSocket连接
type Connection struct {
	NodeCode string
	Conn     *coderws.Conn
}

// Manager 管理节点WebSocket连接
type Manager struct {
	mu          sync.RWMutex
	connections map[string]*Connection
}

// NewManager 创建节点连接管理器
func NewManager() *Manager {
	return &Manager{
		connections: make(map[string]*Connection),
	}
}

// Register 注册节点连接
func (m *Manager) Register(connection *Connection) {
	m.mu.Lock()

	oldConnection := m.connections[connection.NodeCode]
	m.connections[connection.NodeCode] = connection

	m.mu.Unlock()

	// 新连接替换旧连接
	if oldConnection != nil && oldConnection != connection {
		_ = oldConnection.Conn.Close(coderws.StatusNormalClosure, "connection replaced")
	}
}

// Unregister 注销节点连接
func (m *Manager) Unregister(connection *Connection) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 只注销当前连接, 避免旧连接退出时误删新连接
	current, ok := m.connections[connection.NodeCode]
	if ok && current == connection {
		delete(m.connections, connection.NodeCode)
	}
}

// Get 获取节点连接
func (m *Manager) Get(nodeCode string) (*Connection, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	connection, ok := m.connections[nodeCode]
	return connection, ok
}

// IsOnline 判断节点是否在线
func (m *Manager) IsOnline(nodeCode string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	_, ok := m.connections[nodeCode]
	return ok
}

// OnlineCodes 获取当前在线节点编码
func (m *Manager) OnlineCodes() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	codes := make([]string, 0, len(m.connections))
	for nodeCode := range m.connections {
		codes = append(codes, nodeCode)
	}

	return codes
}

// Disconnect 主动断开节点连接
func (m *Manager) Disconnect(nodeCode string) bool {
	m.mu.Lock()

	connection, ok := m.connections[nodeCode]
	if ok {
		delete(m.connections, nodeCode)
	}

	m.mu.Unlock()

	if !ok {
		return false
	}

	_ = connection.Conn.Close(coderws.StatusNormalClosure, "connection closed by server")
	return true
}
