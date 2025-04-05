package manager

import (
	"fmt"
	"github.com/suzmii/ACMBot/internal/fetcher"
	"github.com/suzmii/ACMBot/internal/model/cache"
	"github.com/suzmii/ACMBot/internal/model/provider"
	"github.com/suzmii/ACMBot/internal/model/race"
	"sort"
	"time"
)

const cacheExp = 24 * time.Hour
const updateExp = 5 * time.Hour

var AllResource = race.AllRaceResource

type CachedRace struct {
	source   race.Resource
	provider provider.RaceProvider
	err      error
}

func (r *CachedRace) Update() error {
	races, err := r.provider()
	if err != nil {
		return err
	}
	return cache.SetRace(r.source, races, cacheExp)
}

func (r *CachedRace) Get() ([]race.Race, error) {
	if r.err != nil {
		return nil, r.err
	}
	return cache.GetRace(r.source)
}

type updater struct {
	AllCachedRace map[race.Resource]*CachedRace
	UpdateTicker  *time.Ticker
}

func (r *updater) update() {
	for _, races := range r.AllCachedRace {
		err := races.Update()
		if err != nil {
			races.err = fmt.Errorf("update error: %v", err)
		}
	}
}

func (r *updater) get(resource race.Resource) ([]race.Race, error) {
	races, ok := r.AllCachedRace[resource]
	if !ok {
		return nil, fmt.Errorf("%s not found", resource)
	}
	return races.Get()
}

func (r *updater) start() {
	r.update()
	for range r.UpdateTicker.C {
		r.update()
	}
}

func newUpdater(rp map[race.Resource]provider.RaceProvider, t *time.Ticker) *updater {
	result := &updater{
		AllCachedRace: make(map[race.Resource]*CachedRace),
		UpdateTicker:  t,
	}
	for source, provider_ := range rp {
		result.AllCachedRace[source] = &CachedRace{source, provider_, nil}
	}
	return result
}

var (
	raceAndProvider = map[race.Resource]provider.RaceProvider{
		race.ResourceCodeforces: fetcher.FetchClistCodeforcesContests,
		race.ResourceAtcoder:    fetcher.FetchClistAtcoderContests,
		race.ResourceLeetcode:   fetcher.FetchClistLeetcodeContests,
		race.ResourceLuogu:      fetcher.FetchClistLuoguContests,
		race.ResourceNowcoder:   fetcher.FetchClistNowcoderContests,
	}
	defaultUpdater = newUpdater(raceAndProvider, time.NewTicker(updateExp))
)

func init() {
	go defaultUpdater.start()
}

func GetCachedRacesByResource(resource race.Resource) provider.RaceProvider {
	return func() ([]race.Race, error) {
		return defaultUpdater.get(resource)
	}
}

func GetAllCachedRaces() ([]race.Race, error) {
	var results []race.Race
	for _, s := range AllResource {
		races, err := defaultUpdater.get(s)
		if err != nil {
			continue
		}
		results = append(results, races...)
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].StartTime.Before(results[j].StartTime)
	})
	return results, nil
}
