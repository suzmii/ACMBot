package transer

import (
	"time"

	"github.com/suzmii/ACMBot/api"
	"github.com/suzmii/ACMBot/database"
)

func AtcoderSubmissionApi2DB(userID int64, raw []api.AtcoderSubmission) []database.AtcoderSubmission {
	result := make([]database.AtcoderSubmission, 0, len(raw))

	for _, record := range raw {
		result = append(result, database.AtcoderSubmission{
			UserID:       userID,
			SubmissionID: int64(record.SubmissionId),
			ProblemID:    record.ProblemId,
			Point:        float64(record.Point),
			Status:       record.Status,
			At:           time.Unix(record.SubmissionTime, 0),
		})
	}

	return result
}
