package game_sync

import (
	"context"
	"fmt"
	"log"
	"time"

	"oa.98ent.com/p9/platform-game/common/logger"
	"oa.98ent.com/p9/platform-game/rpc/ent"
	"oa.98ent.com/p9/platform-game/rpc/ent/gamechannel"
	"oa.98ent.com/p9/platform-game/rpc/internal/constant"
	"oa.98ent.com/p9/platform-game/rpc/internal/dao"
	"oa.98ent.com/p9/platform-game/rpc/pb/platform_game"
	platform_game_sync "oa.98ent.com/p9/platform-game/rpc/pb/platform_game"
	"oa.98ent.com/p9/platform-game/rpc/pb/vendors"
)

// ===== 渠道同步服务 =====

// ChannelSyncService 渠道同步服务
type ChannelSyncService struct {
	DAOManager    *dao.Manager
	progressQueue *ProgressQueue
}

// NewChannelSyncService 创建渠道同步服务
func NewChannelSyncService(daoManager *dao.Manager) *ChannelSyncService {
	return &ChannelSyncService{
		DAOManager:    daoManager,
		progressQueue: NewProgressQueue(daoManager),
	}
}

// Preview 预检查渠道同步
func (s *ChannelSyncService) Preview(ctx context.Context, client vendors.VendorGameServiceClient, req *platform_game_sync.SyncPreviewRequest, isGetRemoteClient bool) (*platform_game.SyncPreviewResp, *RemoteChannelResponse, error) {
	// 获取远程渠道数据
	remoteResp, err := getRemoteGameChannel(ctx, client, isGetRemoteClient)
	if err != nil {
		log.Printf("[渠道预检查] 获取测试数据失败: %v", err)
		return nil, nil, err
	}

	// 获取本地渠道数据
	local, localIndex, err := s.fetchLocal(ctx)
	if err != nil {
		return nil, nil, err
	}

	// 转换为接口切片
	remoteItems := make([]interface{}, len(remoteResp.Data))
	for i, v := range remoteResp.Data {
		remoteItems[i] = v
	}

	// 转换本地索引为接口类型
	localIndexInterface := make(map[string]interface{})
	for k, v := range localIndex {
		localIndexInterface[k] = v
	}

	// 当 req 为 nil 时，返回全量结果不做分页
	if req == nil {
		handler := &SyncPageHandler{
			RemoteData: remoteItems,
			LocalIndex: localIndexInterface,
			CompareFunc: func(item interface{}, localIndex map[string]interface{}) *platform_game.SyncDiff {
				remoteChannel := item.(*vendors.GameChannelInfo)
				localIdx := make(map[string]*ent.GameChannel)
				for k, v := range localIndex {
					localIdx[k] = v.(*ent.GameChannel)
				}
				return s.compareAll(remoteChannel, localIdx)
			},
			PageSize: 0,
			Page:     0,
			SkipNoop: false,
		}
		result := handler.Handle(int64(len(remoteResp.Data)), int64(len(local)))
		return result, remoteResp, nil
	}

	// 创建分页处理器
	handler := &SyncPageHandler{
		RemoteData: remoteItems,
		LocalIndex: localIndexInterface,
		CompareFunc: func(item interface{}, localIndex map[string]interface{}) *platform_game.SyncDiff {
			remoteChannel := item.(*vendors.GameChannelInfo)
			localIdx := make(map[string]*ent.GameChannel)
			for k, v := range localIndex {
				localIdx[k] = v.(*ent.GameChannel)
			}
			return s.compareAll(remoteChannel, localIdx)
		},
		PageSize: req.PageSize,
		Page:     req.Page,
		SkipNoop: req.IsSkip,
	}

	// 处理分页逻辑
	result := handler.Handle(int64(len(remoteResp.Data)), int64(len(local)))
	return result, remoteResp, nil
}

