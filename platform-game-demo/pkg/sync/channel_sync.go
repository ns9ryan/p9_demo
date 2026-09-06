package sync

import (
	"context"
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"
	"oa.98ent.com/p9/platform-game/common/logger"
	"oa.98ent.com/p9/platform-game/common/model"
	"oa.98ent.com/p9/platform-game/rpc/pb/vendors"
)

// ===== 渠道同步服务 =====

// ChannelSyncService 渠道同步服务
type ChannelSyncService struct {
	db *gorm.DB
}

// NewChannelSyncService 创建渠道同步服务
func NewChannelSyncService(db *gorm.DB) *ChannelSyncService {
	return &ChannelSyncService{db: db}
}

// Preview 预检查渠道同步
func (s *ChannelSyncService) Preview(ctx context.Context, client vendors.VendorGameServiceClient) (*vendors.SyncPreviewResp, error) {
	// 获取远程渠道数据
	// remoteResp, err := client.GetChannel(ctx, &emptypb.Empty{})
	// if err != nil {
	// 	log.Printf("[渠道预检查] 获取远程数据失败: %v", err)
	// 	return nil, err
	// }
	remoteResp, err := getTestGameChannel()
	if err != nil {
		log.Printf("[渠道预检查] 获取测试数据失败: %v", err)
		return nil, err
	}
	remote := remoteResp.Data

	// 获取本地渠道数据
	local, localIndex, err := s.fetchLocal(ctx)
	if err != nil {
		return nil, err
	}

	// 初始化预检查结果
	result := &vendors.SyncPreviewResp{
		Stats: &vendors.SyncStats{
			RemoteTotal: int64(len(remote)),
			LocalTotal:  int64(len(local)),
		},
		Diffs: make([]*vendors.SyncDiff, 0),
	}

	// 对每条远程渠道进行比较
	for _, remoteItem := range remote {
		diff := s.compareOne(remoteItem, localIndex)
		result.Diffs = append(result.Diffs, diff)

		// 按操作类型统计
		switch diff.Action {
		case "create":
			result.Stats.CreateTotal++
		case "update":
			result.Stats.UpdateTotal++
		case "noop":
			result.Stats.NoopTotal++
		case "conflict":
			result.Stats.ConflictTotal++
		}
	}

	return result, nil
}

// Run 执行渠道同步
// processChannelData 处理渠道数据的创建和更新（从 Run 提取的业务逻辑）
func (s *ChannelSyncService) processChannelData(ctx context.Context, tx *gorm.DB, remoteData []*vendors.GameChannelInfo, applyResult *vendors.SyncApplyResult) error {
	i := 0
	counterMap := make(map[string]int) // 用于记录每个 source_channel_code 的计数，确保唯一性
	for _, remoteChannel := range remoteData {
		if counterMap[remoteChannel.Code] > 0 {
			logger.Errorf("[渠道同步] 检测到重复的远程渠道编码: %s, 计数器: %d, 跳过处理", remoteChannel.Code, counterMap[remoteChannel.Code])
			applyResult.Failed++
			continue
		}
		// 更新计数器
		counterMap[remoteChannel.Code]++
		i++
		// 检查本地是否存在该渠道（使用上游渠道编码）
		var localChannel model.Channel
		exists := tx.Where("source_channel_code = ?", remoteChannel.Code).
			Where("deleted_at IS NULL").
			First(&localChannel).Error == nil

		if !exists {
			// 新增渠道 - 生成 P9 内部的渠道编码
			channelCode := remoteChannel.Code + "_" + fmt.Sprintf("%d", i)
			nameI18n := model.JSONMap{"default": remoteChannel.Name}
			sourceNameI18n := model.JSONMap{"default": remoteChannel.Name}

			newChannel := model.Channel{
				SourceID:          remoteChannel.Id,
				ChannelCode:       channelCode,
				SourceChannelCode: remoteChannel.Code,
				SourceNameI18n:    sourceNameI18n,
				NameI18n:          nameI18n,
				SourceStatus:      remoteChannel.Status,
				Status:            remoteChannel.Status,
				CreatedAt:         time.Now(),
				UpdatedAt:         time.Now(),
			}
			if err := tx.Create(&newChannel).Error; err != nil {
				logger.Errorf("[渠道新增] 失败: %v", err)
				applyResult.Failed++
				continue
			}
			applyResult.Created++
		} else {
			// 更新渠道信息
			channelCode := remoteChannel.Code + "_" + fmt.Sprintf("%d", i)
			sourceNameI18n := model.JSONMap{"default": remoteChannel.Name}
			if err := tx.Model(&localChannel).
				Updates(map[string]interface{}{
					"channel_code":        channelCode,
					"source_channel_code": remoteChannel.Code,
					"source_name_i18n":    sourceNameI18n,
					"source_status":       remoteChannel.Status,
					"updated_at":          time.Now(),
				}).Error; err != nil {
				logger.Errorf("[渠道更新] 失败: %v", err)
				applyResult.Failed++
				continue
			}
			applyResult.Updated++
		}
	}
	return nil
}

