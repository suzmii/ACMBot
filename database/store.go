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
	GetRandomCodeforcesProblem(ctx context.Context, filter CodeforcesProblemFilter) (*CodeforcesProblem, error)

	// Atcoder相关
	CreateAtcoderUser(ctx context.Context, params *CreateAtcoderUserParams) (*AtcoderUserWithRecords, error)
	GetAtcoderUserByUsername(ctx context.Context, username string) (*AtcoderUserWithRecords, error)
	GetAtcoderUserByID(ctx context.Context, userID int64) (*AtcoderUserWithRecords, error)
	UpdateAtcoderSubmissionStatistics(ctx context.Context, userId int) (*AtcoderSubmissionStatistics, error)
	CreateAtcoderSubmissions(ctx context.Context, submissions []AtcoderSubmission) error
}

type CodeforcesProblemFilter struct {
	MinRating *int
	MaxRating *int
	Tags      []string
}
