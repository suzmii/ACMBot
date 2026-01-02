package transer

import (
	"time"

	"github.com/suzmii/ACMBot/api"
	"github.com/suzmii/ACMBot/database"
	"github.com/suzmii/ACMBot/render"
)

func CodeforcesRatingRecordApi2DB(raw []api.RatingRecord) []database.CodeforcesRatingRecord {
	result := make([]database.CodeforcesRatingRecord, 0, len(raw))

	for _, record := range raw {
		result = append(result, database.CodeforcesRatingRecord{
			Rating: record.NewRating,
			At:     time.Unix(record.At, 0),
		})
	}

	return result
}

func CodeforcesRatingRecordDB2Render(raw []database.CodeforcesRatingRecord) []render.CodeforcesRatingChange {
	result := make([]render.CodeforcesRatingChange, 0, len(raw))

	for _, record := range raw {
		result = append(result, render.CodeforcesRatingChange{
			At:        record.At.Unix(),
			NewRating: record.Rating,
		})
	}

	return result
}

func CodeforcesSubmissionApi2DB(userId int64, raw []api.Submission) []database.CodeforcesSubmission {
	result := make([]database.CodeforcesSubmission, 0, len(raw))

	for _, record := range raw {
		result = append(result, database.CodeforcesSubmission{
			UserID: userId,
			Problem: database.CodeforcesProblem{
				ContestID:      record.ContestID,
				ProblemSetName: record.Problem.ProblemSetName,
				Index:          record.Problem.Index,
				Name:           record.Problem.Name,
				Type:           record.Problem.Type,
				Rating:         record.Problem.Rating,
				Tags:           record.Problem.Tags,
			},
			Status: record.Status,
			At:     time.Unix(record.At, 0),
		})
	}

	return result
}
