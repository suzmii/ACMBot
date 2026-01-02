package handler

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/suzmii/ACMBot/consts"
	"github.com/suzmii/ACMBot/database"
)

// GetRace 从数据库获取指定平台的最新比赛记录
// 输入: resource - 平台名称（如 consts.RaceResourceCodeforces）
// 输出: []database.Race - 比赛列表, error - 错误信息
func (h *Handler) GetRace(ctx context.Context, resource string) ([]database.Race, error) {
	if resource == "" {
		return nil, fmt.Errorf("resource cannot be empty")
	}

	races, err := h.store.GetLastRace(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get races from database: %w", err)
	}

	// filter
	races = slices.DeleteFunc(races, func(race database.Race) bool {
		return race.Resource != resource
	})

	if races == nil {
		return []database.Race{}, nil
	}

	return races, nil
}

// GetAllRaces 从数据库获取所有平台的最新比赛记录
// 输入: 无
// 输出: map[string][]database.Race - 按平台分组的比赛列表, error - 错误信息
func (h *Handler) GetAllRaces(ctx context.Context) (map[string][]database.Race, error) {
	resources := []string{
		consts.RaceResourceCodeforces,
		consts.RaceResourceAtcoder,
		consts.RaceResourceLeetcode,
		consts.RaceResourceLuogu,
		consts.RaceResourceNowcoder,
	}

	result := make(map[string][]database.Race)
	for _, resource := range resources {
		races, err := h.GetRace(ctx, resource)
		if err != nil {
			return nil, fmt.Errorf("failed to get races for resource %s: %w", resource, err)
		}
		result[resource] = races
	}

	return result, nil
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
