package repo

import (
	"context"
	"errors"

	"github.com/suzmii/ACMBot/internal/database/dbmodel"
	"github.com/suzmii/ACMBot/internal/database/gen"
	"gorm.io/gorm"
)

func UpdateRatingRecords(ctx context.Context, uid uint, detail *dbmodel.CodeforcesRatingRecords) error {
	err := gen.Q.WithContext(ctx).CodeforcesRatingRecords.
		Where(gen.CodeforcesRatingRecords.UserID.Eq(uid)).
		Save(detail)
	if err != nil {
		return err
	}
	return nil
}

func GetRatingRecords(ctx context.Context, uid uint) (*dbmodel.CodeforcesRatingRecords, error) {
	take, err := gen.Q.WithContext(ctx).CodeforcesRatingRecords.
		Where(gen.CodeforcesRatingRecords.UserID.Eq(uid)).
		Take()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return take, nil
}
