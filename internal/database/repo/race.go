package repo

import (
	"context"
	"errors"
	"time"

	"github.com/suzmii/ACMBot/internal/database/dbmodel"
	"github.com/suzmii/ACMBot/internal/database/gen"
	"gorm.io/gorm"
)

// GetAllRaces returns all races ordered by start time ascending.
func GetAllRaces(ctx context.Context) ([]*dbmodel.Races, error) {
	races, err := gen.Q.WithContext(ctx).
		Races.
		Order(gen.Races.StartAt).
		Find()
	if err != nil {
		return nil, err
	}
	return races, nil
}

// GetLatestRaceUpdatedAt returns the most recent UpdatedAt across all races.
func GetLatestRaceUpdatedAt(ctx context.Context) (time.Time, error) {
	race, err := gen.Q.WithContext(ctx).
		Races.
		Order(gen.Races.UpdatedAt.Desc()).
		Take()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return time.Time{}, nil
		}
		return time.Time{}, err
	}
	return race.UpdatedAt, nil
}

// DeleteFinishedRaces soft-deletes races that have already ended.
func DeleteFinishedRaces(ctx context.Context, now time.Time) error {
	races, err := gen.Q.WithContext(ctx).
		Races.
		Where(gen.Races.EndAt.Lt(now)).
		Find()
	if err != nil {
		return err
	}
	if len(races) == 0 {
		return nil
	}
	_, err = gen.Q.WithContext(ctx).Races.Delete(races...)
	return err
}

// ReplaceRaces replaces all races with the provided list (soft-deleting existing entries then inserting new ones).
func ReplaceRaces(ctx context.Context, newRaces []*dbmodel.Races) error {
	return gen.Q.Transaction(func(tx *gen.Query) error {
		// delete all existing races first
		existing, err := tx.Races.WithContext(ctx).Find()
		if err != nil {
			return err
		}
		if len(existing) > 0 {
			if _, err := tx.Races.WithContext(ctx).Delete(existing...); err != nil {
				return err
			}
		}
		if len(newRaces) == 0 {
			return nil
		}
		return tx.Races.WithContext(ctx).CreateInBatches(newRaces, 500)
	})
}
