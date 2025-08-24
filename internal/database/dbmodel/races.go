package dbmodel

import (
	"fmt"
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

func (r Resource) String() string {
	switch r {
	case ResourceCodeforces:
		return "codeforces"
	case ResourceAtcoder:
		return "atcoder"
	case ResourceLeetcode:
		return "leetcode"
	case ResourceLuogu:
		return "luogu"
	case ResourceNowcoder:
		return "nowcoder"
	default:
		return "unknown"
	}
}

type Races struct {
	gorm.Model

	ID       string `gorm:"primaryKey"`
	Resource Resource
	Title    string
	StartAt  time.Time `gorm:"index:idx_races_start_at"`
	EndAt    time.Time `gorm:"index:idx_races_end_at"`
	Link     string
}

func (r *Races) GenerateID() string {
	return fmt.Sprintf("%s%s", r.Resource.String(), r.Title)
}
