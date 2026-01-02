package render

import (
	"github.com/suzmii/ACMBot/config/subconfig"
	"github.com/suzmii/ACMBot/render/internal"
	"github.com/suzmii/ACMBot/util/logx"
)

var logger = logx.New("render")

type Render struct {
	*internal.Render
}

func NewRender(cfg subconfig.Render) (*Render, error) {
	render, err := internal.New(cfg)
	if err != nil {
		return nil, err
	}
	return &Render{
		Render: render,
	}, nil
}
