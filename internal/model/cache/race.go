package cache

//import (
//	"encoding/json"
//	"github.com/suzmii/ACMBot/internal/model/race"
//	"time"
//)
//
//func keyRace(source race.Resource) string {
//	return "race:" + string(source)
//}
//
//func SetRace(source race.Resource, data []race.Race, exp time.Duration) error {
//	res, err := json.Marshal(data)
//	if err != nil {
//		return err
//	}
//	return rdb.Set(ctx, keyRace(source), res, exp).Err()
//}
//
//func GetRace(source race.Resource) ([]race.Race, error) {
//	races, err := rdb.Get(ctx, keyRace(source)).Result()
//	if err != nil {
//		return nil, err
//	}
//	var result []race.Race
//	err = json.Unmarshal([]byte(races), &result)
//	if err != nil {
//		return nil, err
//	}
//	return result, nil
//}
