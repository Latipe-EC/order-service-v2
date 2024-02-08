package cache

import (
	"fmt"
	"github.com/google/wire"
	"latipe-order-service-v2/config"
	cacheV8 "latipe-order-service-v2/pkg/cache/redisCacheV8"
	cacheV9 "latipe-order-service-v2/pkg/cache/redisCacheV9"
	"time"
)

var Set = wire.NewSet(
	NewCacheEngineV9,
	NewCacheEngineV8,
)

func NewCacheEngineV9(config *config.Config) (*cacheV9.CacheEngine, error) {
	cfg := cacheV9.RedisConfig{
		Address:               fmt.Sprintf("%v:%v", config.Cache.Redis.Address, config.Cache.Redis.Port),
		DB:                    config.Cache.Redis.DbAuth,
		Password:              config.Cache.Redis.Password,
		ContextTimeoutEnabled: true,
		PoolSize:              5,
		PoolTimeout:           5,
		DialTimeout:           5,
		ReadTimeout:           5 * time.Second,
		WriteTimeout:          5 * time.Second,
		ConnectTimeout:        5 * time.Second,
	}
	client, err := cacheV9.NewCacheEngine(cfg)
	if err != nil {
		return nil, err
	}
	return client, err
}

func NewCacheEngineV8(config *config.Config) (*cacheV8.CacheEngine, error) {
	cfg := cacheV8.RedisConfig{
		Address:               fmt.Sprintf("%v:%v", config.Cache.Redis.Address, config.Cache.Redis.Port),
		DB:                    config.Cache.Redis.DbQuery,
		Password:              config.Cache.Redis.Password,
		ContextTimeoutEnabled: true,
		PoolSize:              5,
		PoolTimeout:           5,
		DialTimeout:           5,
		ReadTimeout:           5 * time.Second,
		WriteTimeout:          5 * time.Second,
		ConnectTimeout:        5 * time.Second,
	}
	client, err := cacheV8.NewCacheEngineV8(cfg)
	if err != nil {
		return nil, err
	}
	return client, err
}
