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

	"github.com/jseow5177/pockteer-be/api/middleware"
	"github.com/jseow5177/pockteer-be/api/presenter"
	"github.com/jseow5177/pockteer-be/config"
	"github.com/jseow5177/pockteer-be/dep/repo"
	"github.com/jseow5177/pockteer-be/dep/repo/mongo"
	"github.com/jseow5177/pockteer-be/pkg/logger"
	"github.com/jseow5177/pockteer-be/pkg/router"
	"github.com/jseow5177/pockteer-be/pkg/service"

	bh "github.com/jseow5177/pockteer-be/api/handler/budget"
	ch "github.com/jseow5177/pockteer-be/api/handler/category"
	th "github.com/jseow5177/pockteer-be/api/handler/transaction"
	uh "github.com/jseow5177/pockteer-be/api/handler/user"

	auc "github.com/jseow5177/pockteer-be/usecase/aggr"
	buc "github.com/jseow5177/pockteer-be/usecase/budget"
	cuc "github.com/jseow5177/pockteer-be/usecase/category"
	ttuc "github.com/jseow5177/pockteer-be/usecase/token"
	tuc "github.com/jseow5177/pockteer-be/usecase/transaction"
	uuc "github.com/jseow5177/pockteer-be/usecase/user"
)

type server struct {
	ctx   context.Context
	cfg   *config.Config
	mongo *mongo.Mongo

	categoryRepo    repo.CategoryRepo
	transactionRepo repo.TransactionRepo
	budgetRepo      repo.BudgetRepo
	userRepo        repo.UserRepo

	categoryUseCase    cuc.UseCase
	transactionUseCase tuc.UseCase
	budgetUseCase      buc.UseCase
	aggrUseCase        auc.UseCase
	userUseCase        uuc.UseCase
	tokenUseCase       ttuc.UseCase
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
	s.transactionRepo = mongo.NewTransactionMongo(s.mongo)
	s.budgetRepo = mongo.NewBudgetMongo(s.mongo)
	s.userRepo = mongo.NewUserMongo(s.mongo)

	// init use cases
	s.categoryUseCase = cuc.NewCategoryUseCase(s.categoryRepo)
	s.transactionUseCase = tuc.NewTransactionUseCase(s.categoryUseCase, s.transactionRepo)
	s.budgetUseCase = buc.NewBudgetUseCase(s.budgetRepo, s.categoryRepo)
	s.aggrUseCase = auc.NewAggrUseCase(s.budgetUseCase, s.categoryUseCase)
	s.tokenUseCase = ttuc.NewTokenUseCase(s.cfg.AccessToken, s.cfg.RefreshToken)
	s.userUseCase = uuc.NewUserUseCase(s.userRepo, s.tokenUseCase)

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

	// TODO: Handle inflight requests

	if err := s.mongo.Close(ctx); err != nil {
		log.Ctx(ctx).Error().Msgf("close mongo fail, err: %v", err)
	}

	return nil
}

