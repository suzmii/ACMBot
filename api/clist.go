package api

import (
	"fmt"
	"strings"

	"github.com/suzmii/ACMBot/consts"
	"github.com/suzmii/ACMBot/util"

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

type ClistRace struct {
	Resource string `json:"resource"`
	Event    string `json:"event"`
	Href     string `json:"href"`
	Start    string `json:"start"`
	End      string `json:"end"`
}

type ClistRaceFetcher func() ([]ClistRace, error)

func (api *API) fetchClistContests(sources []string) (races []ClistRace, err error) {
	races, err = fetchClistAPI[[]ClistRace](api.clistClient, "contest", map[string]any{
		"resource__in": strings.Join(sources, ","),
		"order_by":     "start",
		"upcoming":     true,
	})
	logger.Debugf("fetching clist contests, sources=%s, result=%v, err=%v", sources, races, err)
	return
}

func (api *API) FetchClistContests() ([]ClistRace, error) {
	return api.fetchClistContests([]string{
		consts.RaceResourceCodeforces,
		consts.RaceResourceAtcoder,
		consts.RaceResourceLeetcode,
		consts.RaceResourceLuogu,
		consts.RaceResourceNowcoder,
	})
}
