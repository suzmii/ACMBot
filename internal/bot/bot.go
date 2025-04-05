package bot

import (
	"context"

	"github.com/suzmii/ACMBot/conf"
	"github.com/suzmii/ACMBot/internal/logic/manager"
	"github.com/suzmii/ACMBot/internal/logic/tasks"
	"github.com/suzmii/ACMBot/pkg/model/race"
)

var (
	CommandPrefix = conf.GetConfig().Bot.CommandPrefix

	MenuText = `以下是功能列表：所有命令都要加上前缀` + "`" + CommandPrefix + "`" + `哦🥰
0. help(或菜单)，输出本消息

1. cf/at [username]，用于查询codeforces/atcoder用户的基本信息

2. rating(或rt) [username]，用于查询codeforces用户的rating变化曲线

3. 近期[比赛,atc,nk,lg,cf]，用于查询近期的比赛数据，数据来源于clist.by`
)

type CommandHandler struct {
	Commands []string
	Handler  func(context.Context) error
}

var (
	Commands = []CommandHandler{
		{[]string{"近期比赛"}, tasks.RaceHandler(manager.GetAllCachedRaces)},
		{[]string{"近期cf"}, tasks.RaceHandler(manager.GetCachedRacesByResource(race.ResourceCodeforces))},
		{[]string{"近期atc"}, tasks.RaceHandler(manager.GetCachedRacesByResource(race.ResourceAtcoder))},
		{[]string{"近期nk"}, tasks.RaceHandler(manager.GetCachedRacesByResource(race.ResourceNowcoder))},
		{[]string{"近期lg"}, tasks.RaceHandler(manager.GetCachedRacesByResource(race.ResourceLuogu))},

		{[]string{"cf"}, tasks.CodeforcesProfileHandler},
		{[]string{"rt", "rating"}, tasks.CodeforcesRatingHandler},
		{[]string{"at"}, tasks.AtcoderProfileHandler},

		{[]string{"help", "菜单"}, tasks.TextHandler(MenuText)},
	}
)
