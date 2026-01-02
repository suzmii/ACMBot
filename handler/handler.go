package handler

import (
	"github.com/suzmii/ACMBot/api"
	"github.com/suzmii/ACMBot/config/subconfig"
	"github.com/suzmii/ACMBot/database"
	"github.com/suzmii/ACMBot/render"
	"github.com/suzmii/ACMBot/util/logx"
)

var logger = logx.New("handler")

type Handler struct {
	cfg    subconfig.Handler
	store  database.Store
	api    *api.API
	render *render.Render
}

func NewHandler(cfg subconfig.Handler, store database.Store, api *api.API, render *render.Render) *Handler {

	return &Handler{
		cfg:    cfg,
		store:  store,
		api:    api,
		render: render,
	}
}
