package sync

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/gorm"
	"oa.98ent.com/p9/platform-game/rpc/pb/vendors"
)

// ===== 同步服务实现 =====

// SyncServiceImpl 同步服务实现
type SyncServiceImpl struct {
	db                  *gorm.DB
	grpcServerAddr      string
	categorySyncService *CategorySyncService
	providerSyncService *ProviderSyncService
	channelSyncService  *ChannelSyncService
	gameSyncService     *GameSyncService
	currencySyncService *CurrencySyncService
}

// NewSyncServiceImpl 创建同步服务实现
func NewSyncServiceImpl(db *gorm.DB, grpcServerAddr string) *SyncServiceImpl {
	return &SyncServiceImpl{
		db:                  db,
		grpcServerAddr:      grpcServerAddr,
		categorySyncService: NewCategorySyncService(db),
		providerSyncService: NewProviderSyncService(db),
		channelSyncService:  NewChannelSyncService(db),
		gameSyncService:     NewGameSyncService(db),
		currencySyncService: NewCurrencySyncService(db),
	}
}

// GetGRPCClient 获取 gRPC 客户端连接
func (s *SyncServiceImpl) GetGRPCClient(ctx context.Context) (vendors.VendorGameServiceClient, *grpc.ClientConn, error) {
	log.Printf("[GetGRPCClient] 🔍 开始连接 gRPC 服务\n")
	log.Printf("[GetGRPCClient] 📍 目标地址: %s\n", s.grpcServerAddr)

	if s.grpcServerAddr == "" {
		log.Println("[GetGRPCClient] ❌ gRPC 服务器地址为空!")
		return nil, nil, fmt.Errorf("gRPC 服务器地址未配置")
	}

	// 网络诊断：检查能否 DNS 解析
	log.Println("[GetGRPCClient] 🔧 执行网络诊断...")
	host, port, err := net.SplitHostPort(s.grpcServerAddr)
	if err != nil {
		log.Printf("[GetGRPCClient] ⚠️  地址格式解析失败: %v (将直接使用原地址)\n", err)
		host = s.grpcServerAddr
		port = "unknown"
	}

	log.Printf("[GetGRPCClient]    - 主机: %s\n", host)
	log.Printf("[GetGRPCClient]    - 端口: %s\n", port)

	// 尝试 DNS 解析
	log.Printf("[GetGRPCClient] 🔧 尝试 DNS 解析 '%s'...\n", host)
	ips, dnsErr := net.LookupIP(host)
	if dnsErr != nil {
		log.Printf("[GetGRPCClient] ⚠️  DNS 解析失败: %v (可能是本地地址或 127.0.0.1)\n", dnsErr)
	} else {
		log.Printf("[GetGRPCClient] ✓ DNS 解析成功，得到 IP 地址: %v\n", ips)
	}

	// 尝试 TCP 连接测试
	if port != "unknown" {
		log.Printf("[GetGRPCClient] 🔧 测试 TCP 连接到 %s:%s...\n", host, port)
		testConn, tcpErr := net.DialTimeout("tcp", s.grpcServerAddr, 3*time.Second)
		if tcpErr != nil {
			log.Printf("[GetGRPCClient] ⚠️  TCP 连接失败: %v\n", tcpErr)
			log.Printf("[GetGRPCClient]    💡 可能原因：\n")
			log.Printf("[GetGRPCClient]       1. game-vendor-sync 服务未启动\n")
			log.Printf("[GetGRPCClient]       2. 地址或端口配置错误\n")
			log.Printf("[GetGRPCClient]       3. 防火墙阻止连接\n")
		} else {
			log.Printf("[GetGRPCClient] ✓ TCP 连接测试成功\n")
			testConn.Close()
		}
	}

	// 创建 gRPC 连接到游戏供应商服务
	log.Println("[GetGRPCClient] 🔧 执行 grpc.DialContext...")
	log.Println("[GetGRPCClient]    - 超时: 30 秒")
	log.Println("[GetGRPCClient]    - TLS: 禁用（InsecureCredentials）")
	log.Println("[GetGRPCClient]    - 最大接收消息: 50MB")

	// 设置连接超时
	dialCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(
		dialCtx,
		s.grpcServerAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(50*1024*1024)),
		grpc.WithBlock(), // 阻塞模式，直到连接建立才返回
	)
	if err != nil {
		log.Printf("[GetGRPCClient] ❌ gRPC 连接失败!\n")
		log.Printf("[GetGRPCClient]    错误: %v\n", err)
		log.Printf("[GetGRPCClient]    错误类型: %T\n", err)
		log.Printf("[GetGRPCClient]    💡 调试建议：\n")
		log.Printf("[GetGRPCClient]       1. 检查 game-vendor-sync 是否正在运行\n")
		log.Printf("[GetGRPCClient]       2. 检查配置中 GrpcServerAddr 的值: %s\n", s.grpcServerAddr)
		log.Printf("[GetGRPCClient]       3. 尝试手动连接测试: telnet %s\n", s.grpcServerAddr)
		log.Printf("[GetGRPCClient]       4. 查看 game-vendor-sync 的启动日志\n")
		return nil, nil, fmt.Errorf("连接gRPC服务失败: %w", err)
	}

	log.Println("[GetGRPCClient] ✓ gRPC 连接建立成功")
	log.Printf("[GetGRPCClient]    连接状态: %v\n", conn.GetState())

	// 创建客户端
	client := vendors.NewVendorGameServiceClient(conn)
	log.Println("[GetGRPCClient] ✓ gRPC 客户端创建完成")
	log.Printf("[GetGRPCClient] 📍 连接已准备，可开始调用 RPC 方法\n")

	return client, conn, nil
}

