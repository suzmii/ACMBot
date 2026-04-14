package handler

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/suzmii/ACMBot/consts"
	"github.com/suzmii/ACMBot/database"
)

// GetAllRaces 从数据库获取所有平台的最新比赛记录
// 输入: 无
// 输出: map[string][]database.Race - 按平台分组的比赛列表, error - 错误信息
func (h *Handler) GetAllRaces(ctx context.Context) (map[string][]database.Race, error) {
	races, err := h.store.GetLastRace(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get races from database: %w", err)
	}

	result := map[string][]database.Race{
		consts.RaceResourceCodeforces: {},
		consts.RaceResourceAtcoder:    {},
		consts.RaceResourceLeetcode:   {},
		consts.RaceResourceLuogu:      {},
		consts.RaceResourceNowcoder:   {},
	}

	for _, race := range races {
		race.Start = race.Start.In(time.Local)
		race.End = race.End.In(time.Local)

		result[race.Resource] = append(result[race.Resource], race)
	}

	return result, nil
}

// GetUpcomingRace 从数据库获取指定平台的最新比赛记录
// 输入: resource - 平台名称（如 consts.RaceResourceCodeforces）
// 输出: []database.Race - 比赛列表, error - 错误信息
func (h *Handler) GetUpcomingRace(ctx context.Context, resource string) ([]database.Race, error) {
	if resource == "" {
		return nil, fmt.Errorf("resource cannot be empty")
	}

	allRaces, err := h.GetAllRaces(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get races from database: %w", err)
	}

	races := allRaces[resource]

	// filter
	races = slices.DeleteFunc(races, func(race database.Race) bool {
		return race.Resource != resource
	})

	if races == nil {
		return []database.Race{}, nil
	}

	return races, nil
}

// GetUpcomingRaces 获取所有平台的即将开始的比赛（未来7天内）
// 输入: 无
// 输出: []database.Race - 即将开始的比赛列表, error - 错误信息
func (h *Handler) GetUpcomingRaces(ctx context.Context) ([]database.Race, error) {
	allRaces, err := h.GetAllRaces(ctx)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	oneWeekLater := now.Add(7 * 24 * time.Hour)

	var upcomingRaces []database.Race
	for _, races := range allRaces {
		for _, race := range races {
			if race.Start.After(now) && race.Start.Before(oneWeekLater) {
				upcomingRaces = append(upcomingRaces, race)
			}
		}
	}

	return upcomingRaces, nil
}
