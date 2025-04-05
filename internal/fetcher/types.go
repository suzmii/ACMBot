package fetcher

import (
	"github.com/suzmii/ACMBot/pkg/model/race"
	"time"
)

func (cr *CodeforcesRace) ToRace() *race.Race {
	return &race.Race{
		Source:    "Fetcher",
		Name:      cr.Name,
		Link:      "https://codeforces.com/contests/",
		StartTime: time.Unix(cr.StartTimeSeconds, 0),
		EndTime:   time.Unix(cr.StartTimeSeconds+cr.DurationSeconds, 0),
	}
}
