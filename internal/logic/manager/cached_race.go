package manager

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/suzmii/ACMBot/internal/fetcher"
	"github.com/suzmii/ACMBot/internal/model/provider"
	"github.com/suzmii/ACMBot/internal/model/race"
	"time"
)

const updateInterval = 1 * time.Minute // 每分钟检查一次
const maxCacheAge = 4 * time.Hour      // 最大缓存时间，超过4小时需要更新

var raceAndProvider = map[race.Resource]provider.RaceProvider{
	race.ResourceCodeforces: fetcher.FetchClistCodeforcesContests,
	race.ResourceAtcoder:    fetcher.FetchClistAtcoderContests,
	race.ResourceLeetcode:   fetcher.FetchClistLeetcodeContests,
	race.ResourceLuogu:      fetcher.FetchClistLuoguContests,
	race.ResourceNowcoder:   fetcher.FetchClistNowcoderContests,
}

// raceData 存储比赛数据及更新时间
type raceData struct {
	source     race.Resource
	provider   provider.RaceProvider
	cache      []race.Race
	lastUpdate time.Time
}

// Update 更新比赛数据
func (r *raceData) Update() error {
	// 获取最新的数据
	races, err := r.provider()
	if err != nil {
		return fmt.Errorf("failed to fetch races: %v", err)
	}

	// 更新内存中的缓存数据和更新时间
	r.cache = races
	r.lastUpdate = time.Now()

	return nil
}

// updater 用于定期更新比赛数据
type updater struct {
	AllCachedRace map[race.Resource]*raceData
	UpdateTicker  *time.Ticker
}

// startUpdating 每分钟检查并更新比赛数据
func (u *updater) startUpdating() {
	u.updateAll()
	for range u.UpdateTicker.C {
		u.updateAll()
	}
}

// updateAll 检查所有比赛数据的更新时间并更新
func (u *updater) updateAll() {
	for _, raceData := range u.AllCachedRace {
		if time.Since(raceData.lastUpdate) > maxCacheAge {
			logrus.Infof("Updating race %v", raceData.source)
			err := raceData.Update()
			if err != nil {
				logrus.Errorf("failed to update race: %v", err)
				continue
			}
			logrus.Infof("Updated race %v: data: %v", raceData.source, raceData.cache)
		}
	}
}

// newUpdater 创建一个新的 updater 实例
func newUpdater(rp map[race.Resource]provider.RaceProvider) *updater {
	result := &updater{
		AllCachedRace: make(map[race.Resource]*raceData),
		UpdateTicker:  time.NewTicker(updateInterval),
	}
	for source, provider_ := range rp {
		result.AllCachedRace[source] = &raceData{
			source:     source,
			provider:   provider_,
			lastUpdate: time.Time{},
		}
	}
	return result
}

// 初始化 updater
var defaultUpdater = newUpdater(raceAndProvider)

// init 在启动时启动更新协程
func init() {
	go defaultUpdater.startUpdating()
}

// GetRaceProviderByResource 获取指定资源的比赛数据
func GetRaceProviderByResource(resource race.Resource) provider.RaceProvider {
	return func() ([]race.Race, error) {
		raceData, ok := defaultUpdater.AllCachedRace[resource]
		if !ok {
			return nil, fmt.Errorf("resource %s not found", resource)
		}
		return raceData.cache, nil
	}
}

var allCachedRaceProvider = func() ([]race.Race, error) {
	var results []race.Race
	for _, raceData := range defaultUpdater.AllCachedRace {
		results = append(results, raceData.cache...)
	}
	return results, nil
}

// GetAllCachedRaceProvider 获取所有比赛数据
func GetAllCachedRaceProvider() provider.RaceProvider {
	return allCachedRaceProvider
}