// Run 执行渠道同步
// 处理渠道数据的创建和更新（从 Run 提取的业务逻辑）
func (s *ChannelSyncService) processChannelDataWithProgress(ctx context.Context, tx *ent.Tx, remoteData []*vendors.GameChannelInfo, applyResult *Apply, syncCols []string, checkpointID int64) error {
	i := 0
	counterMap := make(map[string]int) // 用于记录每个 source_channel_code 的计数，确保唯一性
	nameMap := make(map[string]string)
	for _, remoteChannel := range remoteData {
		if counterMap[remoteChannel.Code] > 0 {
			logger.Errorf("[渠道同步] 检测到重复的远程渠道编码: %s, 计数器: %d, 跳过处理", remoteChannel.Code, counterMap[remoteChannel.Code])
			applyResult.Failed++
			continue
		}
		nameMap[remoteChannel.Code] = remoteChannel.Name
		// 更新计数器
		counterMap[remoteChannel.Code]++
		i++
		gameChannelRecord, err := tx.GameChannel.Query().
			Where(gamechannel.SourceChannelCodeEQ(remoteChannel.Code), gamechannel.DeletedAtIsNil()).
			Only(ctx)
		channelCode := remoteChannel.Code + "_" + fmt.Sprintf("%d", i)
		if err != nil || gameChannelRecord == nil {
			_, err := tx.GameChannel.Create().
				SetSourceID(remoteChannel.Id).
				SetChannelCode(channelCode).
				SetSourceChannelCode(remoteChannel.Code).
				SetSortNo(remoteChannel.Id).
				SetSourceSortNo(0).
				SetStatus(int64(remoteChannel.Status)).
				SetSourceStatus(int64(remoteChannel.Status)).
				SetSourceLoadType(int64(remoteChannel.LoadType)).
				SetLoadType(int64(remoteChannel.LoadType)).
				SetCreatedAt(time.Now()).
				SetUpdatedAt(time.Now()).
				Save(ctx)
			if err != nil {
				logger.Errorf("[渠道新增] ent 失败: %v", err)
				applyResult.Failed++
				continue
			}
			applyResult.Created++
		} else {
			if s.compare(&platform_game.SyncDiff{}, remoteChannel, gameChannelRecord) {
				logger.Infof("[渠道同步] 检测到需要更新的渠道: %s", remoteChannel.Code)
				update := tx.GameChannel.
					UpdateOneID(gameChannelRecord.ID).
					SetUpdatedAt(time.Now())

				if len(syncCols) == 0 {
					syncCols = append(syncCols, "channel_code")
					syncCols = append(syncCols, "source_channel_code")
					syncCols = append(syncCols, "source_status")
					syncCols = append(syncCols, "source_sort_no")
				}

				for _, col := range syncCols {
					switch col {
					case "channel_code":
						update.SetChannelCode(channelCode)
					case "source_channel_code":
						update.SetSourceChannelCode(remoteChannel.Code)
					case "status":
						update.SetStatus(int64(remoteChannel.Status))
					case "source_status":
						update.SetSourceStatus(int64(remoteChannel.Status))
					case "sort_no":
						update.SetSortNo(remoteChannel.Id)
					case "source_load_type":
						update.SetSourceLoadType(int64(remoteChannel.LoadType))
					case "load_type":
						update.SetLoadType(int64(remoteChannel.LoadType))
					case "source_sort_no":
						update.SetSourceSortNo(0)
					}
				}
				_, err := update.Save(ctx)
				if err != nil {
					logger.Errorf("[渠道更新] ent 失败: %v", err)
					applyResult.Failed++
					continue
				}
				applyResult.Updated++
				logger.Infof("[渠道更新] ✓ 已更新渠道: %s", remoteChannel.Code)
			}

		}
	}
	if err := mergeLocalJSON("game_channel_name_map.json", nameMap, constant.ChannelBiz); err != nil {
		logger.Errorf("[渠道同步] 保存名称映射失败: %v", err)
		// 不返回错误，继续执行后续逻辑
	}
	return nil
}

