package dbmodel

import (
	"time"

	"gorm.io/gorm"
)

type Resource int

const (
	ResourceCodeforces Resource = iota + 1
	ResourceAtcoder
	ResourceLeetcode
	ResourceLuogu
	ResourceNowcoder
)

type Races struct {
	gorm.Model

	Resource Resource
	Title    string
	StartAt  time.Time `gorm:"index:idx_races_start_at"`
	EndAt    time.Time `gorm:"index:idx_races_end_at"`
	Link     string
}
