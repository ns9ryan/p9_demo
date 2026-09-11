package sync

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/gorm"
	"oa.98ent.com/p9/platform-game/rpc/pb/platform_game"
	"oa.98ent.com/p9/platform-game/rpc/pb/vendors"
)

const BatchSize = 10

// ===== 同步服务实现 =====

// SyncServiceImpl 同步服务实现
type SyncServiceImpl struct {
	mu                  sync.Mutex
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
func (s *SyncServiceImpl) Preview(ctx context.Context, req *platform_game.SyncPreviewRequest) (*platform_game.SyncPreviewResp, error) {
	objectType := req.ObjectType
	log.Printf("[Preview] 开始处理 %s 同步预检查\n", objectType)

	var client vendors.VendorGameServiceClient = nil
	isGetRemoteClient := false
	// grpcServerAddr为空时同步本地数据
	if s.grpcServerAddr != "" {
		// 获取 gRPC 客户端
		log.Println("[Preview] 正在获取 gRPC 客户端...")
		cli, conn, err := s.GetGRPCClient(ctx)
		if err != nil {
			log.Printf("[Preview] ✗ 获取 gRPC 客户端失败: %v\n", err)
			return nil, err
		}
		defer conn.Close()
		client = cli
		isGetRemoteClient = true
		log.Println("[Preview] ✓ gRPC 客户端获取成功")
	}

	// 根据对象类型选择对应的同步服务进行预检查
	log.Printf("[Preview] 正在调用 %s 同步服务...\n", objectType)
	switch objectType {
	case "category":
		previewResp, _, err := s.categorySyncService.Preview(ctx, client, req, isGetRemoteClient)
		log.Printf("[Preview] [游戏分类同步] 预检查结果数量: %d", len(previewResp.Diffs))
		return previewResp, err
	case "provider":
		previewResp, _, err := s.providerSyncService.Preview(ctx, client, req, isGetRemoteClient)
		log.Printf("[Preview] [游戏提供商同步] 预检查结果数量: %d", len(previewResp.Diffs))
		return previewResp, err
	case "channel":
		previewResp, _, err := s.channelSyncService.Preview(ctx, client, req, isGetRemoteClient)
		log.Printf("[Preview] [游戏渠道同步] 预检查结果数量: %d", len(previewResp.Diffs))
		return previewResp, err
	case "game":
		previewResp, _, err := s.gameSyncService.Preview(ctx, client, req, isGetRemoteClient)
		log.Printf("[Preview] [游戏同步] 预检查结果数量: %d", len(previewResp.Diffs))
		return previewResp, err
	case "currency":
		previewResp, _, err := s.currencySyncService.Preview(ctx, client, req, isGetRemoteClient)
		log.Printf("[Preview] [游戏货币同步] 预检查结果数量: %d", len(previewResp.Diffs))
		return previewResp, err
	default:
		return nil, fmt.Errorf("未知的对象类型: %s", objectType)
	}
}

// Run 执行同步操作
func (s *SyncServiceImpl) Run(ctx context.Context, objectType string, syncCols []string, localCheckpointID int64) (*platform_game.SyncRunResp, error) {
	// 加锁，确保同一时间仅有一个协程执行同步
	s.mu.Lock()
	defer s.mu.Unlock()
	var client vendors.VendorGameServiceClient = nil
	isGetRemoteClient := false
	// grpcServerAddr为空时同步本地数据
	if s.grpcServerAddr != "" {
		// 获取 gRPC 客户端
		log.Println("[Run] 正在获取 gRPC 客户端...")
		cli, conn, err := s.GetGRPCClient(ctx)
		if err != nil {
			log.Printf("[Run] ✗ 获取 gRPC 客户端失败: %v\n", err)
			return nil, err
		}
		defer conn.Close()
		client = cli
		isGetRemoteClient = true
		log.Println("[Run] ✓ gRPC 客户端获取成功")
	}
	// var client vendors.VendorGameServiceClient = nil

	// 根据对象类型选择对应的同步服务进行同步
	switch objectType {
	case "category":
		err := s.categorySyncService.Run(ctx, client, syncCols, localCheckpointID, isGetRemoteClient)
		if err != nil {
			return nil, err
		}
	case "provider":
		err := s.providerSyncService.Run(ctx, client, syncCols, localCheckpointID, isGetRemoteClient)
		if err != nil {
			return nil, err
		}
	case "channel":
		err := s.channelSyncService.Run(ctx, client, syncCols, localCheckpointID, isGetRemoteClient)
		if err != nil {
			return nil, err
		}
	case "game":
		err := s.gameSyncService.Run(ctx, client, syncCols, localCheckpointID, isGetRemoteClient)
		if err != nil {
			return nil, err
		}
	case "currency":
		err := s.currencySyncService.Run(ctx, client, syncCols, localCheckpointID, isGetRemoteClient)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("未知的对象类型: %s", objectType)
	}

	return &platform_game.SyncRunResp{CheckpointId: localCheckpointID}, nil
}

// SyncPageHandler 分页处理器
type SyncPageHandler struct {
	RemoteData  []interface{}
	LocalIndex  map[string]interface{}
	CompareFunc func(remoteItem interface{}, localIndex map[string]interface{}) *platform_game.SyncDiff
	PageSize    int64
	Page        int64
	SkipNoop    bool
}

// Handle 处理分页逻辑并返回预检查结果
func (h *SyncPageHandler) Handle(totalRemote, totalLocal int64) *platform_game.SyncPreviewResp {
	result := &platform_game.SyncPreviewResp{
		Stats: &platform_game.SyncStats{
			RemoteTotal: int32(totalRemote),
			LocalTotal:  int32(totalLocal),
		},
		Diffs: make([]*platform_game.SyncDiff, 0),
	}

	count := 0
	// 当 Page=0 时，返回全量结果（不做分页）
	if h.Page == 0 {
		for _, item := range h.RemoteData {
			diff := h.CompareFunc(item, h.LocalIndex)

			// 统计各类型操作
			h.recordStats(result.Stats, diff)

			// 如果启用了跳过 noop，则不添加
			if h.SkipNoop && diff.Action == "noop" {
				continue
			}
			count++

			result.Diffs = append(result.Diffs, diff)
		}
		result.Stats.DiffTotal = int32(count)
		return result
	}

	// 计算分页范围
	limit := int(h.PageSize * h.Page)
	offset := limit - int(h.PageSize)

	// 遍历远程数据
	for _, item := range h.RemoteData {
		diff := h.CompareFunc(item, h.LocalIndex)
		// 统计各类型操作
		h.recordStats(result.Stats, diff)

		// 如果启用了跳过 noop，则不计数
		if h.SkipNoop && diff.Action == "noop" {
			continue
		}

		count++

		// 在分页范围内添加差异项
		if count > offset && count <= limit {
			result.Diffs = append(result.Diffs, diff)
		}
	}
	result.Stats.DiffTotal = int32(count)
	return result
}

// recordStats 记录统计信息
func (h *SyncPageHandler) recordStats(stats *platform_game.SyncStats, diff *platform_game.SyncDiff) {
	switch diff.Action {
	case "create":
		stats.CreateTotal++
	case "update":
		stats.UpdateTotal++
	case "noop":
		stats.NoopTotal++
	}
}

func getValidSyncCols(localCols []string, syncCols []string) []string {
	validCols := make([]string, 0)
	for _, col := range localCols {
		for _, c := range syncCols {
			if c == col {
				validCols = append(validCols, col)
				break
			}
		}
	}
	return validCols
}

func GetUpdatedData(allColsData map[string]interface{}, syncCols []string) map[string]interface{} {
	localCols := make([]string, 0, len(allColsData))
	for col := range allColsData {
		localCols = append(localCols, col)
	}
	validCols := getValidSyncCols(localCols, syncCols)
	fmt.Println("validCols ===", validCols)
	updatedData := make(map[string]interface{})
	for _, col := range validCols {
		updatedData[col] = allColsData[col]
	}
	return updatedData
}

type Apply struct {
	Created int32
	Updated int32
	Deleted int32
	Failed  int32
	Skipped int32
}