func (s *server) registerRoutes() http.Handler {
	r := &router.HttpRouter{
		Router: mux.NewRouter(),
	}

	authMiddleware := middleware.NewAuthMiddleware(s.tokenUseCase)

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

	categoryHandler := ch.NewCategoryHandler(s.categoryUseCase)

	// create category
	r.RegisterHttpRoute(&router.HttpRoute{
		Path:   config.PathCreateCategory,
		Method: http.MethodPost,
		Handler: router.Handler{
			Req:       new(presenter.CreateCategoryRequest),
			Res:       new(presenter.CreateCategoryResponse),
			Validator: ch.CreateCategoryValidator,
			HandleFunc: func(ctx context.Context, req, res interface{}) error {
				return categoryHandler.CreateCategory(ctx, req.(*presenter.CreateCategoryRequest), res.(*presenter.CreateCategoryResponse))
			},
		},
		Middlewares: []router.Middleware{authMiddleware},
	})

	// update category
	r.RegisterHttpRoute(&router.HttpRoute{
		Path:   config.PathUpdateCategory,
		Method: http.MethodPost,
		Handler: router.Handler{
			Req:       new(presenter.UpdateCategoryRequest),
			Res:       new(presenter.UpdateCategoryResponse),
			Validator: ch.UpdateCategoryValidator,
			HandleFunc: func(ctx context.Context, req, res interface{}) error {
				return categoryHandler.UpdateCategory(ctx, req.(*presenter.UpdateCategoryRequest), res.(*presenter.UpdateCategoryResponse))
			},
		},
		Middlewares: []router.Middleware{authMiddleware},
	})

	// get category
	r.RegisterHttpRoute(&router.HttpRoute{
		Path:   config.PathGetCategory,
		Method: http.MethodPost,
		Handler: router.Handler{
			Req:       new(presenter.GetCategoryRequest),
			Res:       new(presenter.GetCategoryResponse),
			Validator: ch.GetCategoryValidator,
			HandleFunc: func(ctx context.Context, req, res interface{}) error {
				return categoryHandler.GetCategory(ctx, req.(*presenter.GetCategoryRequest), res.(*presenter.GetCategoryResponse))
			},
		},
		Middlewares: []router.Middleware{authMiddleware},
	})

	// get categories
	r.RegisterHttpRoute(&router.HttpRoute{
		Path:   config.PathGetCategories,
		Method: http.MethodPost,
		Handler: router.Handler{
			Req:       new(presenter.GetCategoriesRequest),
			Res:       new(presenter.GetCategoriesResponse),
			Validator: ch.GetCategoriesValidator,
			HandleFunc: func(ctx context.Context, req, res interface{}) error {
				return categoryHandler.GetCategories(ctx, req.(*presenter.GetCategoriesRequest), res.(*presenter.GetCategoriesResponse))
			},
		},
		Middlewares: []router.Middleware{authMiddleware},
	})

	// ========== Transaction ========== //

	transactionHandler := th.NewTransactionHandler(s.transactionUseCase)

	// create transaction
	r.RegisterHttpRoute(&router.HttpRoute{
		Path:   config.PathCreateTransaction,
		Method: http.MethodPost,
		Handler: router.Handler{
			Req:       new(presenter.CreateTransactionRequest),
			Res:       new(presenter.CreateTransactionResponse),
			Validator: th.CreateTransactionValidator,
			HandleFunc: func(ctx context.Context, req, res interface{}) error {
				return transactionHandler.CreateTransaction(ctx, req.(*presenter.CreateTransactionRequest), res.(*presenter.CreateTransactionResponse))
			},
		},
		Middlewares: []router.Middleware{authMiddleware},
	})

	// update transaction
	r.RegisterHttpRoute(&router.HttpRoute{
		Path:   config.PathUpdateTransaction,
		Method: http.MethodPost,
		Handler: router.Handler{
			Req:       new(presenter.UpdateTransactionRequest),
			Res:       new(presenter.UpdateTransactionResponse),
			Validator: th.UpdateTransactionValidator,
			HandleFunc: func(ctx context.Context, req, res interface{}) error {
				return transactionHandler.UpdateTransaction(ctx, req.(*presenter.UpdateTransactionRequest), res.(*presenter.UpdateTransactionResponse))
			},
		},
		Middlewares: []router.Middleware{authMiddleware},
	})

	// get transaction
	r.RegisterHttpRoute(&router.HttpRoute{
		Path:   config.PathGetTransaction,
		Method: http.MethodPost,
		Handler: router.Handler{
			Req:       new(presenter.GetTransactionRequest),
			Res:       new(presenter.GetTransactionResponse),
			Validator: th.GetTransactionValidator,
			HandleFunc: func(ctx context.Context, req, res interface{}) error {
				return transactionHandler.GetTransaction(ctx, req.(*presenter.GetTransactionRequest), res.(*presenter.GetTransactionResponse))
			},
		},
		Middlewares: []router.Middleware{authMiddleware},
	})

	// aggr transactions
	r.RegisterHttpRoute(&router.HttpRoute{
		Path:   config.PathAggrTransactions,
		Method: http.MethodPost,
		Handler: router.Handler{
			Req:       new(presenter.AggrTransactionsRequest),
			Res:       new(presenter.AggrTransactionsResponse),
			Validator: th.AggrTransactionsValidator,
			HandleFunc: func(ctx context.Context, req, res interface{}) error {
				return transactionHandler.AggrTransactions(ctx, req.(*presenter.AggrTransactionsRequest), res.(*presenter.AggrTransactionsResponse))
			},
		},
		Middlewares: []router.Middleware{authMiddleware},
	})

	// get transactions
	r.RegisterHttpRoute(&router.HttpRoute{
		Path:   config.PathGetTransactions,
		Method: http.MethodPost,
		Handler: router.Handler{
			Req:       new(presenter.GetTransactionsRequest),
			Res:       new(presenter.GetTransactionsResponse),
			Validator: th.GetTransactionsValidator,
			HandleFunc: func(ctx context.Context, req, res interface{}) error {
				return transactionHandler.GetTransactions(ctx, req.(*presenter.GetTransactionsRequest), res.(*presenter.GetTransactionsResponse))
			},
		},
		Middlewares: []router.Middleware{authMiddleware},
	})

	// ========== User ========== //

	userHandler := uh.NewUserHandler(s.userUseCase)

	// get user
	r.RegisterHttpRoute(&router.HttpRoute{
		Path:   config.PathGetUser,
		Method: http.MethodPost,
		Handler: router.Handler{
			Req:       new(presenter.GetUserRequest),
			Res:       new(presenter.GetUserResponse),
			Validator: uh.GetUserValidator,
			HandleFunc: func(ctx context.Context, req, res interface{}) error {
				return userHandler.GetUser(ctx, req.(*presenter.GetUserRequest), res.(*presenter.GetUserResponse))
			},
		},
		Middlewares: []router.Middleware{authMiddleware},
	})

	// sign up
	r.RegisterHttpRoute(&router.HttpRoute{
		Path:   config.PathSignUp,
		Method: http.MethodPost,
		Handler: router.Handler{
			Req:       new(presenter.SignUpRequest),
			Res:       new(presenter.SignUpResponse),
			Validator: uh.SignUpValidator,
			HandleFunc: func(ctx context.Context, req, res interface{}) error {
				return userHandler.SignUp(ctx, req.(*presenter.SignUpRequest), res.(*presenter.SignUpResponse))
			},
		},
	})

	// log in
	r.RegisterHttpRoute(&router.HttpRoute{
		Path:   config.PathLogin,
		Method: http.MethodPost,
		Handler: router.Handler{
			Req:       new(presenter.LogInRequest),
			Res:       new(presenter.LogInResponse),
			Validator: uh.LogInValidator,
			HandleFunc: func(ctx context.Context, req, res interface{}) error {
				return userHandler.LogIn(ctx, req.(*presenter.LogInRequest), res.(*presenter.LogInResponse))
			},
		},
	})

	// ========== Budget ========== //

	budgetHandler := bh.NewBudgetHandler(
		s.budgetUseCase,
		s.aggrUseCase,
	)

	// get budget
	r.RegisterHttpRoute(&router.HttpRoute{
		Path:   config.PathGetBudget,
		Method: http.MethodPost,
		Handler: router.Handler{
			Req:       new(presenter.GetBudgetRequest),
			Res:       new(presenter.GetBudgetResponse),
			Validator: bh.GetBudgetValidator,
			HandleFunc: func(ctx context.Context, req, res interface{}) error {
				return budgetHandler.GetBudget(ctx, req.(*presenter.GetBudgetRequest), res.(*presenter.GetBudgetResponse))
			},
		},
		Middlewares: []router.Middleware{authMiddleware},
	})

	// get budgets
	r.RegisterHttpRoute(&router.HttpRoute{
		Path:   config.PathGetBudgets,
		Method: http.MethodPost,
		Handler: router.Handler{
			Req:       new(presenter.GetBudgetsRequest),
			Res:       new(presenter.GetBudgetsResponse),
			Validator: bh.GetBudgetsValidator,
			HandleFunc: func(ctx context.Context, req, res interface{}) error {
				return budgetHandler.GetBudgets(ctx, req.(*presenter.GetBudgetsRequest), res.(*presenter.GetBudgetsResponse))
			},
		},
		Middlewares: []router.Middleware{authMiddleware},
	})

	// set budget
	r.RegisterHttpRoute(&router.HttpRoute{
		Path:   config.PathSetBudget,
		Method: http.MethodPost,
		Handler: router.Handler{
			Req:       new(presenter.SetBudgetRequest),
			Res:       new(presenter.SetBudgetResponse),
			Validator: bh.SetBudgetValidator,
			HandleFunc: func(ctx context.Context, req, res interface{}) error {
				return budgetHandler.SetBudget(ctx, req.(*presenter.SetBudgetRequest), res.(*presenter.SetBudgetResponse))
			},
		},
		Middlewares: []router.Middleware{authMiddleware},
	})

	return r
}
