package datasync

import (
	"context"
	"sort"
	"time"

	"github.com/suzmii/ACMBot/config"
	"github.com/suzmii/ACMBot/internal/api/client"
	"github.com/suzmii/ACMBot/internal/database/dbmodel"
	"github.com/suzmii/ACMBot/internal/database/repo"
	"github.com/suzmii/ACMBot/internal/model"
)

func GetRaces(ctx context.Context) ([]model.Race, error) {
	// Determine last synced time by the latest race update time in DB
	lastSynced, err := repo.GetLatestRaceUpdatedAt(ctx)
	if err != nil {
		return nil, err
	}

	// If never synced or older than configured hours, refresh from remote and update DB
	refreshHours := config.LoadConfig().Sync.RaceRefreshHours
	if refreshHours <= 0 {
		refreshHours = 24
	}
	if lastSynced.IsZero() || time.Since(lastSynced) > time.Duration(refreshHours)*time.Hour {
		// Fetch from clist for all sources
		var fetched []model.Race
		// Codeforces
		if list, err := client.FetchClistCodeforcesContests(); err == nil {
			fetched = append(fetched, list...)
		} else {
			return nil, err
		}
		// Atcoder
		if list, err := client.FetchClistAtcoderContests(); err == nil {
			fetched = append(fetched, list...)
		} else {
			return nil, err
		}
		// Leetcode
		if list, err := client.FetchClistLeetcodeContests(); err == nil {
			fetched = append(fetched, list...)
		} else {
			return nil, err
		}
		// Luogu
		if list, err := client.FetchClistLuoguContests(); err == nil {
			fetched = append(fetched, list...)
		} else {
			return nil, err
		}
		// Nowcoder
		if list, err := client.FetchClistNowcoderContests(); err == nil {
			fetched = append(fetched, list...)
		} else {
			return nil, err
		}

		// Delete finished races after keep hours
		keepHours := config.LoadConfig().Sync.RaceKeepHours
		if keepHours <= 0 {
			keepHours = 24
		}
		if err := repo.DeleteFinishedRaces(ctx, time.Now(), keepHours); err != nil {
			return nil, err
		}

		// Convert to db models and replace
		dbRaces := make([]*dbmodel.Races, 0, len(fetched))
		for _, r := range fetched {
			dbRaces = append(dbRaces, &dbmodel.Races{
				Resource: convertToDBResource(r.Source),
				Title:    r.Name,
				StartAt:  r.StartTime,
				EndAt:    r.EndTime,
				Link:     r.Link,
			})
		}
		if err := repo.ReplaceRaces(ctx, dbRaces); err != nil {
			return nil, err
		}
	}

	// Read from DB and return
	stored, err := repo.GetAllRaces(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]model.Race, 0, len(stored))
	for _, r := range stored {
		result = append(result, model.Race{
			Source:    convertToModelResource(r.Resource),
			Name:      r.Title,
			Link:      r.Link,
			StartTime: r.StartAt,
			EndTime:   r.EndAt,
		})
	}

	sort.Slice(result, func(i, j int) bool { return result[i].StartTime.Before(result[j].StartTime) })
	return result, nil
}

func convertToDBResource(res model.Resource) dbmodel.Resource {
	switch res {
	case model.ResourceCodeforces:
		return dbmodel.ResourceCodeforces
	case model.ResourceAtcoder:
		return dbmodel.ResourceAtcoder
	case model.ResourceLeetcode:
		return dbmodel.ResourceLeetcode
	case model.ResourceLuogu:
		return dbmodel.ResourceLuogu
	case model.ResourceNowcoder:
		return dbmodel.ResourceNowcoder
	default:
		return dbmodel.ResourceCodeforces
	}
}

func convertToModelResource(res dbmodel.Resource) model.Resource {
	switch res {
	case dbmodel.ResourceCodeforces:
		return model.ResourceCodeforces
	case dbmodel.ResourceAtcoder:
		return model.ResourceAtcoder
	case dbmodel.ResourceLeetcode:
		return model.ResourceLeetcode
	case dbmodel.ResourceLuogu:
		return model.ResourceLuogu
	case dbmodel.ResourceNowcoder:
		return model.ResourceNowcoder
	default:
		return model.ResourceCodeforces
	}
}
