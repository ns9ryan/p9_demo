package task

import (
	"context"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

// markDispatchFailed 标记任务下发失败
func (s *Service) markDispatchFailed(ctx context.Context, taskID int64, runID int64, errorMessage string) {
	// 开启失败状态更新事务
	tx, err := s.db.Tx(ctx)
	if err != nil {
		logx.WithContext(ctx).Errorw(
			"开启任务失败状态事务失败",
			logx.Field("task_id", taskID),
			logx.Field("error", err.Error()),
		)
		return
	}

	// 标记调度任务失败
	err = tx.DispatchTask.
		UpdateOneID(taskID).
		SetStatus(4). // 任务状态: 1待执行, 2执行中, 3成功, 4失败
		Exec(ctx)
	if err != nil {
		_ = tx.Rollback()

		logx.WithContext(ctx).Errorw(
			"更新调度任务失败状态失败",
			logx.Field("task_id", taskID),
			logx.Field("error", err.Error()),
		)
		return
	}

	// 标记任务执行记录失败
	err = tx.DispatchTaskRun.
		UpdateOneID(runID).
		SetStatus(4).                  // 执行状态: 1待执行, 2执行中, 3成功, 4失败
		SetErrorMessage(errorMessage). // 执行失败原因
		SetFinishedAt(time.Now()).     // 执行结束时间
		Exec(ctx)
	if err != nil {
		_ = tx.Rollback()

		logx.WithContext(ctx).Errorw(
			"更新任务执行记录失败状态失败",
			logx.Field("run_id", runID),
			logx.Field("error", err.Error()),
		)
		return
	}

	// 提交失败状态
	if err = tx.Commit(); err != nil {
		logx.WithContext(ctx).Errorw(
			"提交任务失败状态事务失败",
			logx.Field("task_id", taskID),
			logx.Field("error", err.Error()),
		)
	}
}
