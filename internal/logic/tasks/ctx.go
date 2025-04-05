package tasks

import (
	"github.com/suzmii/ACMBot/internal/logic/manager"
	"github.com/suzmii/ACMBot/internal/model/bot"
	"github.com/suzmii/ACMBot/internal/model/provider"
	"github.com/suzmii/ACMBot/internal/model/race"
)

type handles []string

type apiCaller = bot.ApiCaller

type races []race.Race

type textMessage string

type picMessage []byte

type Params []string

type codeforcesUser = *manager.CodeforcesUser

type atcoderUser = *manager.AtcoderUser

type raceProvider = provider.RaceProvider

type platform = bot.Platform