func (s *ChannelSyncService) Run(ctx context.Context, client vendors.VendorGameServiceClient, syncCols []string, checkpointID int64, isGetRemoteClient bool) error {
	logger.Infof("[渠道同步] ===== 开始执行同步 =====")

	// 先执行预检查
	logger.Infof("[渠道同步] 执行预检查")
	previewResp, remoteResp, err := s.Preview(ctx, client, nil, isGetRemoteClient)
	if err != nil {
		logger.Errorf("[渠道同步] 预检查失败: %v", err)
		return err
	}
	// 创建并启动进度队列
	s.progressQueue.Start(ctx)
	if previewResp.Stats.UpdateTotal == 0 && previewResp.Stats.CreateTotal == 0 && previewResp.Stats.DeleteTotal == 0 && previewResp.Stats.RemoteTotal == previewResp.Stats.LocalTotal && len(syncCols) == 0 {
		defer s.progressQueue.Stop()
		logger.Infof("[渠道同步] 无需更新，直接返回")
		logger.Infof("[渠道同步] UpdateTotal=%d, CreateTotal=%d, DeleteTotal=%d, RemoteTotal=%d, LocalTotal=%d",
			previewResp.Stats.UpdateTotal, previewResp.Stats.CreateTotal, previewResp.Stats.DeleteTotal, previewResp.Stats.RemoteTotal, previewResp.Stats.LocalTotal)
		// 更新checkpointID进度为100
		s.progressQueue.Send(&ProgressMessage{
			TableName:      "channel",
			ProcessedCount: previewResp.Stats.RemoteTotal,
			RemoteTotal:    previewResp.Stats.RemoteTotal,
			LocalTotal:     previewResp.Stats.LocalTotal,
			Progress:       100,
			CheckpointID:   checkpointID,
		})
		return nil
	}
	logger.Infof("[渠道同步] 预检查完成: remote_total=%d, local_total=%d",
		previewResp.Stats.RemoteTotal, previewResp.Stats.LocalTotal)

	// 获取远程渠道数据
	logger.Infof("[渠道同步] 获取远程渠道数据")
	logger.Infof("[渠道同步] 获取到远程数据: count=%d", len(remoteResp.Data))

	// 创建初始检查点记录（此时progress=0）
	logger.Infof("[渠道同步] 创建初始检查点")
	logger.Infof("[渠道同步] ✓ 检查点已创建: checkpointID=%d", checkpointID)

	go func() {
		defer s.progressQueue.Stop()
		// 在事务中执行数据库操作
		logger.Infof("[渠道同步] 开始处理渠道数据")
		// 把remoteResp.Data拆分为100条一批进行处理（可根据实际情况调整批次大小）
		batchSize := BatchSize
		totalCount := len(remoteResp.Data)
		apply := &Apply{}
		for i := 0; i < len(remoteResp.Data); i += batchSize {
			end := i + batchSize
			if end > len(remoteResp.Data) {
				end = len(remoteResp.Data)
			}
			batch := remoteResp.Data[i:end]
			logger.Infof("[渠道同步] 处理批次: start=%d, end=%d", i, end)
			// 模拟等待
			logger.Debugf("[进度队列] 模拟处理延迟: 表=channel, checkpointID=%d", checkpointID)
			time.Sleep(1000 * time.Millisecond)
			logger.Debugf("[进度队列] 开始处理消息: 表=channel, checkpointID=%d", checkpointID)
			txn, err := s.DAOManager.DB.Tx(ctx)
			if err != nil {
				logger.Errorf("[渠道同步] 创建事务失败: %v", err)
				return
			}
			err = s.processChannelDataWithProgress(ctx, txn, batch, apply, syncCols, checkpointID)
			if err != nil {
				logger.Errorf("[渠道同步] 处理渠道数据失败: %v", err)
				if rerr := txn.Rollback(); rerr != nil {
					logger.Errorf("[渠道同步] 回滚事务失败: %v", rerr)
				}
				return
			}
			if err := txn.Commit(); err != nil {
				logger.Errorf("[渠道同步] 提交事务失败: %v", err)
				return
			}
			progress := CalculateProgress(int64(i+1), int64(totalCount))
			msg := &ProgressMessage{
				TableName:      "channel",
				ProcessedCount: int32(i + 1),
				RemoteTotal:    previewResp.Stats.RemoteTotal,
				LocalTotal:     previewResp.Stats.LocalTotal,
				Progress:       progress,
				CheckpointID:   checkpointID,
				Created:        apply.Created,
				Updated:        apply.Updated,
				Deleted:        apply.Deleted,
				Failed:         apply.Failed,
				Skipped:        apply.Skipped,
			}
			logger.Infof("[渠道同步] 发送进度消息: 表=channel, 处理数=%d, 总数=%d, 进度=%d%%", i+1, totalCount, progress)
			s.progressQueue.Send(msg)
		}
		if err != nil {
			logger.Errorf("[渠道同步] 处理渠道数据失败: %v", err)
			return
		}

		// 执行逆向同步
		logger.Infof("[渠道同步] 开始执行逆向同步")
		err := s.ReverseSync(ctx, client, remoteResp, apply)
		if err != nil {
			logger.Errorf("❌ [渠道同步] 逆向同步失败: %v", err)
			return
		}
		logger.Infof("[渠道同步] 逆向同步完成: deleted=%d", apply.Deleted)

		// 最后一次进度更新（100%）
		logger.Infof("[渠道同步] 发送最终进度消息（100%%）")
		finalMsg := &ProgressMessage{
			TableName:      "channel",
			ProcessedCount: previewResp.Stats.RemoteTotal,
			RemoteTotal:    previewResp.Stats.RemoteTotal,
			LocalTotal:     previewResp.Stats.LocalTotal,
			Progress:       100,
			CheckpointID:   checkpointID,
			Created:        apply.Created,
			Updated:        apply.Updated,
			Deleted:        apply.Deleted,
			Failed:         apply.Failed,
			Skipped:        apply.Skipped,
		}
		s.progressQueue.Send(finalMsg)

		logger.Infof("[渠道同步] ===== 同步完成 =====")
	}()
	return nil
}

