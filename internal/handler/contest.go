package handler

import (
	"sort"

	"github.com/suzmii/ACMBot/internal/database/repo"
	"github.com/suzmii/ACMBot/internal/datasync"
	"github.com/suzmii/ACMBot/internal/model"
)

func ContestHandler(c *model.Context) error {
	if err := datasync.SyncRacesFromRemoteIfStale(c.Ctx); err != nil {
		return err
	}
	stored, err := repo.GetRaces(c.Ctx)
	if err != nil {
		return err
	}
	contestList := model.RacesFromRepo(stored)
	sort.Slice(contestList, func(i, j int) bool { return contestList[i].StartTime.Before(contestList[j].StartTime) })
	text := "近期比赛:\n" + model.Races(contestList).String()
	c.Messager.Send(model.TextMessage{Text: text})
	return nil
}

func NowcoderContestHandler(c *model.Context) error {
	if err := datasync.SyncRacesFromRemoteIfStale(c.Ctx); err != nil {
		return err
	}
	stored, err := repo.GetRaces(c.Ctx)
	if err != nil {
		return err
	}
	var contestList []model.Race
	for _, r := range model.RacesFromRepo(stored) {
		if r.Source == model.ResourceNowcoder {
			contestList = append(contestList, r)
		}
	}
	text := "近期比赛:\n" + model.Races(contestList).String()
	c.Messager.Send(model.TextMessage{Text: text})
	return nil
}

func LuoguContestHandler(c *model.Context) error {
	if err := datasync.SyncRacesFromRemoteIfStale(c.Ctx); err != nil {
		return err
	}
	stored, err := repo.GetRaces(c.Ctx)
	if err != nil {
		return err
	}
	var list []model.Race
	for _, r := range model.RacesFromRepo(stored) {
		if r.Source == model.ResourceLuogu {
			list = append(list, r)
		}
	}
	text := "近期比赛:\n" + model.Races(list).String()
	c.Messager.Send(model.TextMessage{Text: text})
	return nil
}

func LeetcodeContestHandler(c *model.Context) error {
	if err := datasync.SyncRacesFromRemoteIfStale(c.Ctx); err != nil {
		return err
	}
	stored, err := repo.GetRaces(c.Ctx)
	if err != nil {
		return err
	}
	var list []model.Race
	for _, r := range model.RacesFromRepo(stored) {
		if r.Source == model.ResourceLeetcode {
			list = append(list, r)
		}
	}
	text := "近期比赛:\n" + model.Races(list).String()
	c.Messager.Send(model.TextMessage{Text: text})
	return nil
}

func AtcoderContestHandler(c *model.Context) error {
	if err := datasync.SyncRacesFromRemoteIfStale(c.Ctx); err != nil {
		return err
	}
	stored, err := repo.GetRaces(c.Ctx)
	if err != nil {
		return err
	}
	var list []model.Race
	for _, r := range model.RacesFromRepo(stored) {
		if r.Source == model.ResourceAtcoder {
			list = append(list, r)
		}
	}
	text := "近期比赛:\n" + model.Races(list).String()
	c.Messager.Send(model.TextMessage{Text: text})
	return nil
}

func CodeforcesContestHandler(c *model.Context) error {
	if err := datasync.SyncRacesFromRemoteIfStale(c.Ctx); err != nil {
		return err
	}
	stored, err := repo.GetRaces(c.Ctx)
	if err != nil {
		return err
	}
	var list []model.Race
	for _, r := range model.RacesFromRepo(stored) {
		if r.Source == model.ResourceCodeforces {
			list = append(list, r)
		}
	}
	text := "近期比赛:\n" + model.Races(list).String()
	c.Messager.Send(model.TextMessage{Text: text})
	return nil
}
