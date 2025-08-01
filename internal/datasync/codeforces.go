package datasync

import (
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/suzmii/ACMBot/internal/api/client"
	"github.com/suzmii/ACMBot/internal/database/dbmodel"
	"github.com/suzmii/ACMBot/internal/database/repo"
	"github.com/suzmii/ACMBot/internal/errs"
	"github.com/suzmii/ACMBot/internal/model"
	"gorm.io/gorm"
	"time"
)

// GetOrSyncUser 获取用户，如果是新用户则同步并创建用户
func GetOrSyncUser(ctx context.Context, username string) (model.User, error) {
	repoUser, err := repo.GetUser(ctx, username)
	if err == nil {
		logrus.Debugf("got user: %+v", repoUser)
		return model.UserFromRepo(*repoUser), nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return model.User{}, err
	}
	user, err := User(ctx, username)
	if err != nil {
		return model.User{}, err
	}
	return user, nil
}

func User(ctx context.Context, username string) (model.User, error) {
	logrus.Debugf("syncing user %s", username)
	info, err := client.FetchCodeforcesUserInfo(username, false)
	if err != nil {
		return model.User{}, err
	}

	repoUser := dbmodel.CodeforcesUser{
		Username:      info.Handle,
		AvatarURL:     info.Avatar,
		Organization:  info.Organization,
		FriendOf:      info.FriendCount,
		MaxRating:     info.MaxRating,
		CurrentRating: info.Rating,
	}

	err = repo.SaveUser(ctx, &repoUser)
	if err != nil {
		return model.User{}, err
	}
	logrus.Debugf("synced and saved user: %+v", repoUser)
	return model.UserFromRepo(repoUser), nil
}

func RatingRecords(ctx context.Context, user model.User) (model.RatingRecords, error) {
	logrus.Debug("syncing rating records", user)
	apiRecords, err := client.FetchCodeforcesRatingRecords(user.Name)
	if err != nil {
		return nil, fmt.Errorf("fetch datasync rating records: %w", err)
	}

	if len(apiRecords) == 0 {
		return nil, errs.ErrNoRatingRecords
	}

	repoRecords := dbmodel.CodeforcesRatingRecords{
		UserID: user.ID,
	}
	for _, record := range apiRecords {
		repoRecords.Records = append(repoRecords.Records, record.ToCommon().ToRepo())
	}
	err = repo.UpdateRatingRecords(ctx, user.ID, &repoRecords)
	if err != nil {
		return nil, err
	}

	user.CurrentRating = apiRecords[len(apiRecords)-1].NewRating
	user.MaxRating = apiRecords[0].NewRating
	for _, record := range apiRecords {
		user.MaxRating = max(user.MaxRating, record.NewRating)
	}

	repoUser := user.ToRepo()

	err = repo.SaveUser(ctx, &repoUser)
	if err != nil {
		return nil, err
	}

	logrus.Debugf("synced and saved rating records: %+v", repoRecords)
	return model.RatingRecordsFromRepo(repoRecords), nil
}

// Submissions 因为submission数据需要统计才能用，所以就不返回了; after一定要是当前用户在数据库中最后一个提交记录的时间
func Submissions(ctx context.Context, user model.User) error {
	logrus.Debug("syncing submissions", user)
	submission, err := repo.GetLastSubmission(ctx, user.ID)
	if err != nil {
		return err
	}

	after := time.Time{}

	if submission != nil {
		logrus.Debug("got submission: %v", submission)
		after = submission.CreatedAt
	}

	submissions, err := client.FetchCodeforcesSubmissionsAfter(user.Name, after)
	if err != nil {
		return err
	}
	// 就算submission长度为0也更新
	var repoSubmissions []*dbmodel.CodeforcesSubmission

	for _, submission := range submissions {
		repoSubmissions = append(repoSubmissions, &dbmodel.CodeforcesSubmission{
			Model: gorm.Model{
				CreatedAt: time.Unix(submission.At, 0).Local(),
			},
			UserID:        user.ID,
			Pass:          submission.Status == string(dbmodel.CodeforcesSubmissionStatusOk),
			ProblemID:     submission.Problem.ID(),
			ProblemRating: submission.Problem.Rating,
			Language:      submission.Language,
			Status:        dbmodel.CodeforcesSubmissionStatus(submission.Status),
		})
	}

	err = repo.UpdateSubmissions(ctx, user.ID, repoSubmissions)
	if err != nil {
		return err
	}
	logrus.Debugf("synced and saved submissions: %+v", repoSubmissions)
	return nil
}
