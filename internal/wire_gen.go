// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package server

import (
	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/cors"
	healthcheck2 "github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/swagger"
	"github.com/hellofresh/health-go/v5"
	"latipe-order-service-v2/config"
	order2 "latipe-order-service-v2/internal/api/order"
	"latipe-order-service-v2/internal/common/errors"
	"latipe-order-service-v2/internal/infrastructure/adapter/authserv"
	"latipe-order-service-v2/internal/infrastructure/adapter/deliveryserv"
	"latipe-order-service-v2/internal/infrastructure/adapter/storeserv"
	"latipe-order-service-v2/internal/infrastructure/grpc/deliveryServ"
	"latipe-order-service-v2/internal/infrastructure/grpc/productServ"
	"latipe-order-service-v2/internal/infrastructure/grpc/promotionServ"
	"latipe-order-service-v2/internal/infrastructure/grpc/userServ"
	"latipe-order-service-v2/internal/infrastructure/persistence/commission"
	"latipe-order-service-v2/internal/infrastructure/persistence/db"
	"latipe-order-service-v2/internal/infrastructure/persistence/order"
	"latipe-order-service-v2/internal/middleware"
	"latipe-order-service-v2/internal/middleware/auth"
	"latipe-order-service-v2/internal/publisher"
	"latipe-order-service-v2/internal/router/admin"
	"latipe-order-service-v2/internal/router/delivery"
	"latipe-order-service-v2/internal/router/internalRouter"
	"latipe-order-service-v2/internal/router/statistic"
	"latipe-order-service-v2/internal/router/store"
	"latipe-order-service-v2/internal/router/user"
	"latipe-order-service-v2/internal/services/commands/orderCmd"
	"latipe-order-service-v2/internal/services/queries/orderQuery"
	"latipe-order-service-v2/internal/services/queries/statisticQuery"
	"latipe-order-service-v2/internal/subscriber/purchase"
	"latipe-order-service-v2/internal/subscriber/rating"
	"latipe-order-service-v2/pkg/cache"
	"latipe-order-service-v2/pkg/healthcheck"
	"latipe-order-service-v2/pkg/rabbitclient"
)

// Injectors from server.go:

func New() (*Server, error) {
	configConfig, err := config.NewConfig()
	if err != nil {
		return nil, err
	}
	cacheEngine, err := cache.NewCacheEngineV8(configConfig)
	if err != nil {
		return nil, err
	}
	gorm := db.NewMySQLConnection(configConfig, cacheEngine)
	orderRepository := order.NewGormRepository(gorm)
	commissionRepository := commission.NewCommissionRepository(gorm)
	cacheV9CacheEngine, err := cache.NewCacheEngineV9(configConfig)
	if err != nil {
		return nil, err
	}
	connection := rabbitclient.NewRabbitClientConnection(configConfig)
	publisherTransactionMessage := publisher.NewTransactionProducer(configConfig, connection)
	voucherServiceClient := vouchergrpc.NewVoucherClientGrpcImpl(configConfig)
	productServiceClient := productgrpc.NewProductGrpcClientImpl(configConfig)
	deliveryServiceClient := deliverygrpc.NewDeliveryServiceGRPCClientImpl(configConfig)
	userServiceClient := usergrpc.NewUserServiceClientGRPCImpl(configConfig)
	service := storeserv.NewStoreServiceAdapter(configConfig)
	orderCommandUsecase := orderCmd.NewOrderCommandService(configConfig, orderRepository, commissionRepository, cacheV9CacheEngine, publisherTransactionMessage, voucherServiceClient, productServiceClient, deliveryServiceClient, userServiceClient, service)
	orderQueryUsecase := orderQuery.NewOrderQueryService(orderRepository)
	orderApiHandler := order2.NewOrderHandler(orderCommandUsecase, orderQueryUsecase)
	authservService := authserv.NewAuthServHttpAdapter(configConfig)
	deliveryservService := deliveryserv.NewDeliServHttpAdapter(configConfig)
	authenticationMiddleware := auth.NewAuthMiddleware(authservService, service, deliveryservService, configConfig, cacheV9CacheEngine)
	middlewareMiddleware := middleware.NewMiddleware(authenticationMiddleware)
	adminOrderRouter := adminRouter.NewAdminOrderRouter(orderApiHandler, middlewareMiddleware)
	userOrderRouter := userRouter.NewUserOrderRouter(orderApiHandler, middlewareMiddleware)
	storeOrderRouter := storeRouter.NewStoreOrderRouter(orderApiHandler, middlewareMiddleware)
	deliveryOrderRouter := deliveryRouter.NewDeliveryOrderRouter(orderApiHandler, middlewareMiddleware)
	orderStatisticUsecase := statisticQuery.NewOrderStatisicService(orderRepository)
	orderStatisticApiHandler := order2.NewStatisticHandler(orderStatisticUsecase)
	orderStatisticRouter := statisticRouter.NewStatisticOrderRouter(orderStatisticApiHandler, middlewareMiddleware)
	internalOrderRouter := internalRouter.NewInternalOrderRouter(orderApiHandler, orderStatisticApiHandler, middlewareMiddleware)
	purchaseReplySubscriber := purchase.NewPurchaseReplySubscriber(configConfig, connection, orderCommandUsecase)
	ratingItemSubscriber := rating.NewRatingItemSubscriber(configConfig, connection, orderCommandUsecase)
	server := NewServer(configConfig, adminOrderRouter, userOrderRouter, storeOrderRouter, deliveryOrderRouter, orderStatisticRouter, internalOrderRouter, purchaseReplySubscriber, ratingItemSubscriber)
	return server, nil
}

