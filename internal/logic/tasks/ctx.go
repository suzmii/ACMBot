package tasks

import (
	"github.com/suzmii/ACMBot/internal/model/bot"
	"github.com/suzmii/ACMBot/internal/model/provider"
	"github.com/suzmii/ACMBot/internal/model/race"
	"github.com/suzmii/ACMBot/internal/model/userinfo"
	"github.com/suzmii/ACMBot/internal/renderer"
)

// 此处类型别名主要目的是简化类型出处

type handles []string
type textMessage string

type picMessage []byte
type Params []string

type apiCaller = bot.ApiCaller

type races = []race.Race

type raceProvider = provider.RaceProvider

type platform = bot.Platform

type renderAble = renderer.RenderAble

type user = userinfo.UserInfo
