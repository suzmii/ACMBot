package model

import (
	"fmt"
	"strings"
	"time"

	"github.com/suzmii/ACMBot/internal/database/dbmodel"
	"github.com/suzmii/ACMBot/internal/util"
)

type Resource string

const (
	ResourceCodeforces Resource = "datasync.com"
	ResourceAtcoder    Resource = "atcapi.jp"
	ResourceLeetcode   Resource = "leetcode.com"
	ResourceLuogu      Resource = "luogu.com.cn"
	ResourceNowcoder   Resource = "ac.nowcoder.com"
)

var ResourceName = map[Resource]string{
	ResourceCodeforces: "Codeforces",
	ResourceAtcoder:    "Atcoder",
	ResourceLeetcode:   "力扣",
	ResourceLuogu:      "洛谷",
	ResourceNowcoder:   "牛客",
}

func (r Resource) Name() string {
	return ResourceName[r]
}

var AllRaceResource = []Resource{ResourceCodeforces, ResourceLuogu, ResourceAtcoder, ResourceLeetcode, ResourceNowcoder}

type Race struct {
	ID        string    `json:"id"`
	Source    Resource  `json:"source"`
	Name      string    `json:"name"`
	Link      string    `json:"link"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

type Races []Race

func (rs Races) String() string {
	var b strings.Builder
	for _, r := range rs {
		b.WriteString(r.String())
		b.WriteString("\n\n")
	}
	return b.String()
}

// Conversions between DB models and common models
func ModelResourceFromDB(res dbmodel.Resource) Resource {
	switch res {
	case dbmodel.ResourceCodeforces:
		return ResourceCodeforces
	case dbmodel.ResourceAtcoder:
		return ResourceAtcoder
	case dbmodel.ResourceLeetcode:
		return ResourceLeetcode
	case dbmodel.ResourceLuogu:
		return ResourceLuogu
	case dbmodel.ResourceNowcoder:
		return ResourceNowcoder
	default:
		return ResourceCodeforces
	}
}

func RaceFromRepo(r dbmodel.Races) Race {
	return Race{
		ID:        r.ID,
		Source:    ModelResourceFromDB(r.Resource),
		Name:      r.Title,
		Link:      r.Link,
		StartTime: r.StartAt,
		EndTime:   r.EndAt,
	}
}

func RacesFromRepo(list []*dbmodel.Races) []Race {
	result := make([]Race, 0, len(list))
	for _, r := range list {
		if r == nil {
			continue
		}
		result = append(result, RaceFromRepo(*r))
	}
	return result
}

func (r *Race) String() string {
	d := r.EndTime.Sub(r.StartTime)
	var dStr string
	if h, m := int(d.Hours()), int(d.Minutes())%60; m > 0 {
		dStr = fmt.Sprintf("%d小时%d分钟", h, m)
	} else {
		dStr = fmt.Sprintf("%d小时", h)
	}

	startLeftTime := time.Until(r.StartTime)
	endLeftTime := time.Until(r.EndTime)

	started := startLeftTime.Milliseconds() < 0
	finished := endLeftTime.Milliseconds() < 0

	if !started {
		return fmt.Sprintf(
			""+
				"🕣此比赛尚未开始🕦\n"+
				"比赛来源: %s\n"+
				"比赛名称: %s\n"+
				"距离开始: %s\n"+
				"开始时间: %s\n"+
				"持续时间: %s\n"+
				"传送门🌈: %s",
			r.Source,
			r.Name,
			fmt.Sprintf("%02d天%02d小时%02d分钟", int(startLeftTime.Hours())/24, util.Abs(int(startLeftTime.Hours()))%24, util.Abs(int(startLeftTime.Minutes()))%60),
			r.StartTime.In(time.Local).Format("2006-01-02 15:04:05"),
			dStr,
			r.Link,
		)
	}
	if !finished {
		return fmt.Sprintf(
			""+
				"❗此比赛正在进行中❗\n"+
				"比赛来源: %s\n"+
				"比赛名称: %s\n"+
				"距离结束: %s\n"+
				"开始时间: %s\n"+
				"持续时间: %s\n"+
				"传送门🌈: %s",
			r.Source,
			r.Name,
			fmt.Sprintf("%02d天%02d小时%02d分钟", int(endLeftTime.Hours())/24, util.Abs(int(endLeftTime.Hours()))%24, util.Abs(int(endLeftTime.Minutes()))%60),
			r.StartTime.In(time.Local).Format("2006-01-02 15:04:05"),
			dStr,
			r.Link,
		)
	}
	return fmt.Sprintf(
		""+
			"此比赛已经结束了\n"+
			"比赛来源: %s\n"+
			"比赛名称: %s\n"+
			"开始时间: %s\n"+
			"持续时间: %s\n"+
			"传送门🌈: %s",
		r.Source,
		r.Name,
		r.StartTime.In(time.Local).Format("2006-01-02 15:04:05"),
		dStr,
		r.Link,
	)
}

func (r *Race) NoUrlString() string {
	d := r.EndTime.Sub(r.StartTime)
	var dStr string
	if h, m := int(d.Hours()), int(d.Minutes())%60; m > 0 {
		dStr = fmt.Sprintf("%d小时%d分钟", h, m)
	} else {
		dStr = fmt.Sprintf("%d小时", h)
	}

	startLeftTime := time.Until(r.StartTime)
	endLeftTime := time.Until(r.EndTime)

	started := startLeftTime.Milliseconds() < 0
	finished := endLeftTime.Milliseconds() < 0

	if !started {
		return fmt.Sprintf(
			""+
				"🕣此比赛尚未开始🕦\n"+
				"比赛来源: %s\n"+
				"比赛名称: %s\n"+
				"距离开始: %s\n"+
				"开始时间: %s\n"+
				"持续时间: %s\n",
			r.Source.Name(),
			r.Name,
			fmt.Sprintf("%02d天%02d小时%02d分钟", int(startLeftTime.Hours())/24, util.Abs(int(startLeftTime.Hours()))%24, util.Abs(int(startLeftTime.Minutes()))%60),
			r.StartTime.In(time.Local).Format("2006-01-02 15:04:05"),
			dStr,
		)
	}
	if !finished {
		return fmt.Sprintf(
			""+
				"❗此比赛正在进行中❗\n"+
				"比赛来源: %s\n"+
				"比赛名称: %s\n"+
				"距离结束: %s\n"+
				"开始时间: %s\n"+
				"持续时间: %s\n",
			r.Source.Name(),
			r.Name,
			fmt.Sprintf("%02d天%02d小时%02d分钟", int(endLeftTime.Hours())/24, util.Abs(int(endLeftTime.Hours()))%24, util.Abs(int(endLeftTime.Minutes()))%60),
			r.StartTime.In(time.Local).Format("2006-01-02 15:04:05"),
			dStr,
		)
	}
	return fmt.Sprintf(
		""+
			"此比赛已经结束了\n"+
			"比赛来源: %s\n"+
			"比赛名称: %s\n"+
			"开始时间: %s\n"+
			"持续时间: %s\n",
		r.Source.Name(),
		r.Name,
		r.StartTime.In(time.Local).Format("2006-01-02 15:04:05"),
		dStr,
	)
}

func (r *Race) Markdown() string {
	d := r.EndTime.Sub(r.StartTime)
	var dStr string
	if h, m := int(d.Hours()), int(d.Minutes())%60; m > 0 {
		dStr = fmt.Sprintf("%d小时%d分钟", h, m)
	} else {
		dStr = fmt.Sprintf("%d小时", h)
	}

	startLeftTime := time.Until(r.StartTime)
	endLeftTime := time.Until(r.EndTime)

	started := startLeftTime.Milliseconds() < 0
	finished := endLeftTime.Milliseconds() < 0

	if !started {
		return fmt.Sprintf(
			""+
				"> 🕣此比赛尚未开始🕦\n"+
				"- 比赛来源: %s\n"+
				"- 比赛名称: %s\n"+
				"- 距离开始: %s\n"+
				"- 开始时间: %s\n"+
				"- 持续时间: %s\n"+
				"- 传送门🌈: %s",
			r.Source,
			r.Name,
			fmt.Sprintf("%02d天%02d小时%02d分钟", int(startLeftTime.Hours())/24, util.Abs(int(startLeftTime.Hours()))%24, util.Abs(int(startLeftTime.Minutes()))%60),
			r.StartTime.In(time.Local).Format("2006-01-02 15:04:05"),
			dStr,
			r.Link,
		)
	}
	if !finished {
		return fmt.Sprintf(
			""+
				"> ❗此比赛正在进行中❗\n"+
				"- 比赛来源: %s\n"+
				"- 比赛名称: %s\n"+
				"- 距离结束: %s\n"+
				"- 开始时间: %s\n"+
				"- 持续时间: %s\n"+
				"- 传送门🌈: %s",
			r.Source,
			r.Name,
			fmt.Sprintf("%02d天%02d小时%02d分钟", int(endLeftTime.Hours())/24, util.Abs(int(endLeftTime.Hours()))%24, util.Abs(int(endLeftTime.Minutes()))%60),
			r.StartTime.In(time.Local).Format(time.DateTime),
			dStr,
			r.Link,
		)
	}
	return fmt.Sprintf(
		""+
			"> 此比赛已经结束了\n"+
			"- 比赛来源: %s\n"+
			"- 比赛名称: %s\n"+
			"- 开始时间: %s\n"+
			"- 持续时间: %s\n"+
			"- 传送门🌈: %s",
		r.Source,
		r.Name,
		r.StartTime.In(time.Local).Format("2006-01-02 15:04:05"),
		dStr,
		r.Link,
	)
}
