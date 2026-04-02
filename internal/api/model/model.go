package model

import (
	"fmt"
	"time"

	"github.com/suzmii/ACMBot/internal/model"
)

type User struct {
	Handle       string `json:"handle"`
	Avatar       string `json:"titlePhoto"`
	Rating       int    `json:"rating"`
	MaxRating    int    `json:"maxRating"`
	FriendCount  int    `json:"friendOfCount"`
	Organization string `json:"organization"`
	CreatedAt    int64  `json:"registrationTimeSeconds"`
	MaxRank      string `json:"maxRank"` // 称号
}

type Problem struct {
	ContestID      int      `json:"contestId"`
	ProblemSetName string   `json:"problemsetName"`
	Index          string   `json:"index"`
	Rating         int      `json:"rating"`
	Tags           []string `json:"tags"`
}

func (p Problem) ID() string {
	if p.ContestID == 0 {
		return p.ProblemSetName + p.Index
	}
	return fmt.Sprintf("%d%s", p.ContestID, p.Index)
}

type Submission struct {
	ID       uint    `json:"id"`
	At       int64   `json:"creationTimeSeconds"`
	Status   string  `json:"verdict"`
	Problem  Problem `json:"problem"`
	Language string  `json:"programmingLanguage"`
}

type RatingChange struct {
	At        int64 `json:"ratingUpdateTimeSeconds"`
	NewRating int   `json:"newRating"`
}

func (r *RatingChange) ToCommon() model.RatingRecord {
	return model.RatingRecord{
		At:     time.Unix(r.At, 0).Local(),
		Rating: r.NewRating,
	}
}

type CodeforcesRace struct {
	ID                  int    `json:"id"`
	Name                string `json:"name"`
	Type                string `json:"type"`
	Phase               string `json:"phase"`
	Frozen              bool   `json:"frozen"`
	DurationSeconds     int64  `json:"durationSeconds"`
	StartTimeSeconds    int64  `json:"startTimeSeconds"`
	RelativeTimeSeconds int64  `json:"relativeTimeSeconds"`
}
