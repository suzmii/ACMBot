package client

import (
	"fmt"
	"github.com/suzmii/ACMBot/config"
	"github.com/suzmii/ACMBot/internal/model"
	"github.com/suzmii/ACMBot/internal/util"
	"time"

	"github.com/imroc/req/v3"
)

type clistResponse[T any] struct {
	Meta    any `json:"meta"`
	Objects T   `json:"objects"`
}

func fetchClistAPI[T any](client *req.Client, apiMethod string, args map[string]any) (T, error) {

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

type Contest struct {
	Resource string `json:"resource"`
	Event    string `json:"event"`
	Href     string `json:"href"`
	Start    string `json:"start"`
	End      string `json:"end"`
}

func (c Contest) ToRace() model.Race {
	s, err := time.Parse("2006-01-02T15:04:05", c.Start)
	if err != nil {
		s = time.Unix(0, 0)
	}
	e, err := time.Parse("2006-01-02T15:04:05", c.End)
	if err != nil {
		e = time.Unix(0, 0)
	}
	return model.Race{
		Source:    model.Resource(c.Resource),
		Name:      c.Event,
		Link:      c.Href,
		StartTime: s,
		EndTime:   e,
	}
}

var clistToken = config.LoadConfig().API.ClistAuthenticated

var (
	client *req.Client
)

func init() {
	client = req.C().SetCommonHeader("Authorization", clistToken)
}

func FetchClistContests(source model.Resource) ([]model.Race, error) {
	races, err := fetchClistAPI[[]Contest](client, "contest", map[string]any{
		"resource": source,
		"order_by": "start",
		"upcoming": true,
	})
	if err != nil {
		return nil, err
	}
	result := make([]model.Race, 0, len(races))
	for _, v := range races {
		result = append(result, v.ToRace())
	}
	return result, nil
}

func FetchClistCodeforcesContests() ([]model.Race, error) {
	return FetchClistContests(model.ResourceCodeforces)
}
func FetchClistAtcoderContests() ([]model.Race, error) {
	return FetchClistContests(model.ResourceAtcoder)
}
func FetchClistLeetcodeContests() ([]model.Race, error) {
	return FetchClistContests(model.ResourceLeetcode)
}
func FetchClistLuoguContests() ([]model.Race, error) {
	return FetchClistContests(model.ResourceLuogu)
}
func FetchClistNowcoderContests() ([]model.Race, error) {
	return FetchClistContests(model.ResourceNowcoder)
}
