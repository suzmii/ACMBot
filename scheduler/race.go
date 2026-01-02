package scheduler

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/suzmii/ACMBot/api"
	"github.com/suzmii/ACMBot/database"
	"github.com/suzmii/ACMBot/util/transer"
)

// RaceUpdateTask 比赛更新任务（统一更新所有平台的比赛）
type RaceUpdateTask struct {
	fetcher api.ClistRaceFetcher
	store   database.Store
}

func NewRaceUpdateTask(store database.Store, api *api.API) *RaceUpdateTask {
	return &RaceUpdateTask{
		fetcher: api.FetchClistContests,
		store:   store,
	}
}

func (t *RaceUpdateTask) Name() string {
	return "race-updater"
}

func (t *RaceUpdateTask) Execute(ctx context.Context) error {
	logger.Info("[比赛更新任务] 开始更新比赛信息")

	races, err := t.fetcher()
	if err != nil {
		return fmt.Errorf("failed to fetch races : %w", err)
	}

	dbRaces := make([]database.Race, len(races))
	for i, r := range races {
		dbRaces[i] = transer.RaceApi2DB(r)
	}

	err = t.store.CreateRace(ctx, dbRaces)
	if err != nil {
		return fmt.Errorf("failed to store races: %w", err)
	}

	logger.Infof("[比赛更新任务]: 成功更新%d场比赛", len(races))
	return nil
}

func (t *RaceUpdateTask) CheckIfNeedRunOnStart(ctx context.Context) (bool, error) {
	// 检查所有平台的数据是否超过3小时未更新
	threeHoursAgo := time.Now().Add(-3 * time.Hour)

	createdAt, err := t.store.GetLastRaceCreatedAt(ctx)
	if err != nil {
		if err == pgx.ErrNoRows {
			// 如果没有数据，需要更新
			logger.Infof("[比赛更新任务]无数据，需要更新")
			return true, nil
		}
		return false, fmt.Errorf("failed to check race update time: %w", err)
	}

	// 如果创建时间超过3小时，需要更新
	if createdAt.Time.Before(threeHoursAgo) {
		logger.Infof("[比赛更新任务]数据已过期（最后更新：%v），需要更新", createdAt.Time)
		return true, nil
	}

	logger.Infof("[比赛更新任务] 所有平台数据都在3小时内更新过，暂不更新")
	return false, nil
}

var _ = (CheckableTask)(&RaceUpdateTask{})
