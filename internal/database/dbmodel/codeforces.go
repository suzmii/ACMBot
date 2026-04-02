package dbmodel

import (
	"time"

	"gorm.io/gorm"
)

type CodeforcesSubmissionStatus string

const (
	CodeforcesSubmissionStatusOk CodeforcesSubmissionStatus = "OK"
	// CodeforcesSubmissionStatusFailed                  CodeforcesSubmissionStatus = "FAILED"
	// CodeforcesSubmissionStatusPartial                 CodeforcesSubmissionStatus = "PARTIAL"
	// CodeforcesSubmissionStatusCompilationError        CodeforcesSubmissionStatus = "COMPILATION_ERROR"
	// CodeforcesSubmissionStatusRuntimeError            CodeforcesSubmissionStatus = "RUNTIME_ERROR"
	// CodeforcesSubmissionStatusWrongAnswer             CodeforcesSubmissionStatus = "WRONG_ANSWER"
	// CodeforcesSubmissionStatusPresentationError       CodeforcesSubmissionStatus = "PRESENTATION_ERROR"
	// CodeforcesSubmissionStatusTimeLimitExceeded       CodeforcesSubmissionStatus = "TIME_LIMIT_EXCEEDED"
	// CodeforcesSubmissionStatusMemoryLimitExceeded     CodeforcesSubmissionStatus = "MEMORY_LIMIT_EXCEEDED"
	// CodeforcesSubmissionStatusIdlenessLimitExceeded   CodeforcesSubmissionStatus = "IDLENESS_LIMIT_EXCEEDED"
	// CodeforcesSubmissionStatusSecurityViolated        CodeforcesSubmissionStatus = "SECURITY_VIOLATED"
	// CodeforcesSubmissionStatusCrashed                 CodeforcesSubmissionStatus = "CRASHED"
	// CodeforcesSubmissionStatusInputPreparationCrashed CodeforcesSubmissionStatus = "INPUT_PREPARATION_CRASHED"
	// CodeforcesSubmissionStatusChallenged              CodeforcesSubmissionStatus = "CHALLENGED"
	// CodeforcesSubmissionStatusSkipped                 CodeforcesSubmissionStatus = "SKIPPED"
	// CodeforcesSubmissionStatusTesting                 CodeforcesSubmissionStatus = "TESTING"
	// CodeforcesSubmissionStatusRejected                CodeforcesSubmissionStatus = "REJECTED"
)

type CodeforcesUser struct {
	gorm.Model

	Username     string `gorm:"uniqueIndex"`
	AvatarURL    string
	Organization string
	FriendOf     int

	MaxRating     int
	CurrentRating int

	SubmissionUpdatedAt time.Time
}

type CodeforcesSubmission struct {
	gorm.Model
	// createdAt 替换为实际提交时间
	UserID        uint   `gorm:"index:idx_user_pass_problem,priority:1"`
	Pass          bool   `gorm:"index:idx_user_pass_problem,priority:2"`
	ProblemID     string `gorm:"index:idx_user_pass_problem,priority:3"`
	ProblemRating int

	Language string
	Status   CodeforcesSubmissionStatus
}

type CodeforcesUserPassedProblem struct {
	gorm.Model
	UserID        uint   `gorm:"uniqueIndex:idx_user_problem,priority:1;index:idx_user_problem_rating,priority:1"`
	ProblemID     string `gorm:"uniqueIndex:idx_user_problem,priority:2"`
	ProblemRating int    `gorm:"index:idx_user_problem_rating,priority:2"`
}

type CodeforcesRatingRecords struct {
	UserID    uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt           `gorm:"index"`
	Records   []CodeforcesRatingRecord `gorm:"type:json;serializer:json"`
}

type CodeforcesRatingRecord struct {
	At     time.Time
	Rating int
}

var CodeforcesModels = []interface{}{
	new(CodeforcesUser),
	new(CodeforcesSubmission),
	new(CodeforcesUserPassedProblem),
	new(CodeforcesRatingRecords),
}
