package database

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/suzmii/ACMBot/database/sqlc"
	"github.com/suzmii/ACMBot/util/jsonx"
)

type AtcoderProblemForQuery struct {
	ID    string  `json:"id"`
	Point float64 `json:"point"`
}

type AtcoderSubmissionStatistics struct {
	TotalCount int `json:"total_count"`

	Ac []AtcoderProblemForQuery `json:"ac"`

	LastSubmissionAt time.Time `json:"last_submission_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	Version          int       `json:"version"`
}

type AtcoderSubmission struct {
	UserID       int64
	SubmissionID int64
	ProblemID    string
	Point        float64
	Status       string
	At           time.Time
}

func (db *dbStore) CreateAtcoderSubmissions(ctx context.Context, submissions []AtcoderSubmission) error {
	if len(submissions) == 0 {
		return nil
	}

	userIDs := make([]int64, len(submissions))
	submissionIDs := make([]int64, len(submissions))
	problemIDs := make([]string, len(submissions))
	points := make([]float64, len(submissions))
	statuses := make([]string, len(submissions))
	ats := make([]pgtype.Timestamptz, len(submissions))

	slices.SortFunc(submissions, func(a, b AtcoderSubmission) int {
		return int(a.At.Sub(b.At))
	})

	for i, submission := range submissions {
		userIDs[i] = submission.UserID
		submissionIDs[i] = submission.SubmissionID
		problemIDs[i] = submission.ProblemID
		points[i] = submission.Point
		statuses[i] = submission.Status
		ats[i] = pgtype.Timestamptz{Time: submission.At, Valid: true}
	}

	err := db.Querier.CreateAtcoderSubmissionsRaw(ctx, sqlc.CreateAtcoderSubmissionsRawParams{
		Column1: userIDs,
		Column2: submissionIDs,
		Column3: problemIDs,
		Column4: points,
		Column5: statuses,
		Column6: ats,
	})
	if err != nil {
		return fmt.Errorf("failed to create atcoder submissions: %v", err)
	}

	return nil
}

func (db *dbStore) UpdateAtcoderSubmissionStatistics(ctx context.Context, userId int) (*AtcoderSubmissionStatistics, error) {
	user, err := db.GetAtcoderUserByID(ctx, int64(userId))
	if err != nil {
		return nil, fmt.Errorf("failed to get atcoder user: %v", err)
	}

	lastSubmissionAt := user.SubmissionStatistics.LastSubmissionAt
	newSubmissions, err := db.Querier.GetAtcoderSubmissionsAfter(ctx, sqlc.GetAtcoderSubmissionsAfterParams{
		UserID: user.ID,
		At:     pgtype.Timestamptz{Time: lastSubmissionAt, Valid: true},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get submissions after timestamp: %v", err)
	}

	acProblemMap := make(map[string]AtcoderProblemForQuery)
	for _, acProblem := range user.SubmissionStatistics.Ac {
		acProblemMap[acProblem.ID] = acProblem
	}

	newTotalCount := user.SubmissionStatistics.TotalCount
	for _, submission := range newSubmissions {
		newTotalCount++

		if submission.At.Time.Sub(lastSubmissionAt) > 0 {
			lastSubmissionAt = submission.At.Time
		}

		if submission.Status == "AC" {
			acProblemMap[submission.ProblemID] = AtcoderProblemForQuery{
				ID:    submission.ProblemID,
				Point: submission.Point,
			}
		}
	}

	acList := make([]AtcoderProblemForQuery, 0, len(acProblemMap))
	for _, problem := range acProblemMap {
		acList = append(acList, problem)
	}

	newStats := AtcoderSubmissionStatistics{
		TotalCount:       newTotalCount,
		Ac:               acList,
		LastSubmissionAt: lastSubmissionAt,
		UpdatedAt:        time.Now(),
		Version:          user.SubmissionStatistics.Version,
	}

	jsonb, err := jsonx.Marshal(newStats)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal submission statistics: %v", err)
	}

	err = db.Querier.UpdateAtcoderSubmissionStatisticsRaw(ctx, sqlc.UpdateAtcoderSubmissionStatisticsRawParams{
		ID:                   user.ID,
		SubmissionStatistics: jsonb,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update submission statistics: %v", err)
	}

	return &newStats, nil
}

type AtcoderUserWithRecords struct {
	sqlc.AtcoderUser
	SubmissionStatistics *AtcoderSubmissionStatistics
}

type CreateAtcoderUserParams struct {
	Username         string `json:"username"`
	AvatarUrl        string `json:"avatar_url"`
	Rank             string `json:"rank"`
	Rating           int32  `json:"rating"`
	HighestRating    int32  `json:"highest_rating"`
	PromotionMessage string `json:"promotion_message"`
}

func (db *dbStore) CreateAtcoderUser(ctx context.Context, params *CreateAtcoderUserParams) (*AtcoderUserWithRecords, error) {
	result := &AtcoderUserWithRecords{
		SubmissionStatistics: &AtcoderSubmissionStatistics{
			TotalCount:       0,
			Ac:               []AtcoderProblemForQuery{},
			LastSubmissionAt: time.Time{},
			UpdatedAt:        time.Time{},
			Version:          1,
		},
	}

	ssjsonb, err := jsonx.Marshal(result.SubmissionStatistics)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal submission statistics: %w", err)
	}

	user, err := db.CreateAtcoderUserRaw(ctx, sqlc.CreateAtcoderUserRawParams{
		Username:             params.Username,
		AvatarUrl:            params.AvatarUrl,
		Rank:                 params.Rank,
		Rating:               params.Rating,
		HighestRating:        params.HighestRating,
		PromotionMessage:     params.PromotionMessage,
		SubmissionStatistics: ssjsonb,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	result.AtcoderUser = user
	return result, nil
}

func (db *dbStore) parseAtcoderUser(user sqlc.AtcoderUser) (*AtcoderUserWithRecords, error) {
	result := &AtcoderUserWithRecords{AtcoderUser: user}

	if len(user.SubmissionStatistics) > 0 {
		var stats AtcoderSubmissionStatistics
		if err := jsonx.Unmarshal(user.SubmissionStatistics, &stats); err != nil {
			return nil, fmt.Errorf("failed to unmarshal submission statistics: %w", err)
		}
		result.SubmissionStatistics = &stats
	}

	if result.SubmissionStatistics == nil {
		result.SubmissionStatistics = &AtcoderSubmissionStatistics{
			TotalCount:       0,
			Ac:               []AtcoderProblemForQuery{},
			LastSubmissionAt: time.Time{},
			UpdatedAt:        time.Time{},
			Version:          1,
		}
	}

	return result, nil
}

func (db *dbStore) GetAtcoderUserByID(ctx context.Context, userID int64) (*AtcoderUserWithRecords, error) {
	user, err := db.Querier.GetAtcoderUserByIDRaw(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get atcoder user by id: %w", err)
	}
	return db.parseAtcoderUser(user)
}

func (db *dbStore) GetAtcoderUserByUsername(ctx context.Context, username string) (*AtcoderUserWithRecords, error) {
	user, err := db.Querier.GetAtcoderUserByUsernameRaw(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("failed to get atcoder user by username: %w", err)
	}
	return db.parseAtcoderUser(user)
}
