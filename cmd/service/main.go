package main

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"

	"github.com/jseow5177/pockteer-be/api/handler/budget"
	"github.com/jseow5177/pockteer-be/api/handler/category"
	"github.com/jseow5177/pockteer-be/api/middleware"
	"github.com/jseow5177/pockteer-be/config"
	"github.com/jseow5177/pockteer-be/model"
	"github.com/jseow5177/pockteer-be/pkg/logger"
	"github.com/jseow5177/pockteer-be/pkg/router"
	"github.com/jseow5177/pockteer-be/pkg/service"
)

type server struct {
	cfg *config.Config
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	s := new(server)
	if err := service.Run(s); err != nil {
		log.Fatal().Msg(err.Error())
	}
}

func (s *server) Init() error {
	s.cfg = config.NewConfig()
	return nil
}

func (s *server) Start() error {
	ctx := context.Background()

	// init logger
	ctx = logger.InitZeroLog(ctx, s.cfg.Server.LogLevel)

	// init rate limiter
	limiter := middleware.NewRateLimiter(s.cfg.Server.RateLimits)

	// TODO: init Firestore

	// TODO: init repositories

	// TODO: init use cases

	// start server
	addr := fmt.Sprintf(":%d", s.cfg.Server.Port)
	go func() {
		log.Info().Msgf("starting HTTP server at %s", addr)

		httpServer := &http.Server{
			BaseContext: func(_ net.Listener) context.Context {
				return ctx
			},
			Addr:    addr,
			Handler: middleware.RateLimit(ctx, limiter, middleware.Log(s.registerRoutes())),
		}
		err := httpServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal().Msgf("fail to start HTTP server, err: %v", err)
		}
	}()

	return nil
}

func (s *server) Stop() error {
	// close dependencies
	return nil
}

func (s *server) registerRoutes() http.Handler {
	r := &router.HttpRouter{
		Router: mux.NewRouter(),
	}

	// healthcheck
	r.RegisterHttpRoute(&router.HttpRoute{
		Path:   config.PathHealthCheck,
		Method: http.MethodGet,
		Handler: router.Handler{
			Req:       new(model.HealthCheckRequest),
			Res:       new(model.HealthCheckResponse),
			Validator: nil,
			HandleFunc: func(ctx context.Context, req, res interface{}) error {
				return nil
			},
		},
	})

	// create category
	r.RegisterHttpRoute(&router.HttpRoute{
		Path:   config.PathCreateCategory,
		Method: http.MethodPost,
		Handler: router.Handler{
			Req:       new(model.CreateCategoryRequest),
			Res:       new(model.CreateCategoryResponse),
			Validator: category.CreateCategoryValidator,
			HandleFunc: func(ctx context.Context, req, res interface{}) error {
				return category.CreateCategory(ctx, req.(*model.CreateCategoryRequest), res.(*model.CreateCategoryResponse))
			},
		},
	})

	// create budget
	r.RegisterHttpRoute(&router.HttpRoute{
		Path:   config.PathCreateBudget,
		Method: http.MethodPost,
		Handler: router.Handler{
			Req:       new(model.CreateBudgetRequest),
			Res:       new(model.CreateCategoryResponse),
			Validator: budget.CreateBudgetValidator,
			HandleFunc: func(ctx context.Context, req, res interface{}) error {
				return budget.CreateBudget(ctx, req.(*model.CreateBudgetRequest), res.(*model.CreateBudgetResponse))
			},
		},
	})

	return r
}
