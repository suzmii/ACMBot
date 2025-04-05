package fetcher

import (
	"time"

	"github.com/suzmii/ACMBot/app/model"
)

func (cr *CodeforcesRace) ToRace() *model.Race {
	return &model.Race{
		Source:    "Fetcher",
		Name:      cr.Name,
		Link:      "https://codeforces.com/contests/",
		StartTime: time.Unix(cr.StartTimeSeconds, 0),
		EndTime:   time.Unix(cr.StartTimeSeconds+cr.DurationSeconds, 0),
	}
}
