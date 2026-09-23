package cache

import (
	"strconv"
	"sync"
	"time"

	"github.com/patrickmn/go-cache"
)

var FifteenSecondsCache DataSetCache
var ThirtySecondsCache DataSetCache
var TwoMinuteCache DataSetCache
var TenMinuteCache DataSetCache
var OneDayCache DataSetCache
var OneWeekCache DataSetCache
var OneSecondsCache DataSetCache

var cMap sync.Map

func init() {

	t1 := 1 * time.Second
	OneSecondsCache.Cache = cache.New(t1, t1)
	OneSecondsCache.Expiration = t1
	OneSecondsCache.UpdateInterval = t1

	t15 := 15 * time.Second
	FifteenSecondsCache.Cache = cache.New(t15, t15)
	FifteenSecondsCache.Expiration = t15
	FifteenSecondsCache.UpdateInterval = t15 / 2

	t30 := 30 * time.Second
	ThirtySecondsCache.Cache = cache.New(t30, t30)
	ThirtySecondsCache.Expiration = t30
	ThirtySecondsCache.UpdateInterval = t30 / 2

	t120 := 2 * time.Minute
	TwoMinuteCache.Cache = cache.New(t120, t120)
	TwoMinuteCache.Expiration = t120
	TwoMinuteCache.UpdateInterval = t120 / 2

	t600 := 10 * time.Minute
	TenMinuteCache.Cache = cache.New(t600, t600)
	TenMinuteCache.Expiration = t600
	TenMinuteCache.UpdateInterval = t600 / 2

	t1440 := 1440 * time.Minute
	OneDayCache.Cache = cache.New(t1440, t1440)
	OneDayCache.Expiration = t1440
	OneDayCache.UpdateInterval = t1440 / 2

	tWeek := 7 * 1440 * time.Minute
	OneWeekCache.Cache = cache.New(tWeek, tWeek)
	OneWeekCache.Expiration = tWeek
	OneWeekCache.UpdateInterval = tWeek
}

type DataSetCache struct {
	Cache          *cache.Cache
	UpdateInterval time.Duration
	Expiration     time.Duration
}

func (dataSet *DataSetCache) GetC(key string, callback Callback, params map[string]interface{}, noCache bool) (any, error) {
	data, eTime, found := dataSet.Cache.GetWithExpiration(key)
	if !found || noCache {
		d, err := callback(params)
		if err == nil {
			data = d
			dataSet.SetCache(key, data)
		} else {
			return nil, err
		}
	}

	if found && time.Until(eTime).Seconds() <= dataSet.UpdateInterval.Seconds() {
		//is lock or lock Expiration
		tKey := "cMapT_" + key
		if v, ok := cMap.Load(tKey); !ok || (v != nil && time.Until(v.(time.Time)) <= 0) {
			go func() {
				cMap.Store(tKey, time.Now().Add(dataSet.UpdateInterval)) //lock
				defer cMap.Delete(tKey)                                  //unlock

				td, err := callback(params)
				if err == nil {
					dataSet.SetCache(key, td)
				}
			}()

		}
	}
	return data, nil
}

