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
	StartAt  time.Time
	EndAt    time.Time
	Link     string
}
