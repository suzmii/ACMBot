package provider

import (
	"github.com/suzmii/ACMBot/pkg/model/race"
)

type RaceProvider func() ([]race.Race, error)

type PicProvider func() ([]byte, error)