// Preview 预检查同步
func (s *SyncServiceImpl) Preview(ctx context.Context, objectType string) (*vendors.SyncPreviewResp, error) {
	log.Printf("[Preview] 开始处理 %s 同步预检查\n", objectType)
	// 获取 gRPC 客户端
	// log.Println("[Preview] 正在获取 gRPC 客户端...")
	// client, conn, err := s.GetGRPCClient(ctx)
	// if err != nil {
	// 	log.Printf("[Preview] ✗ 获取 gRPC 客户端失败: %v\n", err)
	// 	return nil, err
	// }
	// defer conn.Close()
	// log.Println("[Preview] ✓ gRPC 客户端获取成功")
	var client vendors.VendorGameServiceClient = nil

	// 根据对象类型选择对应的同步服务进行预检查
	log.Printf("[Preview] 正在调用 %s 同步服务...\n", objectType)
	switch objectType {
	case "category":
		return s.categorySyncService.Preview(ctx, client)
	case "provider":
		return s.providerSyncService.Preview(ctx, client)
	case "channel":
		return s.channelSyncService.Preview(ctx, client)
	case "game":
		return s.gameSyncService.Preview(ctx, client)
	case "currency":
		return s.currencySyncService.Preview(ctx, client)
	default:
		return nil, fmt.Errorf("未知的对象类型: %s", objectType)
	}
}

// Run 执行同步操作
func (s *SyncServiceImpl) Run(ctx context.Context, objectType string) (*vendors.SyncRunResp, error) {
	// 获取 gRPC 客户端
	// client, conn, err := s.GetGRPCClient(ctx)
	// if err != nil {
	// 	return nil, err
	// }
	// defer conn.Close()
	var client vendors.VendorGameServiceClient = nil

	// 根据对象类型选择对应的同步服务进行同步
	switch objectType {
	case "category":
		return s.categorySyncService.Run(ctx, client)
	case "provider":
		return s.providerSyncService.Run(ctx, client)
	case "channel":
		return s.channelSyncService.Run(ctx, client)
	case "game":
		return s.gameSyncService.Run(ctx, client)
	case "currency":
		return s.currencySyncService.Run(ctx, client)
	default:
		return nil, fmt.Errorf("未知的对象类型: %s", objectType)
	}
}

// SyncAll 全量同步所有类型的数据
func (s *SyncServiceImpl) SyncAll(ctx context.Context) (*vendors.SyncRunResp, error) {
	// 创建聚合结果
	totalResult := &vendors.SyncRunResp{
		Preview: &vendors.SyncPreviewResp{
			Stats: &vendors.SyncStats{},
			Diffs: make([]*vendors.SyncDiff, 0),
		},
		Apply: &vendors.SyncApplyResult{},
	}

	// 按依赖顺序同步：分类 -> 提供商 -> 渠道 -> 游戏 -> 货币
	syncOrder := []string{"category", "provider", "channel", "game", "currency"}

	for _, objectType := range syncOrder {
		log.Printf("[全量同步] 开始同步%s", objectType)

		// 执行同步
		result, err := s.Run(ctx, objectType)
		if err != nil {
			log.Printf("[全量同步] 同步%s失败: %v", objectType, err)
			return nil, fmt.Errorf("同步%s失败: %w", objectType, err)
		}

		// 汇总统计信息
		if result.Preview != nil && result.Preview.Stats != nil {
			totalResult.Preview.Stats.RemoteTotal += result.Preview.Stats.RemoteTotal
			totalResult.Preview.Stats.LocalTotal += result.Preview.Stats.LocalTotal
			totalResult.Preview.Stats.CreateTotal += result.Preview.Stats.CreateTotal
			totalResult.Preview.Stats.UpdateTotal += result.Preview.Stats.UpdateTotal
			totalResult.Preview.Stats.DeleteTotal += result.Preview.Stats.DeleteTotal
			totalResult.Preview.Stats.NoopTotal += result.Preview.Stats.NoopTotal
			totalResult.Preview.Stats.ConflictTotal += result.Preview.Stats.ConflictTotal
			totalResult.Preview.Stats.ErrorTotal += result.Preview.Stats.ErrorTotal
		}

		// 汇总差异信息
		if result.Preview != nil && result.Preview.Diffs != nil {
			totalResult.Preview.Diffs = append(totalResult.Preview.Diffs, result.Preview.Diffs...)
		}

		// 汇总应用结果
		if result.Apply != nil {
			totalResult.Apply.Created += result.Apply.Created
			totalResult.Apply.Updated += result.Apply.Updated
			totalResult.Apply.Deleted += result.Apply.Deleted
			totalResult.Apply.Failed += result.Apply.Failed
			totalResult.Apply.Skipped += result.Apply.Skipped
		}

		log.Printf("[全量同步] 完成%s同步 (新增:%d, 更新:%d, 失败:%d)", objectType,
			result.Apply.Created, result.Apply.Updated, result.Apply.Failed)
	}

	return totalResult, nil
}