// 获取本地所有渠道
func (s *ChannelSyncService) fetchLocal(ctx context.Context) ([]*ent.GameChannel, map[string]*ent.GameChannel, error) {
	channels, err := s.DAOManager.GameChannel.GetAllGameChannel(ctx)
	if err != nil {
		return nil, nil, err
	}

	// 构建本地渠道索引（按上游渠道编码）
	index := make(map[string]*ent.GameChannel)
	for _, channel := range channels {
		if channel.SourceChannelCode != "" {
			index[channel.SourceChannelCode] = channel
		}
	}

	return channels, index, nil
}

// 比较单个渠道的本地和远程数据
func (s *ChannelSyncService) compareAll(remote *vendors.GameChannelInfo, localIndex map[string]*ent.GameChannel) *platform_game.SyncDiff {
	local, exists := localIndex[remote.Code]

	// 初始化差异记录
	diff := &platform_game.SyncDiff{
		ObjectType: "channel",
		ObjectId:   remote.Id,
		ObjectCode: remote.Code,
		RemoteId:   remote.Id,
		RemoteCode: remote.Code,
	}

	// 如果本地不存在该渠道
	if !exists {
		diff.Action = constant.SyncActionCreate
		diff.Reason = "本地不存在该渠道"
		return diff
	}

	// 比较本地和远程数据
	s.compare(diff, remote, local)

	return diff
}

