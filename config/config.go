package config

import (
	"errors"
	"fmt"
	"github.com/google/wire"
	"github.com/spf13/viper"
	"log"
	"os"
	"time"
)

var Set = wire.NewSet(NewConfig)

type Config struct {
	Server Server
	DB     DB
	Cache  Cache
	GRPC   GRPC
	//Adapters Adapters
	AdapterService AdapterService
	RabbitMQ       RabbitMQ
	Metrics        Metrics
}

type Server struct {
	KeyID               string
	Name                string
	ApiHeaderKey        string
	AppVersion          string
	Port                string
	BaseURI             string
	Mode                string
	ReadTimeout         time.Duration
	WriteTimeout        time.Duration
	SSL                 bool
	CtxDefaultTimeout   int
	CSRF                bool
	Debug               bool
	MaxCountRequest     int           // max count of connections
	ExpirationLimitTime time.Duration //  expiration time of the limit
}

type GRPC struct {
	DeliveryServiceGrpc DeliveryServiceGrpc
	VoucherServiceGrpc  VoucherServiceGrpc
	ProductServiceGrpc  ProductServiceGrpc
	UserServiceGrpc     UserServiceGrpc
}

type DeliveryServiceGrpc struct {
	Connection string
	APIKey     string
}
type VoucherServiceGrpc struct {
	Connection string
	APIKey     string
}
type ProductServiceGrpc struct {
	Connection string
	APIKey     string
}
type UserServiceGrpc struct {
	Connection string
	APIKey     string
}

type DB struct {
	Mysql Mysql
}

type Cache struct {
	Redis Redis
}

type Redis struct {
	Address  string
	Port     int
	Password string
	DbQuery  int
	DbAuth   int
}

type Mysql struct {
	Host            string
	Port            int
	UserName        string
	Password        string
	Database        string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
	Options         string
}

type Mongodb struct {
	Address         string
	Username        string
	Password        string
	DbName          string
	ConnectTimeout  time.Duration
	MaxConnIdleTime int
	MinPoolSize     uint64
	MaxPoolSize     uint64
}

type RabbitMQ struct {
	SagaOrderEvent    SagaOrderEvent
	RatingReplyEvent  RatingReplyEvent
	NotificationEvent NotificationEvent
	Connection        string
	NameService       string
}

type RatingReplyEvent struct {
	Exchange   string
	RoutingKey string
	Queue      string
}

type NotificationEvent struct {
	Exchange   string
	RoutingKey string
	Queue      string
}

type SagaOrderEvent struct {
	Exchange          string
	PublishRoutingKey string
	CancelRoutingKey  string
	ReplyRoutingKey   string
	Queue             string
}

type AdapterService struct {
	AuthService      AuthService
	UserService      UserService
	ProductService   ProductService
	EmailService     EmailService
	DeliveryService  DeliveryService
	PromotionService PromotionService
	StoreService     StoreService
}
type AuthService struct {
	BaseURL     string
	InternalKey string
}

type UserService struct {
	AuthURL     string
	UserURL     string
	InternalKey string
}

type ProductService struct {
	BaseURL     string
	InternalKey string
}

type StoreService struct {
	BaseURL     string
	InternalKey string
}

type EmailService struct {
	Email string
	Host  string
	Key   string
}

type DeliveryService struct {
	BaseURL     string
	InternalKey string
}

type PromotionService struct {
	BaseURL     string
	InternalKey string
}

type Metrics struct {
	Username      string
	Password      string
	PrometheusURL string
	FiberURL      string
}

// Get config path for local or docker
func getDefaultConfig() string {
	return "/config/config"
}

// Load config file from given path
func NewConfig() (*Config, error) {
	config := Config{}
	path := os.Getenv("cfgPath")
	if path == "" {
		path = getDefaultConfig()
	}
	fmt.Printf("config path:%s\n", path)

	v := viper.New()

	v.SetConfigName(path)
	v.AddConfigPath(".")
	v.AutomaticEnv()
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, errors.New("config file not found")
		}
		return nil, err
	}

	err := v.Unmarshal(&config)
	if err != nil {
		log.Printf("unable to decode into struct, %v", err)
		return nil, err
	}

	return &config, nil
}
