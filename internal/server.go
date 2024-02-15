//go:build wireinject
// +build wireinject

// this code to enable wire inject
package server

import (
	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/google/wire"
	"github.com/hellofresh/health-go/v5"
	"latipe-order-service-v2/config"
	"latipe-order-service-v2/internal/api"
	"latipe-order-service-v2/internal/common/errors"
	"latipe-order-service-v2/internal/infrastructure/adapter/authserv"
	"latipe-order-service-v2/internal/infrastructure/adapter/deliveryserv"
	"latipe-order-service-v2/internal/infrastructure/adapter/productserv"
	"latipe-order-service-v2/internal/infrastructure/adapter/storeserv"
	"latipe-order-service-v2/internal/infrastructure/adapter/userserv"
	"latipe-order-service-v2/internal/infrastructure/adapter/vouchersev"
	grpc_adapt "latipe-order-service-v2/internal/infrastructure/grpc"
	"latipe-order-service-v2/internal/infrastructure/persistence"
	"latipe-order-service-v2/internal/middleware"
	producer "latipe-order-service-v2/internal/publisher"
	"latipe-order-service-v2/internal/router"
	"latipe-order-service-v2/internal/services"
	"latipe-order-service-v2/internal/subscriber"
	"latipe-order-service-v2/pkg/cache"
	healthcheckService "latipe-order-service-v2/pkg/healthcheck"
	"latipe-order-service-v2/pkg/rabbitclient"
)

type Server struct {
	app       *fiber.App
	cfg       *config.Config
	orderSubs *subscriber.PurchaseReplySubscriber
}

func New() (*Server, error) {
	panic(wire.Build(wire.NewSet(
		NewServer,
		config.Set,
		api.Set,
		router.Set,
		rabbitclient.Set,
		persistence.Set,
		grpc_adapt.Set,
		userserv.Set,
		authserv.Set,
		deliveryserv.Set,
		productserv.Set,
		services.Set,
		middleware.Set,
		cache.Set,
		voucherserv.Set,
		storeserv.Set,
		producer.Set,
		subscriber.Set,
	)))
}

func NewServer(
	cfg *config.Config,
	orderRouter router.OrderRouter,
	orderSubs *subscriber.PurchaseReplySubscriber) *Server {

	app := fiber.New(fiber.Config{
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		JSONDecoder:  sonic.Unmarshal,
		JSONEncoder:  sonic.Marshal,
		ErrorHandler: errors.CustomErrorHandler,
	})

	//providing basic authentication for metrics endpoints
	basicAuthConfig := basicauth.Config{
		Users: map[string]string{
			cfg.Metrics.Username: cfg.Metrics.Password,
		},
	}

	// Healthcheck
	h, _ := healthcheckService.NewHealthCheckService(cfg)
	app.Get("/status", basicauth.New(basicAuthConfig), adaptor.HTTPHandlerFunc(h.HandlerFunc))
	app.Use(healthcheck.New())
	app.Use(healthcheck.New(healthcheck.Config{
		LivenessProbe: func(c *fiber.Ctx) bool {
			return true
		},
		LivenessEndpoint: "/live",
		ReadinessProbe: func(c *fiber.Ctx) bool {
			result := h.Measure(c.Context())
			return result.Status == health.StatusOK
		},
		ReadinessEndpoint: "/ready",
	}))

	//fiber dashboard
	app.Get(cfg.Metrics.FiberURL, basicauth.New(basicAuthConfig),
		monitor.New(monitor.Config{Title: "Orders Services Metrics Page (Fiber)"}))

	prometheus := fiberprometheus.New("latipe-order-service-v2")
	prometheus.RegisterAt(app, cfg.Metrics.PrometheusURL, basicauth.New(basicAuthConfig))
	app.Use(prometheus.Middleware)

	// Initialize default config
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

	orderRouter.Init(&v2)

	return &Server{
		cfg:       cfg,
		app:       app,
		orderSubs: orderSubs,
	}
}

func (serv Server) App() *fiber.App {
	return serv.app
}

func (serv Server) Config() *config.Config {
	return serv.cfg
}

func (serv Server) OrderTransactionSubscriber() *subscriber.PurchaseReplySubscriber {
	return serv.orderSubs
}
