package gorm

import (
	"context"
	"database/sql"
	"errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"latipe-order-service-v2/config"
	cacheV8 "latipe-order-service-v2/pkg/cache/redisCacheV8"
	"strings"
	"time"
)

// Gorm defines a interface for access the database.
type Gorm interface {
	DB() *gorm.DB
	SqlDB() *sql.DB
	Exec(fc func(tx *gorm.DB) error, ctx context.Context) (err error)
	Transaction(fc func(tx *gorm.DB) error) (err error)
	Close() error
	DropTableIfExists(value interface{}) error
}

// Config GORM Config
type Config struct {
	Debug           bool
	DBType          string
	DSN             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
	TablePrefix     string
}

// _gorm gorm struct
type _gorm struct {
	db    *gorm.DB
	sqlDB *sql.DB
	cfg   *config.Config
}

// New Create gorm.DB and  instance
func New(c Config, redisClient *cacheV8.CacheEngine) (Gorm, error) {
	var dial gorm.Dialector

	switch strings.ToLower(c.DBType) {
	case "mysql":
		dial = mysql.Open(c.DSN)
	default:
		return nil, errors.New("DBType does not support")
	}

	gConfig := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   c.TablePrefix,
			SingularTable: true,
		},
	}

	db, err := gorm.Open(dial, gConfig)
	if err != nil {
		return nil, err
	}

	if c.Debug {
		db = db.Debug()
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	if c.MaxOpenConns != 0 {
		sqlDB.SetMaxOpenConns(c.MaxOpenConns)
	}
	if c.ConnMaxLifetime != 0 {
		sqlDB.SetConnMaxLifetime(c.ConnMaxLifetime)
	}
	if c.MaxIdleConns != 0 {
		sqlDB.SetMaxIdleConns(c.MaxIdleConns)
	}
	if c.ConnMaxIdleTime != 0 {
		sqlDB.SetConnMaxIdleTime(c.ConnMaxIdleTime)
	}

	//_ = db.Use(NewCacheGormPlugin(redisClient))
	return &_gorm{
		db:    db,
		sqlDB: sqlDB,
	}, nil
}
