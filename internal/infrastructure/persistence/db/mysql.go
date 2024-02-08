package db

import (
	"fmt"
	"latipe-order-service-v2/config"
	cacheV8 "latipe-order-service-v2/pkg/cache/redisCacheV8"
	"latipe-order-service-v2/pkg/db/gorm"
	"log"
)

func NewMySQLConnection(configuration *config.Config, redisClient *cacheV8.CacheEngine) gorm.Gorm {
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=True&loc=Local",
		configuration.DB.Mysql.UserName,
		configuration.DB.Mysql.Password,
		configuration.DB.Mysql.Host,
		configuration.DB.Mysql.Port,
		configuration.DB.Mysql.Database,
	)
	cfg := gorm.Config{
		DSN:             dataSourceName,
		MaxOpenConns:    configuration.DB.Mysql.MaxOpenConns,
		MaxIdleConns:    configuration.DB.Mysql.MaxIdleConns,
		ConnMaxLifetime: configuration.DB.Mysql.ConnMaxLifetime,
		ConnMaxIdleTime: configuration.DB.Mysql.ConnMaxIdleTime,
		DBType:          "mysql",
		Debug:           true,
	}
	conn, err := gorm.New(cfg, redisClient)
	if err != nil {
		panic(err)
	}

	log.Printf("[%s] Gorm has created database connection", "INFO")
	return conn

}
