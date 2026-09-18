package game_sync

import (
	"context"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"oa.98ent.com/p9/platform-game/rpc/ent"
	"oa.98ent.com/p9/platform-game/rpc/internal/config"
	"oa.98ent.com/p9/platform-game/rpc/internal/constant"
	"oa.98ent.com/p9/platform-game/rpc/internal/dao"
	"oa.98ent.com/p9/platform-game/rpc/internal/locales"
	"oa.98ent.com/p9/platform-game/rpc/pb/platform_game"
	"oa.98ent.com/p9/platform-game/rpc/pb/vendors"
)

// ===== 提供商同步服务 =====

// ProviderSyncService 提供商同步服务
type ProviderSyncService struct {
	DAOManager    *dao.Manager
	config        config.Config
	progressQueue *ProgressQueue
	logx.Logger
}

// NewProviderSyncService 创建提供商同步服务
func NewProviderSyncService(ctx context.Context, config config.Config, daoManager *dao.Manager) *ProviderSyncService {
	return &ProviderSyncService{
		DAOManager:    daoManager,
		config:        config,
		progressQueue: NewProgressQueue(daoManager),
		Logger:        logx.WithContext(ctx),
	}
}

// Preview 预检查提供商同步
func (s *ProviderSyncService) Preview(ctx context.Context, client vendors.VendorGameServiceClient, req *platform_game.SyncPreviewRequest, isGetRemoteClient bool) (*platform_game.SyncPreviewResp, *RemoteProviderResponse, error) {
	// 获取远程提供商数据
	remoteResp, err := s.getRemoteGameProviders(ctx, client, isGetRemoteClient)
	if err != nil {
		s.Errorf("[提供商预检查] 获取测试数据失败: %v", err)
		return nil, nil, err
	}
	remote := remoteResp.Data

	// 获取本地提供商数据
	local, localIndex, err := s.fetchLocal(ctx)
	if err != nil {
		return nil, nil, err
	}

	// 转换为接口切片
	remoteItems := make([]interface{}, len(remote))
	for i, v := range remote {
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
				remoteProvider := item.(*vendors.VendorInfo)
				localIdx := make(map[string]*ent.GameProvider)
				for k, v := range localIndex {
					localIdx[k] = v.(*ent.GameProvider)
				}
				return s.compareAll(remoteProvider, localIdx)
			},
			PageSize: 0,
			Page:     0,
			SkipNoop: false,
		}
		result := handler.Handle(int64(len(remote)), int64(len(local)))
		return result, remoteResp, nil
	}

	// 创建分页处理器
	handler := &SyncPageHandler{
		RemoteData: remoteItems,
		LocalIndex: localIndexInterface,
		CompareFunc: func(item interface{}, localIndex map[string]interface{}) *platform_game.SyncDiff {
			remoteProvider := item.(*vendors.VendorInfo)
			localIdx := make(map[string]*ent.GameProvider)
			for k, v := range localIndex {
				localIdx[k] = v.(*ent.GameProvider)
			}
			return s.compareAll(remoteProvider, localIdx)
		},
		PageSize: req.PageSize,
		Page:     req.Page,
		SkipNoop: req.IsSkip,
	}

	// 处理分页逻辑
	result := handler.Handle(int64(len(remote)), int64(len(local)))

	return result, remoteResp, nil
}

