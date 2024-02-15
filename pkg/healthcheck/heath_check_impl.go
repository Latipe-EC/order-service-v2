package healthcheck

import (
	"fmt"
	"github.com/hellofresh/health-go/v5"
	healthMysql "github.com/hellofresh/health-go/v5/checks/mysql"
	healthRabbit "github.com/hellofresh/health-go/v5/checks/rabbitmq"
	healthRedis "github.com/hellofresh/health-go/v5/checks/redis"
	"latipe-order-service-v2/config"
	"time"
)

func NewHealthCheckService(config *config.Config) (*health.Health, error) {
	// add some checks on instance creation
	h, err := health.New(health.WithComponent(health.Component{
		Name:    "order-service-v2",
		Version: "v2.0",
	}))
	if err != nil {
		return nil, err
	}

	//mysql check
	err = h.Register(health.Config{
		Name:      "mysql",
		Timeout:   time.Second * 2,
		SkipOnErr: false,
		Check: healthMysql.New(healthMysql.Config{
			DSN: fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4",
				config.DB.Mysql.UserName,
				config.DB.Mysql.Password,
				config.DB.Mysql.Host,
				config.DB.Mysql.Port,
				config.DB.Mysql.Database,
			),
		}),
	})
	if err != nil {
		return nil, err
	}

	//rabbitMQ check
	err = h.Register(health.Config{
		Name:      "rabbitMQ",
		Timeout:   time.Second * 2,
		SkipOnErr: false,
		Check: healthRabbit.New(healthRabbit.Config{
			DSN: config.RabbitMQ.Connection,
		}),
	})

	//redis check
	err = h.Register(health.Config{
		Name:      "redis-cache-query",
		Timeout:   time.Second * 2,
		SkipOnErr: false,
		Check: healthRedis.New(healthRedis.Config{
			DSN: fmt.Sprintf("redis://%s:%d/%d", config.Cache.Redis.Address, config.Cache.Redis.Port, config.Cache.Redis.DbQuery),
		}),
	})
	if err != nil {
		return nil, err
	}

	//redis check
	err = h.Register(health.Config{
		Name:      "redis-cache-middleware",
		Timeout:   time.Second * 2,
		SkipOnErr: false,
		Check: healthRedis.New(healthRedis.Config{
			DSN: fmt.Sprintf("redis://%s:%d/%d", config.Cache.Redis.Address, config.Cache.Redis.Port, config.Cache.Redis.DbAuth),
		}),
	})
	if err != nil {
		return nil, err
	}

	return h, nil
}
