package fetcher

import (
	"errors"
	"fmt"
	"github.com/suzmii/ACMBot/internal/model/race"
	"github.com/suzmii/ACMBot/internal/util"
	"time"

	"github.com/imroc/req/v3"
	"github.com/suzmii/ACMBot/conf"
)

var (
	apiKey = conf.GetConfig().Fetcher.ClistAuthenticated
	client = req.C().SetCommonHeader("Authorization", apiKey)
)

/*
   "duration": 5400,
   "end": "2020-06-14T04:00:00",
   "event": "Weekly Contest 193",
   "host": "leetcode.com",
   "href": "https://leetcode.com/contest/weekly-contest-193",
   "id": 20198406,
   "n_problems": 4,
   "n_statistics": 10545,
   "parsed_at": "2023-10-10T11:31:48.984866",
   "problems": null,
   "resource": "leetcode.com",
   "resource_id": 102,
   "start": "2020-06-14T02:30:00"
*/

type clistResponse[T any] struct {
	Meta    any `json:"meta"`
	Objects T   `json:"objects"`
}

func fetchClistAPI[T any](apiMethod string, args map[string]any) (T, error) {
	if apiKey == "" {
		return util.Zero[T](), errors.New("api key empty")
	}

	c := client.Clone()

	for k, v := range args {
		c.SetCommonQueryParam(k, fmt.Sprint(v))
	}

	const baseURL = "https://clist.by/api/v4/"
	fullURL := baseURL + apiMethod
	res, err := c.R().Get(fullURL)
	if err != nil {
		return util.Zero[T](), err
	}
	var result clistResponse[T]
	err = res.UnmarshalJson(&result)
	if err != nil {
		return util.Zero[T](), err
	}
	return result.Objects, nil
}

type ClistContest struct {
	Resource string `json:"resource"`
	Event    string `json:"event"`
	Href     string `json:"href"`
	Start    string `json:"start"`
	End      string `json:"end"`
}

func (c ClistContest) ToRace() race.Race {
	s, err := time.Parse("2006-01-02T15:04:05", c.Start)
	if err != nil {
		s = time.Unix(0, 0)
	}
	e, err := time.Parse("2006-01-02T15:04:05", c.End)
	if err != nil {
		e = time.Unix(0, 0)
	}
	return race.Race{
		Source:    race.Resource(c.Resource),
		Name:      c.Event,
		Link:      c.Href,
		StartTime: s,
		EndTime:   e,
	}
}

func FetchClistContests(source race.Resource) ([]race.Race, error) {
	races, err := fetchClistAPI[[]ClistContest]("contest", map[string]any{
		"resource": source,
		"order_by": "start",
		"upcoming": true,
	})
	if err != nil {
		return nil, err
	}
	result := make([]race.Race, 0, len(races))
	for _, v := range races {
		result = append(result, v.ToRace())
	}
	return result, nil
}

func FetchClistCodeforcesContests() ([]race.Race, error) {
	return FetchClistContests(race.ResourceCodeforces)
}
func FetchClistAtcoderContests() ([]race.Race, error) {
	return FetchClistContests(race.ResourceAtcoder)
}
func FetchClistLeetcodeContests() ([]race.Race, error) {
	return FetchClistContests(race.ResourceLeetcode)
}
func FetchClistLuoguContests() ([]race.Race, error) {
	return FetchClistContests(race.ResourceLuogu)
}
func FetchClistNowcoderContests() ([]race.Race, error) {
	return FetchClistContests(race.ResourceNowcoder)
}
