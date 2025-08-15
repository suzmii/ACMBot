package datasync

import (
	"context"
	"time"

	"github.com/suzmii/ACMBot/config"
	"github.com/suzmii/ACMBot/internal/api/client"
	"github.com/suzmii/ACMBot/internal/database/dbmodel"
	"github.com/suzmii/ACMBot/internal/database/repo"
	"github.com/suzmii/ACMBot/internal/model"
)

// SyncRacesFromRemoteIfStale refreshes races from remote and updates DB if needed
func SyncRacesFromRemoteIfStale(ctx context.Context) error {
	// Determine last synced time by the latest race update time in DB
	lastSynced, err := repo.GetLatestRaceUpdatedAt(ctx)
	if err != nil {
		return err
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
			return err
		}
		// Atcoder
		if list, err := client.FetchClistAtcoderContests(); err == nil {
			fetched = append(fetched, list...)
		} else {
			return err
		}
		// Leetcode
		if list, err := client.FetchClistLeetcodeContests(); err == nil {
			fetched = append(fetched, list...)
		} else {
			return err
		}
		// Luogu
		if list, err := client.FetchClistLuoguContests(); err == nil {
			fetched = append(fetched, list...)
		} else {
			return err
		}
		// Nowcoder
		if list, err := client.FetchClistNowcoderContests(); err == nil {
			fetched = append(fetched, list...)
		} else {
			return err
		}

		// Delete finished races after keep hours
		keepHours := config.LoadConfig().Sync.RaceKeepHours
		if keepHours <= 0 {
			keepHours = 24
		}
		if err := repo.DeleteFinishedRaces(ctx, time.Now(), keepHours); err != nil {
			return err
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
			return err
		}
	}

	return nil
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
