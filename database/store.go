package database

import (
	"context"

	"github.com/suzmii/ACMBot/database/sqlc"
)

type Store interface {
	sqlc.Querier

	CreateRace(ctx context.Context, races []Race) error
	GetLastRace(ctx context.Context) ([]Race, error)

	// Codeforces相关
	CreateCodeforcesUser(ctx context.Context, params *CreateCodeforcesUserParams) (*CodeforcesUserWithRecords, error)
	GetCodeforcesUserByUsername(ctx context.Context, username string) (*CodeforcesUserWithRecords, error)
	GetCodeforcesUserByID(ctx context.Context, userID int64) (*CodeforcesUserWithRecords, error)
	UpdateCodeforcesRatingRecords(ctx context.Context, userId int, originRecords []CodeforcesRatingRecord) (*CodeforcesRatingRecords, error)
	UpdateCodeforcesSubmissionStatistics(ctx context.Context, userId int) (*CodeforcesSubmissionStatistics, error)
	CreateCodeforcesSubmissions(ctx context.Context, submissions []CodeforcesSubmission) error
}