func (s *ChannelSyncService) compare(diff *platform_game.SyncDiff, remote *vendors.GameChannelInfo, local *ent.GameChannel) bool {
	if local.SourceStatus != int64(remote.Status) {
		diff.Action = constant.SyncActionUpdate
		diff.Reason = fmt.Sprintf("渠道状态变更: %d -> %d", local.SourceStatus, remote.Status)
		diff.ConflictType = "status_changed"
		return true
	}
	if local.SourceLoadType != int64(remote.LoadType) {
		diff.Action = constant.SyncActionUpdate
		diff.Reason = fmt.Sprintf("渠道加载类型变更: %d -> %d", local.SourceLoadType, remote.LoadType)
		diff.ConflictType = "load_type_changed"
		return true
	}
	diff.Action = constant.SyncActionNoop
	diff.Reason = "本地和远程数据一致"
	return false
}

// ReverseSync 逆向同步：检查本地数据在远程是否存在，不存在则软删除
func (s *ChannelSyncService) ReverseSync(ctx context.Context, client vendors.VendorGameServiceClient, remoteResp *RemoteChannelResponse, apply *Apply) error {
	log.Println("[渠道逆向同步] 开始执行逆向同步...")

	// 构建远程渠道编码索引
	remoteIndex := make(map[string]bool)
	for _, remoteChannel := range remoteResp.Data {
		remoteIndex[remoteChannel.Code] = true
	}

	// 获取本地渠道数据
	localChannels, err := s.DAOManager.GameChannel.GetAllGameChannel(ctx)
	if err != nil {
		log.Printf("[渠道逆向同步] ✗ 获取本地数据失败: %v\n", err)
		return err
	}
	log.Printf("[渠道逆向同步] ✓ 获取本地数据成功, 共 %d 条\n", len(localChannels))

	for _, localChannel := range localChannels {
		// 如果本地渠道在远程不存在，则软删除
		if !remoteIndex[localChannel.SourceChannelCode] {
			_, err := s.DAOManager.GameChannel.UpdateGameChannel(ctx, localChannel.ID, map[string]interface{}{"deleted_at": time.Now()})

			if err != nil {
				logger.Errorf("❌ [渠道逆向同步] 软删除失败: %v", err)
				apply.Failed++
				continue
			}
			apply.Deleted++
			logger.Infof("[渠道逆向同步] ✓ 已软删除渠道: %s (ID: %d)", localChannel.SourceChannelCode, localChannel.ID)
		}
	}
	if err != nil {
		return err
	}

	log.Printf("[渠道逆向同步] ✓ 完成，删除: %d\n", apply.Deleted)
	return nil
}

func getRemoteGameChannel(ctx context.Context, client vendors.VendorGameServiceClient, isGetRemoteClient bool) (*RemoteChannelResponse, error) {
	if isGetRemoteClient {
		logger.Infof("[同步数据源] 使用远程客户端获取数据")
		remoteResp, err := client.GetChannel(ctx, &vendors.Empty{})
		if err != nil {
			logger.Errorf("❌ [同步数据源] 获取远程数据失败: %v", err)
			return nil, err
		}
		logger.Infof("[同步数据源] ✓ 获取远程数据成功, 共 %d 条", len(remoteResp.ChannelList))
		return &RemoteChannelResponse{Data: remoteResp.ChannelList}, nil
	} else {
		logger.Infof("[同步数据源] 使用本地 JSON 数据获取渠道信息")
		var list []*vendors.GameChannelInfo
		if err := loadLocalJSON("channel.json", &list); err != nil {
			logger.Errorf("[同步数据源] 读取 channel.json 失败: %v", err)
			return nil, err
		}
		logger.Infof("[同步数据源] ✓ 读取本地 JSON 数据成功, 共 %d 条", len(list))
		return &RemoteChannelResponse{Data: list}, nil
	}
}

type RemoteChannelResponse struct {
	Data []*vendors.GameChannelInfo `json:"data"`
}
