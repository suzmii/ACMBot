package scheduler

import (
	"context"
	"fmt"
	"time"

	"github.com/suzmii/ACMBot/config/subconfig"
	"github.com/suzmii/ACMBot/database"
)

// Task 定义调度任务的接口
type Task interface {
	// Name 返回任务名称
	Name() string
	// Execute 执行任务逻辑
	Execute(ctx context.Context) error
}

// CheckableTask 定义可以检查是否需要更新的任务接口
type CheckableTask interface {
	Task
	// CheckIfNeedUpdate 检查是否需要更新（用于run_on_start: auto模式）
	CheckIfNeedRunOnStart(ctx context.Context) (bool, error)
}

// TaskWrapper 包装Task，提供重试机制
type TaskWrapper struct {
	task  Task
	cfg   subconfig.TaskConfig
	store database.Store
}

func NewTaskWrapper(task Task, cfg subconfig.TaskConfig, store database.Store) *TaskWrapper {
	return &TaskWrapper{
		task:  task,
		cfg:   cfg,
		store: store,
	}
}

func (tw *TaskWrapper) ExecuteWithRetry(ctx context.Context) error {
	var err error
	for i := 1; i <= tw.cfg.RetryCount; i++ {
		err = tw.task.Execute(ctx)
		if err == nil {
			logger.Infof("[%s] 任务执行成功", tw.task.Name())
			return nil
		}

		if i < tw.cfg.RetryCount {
			logger.Warnf("[%s] 第%d次执行失败，将在%v后重试: %v", tw.task.Name(), i, tw.cfg.RetryWait, err)
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(tw.cfg.RetryWait):
			}
		}
	}

	return fmt.Errorf("[%s] 任务执行失败，已重试%d次: %w", tw.task.Name(), tw.cfg.RetryCount, err)
}