// server.go:

type Server struct {
	app        *fiber.App
	cfg        *config.Config
	orderSubs  *purchase.PurchaseReplySubscriber
	ratingSubs *rating.RatingItemSubscriber
}

func NewServer(
	cfg *config.Config, adminRouter2 adminRouter.AdminOrderRouter, userRouter2 userRouter.UserOrderRouter, storeRouter2 storeRouter.StoreOrderRouter, deliveryRouter2 deliveryRouter.DeliveryOrderRouter, statisticRouter2 statisticRouter.OrderStatisticRouter, internalRouter2 internalRouter.InternalOrderRouter,

	orderSubs *purchase.PurchaseReplySubscriber,
	ratingSubs *rating.RatingItemSubscriber) *Server {

	app := fiber.New(fiber.Config{
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		JSONDecoder:  sonic.Unmarshal,
		JSONEncoder:  sonic.Marshal,
		ErrorHandler: errors.CustomErrorHandler,
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://127.0.0.1:5500",
		AllowHeaders: "Origin, X-Requested-With, Content-Type, Accept, Authorization",
		AllowMethods: "GET,HEAD,OPTIONS,POST,PUT",
	}))

	basicAuthConfig := basicauth.Config{
		Users: map[string]string{
			cfg.Metrics.Username: cfg.Metrics.Password,
		},
	}

	h, _ := healthcheck.NewHealthCheckService(cfg)
	app.Get("/health", basicauth.New(basicAuthConfig), adaptor.HTTPHandlerFunc(h.HandlerFunc))
	app.Use(healthcheck2.New())
	app.Use(healthcheck2.New(healthcheck2.Config{
		LivenessProbe: func(c *fiber.Ctx) bool {
			return true
		},
		LivenessEndpoint: "/liveness",
		ReadinessProbe: func(c *fiber.Ctx) bool {
			result := h.Measure(c.Context())
			return result.Status == health.StatusOK
		},
		ReadinessEndpoint: "/readiness",
	}))

	app.Get("/swagger/*", basicauth.New(basicAuthConfig), swagger.HandlerDefault)

	app.Get(cfg.Metrics.FiberURL, basicauth.New(basicAuthConfig), monitor.New(monitor.Config{Title: "Orders Services Metrics Page (Fiber)"}))

	prometheus := fiberprometheus.New("latipe-order-service-v2")
	prometheus.RegisterAt(app, cfg.Metrics.PrometheusURL, basicauth.New(basicAuthConfig))
	app.Use(prometheus.Middleware)

	app.Use(logger.New())

	app.Get("", func(ctx *fiber.Ctx) error {
		s := struct {
			Message string `json:"message"`
			Version string `json:"version"`
		}{
			Message: "the orders service was developed by tdat.it",
			Version: "v2.0.0",
		}
		return ctx.JSON(s)
	})

	api := app.Group("/api")
	v2 := api.Group("/v2")
	orders := v2.Group("/orders")
	userRouter2.
		Init(&orders)
	storeRouter2.
		Init(&orders)
	adminRouter2.
		Init(&orders)
	deliveryRouter2.
		Init(&orders)
	statisticRouter2.
		Init(&orders)
	internalRouter2.
		Init(&orders)

	return &Server{
		cfg:        cfg,
		app:        app,
		orderSubs:  orderSubs,
		ratingSubs: ratingSubs,
	}
}

func (serv Server) App() *fiber.App {
	return serv.app
}

func (serv Server) Config() *config.Config {
	return serv.cfg
}

func (serv Server) OrderTransactionSubscriber() *purchase.PurchaseReplySubscriber {
	return serv.orderSubs
}

func (serv Server) RatingItemSubscriber() *rating.RatingItemSubscriber {
	return serv.ratingSubs
}