func (s *ChannelSyncService) Run(ctx context.Context, client vendors.VendorGameServiceClient) (*vendors.SyncRunResp, error) {
	logger.Infof("[渠道同步] ===== 开始执行同步 =====")

	// 先执行预检查
	logger.Infof("[渠道同步] 执行预检查")
	previewResp, err := s.Preview(ctx, client)
	if err != nil {
		logger.Errorf("[渠道同步] 预检查失败: %v", err)
		return nil, err
	}
	logger.Infof("[渠道同步] 预检查完成: remote_total=%d, local_total=%d",
		previewResp.Stats.RemoteTotal, previewResp.Stats.LocalTotal)

	// 构建执行结果
	result := &vendors.SyncRunResp{
		Preview: previewResp,
		Apply:   &vendors.SyncApplyResult{},
	}

	// 获取远程渠道数据
	logger.Infof("[渠道同步] 获取远程渠道数据")
	// remoteResp, err := client.GetChannel(ctx, &emptypb.Empty{})
	// if err != nil {
	// 	logger.Errorf("[渠道同步] 获取远程数据失败: %v", err)
	// 	return nil, err
	// }
	remoteResp, err := getTestGameChannel()
	if err != nil {
		logger.Errorf("[渠道同步] 获取远程数据失败: %v", err)
		return nil, err
	}
	logger.Infof("[渠道同步] 获取到远程数据: count=%d", len(remoteResp.Data))

	// 在事务中执行数据库操作
	logger.Infof("[渠道同步] 开始处理渠道数据")
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return s.processChannelData(ctx, tx, remoteResp.Data, result.Apply)
	})
	if err != nil {
		logger.Errorf("[渠道同步] 处理渠道数据失败: %v", err)
		return nil, err
	}
	logger.Infof("[渠道同步] 处理完成: created=%d, updated=%d, deleted=%d, failed=%d",
		result.Apply.Created, result.Apply.Updated, result.Apply.Deleted, result.Apply.Failed)

	// 执行逆向同步
	logger.Infof("[渠道同步] 开始执行逆向同步")
	reverseResp, err := s.ReverseSync(ctx, client)
	if err != nil {
		logger.Errorf("[渠道同步] 逆向同步失败: %v", err)
		return nil, err
	}
	result.Apply.Deleted = reverseResp.Apply.Deleted
	logger.Infof("[渠道同步] 逆向同步完成: deleted=%d", result.Apply.Deleted)

	// 保存同步检查点
	logger.Infof("[渠道同步] 准备保存同步检查点")
	checkpointMgr := NewCheckpointManager(s.db)
	checkpointValue := previewResp.Stats.RemoteTotal // 使用远程总数作为检查点值
	if err := checkpointMgr.UpdateCheckpoint(ctx, "CHANNEL",
		fmt.Sprintf("%d", checkpointValue), result, nil); err != nil {
		logger.Errorf("[渠道同步] 保存检查点失败: %v", err)
		return nil, err
	}
	logger.Infof("[渠道同步] ✓ 检查点已保存")

	logger.Infof("[渠道同步] ===== 同步完成 =====")
	return result, nil
}

// fetchLocal 获取本地所有渠道
func (s *ChannelSyncService) fetchLocal(ctx context.Context) ([]*model.Channel, map[string]*model.Channel, error) {
	var channels []*model.Channel
	if err := s.db.WithContext(ctx).Where("deleted_at IS NULL").Find(&channels).Error; err != nil {
		return nil, nil, err
	}

	// 构建本地渠道索引（按上游渠道编码）
	index := make(map[string]*model.Channel)
	for _, channel := range channels {
		if channel.SourceChannelCode != "" {
			index[channel.SourceChannelCode] = channel
		}
	}

	return channels, index, nil
}

