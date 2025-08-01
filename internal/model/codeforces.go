package model

import (
	"github.com/suzmii/ACMBot/internal/database/dbmodel"
	"gorm.io/gorm"
	"time"
)

type RatingRecord struct {
	At     time.Time
	Rating int
}

func RatingRecordFromRepo(repo dbmodel.CodeforcesRatingRecord) RatingRecord {
	return RatingRecord{
		At:     repo.At,
		Rating: repo.Rating,
	}
}

func (r RatingRecord) ToRepo() dbmodel.CodeforcesRatingRecord {
	return dbmodel.CodeforcesRatingRecord{
		At:     r.At,
		Rating: r.Rating,
	}
}

type RatingRecords []RatingRecord

func RatingRecordsFromRepo(repo dbmodel.CodeforcesRatingRecords) RatingRecords {
	var ratingRecords RatingRecords
	for _, r := range repo.Records {
		ratingRecords = append(ratingRecords, RatingRecordFromRepo(r))
	}
	return ratingRecords
}

func (r RatingRecords) ToRepo() dbmodel.CodeforcesRatingRecords {
	if len(r) == 0 {
		return dbmodel.CodeforcesRatingRecords{}
	}
	result := dbmodel.CodeforcesRatingRecords{
		Records: make([]dbmodel.CodeforcesRatingRecord, 0, len(r)),
	}
	for _, r := range r {
		result.Records = append(result.Records, r.ToRepo())
	}
	return result
}

type User struct {
	ID           uint
	Name         string
	Avatar       string
	Organization string
	FriendOf     int

	MaxRating     int
	CurrentRating int

	SubmissionUpdatedAt time.Time
	UpdatedAt           time.Time
}

func UserFromRepo(repo dbmodel.CodeforcesUser) User {
	return User{
		ID:                  repo.ID,
		Name:                repo.Username,
		Avatar:              repo.AvatarURL,
		Organization:        repo.Organization,
		FriendOf:            repo.FriendOf,
		MaxRating:           repo.MaxRating,
		CurrentRating:       repo.CurrentRating,
		SubmissionUpdatedAt: repo.SubmissionUpdatedAt,
		UpdatedAt:           repo.UpdatedAt,
	}
}

func (u User) ToRepo() dbmodel.CodeforcesUser {
	return dbmodel.CodeforcesUser{
		Model: gorm.Model{
			ID:        u.ID,
			UpdatedAt: u.UpdatedAt,
		},
		Username:            u.Name,
		AvatarURL:           u.Avatar,
		Organization:        u.Organization,
		FriendOf:            u.FriendOf,
		MaxRating:           u.MaxRating,
		CurrentRating:       u.CurrentRating,
		SubmissionUpdatedAt: u.SubmissionUpdatedAt,
	}
}

type SolvedDetail struct {
	Range800to1400  int
	Range1400to2000 int
	Range2000to2600 int
	Above2600       int
}
