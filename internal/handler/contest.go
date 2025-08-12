package handler

import (
	"fmt"
	"sort"

	"github.com/suzmii/ACMBot/internal/api/client"
	"github.com/suzmii/ACMBot/internal/datasync"
	"github.com/suzmii/ACMBot/internal/model"
)

func FetchContest(c *model.Context) error {
	contestList, err := datasync.GetRaces(c.Ctx)
	if err != nil {
		return err
	}
	sort.Slice(contestList, func(i, j int) bool {
		return contestList[i].StartTime.Before(contestList[j].StartTime)
	})
	c.Messager.Send(model.TextMessage{Text: fmt.Sprintf("近期比赛: %v", contestList)})
	return nil
}

func FetchNowcoderContest(c *model.Context) error {
	contestList, err := client.FetchClistNowcoderContests()
	if err != nil {
		return err
	}
	c.Messager.Send(model.TextMessage{Text: fmt.Sprintf("近期比赛: %v", contestList)})
	return nil
}

func FetchLuoguContest(c *model.Context) error {
	contestList, err := client.FetchClistLuoguContests()
	if err != nil {
		return err
	}
	c.Messager.Send(model.TextMessage{Text: fmt.Sprintf("近期比赛: %v", contestList)})
	return nil
}

func FetchLeetcodeContest(c *model.Context) error {
	contestList, err := client.FetchClistLeetcodeContests()
	if err != nil {
		return err
	}
	c.Messager.Send(model.TextMessage{Text: fmt.Sprintf("近期比赛: %v", contestList)})
	return nil
}

func FetchAtcoderContest(c *model.Context) error {
	contestList, err := client.FetchClistAtcoderContests()
	if err != nil {
		return err
	}
	c.Messager.Send(model.TextMessage{Text: fmt.Sprintf("近期比赛: %v", contestList)})
	return nil
}

func FetchCodeforcesContest(c *model.Context) error {
	contestList, err := client.FetchClistCodeforcesContests()
	if err != nil {
		return err
	}
	c.Messager.Send(model.TextMessage{Text: fmt.Sprintf("近期比赛: %v", contestList)})
	return nil
}
