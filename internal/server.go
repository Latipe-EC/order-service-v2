//go:build wireinject
// +build wireinject

package server

import (
	"encoding/json"
	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/google/wire"
	"latipe-order-service-v2/config"
	"latipe-order-service-v2/internal/api"
	"latipe-order-service-v2/internal/app"
	"latipe-order-service-v2/internal/common/errors"
	"latipe-order-service-v2/internal/infrastructure/adapter/authserv"
	"latipe-order-service-v2/internal/infrastructure/adapter/deliveryserv"
	"latipe-order-service-v2/internal/infrastructure/adapter/productserv"
	"latipe-order-service-v2/internal/infrastructure/adapter/storeserv"
	"latipe-order-service-v2/internal/infrastructure/adapter/userserv"
	voucherserv "latipe-order-service-v2/internal/infrastructure/adapter/vouchersev"
	"latipe-order-service-v2/internal/infrastructure/persistence"
	"latipe-order-service-v2/internal/middleware"
	producer "latipe-order-service-v2/internal/msgqueue"
	router2 "latipe-order-service-v2/internal/router"
	"latipe-order-service-v2/internal/worker"
	"latipe-order-service-v2/pkg/cache"
)

type Server struct {
	app       *fiber.App
	cfg       *config.Config
	orderSubs *worker.OrderTransactionSubscriber
}

func New() (*Server, error) {
	panic(wire.Build(wire.NewSet(
		NewServer,
		config.Set,
		api.Set,
		router2.Set,
		persistence.Set,
		userserv.Set,
		authserv.Set,
		deliveryserv.Set,
		productserv.Set,
		app.Set,
		middleware.Set,
		cache.Set,
		voucherserv.Set,
		storeserv.Set,
		producer.Set,
		worker.Set,
	)))
}

func NewServer(
	cfg *config.Config,
	orderRouter router2.OrderRouter,
	orderSubs *worker.OrderTransactionSubscriber) *Server {

	app := fiber.New(fiber.Config{
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		JSONDecoder:  json.Unmarshal,
		JSONEncoder:  json.Marshal,
		ErrorHandler: errors.CustomErrorHandler,
	})

	prometheus := fiberprometheus.New("latipe-order-service-v2")
	prometheus.RegisterAt(app, "/metrics")
	app.Use(prometheus.Middleware)

	// Initialize default config
	app.Use(logger.New())

	app.Get("", func(ctx *fiber.Ctx) error {
		s := struct {
			Message string `json:"message"`
			Version string `json:"version"`
		}{
			Message: "Order rest-api was developed by TienDat",
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

func (serv Server) OrderTransactionSubscriber() *worker.OrderTransactionSubscriber {
	return serv.orderSubs
}
