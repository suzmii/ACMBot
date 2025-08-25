package repo

import (
	"context"
	"errors"
	"time"

	"github.com/suzmii/ACMBot/internal/database/dbmodel"
	"github.com/suzmii/ACMBot/internal/database/gen"
	"gorm.io/gorm"
)

// GetRaces returns all races ordered by start time ascending.
func GetRaces(ctx context.Context) ([]*dbmodel.Races, error) {
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

// DeleteFinishedRaces soft-deletes races that have already ended and exceeded keep hours.
func DeleteFinishedRaces(ctx context.Context, now time.Time, keepHours int) error {
	races, err := gen.Q.WithContext(ctx).
		Races.
		Where(gen.Races.EndAt.Lt(now.Add(-time.Duration(keepHours) * time.Hour))).
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
	races, err := GetRaces(ctx)
	if err != nil {
		return err
	}
	ids := make(map[string]bool)
	newIds := make(map[string]bool)
	for _, race := range races {
		ids[race.ID] = true
	}
	for _, race := range newRaces {
		newIds[race.ID] = true
	}
	for _, race := range newRaces {
		if ids[race.ID] == false {
			gen.Q.WithContext(ctx).Races.Create(race)
		}
	}
	for _, race := range races {
		if newIds[race.ID] == false {
			gen.Q.WithContext(ctx).Races.Delete(race)
		}
	}
	return nil
}
