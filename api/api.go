package api

import (
	"errors"

	"github.com/gocolly/colly/v2"
	"github.com/imroc/req/v3"
	"github.com/suzmii/ACMBot/config/subconfig"
	"github.com/suzmii/ACMBot/util/logx"
)

var logger = logx.New("api")

type API struct {
	cfg           subconfig.API
	clistClient   *req.Client
	atcoderClient *colly.Collector
}

func New(cfg subconfig.API) (*API, error) {
	if cfg.ClistAuthenticated == "" {
		return nil, errors.New("clist token is empty")
	}
	if cfg.CodeforcesKey == "" || cfg.CodeforcesSecret == "" {
		return nil, errors.New("codeforces token or secret is empty")
	}

	atcoderClient := colly.NewCollector(
		colly.AllowedDomains("atcoder.jp", "www.atcoder.jp", "img.atcoder.jp"),
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3"),
	)

	return &API{
		cfg:           cfg,
		clistClient:   req.C().SetCommonHeader("Authorization", cfg.ClistAuthenticated),
		atcoderClient: atcoderClient,
	}, nil
}
