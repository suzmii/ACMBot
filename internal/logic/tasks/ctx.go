package tasks

import (
	"github.com/suzmii/ACMBot/internal/logic/manager"
	"github.com/suzmii/ACMBot/pkg/model/bot"
	"github.com/suzmii/ACMBot/pkg/model/provider"
	"github.com/suzmii/ACMBot/pkg/model/race"
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
