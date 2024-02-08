package auth

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"latipe-order-service-v2/config"
	"latipe-order-service-v2/internal/common/errors"
	"latipe-order-service-v2/internal/infrastructure/adapter/authserv"
	"latipe-order-service-v2/internal/infrastructure/adapter/authserv/dto"
	"latipe-order-service-v2/internal/infrastructure/adapter/deliveryserv"
	deliDto "latipe-order-service-v2/internal/infrastructure/adapter/deliveryserv/dto"
	"latipe-order-service-v2/internal/infrastructure/adapter/storeserv"
	storeDTO "latipe-order-service-v2/internal/infrastructure/adapter/storeserv/dto"
	cacheEngine "latipe-order-service-v2/pkg/cache/redisCacheV9"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2/log"
)

type AuthenticationMiddleware struct {
	authServ    authserv.Service
	storeServ   storeserv.Service
	delivery    deliveryserv.Service
	cacheEngine *cacheEngine.CacheEngine
	cfg         *config.Config
}

func NewAuthMiddleware(authServ authserv.Service, storeServ storeserv.Service,
	deliServ deliveryserv.Service, config *config.Config, cacheEngine *cacheEngine.CacheEngine) *AuthenticationMiddleware {
	return &AuthenticationMiddleware{
		authServ:    authServ,
		storeServ:   storeServ,
		delivery:    deliServ,
		cfg:         config,
		cacheEngine: cacheEngine,
	}
}

func (a AuthenticationMiddleware) RequiredAuthentication() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		bearToken := ctx.Get("Authorization")
		if bearToken == "" {
			return errors.ErrUnauthenticated
		}

		bearToken = strings.Split(bearToken, " ")[1]

		var resp *dto.AuthorizationResponse
		//get cache
		data, err := a.cacheEngine.Get(fmt.Sprintf("auth:%v", bearToken))
		if data != nil {
			err = json.Unmarshal(data, &resp)
			if err != nil {
				log.Error(err)
			}
			log.Info("get cache auth data")
		} else {
			req := dto.AuthorizationRequest{
				Token: bearToken,
			}
			resp, err = a.authServ.Authorization(ctx.Context(), &req)
			if err != nil {
				return err
			}
			log.Info("set cache auth data")
			err = a.cacheEngine.Set(fmt.Sprintf("auth:%v", bearToken), resp, 5*time.Minute)
		}

		ctx.Locals(USERNAME, resp.Email)
		ctx.Locals(USER_ID, resp.Id)
		ctx.Locals(ROLE, resp.Role)
		ctx.Locals(BEARER_TOKEN, bearToken)
		return ctx.Next()
	}
}

func (a AuthenticationMiddleware) RequiredStoreAuthentication() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		bearToken := ctx.Get("Authorization")
		if bearToken == "" {
			return errors.ErrUnauthenticated
		}

		bearToken = strings.Split(bearToken, " ")[1]

		var resp *dto.AuthorizationResponse
		//get cache
		auth, err := a.cacheEngine.Get(fmt.Sprintf("auth:%v", bearToken))
		if auth != nil {
			err = json.Unmarshal(auth, &resp)
			if err != nil {
				log.Error(err)
			}
			log.Info("get cache auth data")
		} else {
			req := dto.AuthorizationRequest{
				Token: bearToken,
			}
			resp, err = a.authServ.Authorization(ctx.Context(), &req)
			if err != nil {
				return err
			}
			log.Info("set cache auth data")
			err = a.cacheEngine.Set(fmt.Sprintf("auth:%v", bearToken), resp, 5*time.Minute)
		}

		//get store cache
		var storeResp *storeDTO.GetStoreIdByUserResponse
		store, err := a.cacheEngine.Get(fmt.Sprintf("store:%v", resp.Id))
		if store != nil {
			err = json.Unmarshal(store, &resp)
			if err != nil {
				log.Error(err)
			}
			log.Info("get cache auth data")
		} else {
			//validate store
			storeReq := storeDTO.GetStoreIdByUserRequest{UserID: resp.Id}
			storeReq.BaseHeader.BearToken = bearToken

			storeResp, err = a.storeServ.GetStoreByUserId(ctx.Context(), &storeReq)
			if err != nil {
				return err
			}

			if storeResp.StoreID == "" {
				return errors.ErrPermissionDenied
			}
			log.Info("set cache auth data")
			err = a.cacheEngine.Set(fmt.Sprintf("store:%v", storeResp.StoreID), resp, 5*time.Minute)
		}

		ctx.Locals(USERNAME, resp.Email)
		ctx.Locals(USER_ID, resp.Id)
		ctx.Locals(ROLE, resp.Role)
		ctx.Locals(BEARER_TOKEN, bearToken)
		ctx.Locals(STORE_ID, storeResp.StoreID)

		return ctx.Next()
	}
}

