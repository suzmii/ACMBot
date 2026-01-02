package transer

import (
	"time"

	"github.com/suzmii/ACMBot/api"
	"github.com/suzmii/ACMBot/database"
)

func RaceApi2DB(c api.ClistRace) database.Race {
	s, err := time.Parse("2006-01-02T15:04:05", c.Start)
	if err != nil {
		s = time.Unix(0, 0)
	}
	e, err := time.Parse("2006-01-02T15:04:05", c.End)
	if err != nil {
		e = time.Unix(0, 0)
	}

	return database.Race{
		Title:    c.Event,
		Resource: c.Resource,
		Start:    s,
		End:      e,
		Link:     c.Href,
	}
}
