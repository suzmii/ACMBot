package api

import (
	"errors"

	"github.com/imroc/req/v3"
	"github.com/suzmii/ACMBot/config/subconfig"
	"github.com/suzmii/ACMBot/util/logx"
)

var logger = logx.New("api")

type API struct {
	cfg         subconfig.API
	clistClient *req.Client
}

func New(cfg subconfig.API) (*API, error) {
	if cfg.ClistAuthenticated == "" {
		return nil, errors.New("clist token is empty")
	}
	if cfg.CodeforcesKey == "" || cfg.CodeforcesSecret == "" {
		return nil, errors.New("codeforces token or secret is empty")
	}

	return &API{
		cfg:         cfg,
		clistClient: req.C().SetCommonHeader("Authorization", cfg.ClistAuthenticated),
	}, nil
}
