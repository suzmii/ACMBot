package repo

import (
	"context"
	"errors"
	"github.com/suzmii/ACMBot/internal/database/dbmodel"
	"github.com/suzmii/ACMBot/internal/database/gen"
	"github.com/suzmii/ACMBot/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"sync"
	"time"
)

//func (r *Repo) CountSolvedProblem(ctx context.Context, uid uint) (int64, error) {
//	count, err := gen.Q.WithContext(ctx).CodeforcesUserPassedProblem.
//		Where(gen.CodeforcesUserPassedProblem.UserID.Eq(uid)).
//		Count()
//	if err != nil {
//		return 0, err
//	}
//	return count, nil
//}

func CountSolvedDetail(ctx context.Context, uid uint) (*model.SolvedDetail, error) {
	var solvedDetail model.SolvedDetail
	errCh := make(chan error, 4)
	wg := sync.WaitGroup{}

	queryRange := func(low, high int, dest *int) {
		defer wg.Done()
		c, err := gen.Q.WithContext(ctx).CodeforcesUserPassedProblem.
			Where(
				gen.CodeforcesUserPassedProblem.UserID.Eq(uid),
				gen.CodeforcesUserPassedProblem.ProblemRating.Between(low, high),
			).Count()
		if err != nil {
			select { // 避免阻塞
			case errCh <- err:
			default:
			}
			return
		}
		*dest = int(c)
	}

	wg.Add(4)
	go queryRange(800, 1400, &solvedDetail.Range800to1400)
	go queryRange(1400, 2000, &solvedDetail.Range1400to2000)
	go queryRange(2000, 2600, &solvedDetail.Range2000to2600)
	go queryRange(2600, 9999, &solvedDetail.Above2600)

	wg.Wait()
	close(errCh)

	if err, ok := <-errCh; ok {
		return nil, err
	}
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	return &solvedDetail, nil
}

func UpdateSubmissions(ctx context.Context, uid uint, submissions []*dbmodel.CodeforcesSubmission) error {

	if len(submissions) == 0 {
		_, err := gen.Q.WithContext(ctx).
			CodeforcesUser.
			Where(gen.CodeforcesUser.ID.Eq(uid)).
			UpdateColumnSimple(gen.CodeforcesUser.SubmissionUpdatedAt.Value(time.Now()))
		return err
	}

	passedSubmission := make([]*dbmodel.CodeforcesUserPassedProblem, 0, len(submissions))

	for _, submission := range submissions {
		if submission.Status == dbmodel.CodeforcesSubmissionStatusOk {
			passedSubmission = append(passedSubmission, &dbmodel.CodeforcesUserPassedProblem{
				Model: gorm.Model{
					CreatedAt: submission.CreatedAt,
				},
				UserID:        submission.UserID,
				ProblemID:     submission.ProblemID,
				ProblemRating: submission.ProblemRating,
			})
		}
	}

	err := gen.Q.Transaction(func(tx *gen.Query) error {
		_, err := tx.CodeforcesUser.
			WithContext(ctx).
			Where(tx.CodeforcesUser.ID.Eq(submissions[0].UserID)).
			UpdateSimple(tx.CodeforcesUser.SubmissionUpdatedAt.Value(time.Now()))
		if err != nil {
			return err
		}

		err = tx.CodeforcesSubmission.
			WithContext(ctx).
			Clauses(clause.OnConflict{DoNothing: true}).
			CreateInBatches(submissions, 500)
		if err != nil {
			return err
		}

		err = tx.CodeforcesUserPassedProblem.
			WithContext(ctx).
			Clauses(clause.OnConflict{DoNothing: true}).
			CreateInBatches(passedSubmission, 500)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func GetLastSubmission(ctx context.Context, uid uint) (*dbmodel.CodeforcesSubmission, error) {
	submission, err := gen.Q.WithContext(ctx).CodeforcesSubmission.
		Where(gen.CodeforcesSubmission.UserID.Eq(uid)).
		Order(gen.CodeforcesSubmission.CreatedAt.Desc()).
		First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return submission, nil
}
