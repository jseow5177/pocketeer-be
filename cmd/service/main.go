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

	"github.com/jseow5177/pockteer-be/api/handler/category"
	"github.com/jseow5177/pockteer-be/api/middleware"
	"github.com/jseow5177/pockteer-be/config"
	"github.com/jseow5177/pockteer-be/data/presenter"
	"github.com/jseow5177/pockteer-be/dep/repo"
	"github.com/jseow5177/pockteer-be/dep/repo/mongo"
	"github.com/jseow5177/pockteer-be/pkg/logger"
	"github.com/jseow5177/pockteer-be/pkg/router"
	"github.com/jseow5177/pockteer-be/pkg/service"
)

type server struct {
	ctx   context.Context
	cfg   *config.Config
	mongo *mongo.Mongo

	categoryRepo repo.CategoryRepo
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
	var err error

	// init logger
	s.ctx = logger.InitZeroLog(context.Background(), s.cfg.Server.LogLevel)

	// init rate limiter
	limiter := middleware.NewRateLimiter(s.cfg.Server.RateLimits)

	// init mongo
	s.mongo, err = mongo.NewMongo(s.ctx, s.cfg.Mongo)
	if err != nil {
		log.Ctx(s.ctx).Error().Msgf("fail to init mongo client, err: %v", err)
		return err
	}
	defer func() {
		if err != nil {
			_ = s.mongo.Close(s.ctx)
		}
	}()

	// init repos
	s.categoryRepo = mongo.NewCategoryMongo(s.mongo)

	// start server
	addr := fmt.Sprintf(":%d", s.cfg.Server.Port)
	go func() {
		log.Info().Msgf("starting HTTP server at %s", addr)

		httpServer := &http.Server{
			BaseContext: func(_ net.Listener) context.Context {
				return s.ctx
			},
			Addr:    addr,
			Handler: middleware.RateLimit(s.ctx, limiter, middleware.Log(s.registerRoutes())),
		}
		err := httpServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal().Msgf("fail to start HTTP server, err: %v", err)
		}
	}()

	return nil
}

func (s *server) Stop() error {
	ctx, cancel := context.WithTimeout(s.ctx, 30*time.Second)
	defer cancel()

	if err := s.mongo.Close(ctx); err != nil {
		log.Ctx(ctx).Error().Msgf("close mongo fail, err: %v", err)
	}

	return nil
}

func (s *server) registerRoutes() http.Handler {
	r := &router.HttpRouter{
		Router: mux.NewRouter(),
	}

	// ========== Healthcheck ========== //

	// healthcheck
	r.RegisterHttpRoute(&router.HttpRoute{
		Path:   config.PathHealthCheck,
		Method: http.MethodGet,
		Handler: router.Handler{
			Req:       new(presenter.HealthCheckRequest),
			Res:       new(presenter.HealthCheckResponse),
			Validator: nil,
			HandleFunc: func(ctx context.Context, req, res interface{}) error {
				return nil
			},
		},
	})

	// ========== Category ========== //

	catHandler := category.NewCategoryHandler(s.categoryRepo)

	// create category
	r.RegisterHttpRoute(&router.HttpRoute{
		Path:   config.PathCreateCategory,
		Method: http.MethodPost,
		Handler: router.Handler{
			Req:       new(presenter.CreateCategoryRequest),
			Res:       new(presenter.CreateCategoryResponse),
			Validator: category.CreateCategoryValidator,
			HandleFunc: func(ctx context.Context, req, res interface{}) error {
				return catHandler.CreateCategory(ctx, req.(*presenter.CreateCategoryRequest), res.(*presenter.CreateCategoryResponse))
			},
		},
	})

	// get categories
	r.RegisterHttpRoute(&router.HttpRoute{
		Path:   config.PathGetCategories,
		Method: http.MethodPost,
		Handler: router.Handler{
			Req:       new(presenter.GetCategoriesRequest),
			Res:       new(presenter.GetCategoriesResponse),
			Validator: category.GetCategoriesValidator,
			HandleFunc: func(ctx context.Context, req, res interface{}) error {
				return catHandler.GetCategories(ctx, req.(*presenter.GetCategoriesRequest), res.(*presenter.GetCategoriesResponse))
			},
		},
	})

	return r
}
