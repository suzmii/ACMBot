package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/suzmii/ACMBot/util/jsonx"
)

type Race struct {
	Title    string    `json:"title"`
	Resource string    `json:"resource"`
	Start    time.Time `json:"start"`
	End      time.Time `json:"end"`
	Link     string    `json:"link"`
}

func (s *dbStore) CreateRace(ctx context.Context, races []Race) error {
	jsonData, err := jsonx.Marshal(races)
	if err != nil {
		return fmt.Errorf("failed to prase races: %w", err)
	}
	return s.Querier.CreateRaceRaw(ctx, jsonData)
}

func (s *dbStore) GetLastRace(ctx context.Context) ([]Race, error) {
	raw, err := s.Querier.GetLastRaceRaw(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	var result []Race

	err = jsonx.Unmarshal(raw, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal races: %w", err)
	}

	return result, nil
}