// compareOne 比较单个渠道的本地和远程数据
func (s *ChannelSyncService) compareOne(remote *vendors.GameChannelInfo, localIndex map[string]*model.Channel) *vendors.SyncDiff {
	local, exists := localIndex[remote.Code]

	// 初始化差异记录
	diff := &vendors.SyncDiff{
		ObjectType: "channel",
		ObjectId:   remote.Id,
		ObjectCode: remote.Code,
		RemoteId:   remote.Id,
		RemoteCode: remote.Code,
	}

	// 如果本地不存在该渠道
	if !exists {
		diff.Action = "create"
		diff.Reason = "本地不存在该渠道"
		return diff
	}

	// 对比渠道名称是否变更
	var localSourceName string
	if sourceNameI18n, ok := local.SourceNameI18n["default"]; ok {
		localSourceName = sourceNameI18n.(string)
	}
	if localSourceName != remote.Name {
		diff.Action = "update"
		diff.Reason = fmt.Sprintf("渠道名称变更: %s -> %s", localSourceName, remote.Name)
		return diff
	}

	// 本地和远程一致，无需操作
	diff.Action = "noop"
	diff.Reason = "本地和远程数据一致"
	return diff
}

// ReverseSync 逆向同步：检查本地数据在远程是否存在，不存在则软删除
func (s *ChannelSyncService) ReverseSync(ctx context.Context, client vendors.VendorGameServiceClient) (*vendors.SyncRunResp, error) {
	log.Println("[渠道逆向同步] 开始执行逆向同步...")

	// 获取远程渠道数据
	// remoteResp, err := client.GetChannel(ctx, &emptypb.Empty{})
	// if err != nil {
	// 	log.Printf("[渠道逆向同步] ✗ 获取远程数据失败: %v\n", err)
	// 	return nil, err
	// }
	remoteResp, err := getTestGameChannel()
	if err != nil {
		log.Printf("[渠道逆向同步] ✗ 获取测试数据失败: %v\n", err)
		return nil, err
	}
	log.Printf("[渠道逆向同步] ✓ 获取远程数据成功, 共 %d 条\n", len(remoteResp.Data))

	// 构建远程渠道编码索引
	remoteIndex := make(map[string]bool)
	for _, remoteChannel := range remoteResp.Data {
		remoteIndex[remoteChannel.Code] = true
	}

	// 获取本地渠道数据
	var localChannels []*model.Channel
	if err := s.db.WithContext(ctx).Where("deleted_at IS NULL").Find(&localChannels).Error; err != nil {
		log.Printf("[渠道逆向同步] ✗ 获取本地数据失败: %v\n", err)
		return nil, err
	}
	log.Printf("[渠道逆向同步] ✓ 获取本地数据成功, 共 %d 条\n", len(localChannels))

	// 构建执行结果
	result := &vendors.SyncRunResp{
		Preview: &vendors.SyncPreviewResp{
			Stats: &vendors.SyncStats{
				LocalTotal: int64(len(localChannels)),
			},
		},
		Apply: &vendors.SyncApplyResult{},
	}

	// 在事务中执行软删除
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, localChannel := range localChannels {
			// 如果本地渠道在远程不存在，则软删除
			if !remoteIndex[localChannel.SourceChannelCode] {
				// 检查唯一性：是否存在其他 deleted_at IS NULL 的记录有相同的 source_channel_code
				var duplicateCount int64
				if err := tx.Model(&model.Channel{}).
					Where("source_channel_code = ?", localChannel.SourceChannelCode).
					Where("deleted_at IS NULL").
					Count(&duplicateCount).Error; err != nil {
					log.Printf("[渠道逆向同步] 唯一性检查失败: %v", err)
					result.Apply.Failed++
					continue
				}

				if duplicateCount > 1 {
					logger.Errorf("[渠道逆向同步] 存在 %d 条记录有相同的 source_channel_code=%s, 跳过本条删除以保护数据", duplicateCount, localChannel.SourceChannelCode)
					result.Apply.Failed++
					continue
				}

				if err := tx.Model(localChannel).Update("deleted_at", time.Now()).Error; err != nil {
					log.Printf("[渠道逆向同步] 软删除失败: %v", err)
					result.Apply.Failed++
					continue
				}
				result.Apply.Deleted++
				log.Printf("[渠道逆向同步] ✓ 已软删除渠道: %s\n", localChannel.SourceChannelCode)
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	log.Printf("[渠道逆向同步] ✓ 完成，删除: %d\n", result.Apply.Deleted)
	return result, nil
}

func getTestGameChannel() (*vendors.GetGameChannelResp, error) {
	var list []*vendors.GameChannelInfo
	if err := loadTestJSON("channel.json", &list); err != nil {
		log.Printf("[渠道测试数据] 读取 channel.json 失败: %v", err)
		return &vendors.GetGameChannelResp{Data: []*vendors.GameChannelInfo{}}, err
	}

	return &vendors.GetGameChannelResp{Data: list}, nil
}
