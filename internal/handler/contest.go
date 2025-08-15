package handler

import (
	"fmt"

	"github.com/suzmii/ACMBot/internal/datasync"
	"github.com/suzmii/ACMBot/internal/model"
)

func ContestHandler(c *model.Context) error {
	contestList, err := datasync.GetRaces(c.Ctx)
	if err != nil {
		return err
	}
	c.Messager.Send(model.TextMessage{Text: fmt.Sprintf("近期比赛: %v", contestList)})
	return nil
}

func NowcoderContestHandler(c *model.Context) error {
	list, err := datasync.GetRaces(c.Ctx)
	if err != nil {
		return err
	}
	var contestList []model.Race
	for _, r := range list {
		if r.Source == model.ResourceNowcoder {
			contestList = append(contestList, r)
		}
	}
	c.Messager.Send(model.TextMessage{Text: fmt.Sprintf("近期比赛: %v", contestList)})
	return nil
}

func LuoguContestHandler(c *model.Context) error {
	list, err := datasync.GetRaces(c.Ctx)
	if err != nil {
		return err
	}
	var contestList []model.Race
	for _, r := range list {
		if r.Source == model.ResourceLuogu {
			contestList = append(contestList, r)
		}
	}
	c.Messager.Send(model.TextMessage{Text: fmt.Sprintf("近期比赛: %v", contestList)})
	return nil
}

func LeetcodeContestHandler(c *model.Context) error {
	list, err := datasync.GetRaces(c.Ctx)
	if err != nil {
		return err
	}
	var contestList []model.Race
	for _, r := range list {
		if r.Source == model.ResourceLeetcode {
			contestList = append(contestList, r)
		}
	}
	c.Messager.Send(model.TextMessage{Text: fmt.Sprintf("近期比赛: %v", contestList)})
	return nil
}

func AtcoderContestHandler(c *model.Context) error {
	list, err := datasync.GetRaces(c.Ctx)
	if err != nil {
		return err
	}
	var contestList []model.Race
	for _, r := range list {
		if r.Source == model.ResourceAtcoder {
			contestList = append(contestList, r)
		}
	}
	c.Messager.Send(model.TextMessage{Text: fmt.Sprintf("近期比赛: %v", contestList)})
	return nil
}

func CodeforcesContestHandler(c *model.Context) error {
	list, err := datasync.GetRaces(c.Ctx)
	if err != nil {
		return err
	}
	var contestList []model.Race
	for _, r := range list {
		if r.Source == model.ResourceCodeforces {
			contestList = append(contestList, r)
		}
	}
	c.Messager.Send(model.TextMessage{Text: fmt.Sprintf("近期比赛: %v", contestList)})
	return nil
}
