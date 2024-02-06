package gorm

import (
	"github.com/Pacific73/gorm-cache/cache"
	"github.com/Pacific73/gorm-cache/config"
	"log"
)

func NewCacheGormPlugin() *cache.Gorm2Cache {
	cachePlugin, err := cache.NewGorm2Cache(&config.CacheConfig{
		CacheLevel:           config.CacheLevelAll,
		CacheStorage:         config.CacheStorageMemory,
		InvalidateWhenUpdate: true, // when you create/update/delete objects, invalidate cache
		CacheTTL:             5000, // 5000 ms
		CacheMaxItemCnt:      5,    // if length of objects retrieved one single time
		// exceeds this number, then don't cache
	})
	if err != nil {
		log.Panic(err)
	}

	return cachePlugin
}
