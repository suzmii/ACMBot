package database

import (
	"github.com/suzmii/ACMBot/internal/database/dbmodel"
)

var AllModels []interface{}

func init() {
	AllModels = append(AllModels, dbmodel.CodeforcesModels...)
	AllModels = append(AllModels, new(dbmodel.Races))
}