// 处理提供商数据的创建和更新（从 Run 提取的业务逻辑）
func (s *ProviderSyncService) processProviderDataWithProgress(ctx context.Context, remoteData []*vendors.VendorInfo, applyResult *Apply, syncCols []string, checkpointID int64) error {
	i := 0
	counterMap := make(map[string]int) // 用于记录每个 source_provider_code 的计数，确保唯一性
	nameMap := make(map[string]string)
	createList := make([]*ent.GameProviderCreate, 0, len(remoteData))
	for _, remoteProvider := range remoteData {
		if counterMap[remoteProvider.Code] > 0 {
			s.Errorf("[提供商同步] 检测到重复的远程提供商编码: %s, 计数器: %d, 跳过处理", remoteProvider.Code, counterMap[remoteProvider.Code])
			applyResult.Failed++
			continue
		}
		nameMap[remoteProvider.Code] = remoteProvider.Name
		// 更新计数器
		counterMap[remoteProvider.Code]++
		i++
		// 检查本地是否存在该提供商（使用上游提供商编码）
		localProvider, err := s.DAOManager.GameProvider.GetGameProviderBySourceCode(ctx, remoteProvider.Code)
		// 生成 P9 内部的提供商编码
		providerCode := remoteProvider.Code + "_" + fmt.Sprintf("%d", i)
		if err != nil || localProvider == nil {
			createList = append(createList, s.DAOManager.DB.GameProvider.Create().
				SetSourceID(remoteProvider.Id).
				SetProviderCode(providerCode).
				SetSourceProviderCode(remoteProvider.Code).
				SetSortNo(int64(remoteProvider.Id)).
				SetSourceSortNo(0).
				SetSourceStatus(int64(remoteProvider.Status)).
				SetStatus(int64(remoteProvider.Status)).
				SetCreatedAt(time.Now()).
				SetUpdatedAt(time.Now()))
		} else {
			if s.compare(&platform_game.SyncDiff{}, remoteProvider, localProvider) || len(syncCols) > 0 {
				s.Infof("[提供商同步] 检测到需要更新的提供商: %s", remoteProvider.Code)
				update := s.DAOManager.DB.GameProvider.
					UpdateOneID(localProvider.ID).
					SetUpdatedAt(time.Now())

				if len(syncCols) == 0 {
					syncCols = append(syncCols, "provider_code")
					syncCols = append(syncCols, "source_provider_code")
					syncCols = append(syncCols, "source_status")
					syncCols = append(syncCols, "source_sort_no")
				}

				for _, col := range syncCols {
					switch col {
					case "provider_code":
						update.SetProviderCode(providerCode)
					case "source_provider_code":
						update.SetSourceProviderCode(remoteProvider.Code)
					case "status":
						update.SetStatus(int64(remoteProvider.Status))
					case "source_status":
						update.SetSourceStatus(int64(remoteProvider.Status))
					case "sort_no":
						update.SetSortNo(int64(remoteProvider.Id))
					case "source_sort_no":
						update.SetSourceSortNo(0)
					}
				}

				if _, err := update.Save(ctx); err != nil {
					s.Errorf("[提供商更新] 失败: %v", err)
					applyResult.Failed++
					continue
				}
				applyResult.Updated++
				s.Infof("[提供商更新] ✓ 已更新提供商: %s", remoteProvider.Code)
			}

		}
	}
	if len(createList) > 0 {
		_, err := s.DAOManager.GameProvider.BatchCreateGameProvider(ctx, createList)
		if err != nil {
			s.Errorf("[提供商新增] 批量创建失败: %v", err)
			applyResult.Failed += int32(len(createList))
		} else {
			applyResult.Created += int32(len(createList))
		}
	}
	// 将 nameMap 写入本地 JSON 文件
	if err := locales.MergeLocalJSON("game_provider_name_map.json", nameMap, constant.ProviderBiz); err != nil {
		s.Errorf("[提供商同步] 保存名称映射失败: %v", err)
		// 不返回错误，继续执行后续逻辑
	}
	return nil
}

