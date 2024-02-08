package gorm

import (
	"github.com/Pacific73/gorm-cache/cache"
	"github.com/Pacific73/gorm-cache/config"
	cacheV8 "latipe-order-service-v2/pkg/cache/redisCacheV8"
	"log"
)

func NewCacheGormPlugin(redisClient *cacheV8.CacheEngine) *cache.Gorm2Cache {
	cachePlugin, err := cache.NewGorm2Cache(&config.CacheConfig{
		CacheLevel:           config.CacheLevelAll,
		RedisConfig:          cache.NewRedisConfigWithClient(redisClient.Client()),
		CacheStorage:         config.CacheStorageRedis,
		InvalidateWhenUpdate: true,          // when you create/update/delete objects, invalidate cache
		CacheTTL:             5 * 60 * 1000, // 5m
		CacheMaxItemCnt:      10,            // if length of objects retrieved one single time
		DebugMode:            false,
		// exceeds this number, then don't cache
	})
	if err != nil {
		log.Panic(err)
	}

	return cachePlugin
}
