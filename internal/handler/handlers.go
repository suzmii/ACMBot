package handler

import (
	"github.com/suzmii/ACMBot/internal/model"
)

//var (
//	Commands = []Event{
//		{[]string{"近期比赛"}, tasks.RaceHandler(manager.GetAllCachedRaceProvider())},
//		{[]string{"近期cf"}, tasks.RaceHandler(manager.GetRaceProviderByResource(racemodel.ResourceCodeforces))},
//		{[]string{"近期atc"}, tasks.RaceHandler(manager.GetRaceProviderByResource(racemodel.ResourceAtcoder))},
//		{[]string{"近期nk"}, tasks.RaceHandler(manager.GetRaceProviderByResource(racemodel.ResourceNowcoder))},
//		{[]string{"近期lg"}, tasks.RaceHandler(manager.GetRaceProviderByResource(racemodel.ResourceLuogu))},
//
//		{[]string{"cf"}, tasks.CodeforcesProfileHandler},
//		{[]string{"rt", "rating"}, tasks.CodeforcesRatingHandler},
//		{[]string{"at"}, tasks.AtcoderProfileHandler},
//
//		{[]string{"help", "菜单"}, tasks.TextHandler(MenuText)},
//	}
//)

var Events = []model.Event{
	{
		Commands: []string{"rating"},
		Handler:  RatingDetailHandler,
	},
	{
		Commands: []string{"cf"},
		Handler:  ProfileHandler,
	},
	{
		Commands: []string{"近期比赛"},
		Handler:  FetchContest,
	},
	{
		Commands: []string{"近期nowcoder", "近期nc", "近期牛客"},
		Handler:  FetchNowcoderContest,
	},
	{
		Commands: []string{"近期cf"},
		Handler:  FetchCodeforcesContest,
	},
	{
		Commands: []string{"近期luogu", "近期lg"},
		Handler:  FetchLuoguContest,
	},
	{
		Commands: []string{"近期leetcode", "近期lc"},
		Handler:  FetchLeetcodeContest,
	},
	{
		Commands: []string{"近期atc", "近期atcoder"},
		Handler:  FetchAtcoderContest,
	},
}