func (s *ProviderSyncService) Run(ctx context.Context, client vendors.VendorGameServiceClient, syncCols []string, checkpointID int64, isGetRemoteClient bool) error {
	s.Infof("[提供商同步] ===== 开始执行同步 =====")

	// 先执行预检查
	s.Infof("[提供商同步] 执行预检查")
	previewResp, remoteResp, err := s.Preview(ctx, client, nil, isGetRemoteClient)
	if err != nil {
		s.Errorf("[提供商同步] 预检查失败: %v", err)
		return err
	}
	// 创建并启动进度队列
	s.progressQueue.Start(ctx)
	defer s.progressQueue.Stop()
	if previewResp.Stats.UpdateTotal == 0 && previewResp.Stats.CreateTotal == 0 && previewResp.Stats.DeleteTotal == 0 && previewResp.Stats.RemoteTotal == previewResp.Stats.LocalTotal && len(syncCols) == 0 {
		s.Infof("[提供商同步] 无需更新，直接返回")
		s.Infof("[提供商同步] UpdateTotal=%d, CreateTotal=%d, DeleteTotal=%d, RemoteTotal=%d, LocalTotal=%d",
			previewResp.Stats.UpdateTotal, previewResp.Stats.CreateTotal, previewResp.Stats.DeleteTotal, previewResp.Stats.RemoteTotal, previewResp.Stats.LocalTotal)
		// 更新checkpointID进度为100
		s.progressQueue.Send(&ProgressMessage{
			TableName:      "provider",
			ProcessedCount: previewResp.Stats.RemoteTotal,
			RemoteTotal:    previewResp.Stats.RemoteTotal,
			LocalTotal:     previewResp.Stats.LocalTotal,
			Progress:       100,
			CheckpointID:   checkpointID,
		})
		return nil
	}
	// 把remoteResp.Data拆分为多条一批进行处理（可根据实际情况调整批次大小）
	batchSize := s.config.SyncBatchSize
	s.Infof("[提供商同步] 预检查完成: remote_total=%d, local_total=%d, 检查点已创建: checkpointID=%d, batch_size=%d",
		previewResp.Stats.RemoteTotal, previewResp.Stats.LocalTotal, checkpointID, batchSize)
	totalCount := len(remoteResp.Data)
	apply := &Apply{}
	for i := 0; i < len(remoteResp.Data); i += batchSize {
		end := i + batchSize
		if end > len(remoteResp.Data) {
			end = len(remoteResp.Data)
		}
		batch := remoteResp.Data[i:end]
		s.Infof("[提供商同步] 处理批次: start=%d, end=%d", i, end)
		// 模拟等待
		s.Debugf("[进度队列] 模拟处理延迟: 表=provider, checkpointID=%d", checkpointID)
		time.Sleep(1000 * time.Millisecond)
		s.Debugf("[进度队列] 开始处理消息: 表=provider, checkpointID=%d", checkpointID)
		err = s.processProviderDataWithProgress(ctx, batch, apply, syncCols, checkpointID)
		if err != nil {
			s.Errorf("[提供商同步] 处理数据失败: %v", err)
			return err
		}
		progress := CalculateProgress(int64(i+1), int64(totalCount))
		msg := &ProgressMessage{
			TableName:      "provider",
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
		s.Infof("[提供商同步] 发送进度消息: 表=provider, 处理数=%d, 总数=%d, 进度=%d%%", i+1, totalCount, progress)
		s.progressQueue.Send(msg)
	}
	if err != nil {
		s.Errorf("[提供商同步] 处理提供商数据失败: %v", err)
		return err
	}

	// 执行逆向同步
	s.Infof("[提供商同步] 开始执行逆向同步")
	err = s.ReverseSync(ctx, client, remoteResp, apply)
	if err != nil {
		s.Errorf("❌ [提供商同步] 逆向同步失败: %v", err)
		return err
	}
	s.Infof("[提供商同步] 逆向同步完成: deleted=%d", apply.Deleted)

	// 最后一次进度更新（100%）
	s.Infof("[提供商同步] 发送最终进度消息（100%%）")
	finalMsg := &ProgressMessage{
		TableName:      "provider",
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

	s.Infof("[提供商同步] ===== 同步完成 =====")
	// 构建执行结果
	return nil
}

// fetchLocal 获取本地所有提供商
func (s *ProviderSyncService) fetchLocal(ctx context.Context) ([]*ent.GameProvider, map[string]*ent.GameProvider, error) {
	providers, err := s.DAOManager.GameProvider.GetAllGameProviders(ctx)
	if err != nil {
		return nil, nil, err
	}

	// 构建本地提供商索引（按上游提供商编码）
	index := make(map[string]*ent.GameProvider)
	for _, provider := range providers {
		if provider.SourceProviderCode != "" {
			index[provider.SourceProviderCode] = provider
		}
	}

	return providers, index, nil
}

// 比较单个提供商的本地和远程数据
func (s *ProviderSyncService) compareAll(remote *vendors.VendorInfo, localIndex map[string]*ent.GameProvider) *platform_game.SyncDiff {
	local, exists := localIndex[remote.Code]

	// 初始化差异记录
	diff := &platform_game.SyncDiff{
		ObjectType: "provider",
		ObjectId:   remote.Id,
		ObjectCode: remote.Code,
		RemoteId:   remote.Id,
		RemoteCode: remote.Code,
	}

	// 如果本地不存在该提供商
	if !exists {
		diff.Action = "create"
		diff.Reason = "本地不存在该提供商"
		return diff
	}

	s.compare(diff, remote, local)
	return diff
}

func (s *ProviderSyncService) compare(diff *platform_game.SyncDiff, remote *vendors.VendorInfo, local *ent.GameProvider) bool {
	if local.SourceStatus != int64(remote.Status) {
		diff.Action = constant.SyncActionUpdate
		diff.Reason = fmt.Sprintf("提供商状态变更: %d -> %d", local.SourceStatus, remote.Status)
		diff.ConflictType = "status_changed"
		s.Infof("[提供商同步] 检测到提供商状态变更: %d -> %d", local.SourceStatus, remote.Status)
		return true
	}
	diff.Action = constant.SyncActionNoop
	diff.Reason = "本地和远程数据一致"
	return false
}

// 逆向同步：检查本地数据在远程是否存在，不存在则软删除
func (s *ProviderSyncService) ReverseSync(ctx context.Context, client vendors.VendorGameServiceClient, remoteResp *RemoteProviderResponse, apply *Apply) error {
	s.Infof("[提供商逆向同步] 开始执行逆向同步")

	// 构建远程提供商编码索引
	remoteIndex := make(map[string]bool)
	for _, remoteProvider := range remoteResp.Data {
		remoteIndex[remoteProvider.Code] = true
	}

	// 获取本地提供商数据
	localProviders, err := s.DAOManager.GameProvider.GetAllGameProviders(ctx)
	if err != nil {
		s.Errorf("[提供商逆向同步] ✗ 获取本地数据失败: %v", err)
		return err
	}
	s.Infof("[提供商逆向同步] ✓ 获取本地数据成功, 共 %d 条", len(localProviders))

	// 在事务中执行软删除
	tx, err := s.DAOManager.DB.Tx(ctx)
	if err != nil {
		return fmt.Errorf("starting transaction: %w", err)
	}

	_ = func(tx *ent.Tx, err error) error {
		if rerr := tx.Rollback(); rerr != nil {
			err = fmt.Errorf("%w: rolling back transaction: %v", err, rerr)
		}
		return err
	}

	for _, localProvider := range localProviders {
		// 如果本地提供商在远程不存在，则软删除
		if !remoteIndex[localProvider.SourceProviderCode] {
			err := tx.GameProvider.UpdateOne(localProvider).
				SetDeletedAt(time.Now()).
				Exec(ctx)

			if err != nil {
				s.Errorf("[提供商逆向同步] 软删除失败: %v", err)
				apply.Failed++
				continue
			}
			apply.Deleted++
			s.Infof("[提供商逆向同步] ✓ 已软删除提供商: %s (ID: %d)", localProvider.SourceProviderCode, localProvider.ID)
		}
	}

	err = tx.Commit()

	if err != nil {
		return err
	}

	s.Infof("[提供商逆向同步] ✓ 完成，删除: %d", apply.Deleted)
	return nil
}

func (s *ProviderSyncService) getRemoteGameProviders(ctx context.Context, client vendors.VendorGameServiceClient, isGetRemoteClient bool) (*RemoteProviderResponse, error) {
	if isGetRemoteClient {
		s.Infof("[同步数据源] 使用远程客户端获取数据")
		remoteResp, err := client.GetVendor(ctx, &vendors.Empty{})
		if err != nil {
			s.Errorf("❌ [同步数据源] 获取远程数据失败: %v", err)
			return nil, err
		}
		s.Infof("[同步数据源] ✓ 获取远程数据成功, 共 %d 条", len(remoteResp.VendorList))
		return &RemoteProviderResponse{Data: remoteResp.VendorList}, nil
	} else {
		s.Infof("[同步数据源] 使用本地 JSON 数据获取提供商信息")
		var list []*vendors.VendorInfo
		if err := locales.LoadLocalVendorRemoteJSON("provider.json", &list); err != nil {
			s.Errorf("[同步数据源] 读取 provider.json 失败: %v", err)
			return nil, err
		}
		s.Infof("[同步数据源] ✓ 读取本地 JSON 数据成功, 共 %d 条", len(list))
		return &RemoteProviderResponse{Data: list}, nil
	}
}

type RemoteProviderResponse struct {
	Data []*vendors.VendorInfo `json:"data"`
}
