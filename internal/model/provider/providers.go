package provider

import (
	"github.com/suzmii/ACMBot/internal/model/race"
)

type RaceProvider func() ([]race.Race, error)

type PicProvider func() ([]byte, error)
