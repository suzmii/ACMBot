package scheduler

import (
	"context"
	"fmt"

	"github.com/robfig/cron/v3"
	"github.com/suzmii/ACMBot/config/subconfig"
	"github.com/suzmii/ACMBot/database"
	"github.com/suzmii/ACMBot/util/logx"
)

var logger = logx.New("scheduler")

// Scheduler 统一调度器
type Scheduler struct {
	cfg   subconfig.Scheduler
	store database.Store
	cron  *cron.Cron
	tasks map[string]TaskWrapper
}

// NewScheduler 创建新的调度器
func NewScheduler(cfg subconfig.Scheduler, store database.Store, tasks []Task) (*Scheduler, error) {
	s := &Scheduler{
		cfg:   cfg,
		store: store,
		cron:  cron.New(),
		tasks: make(map[string]TaskWrapper),
	}

	// 注册所有任务
	for _, task := range tasks {
		taskCfg, ok := cfg.Tasks[task.Name()]
		if !ok {
			logger.Warnf("任务 %s 未找到配置，使用默认配置", task.Name())
			continue
		}

		wrapper := NewTaskWrapper(task, taskCfg, store)
		s.tasks[task.Name()] = *wrapper

		// 添加到cron
		_, err := s.cron.AddFunc(taskCfg.Spec, func() {
			ctx := context.Background()
			err := wrapper.ExecuteWithRetry(ctx)
			if err != nil {
				logger.Errorf("任务 %s 执行失败: %v", task.Name(), err)
			}
		})
		if err != nil {
			return nil, fmt.Errorf("failed to add task %s: %w", task.Name(), err)
		}

		logger.Infof("已注册任务: %s, 调度规则: %s", task.Name(), taskCfg.Spec)
	}

	return s, nil
}

// Start 启动调度器
func (s *Scheduler) Start() {
	ctx := context.Background()

	// 根据配置决定是否在启动时执行任务
	for name, wrapper := range s.tasks {
		switch wrapper.cfg.RunOnStart {
		case subconfig.RunOnStartTrue:
			// true: 启动时立即执行
			go s.executeTaskOnStart(name, wrapper)
		case subconfig.RunOnStartFalse:
			// false: 不在启动时执行
			logger.Infof("任务 %s 配置为不在启动时执行", name)
		case subconfig.RunOnStartAuto:
			// auto: 检查任务是否实现了CheckableTask接口
			if checkable, ok := wrapper.task.(CheckableTask); ok {
				needUpdate, err := checkable.CheckIfNeedRunOnStart(ctx)
				if err != nil {
					logger.Errorf("任务 %s 检查是否需要更新失败: %v", name, err)
				} else if needUpdate {
					go s.executeTaskOnStart(name, wrapper)
				} else {
					logger.Infof("任务 %s 的check函数表示不需要现在就执行", name)
				}
			} else {
				// 如果没有实现CheckableTask接口，默认不执行
				logger.Infof("任务 %s 未实现CheckableTask接口，跳过启动时执行", name)
			}
		}
	}

	s.cron.Start()
	logger.Info("调度器已启动")
}

// executeTaskOnStart 在启动时执行任务
func (s *Scheduler) executeTaskOnStart(name string, wrapper TaskWrapper) {
	ctx := context.Background()
	err := wrapper.ExecuteWithRetry(ctx)
	if err != nil {
		logger.Errorf("任务 %s 启动时执行失败: %v", name, err)
	} else {
		logger.Infof("任务 %s 启动时执行成功", name)
	}
}

// Stop 停止调度器
func (s *Scheduler) Stop() {
	s.cron.Stop()
	logger.Info("调度器已停止")
}
