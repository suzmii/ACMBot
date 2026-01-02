package subconfig

import "time"

type RunOnStartMode string

const (
	RunOnStartTrue  RunOnStartMode = "true"
	RunOnStartFalse RunOnStartMode = "false"
	RunOnStartAuto  RunOnStartMode = "auto"
)

type TaskConfig struct {
	Spec       string         `mapstructure:"spec" toml:"spec"`
	RunOnStart RunOnStartMode `mapstructure:"run_on_start" toml:"run_on_start"`
	RetryCount int            `mapstructure:"retry_count" toml:"retry_count"`
	RetryWait  time.Duration  `mapstructure:"retry_wait" toml:"retry_wait"`
}

type Scheduler struct {
	Tasks map[string]TaskConfig `mapstructure:"tasks" toml:"tasks"`
}

var DefaultScheduler = Scheduler{
	Tasks: map[string]TaskConfig{
		"race-updater": {
			Spec:       "@every 1h",
			RunOnStart: RunOnStartAuto,
			RetryCount: 3,
			RetryWait:  time.Minute,
		},
	},
}