func (dataSet *DataSetCache) MGetCKeyString(keys []string, callback CallbackMapString, params map[string]any, prefix string, noCache bool) (map[string]any, error) {
	unFounds := []string{}
	needs := []string{}
	founds := map[string]any{}

	// 根据key获取缓存 区分未缓存和需要重新缓存的key
	for _, key := range keys {

		data, eTime, found := dataSet.Cache.GetWithExpiration(prefix + key)

		if !found || noCache {
			unFounds = append(unFounds, key)
		} else if data != nil {
			founds[key] = data
		}

		if found && time.Until(eTime).Seconds() <= dataSet.UpdateInterval.Seconds() {
			needs = append(needs, key)
		}
	}

	// 未缓存的根据callback返回内容
	if len(unFounds) != 0 {
		list, err := callback(unFounds, params)
		if err == nil {
			for _, key := range unFounds {
				if val, ok := list[key]; ok {
					dataSet.SetCache(prefix+key, val)
					founds[key] = val
				} else {
					dataSet.SetCache(prefix+key, nil)
				}
			}
		} else {
			return founds, err
		}
	}

	// 需要重新缓存的
	if len(needs) != 0 {
		//is lock or lock Expiration
		var realNeeds []string
		for _, key := range needs {
			tKey := "t_" + prefix + "_" + key
			if v, ok := cMap.Load(tKey); !ok || (v != nil && time.Until(v.(time.Time)) <= 0) {
				cMap.Store(tKey, time.Now().Add(dataSet.UpdateInterval)) //lock
				realNeeds = append(realNeeds, key)
			}
		}

		if realNeeds != nil {
			go func() {
				for _, key := range realNeeds {
					tKey := "t_" + prefix + "_" + key
					defer cMap.Delete(tKey) //unlock
				}

				list, err := callback(needs, params)

				if err == nil {
					for _, key := range realNeeds {
						if val, ok := list[key]; ok {
							dataSet.SetCache(prefix+key, val)
						} else {
							dataSet.SetCache(prefix+key, nil)
						}
					}
				}
			}()
		}
	}

	return founds, nil
}

func (dataSet *DataSetCache) MGetCKeyInt(keys []int, callback CallbackMapInt, params map[string]any, prefix string, noCache bool) (map[int]any, error) {
	unFounds := []int{}
	needs := []int{}
	founds := map[int]any{}

	// 根据key获取缓存 区分未缓存和需要重新缓存的key
	for _, key := range keys {

		data, eTime, found := dataSet.Cache.GetWithExpiration(prefix + strconv.Itoa(key))

		if !found || noCache {
			unFounds = append(unFounds, key)
		} else if data != nil {
			founds[key] = data
		}

		if found && time.Until(eTime).Seconds() <= dataSet.UpdateInterval.Seconds() {
			needs = append(needs, key)
		}
	}

	// 未缓存的根据callback返回内容
	if len(unFounds) != 0 {
		list, err := callback(unFounds, params)
		if err == nil {
			for _, key := range unFounds {
				if val, ok := list[key]; ok {
					dataSet.SetCache(prefix+strconv.Itoa(key), val)
					founds[key] = val
				} else {
					dataSet.SetCache(prefix+strconv.Itoa(key), nil)
				}
			}
		} else {
			return founds, err
		}
	}

	// 需要重新缓存的
	if len(needs) != 0 {

		//is lock or lock Expiration
		var realNeeds []int
		for _, key := range needs {
			tKey := "t_" + prefix + "_" + strconv.Itoa(key)
			if v, ok := cMap.Load(tKey); !ok || (v != nil && time.Until(v.(time.Time)) <= 0) {
				cMap.Store(tKey, time.Now().Add(dataSet.UpdateInterval)) //lock
				realNeeds = append(realNeeds, key)
			}
		}

		if realNeeds != nil {
			go func() {
				for _, key := range realNeeds {
					tKey := "t_" + prefix + "_" + strconv.Itoa(key)
					defer cMap.Delete(tKey) //unlock
				}

				list, err := callback(needs, params)

				if err == nil {
					for _, key := range realNeeds {
						if val, ok := list[key]; ok {
							dataSet.SetCache(prefix+strconv.Itoa(key), val)
						} else {
							dataSet.SetCache(prefix+strconv.Itoa(key), nil)
						}
					}
				}
			}()
		}
	}

	return founds, nil
}

func (dataSet *DataSetCache) GetCache(key string) (any, bool) {
	return dataSet.Cache.Get(key)
}

func (dataSet *DataSetCache) SetCache(key string, val any) {
	dataSet.Cache.Set(key, val, dataSet.Expiration)
}

type Callback func(params map[string]any) (any, error)
type CallbackMapString func(keys []string, params map[string]any) (map[string]interface{}, error)
type CallbackMapInt func(keys []int, params map[string]interface{}) (map[int]interface{}, error)
