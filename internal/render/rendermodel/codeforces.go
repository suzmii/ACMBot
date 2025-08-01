package rendermodel

import (
	"github.com/suzmii/ACMBot/internal/model"
	"html/template"
)

type CodeforcesRatingChange struct {
	At        int64 `json:"at"`
	NewRating int   `json:"newRating"`
}

func CodeforcesRatingChangeFromCommon(r model.RatingRecord) CodeforcesRatingChange {
	var c CodeforcesRatingChange
	c.NewRating = r.Rating
	c.At = r.At.Unix()
	return c
}

type CodeforcesRatingRecords struct {
	Data   []CodeforcesRatingChange
	Handle string

	EchartsJS template.JS
}

func CodeforcesRatingRecordsFromCommon(username string, records []model.RatingRecord) CodeforcesRatingRecords {
	var r CodeforcesRatingRecords
	r.Handle = username
	for _, record := range records {
		r.Data = append(r.Data, CodeforcesRatingChangeFromCommon(record))
	}

	return r
}

type CodeforcesRatingLevel string

const (
	CodeforcesRatingLevelNewbie                   CodeforcesRatingLevel = "Newbie"
	CodeforcesRatingLevelPupil                    CodeforcesRatingLevel = "Pupil"
	CodeforcesRatingLevelSpecialist               CodeforcesRatingLevel = "Specialist"
	CodeforcesRatingLevelExpert                   CodeforcesRatingLevel = "Expert"
	CodeforcesRatingLevelCandidateMaster          CodeforcesRatingLevel = "CM"
	CodeforcesRatingLevelMaster                   CodeforcesRatingLevel = "Master"
	CodeforcesRatingLevelInternationalMaster      CodeforcesRatingLevel = "IM"
	CodeforcesRatingLevelGrandmaster              CodeforcesRatingLevel = "GM"
	CodeforcesRatingLevelInternationalGrandmaster CodeforcesRatingLevel = "IGM"
	CodeforcesRatingLevelLegendaryGrandmaster     CodeforcesRatingLevel = "LGM"
	CodeforcesRatingLevelTourist                  CodeforcesRatingLevel = "Tourist"
)

func Rating2Level(rating int) CodeforcesRatingLevel {
	switch {
	case rating < 1200:
		return CodeforcesRatingLevelNewbie
	case rating < 1400:
		return CodeforcesRatingLevelPupil
	case rating < 1600:
		return CodeforcesRatingLevelSpecialist
	case rating < 1900:
		return CodeforcesRatingLevelExpert
	case rating < 2100:
		return CodeforcesRatingLevelCandidateMaster
	case rating < 2300:
		return CodeforcesRatingLevelMaster
	case rating < 2400:
		return CodeforcesRatingLevelInternationalMaster
	case rating < 2600:
		return CodeforcesRatingLevelGrandmaster
	case rating < 3000:
		return CodeforcesRatingLevelInternationalGrandmaster
	case rating < 4000:
		return CodeforcesRatingLevelLegendaryGrandmaster
	default:
		return CodeforcesRatingLevelTourist
	}
}

type CodeforcesUserSolvedData struct {
	Range   string
	Percent float32
}

type CodeforcesUserProfile struct {
	Avatar    string
	Handle    string
	MaxRating int
	FriendOf  int
	Rating    int
	Level     CodeforcesRatingLevel
	Solved    int

	SolvedData []CodeforcesUserSolvedData

	TailwindJS template.JS
	FontCSS    template.CSS
}

func CodeforcesUserProfileFromCommon(user model.User, solved model.SolvedDetail) CodeforcesUserProfile {
	solveCnt := solved.Range800to1400 + solved.Range1400to2000 + solved.Range2000to2600 + solved.Above2600
	solveCntFloat := float32(solveCnt)
	return CodeforcesUserProfile{
		Avatar:    user.Avatar,
		Handle:    user.Name,
		MaxRating: user.MaxRating,
		FriendOf:  user.FriendOf,
		Rating:    user.CurrentRating,
		Level:     Rating2Level(user.CurrentRating),
		Solved:    solveCnt,
		SolvedData: []CodeforcesUserSolvedData{
			{
				Range:   "800+",
				Percent: float32(solved.Range800to1400) / solveCntFloat * 100,
			},
			{
				Range:   "1400+",
				Percent: float32(solved.Range1400to2000) / solveCntFloat * 100,
			},
			{
				Range:   "2000+",
				Percent: float32(solved.Range2000to2600) / solveCntFloat * 100,
			},
			{
				Range:   "2600+",
				Percent: float32(solved.Above2600) / solveCntFloat * 100,
			},
		},
	}
}