func (a AuthenticationMiddleware) RequiredDeliveryAuthentication() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		bearToken := ctx.Get("Authorization")
		if bearToken == "" {
			return errors.ErrUnauthenticated
		}

		bearToken = strings.Split(bearToken, " ")[1]
		var resp *dto.AuthorizationResponse

		//get cache
		auth, err := a.cacheEngine.Get(fmt.Sprintf("auth:%v", bearToken))
		if auth != nil {
			err = json.Unmarshal(auth, &resp)
			if err != nil {
				log.Error(err)
			}
			log.Info("get cache auth data")
		} else {
			req := dto.AuthorizationRequest{
				Token: bearToken,
			}
			resp, err = a.authServ.Authorization(ctx.Context(), &req)
			if err != nil {
				return err
			}
			log.Info("set cache auth data")
			err = a.cacheEngine.Set(fmt.Sprintf("auth:%v", bearToken), resp, 5*time.Minute)
		}

		//get delivery cache

		var deliResp *deliDto.GetDeliveryByTokenResponse
		store, err := a.cacheEngine.Get(fmt.Sprintf("deli:%v", resp.Id))
		if store != nil {
			err = json.Unmarshal(store, &resp)
			if err != nil {
				log.Error(err)
			}

		} else {
			deliReq := deliDto.GetDeliveryByTokenRequest{BearerToken: bearToken}

			deliResp, err = a.delivery.GetDeliveryByToken(ctx.Context(), &deliReq)
			if err != nil {
				return err
			}

			if deliResp.ID == "" {
				return errors.ErrPermissionDenied
			}

			err = a.cacheEngine.Set(fmt.Sprintf("deli:%v", deliResp.ID), resp, 5*time.Minute)
		}

		ctx.Locals(USERNAME, resp.Email)
		ctx.Locals(USER_ID, resp.Id)
		ctx.Locals(ROLE, resp.Role)
		ctx.Locals(BEARER_TOKEN, bearToken)
		ctx.Locals(DELIVERY_ID, deliResp.ID)

		return ctx.Next()
	}
}

func (a AuthenticationMiddleware) RequiredRole(roles []string) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		bearToken := ctx.Get("Authorization")
		if bearToken == "" {
			return errors.ErrUnauthenticated
		}

		bearToken = strings.Split(bearToken, " ")[1]
		var resp *dto.AuthorizationResponse
		//get cache
		auth, err := a.cacheEngine.Get(fmt.Sprintf("auth:%v", bearToken))
		if auth != nil {
			err = json.Unmarshal(auth, &resp)
			if err != nil {
				log.Error(err)
			}
			log.Info("get cache auth data")
		} else {
			req := dto.AuthorizationRequest{
				Token: bearToken,
			}
			resp, err = a.authServ.Authorization(ctx.Context(), &req)
			if err != nil {
				return err
			}
			log.Info("set cache auth data")
			err = a.cacheEngine.Set(fmt.Sprintf("auth:%v", bearToken), resp, 5*time.Minute)
		}

		for _, i := range roles {
			if i == resp.Role {
				ctx.Locals(USERNAME, resp.Email)
				ctx.Locals(USER_ID, resp.Id)
				ctx.Locals(ROLE, resp.Role)
				ctx.Locals(BEARER_TOKEN, bearToken)
				return ctx.Next()
			}
		}
		return errors.ErrPermissionDenied
	}
}

func (a AuthenticationMiddleware) RequiredInternalService() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		token := ctx.Get("X-API-KEY")
		if token == "" || token != a.cfg.Server.ApiHeaderKey {
			return errors.ErrUnauthenticated
		}

		return ctx.Next()
	}
}
