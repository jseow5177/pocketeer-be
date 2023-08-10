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
	"github.com/jseow5177/pockteer-be/dep/api"
	"github.com/jseow5177/pockteer-be/dep/api/finnhub"
	"github.com/jseow5177/pockteer-be/dep/mailer"
	"github.com/jseow5177/pockteer-be/dep/mailer/brevo"
	"github.com/jseow5177/pockteer-be/dep/repo"
	"github.com/jseow5177/pockteer-be/dep/repo/mem"
	"github.com/jseow5177/pockteer-be/dep/repo/mongo"
	"github.com/jseow5177/pockteer-be/pkg/logger"
	"github.com/jseow5177/pockteer-be/pkg/router"
	"github.com/jseow5177/pockteer-be/pkg/service"

	ach "github.com/jseow5177/pockteer-be/api/handler/account"
	bh "github.com/jseow5177/pockteer-be/api/handler/budget"
	ch "github.com/jseow5177/pockteer-be/api/handler/category"
	hh "github.com/jseow5177/pockteer-be/api/handler/holding"
	lh "github.com/jseow5177/pockteer-be/api/handler/lot"
	sh "github.com/jseow5177/pockteer-be/api/handler/security"
	th "github.com/jseow5177/pockteer-be/api/handler/transaction"
	uh "github.com/jseow5177/pockteer-be/api/handler/user"

	acuc "github.com/jseow5177/pockteer-be/usecase/account"
	buc "github.com/jseow5177/pockteer-be/usecase/budget"
	cuc "github.com/jseow5177/pockteer-be/usecase/category"
	huc "github.com/jseow5177/pockteer-be/usecase/holding"
	luc "github.com/jseow5177/pockteer-be/usecase/lot"
	suc "github.com/jseow5177/pockteer-be/usecase/security"
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
	accountRepo     repo.AccountRepo
	holdingRepo     repo.HoldingRepo
	lotRepo         repo.LotRepo
	securityRepo    repo.SecurityRepo
	quoteRepo       repo.QuoteRepo

	securityAPI api.SecurityAPI

	mailer mailer.Mailer

	categoryUseCase    cuc.UseCase
	transactionUseCase tuc.UseCase
	budgetUseCase      buc.UseCase
	userUseCase        uuc.UseCase
	tokenUseCase       ttuc.UseCase
	accountUseCase     acuc.UseCase
	securityUseCase    suc.UseCase
	holdingUseCase     huc.UseCase
	lotUseCase         luc.UseCase
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

	// init mongo repos
	s.categoryRepo = mongo.NewCategoryMongo(s.mongo)
	s.transactionRepo = mongo.NewTransactionMongo(s.mongo)
	s.budgetRepo = mongo.NewBudgetMongo(s.mongo)
	s.userRepo = mongo.NewUserMongo(s.mongo)
	s.accountRepo = mongo.NewAccountMongo(s.mongo)
	s.holdingRepo = mongo.NewHoldingMongo(s.mongo)
	s.lotRepo = mongo.NewLotMongo(s.mongo)
	s.securityRepo = mongo.NewSecurityMongo(s.mongo)

	// init apis
	s.securityAPI = finnhub.NewFinnHubMgr(s.cfg.FinnHub)

	// init mailer
	s.mailer, err = brevo.NewBrevoMgr(s.cfg.Brevo)
	if err != nil {
		log.Ctx(s.ctx).Error().Msgf("fail to init brevo mailer, err: %v", err)
		return err
	}

	// init mem repos
	s.quoteRepo, err = mem.NewQuoteMemCache(s.cfg.QuoteMemCache, s.securityAPI)
	if err != nil {
		log.Ctx(s.ctx).Error().Msgf("fail to init quote repo, err: %v", err)
		return err
	}

	// init use cases
	s.transactionUseCase = tuc.NewTransactionUseCase(s.mongo, s.categoryRepo, s.accountRepo, s.transactionRepo, s.budgetRepo)
	s.budgetUseCase = buc.NewBudgetUseCase(s.mongo, s.budgetRepo, s.categoryRepo, s.transactionRepo)
	s.categoryUseCase = cuc.NewCategoryUseCase(s.categoryRepo, s.transactionRepo, s.budgetUseCase)
	s.tokenUseCase = ttuc.NewTokenUseCase(s.cfg.Tokens)
	s.userUseCase = uuc.NewUserUseCase(s.mongo, s.userRepo, s.tokenUseCase, s.mailer)
	s.securityUseCase = suc.NewSecurityUseCase(s.securityRepo)
	s.holdingUseCase = huc.NewHoldingUseCase(s.accountRepo, s.holdingRepo, s.lotRepo, s.securityRepo, s.quoteRepo)
	s.lotUseCase = luc.NewLotUseCase(s.lotRepo, s.holdingRepo)
	s.accountUseCase = acuc.NewAccountUseCase(s.mongo, s.accountRepo, s.transactionRepo, s.holdingUseCase)

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

	authMiddleware := middleware.NewAuthMiddleware(s.userUseCase)

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

	// get category budget
	r.RegisterHttpRoute(&router.HttpRoute{
		Path:   config.PathGetCategoryBudget,
		Method: http.MethodPost,
		Handler: router.Handler{
			Req:       new(presenter.GetCategoryBudgetRequest),
			Res:       new(presenter.GetCategoryBudgetResponse),
			Validator: ch.GetCategoryBudgetValidator,
			HandleFunc: func(ctx context.Context, req, res interface{}) error {
				return categoryHandler.GetCategoryBudget(ctx, req.(*presenter.GetCategoryBudgetRequest), res.(*presenter.GetCategoryBudgetResponse))
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

	// get categories budget
	r.RegisterHttpRoute(&router.HttpRoute{
		Path:   config.PathGetCategoriesBudget,
		Method: http.MethodPost,
		Handler: router.Handler{
			Req:       new(presenter.GetCategoriesBudgetRequest),
			Res:       new(presenter.GetCategoriesBudgetResponse),
			Validator: ch.GetCategoriesBudgetValidator,
			HandleFunc: func(ctx context.Context, req, res interface{}) error {
				return categoryHandler.GetCategoriesBudget(ctx, req.(*presenter.GetCategoriesBudgetRequest), res.(*presenter.GetCategoriesBudgetResponse))
			},
		},
		Middlewares: []router.Middleware{authMiddleware},
	})

	// ========== Account ========== //

	accountHandler := ach.NewAccountHandler(s.accountUseCase)

	// create account
	r.RegisterHttpRoute(&router.HttpRoute{
		Path:   config.PathCreateAccount,
		Method: http.MethodPost,
		Handler: router.Handler{
			Req:       new(presenter.CreateAccountRequest),
			Res:       new(presenter.CreateAccountResponse),
			Validator: ach.CreateAccountValidator,
			HandleFunc: func(ctx context.Context, req, res interface{}) error {
				return accountHandler.CreateAccount(ctx, req.(*presenter.CreateAccountRequest), res.(*presenter.CreateAccountResponse))
			},
		},
		Middlewares: []router.Middleware{authMiddleware},
	})

	// update account
	r.RegisterHttpRoute(&router.HttpRoute{
		Path:   config.PathUpdateAccount,
		Method: http.MethodPost,
		Handler: router.Handler{
			Req:       new(presenter.UpdateAccountRequest),
			Res:       new(presenter.UpdateAccountResponse),
			Validator: ach.UpdateAccountValidator,
			HandleFunc: func(ctx context.Context, req, res interface{}) error {
				return accountHandler.UpdateAccount(ctx, req.(*presenter.UpdateAccountRequest), res.(*presenter.UpdateAccountResponse))
			},
		},
		Middlewares: []router.Middleware{authMiddleware},
	})

	// get account
	r.RegisterHttpRoute(&router.HttpRoute{
		Path:   config.PathGetAccount,
		Method: http.MethodPost,
		Handler: router.Handler{
			Req:       new(presenter.GetAccountRequest),
			Res:       new(presenter.GetAccountResponse),
			Validator: ach.GetAccountValidator,
			HandleFunc: func(ctx context.Context, req, res interface{}) error {
				return accountHandler.GetAccount(ctx, req.(*presenter.GetAccountRequest), res.(*presenter.GetAccountResponse))
			},
		},
		Middlewares: []router.Middleware{authMiddleware},
	})

	// gets accounts
	r.RegisterHttpRoute(&router.HttpRoute{
		Path:   config.PathGetAccounts,
		Method: http.MethodPost,
		Handler: router.Handler{
			Req:       new(presenter.GetAccountsRequest),
			Res:       new(presenter.GetAccountsResponse),
			Validator: ach.GetAccountsValidator,
			HandleFunc: func(ctx context.Context, req, res interface{}) error {
				return accountHandler.GetAccounts(ctx, req.(*presenter.GetAccountsRequest), res.(*presenter.GetAccountsResponse))
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

	// delete transactions
	r.RegisterHttpRoute(&router.HttpRoute{
		Path:   config.PathDeleteTransaction,
		Method: http.MethodPost,
		Handler: router.Handler{
			Req:       new(presenter.DeleteTransactionRequest),
			Res:       new(presenter.DeleteTransactionResponse),
			Validator: th.DeleteTransactionValidator,
			HandleFunc: func(ctx context.Context, req, res interface{}) error {
				return transactionHandler.DeleteTransaction(ctx, req.(*presenter.DeleteTransactionRequest), res.(*presenter.DeleteTransactionResponse))
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

	// verify email
	r.RegisterHttpRoute(&router.HttpRoute{
		Path:   config.PathVerifyEmail,
		Method: http.MethodPost,
		Handler: router.Handler{
			Req:       new(presenter.VerifyEmailRequest),
			Res:       new(presenter.VerifyEmailResponse),
			Validator: uh.VerifyEmailValidator,
			HandleFunc: func(ctx context.Context, req, res interface{}) error {
				return userHandler.VerifyEmail(ctx, req.(*presenter.VerifyEmailRequest), res.(*presenter.VerifyEmailResponse))
			},
		},
	})

	// ========== Budget ========== //

	budgetHandler := bh.NewBudgetHandler(s.budgetUseCase)

	// create budget
	r.RegisterHttpRoute(&router.HttpRoute{
		Path:   config.PathCreateBudget,
		Method: http.MethodPost,
		Handler: router.Handler{
			Req:       new(presenter.CreateBudgetRequest),
			Res:       new(presenter.CreateBudgetResponse),
			Validator: bh.CreateBudgetValidator,
			HandleFunc: func(ctx context.Context, req, res interface{}) error {
				return budgetHandler.CreateBudget(ctx, req.(*presenter.CreateBudgetRequest), res.(*presenter.CreateBudgetResponse))
			},
		},
		Middlewares: []router.Middleware{authMiddleware},
	})

	// update budget
	r.RegisterHttpRoute(&router.HttpRoute{
		Path:   config.PathUpdateBudget,
		Method: http.MethodPost,
		Handler: router.Handler{
			Req:       new(presenter.UpdateBudgetRequest),
			Res:       new(presenter.UpdateBudgetResponse),
			Validator: bh.UpdateBudgetValidator,
			HandleFunc: func(ctx context.Context, req, res interface{}) error {
				return budgetHandler.UpdateBudget(ctx, req.(*presenter.UpdateBudgetRequest), res.(*presenter.UpdateBudgetResponse))
			},
		},
		Middlewares: []router.Middleware{authMiddleware},
	})

	// delete budget
	r.RegisterHttpRoute(&router.HttpRoute{
		Path:   config.PathDeleteBudget,
		Method: http.MethodPost,
		Handler: router.Handler{
			Req:       new(presenter.DeleteBudgetRequest),
			Res:       new(presenter.DeleteBudgetResponse),
			Validator: bh.DeleteBudgetValidator,
			HandleFunc: func(ctx context.Context, req, res interface{}) error {
				return budgetHandler.DeleteBudget(ctx, req.(*presenter.DeleteBudgetRequest), res.(*presenter.DeleteBudgetResponse))
			},
		},
		Middlewares: []router.Middleware{authMiddleware},
	})

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

	// update budget
	r.RegisterHttpRoute(&router.HttpRoute{
		Path:   config.PathUpdateBudget,
		Method: http.MethodPost,
		Handler: router.Handler{
			Req:       new(presenter.UpdateBudgetRequest),
			Res:       new(presenter.UpdateBudgetResponse),
			Validator: bh.UpdateBudgetValidator,
			HandleFunc: func(ctx context.Context, req, res interface{}) error {
				return budgetHandler.UpdateBudget(ctx, req.(*presenter.UpdateBudgetRequest), res.(*presenter.UpdateBudgetResponse))
			},
		},
		Middlewares: []router.Middleware{authMiddleware},
	})

	// ========== Security ========== //

	securityHandler := sh.NewSecurityHandler(s.securityUseCase)

	// search securities
	r.RegisterHttpRoute(&router.HttpRoute{
		Path:   config.PathSearchSecurities,
		Method: http.MethodPost,
		Handler: router.Handler{
			Req:       new(presenter.SearchSecuritiesRequest),
			Res:       new(presenter.SearchSecuritiesResponse),
			Validator: sh.SearchSecuritiesValidator,
			HandleFunc: func(ctx context.Context, req, res interface{}) error {
				return securityHandler.SearchSecurities(ctx, req.(*presenter.SearchSecuritiesRequest), res.(*presenter.SearchSecuritiesResponse))
			},
		},
		Middlewares: []router.Middleware{authMiddleware},
	})

	// ========== Holding ========== //

	holdingHandler := hh.NewHoldingHandler(s.holdingUseCase)

	// create holding
	r.RegisterHttpRoute(&router.HttpRoute{
		Path:   config.PathCreateHolding,
		Method: http.MethodPost,
		Handler: router.Handler{
			Req:       new(presenter.CreateHoldingRequest),
			Res:       new(presenter.CreateHoldingResponse),
			Validator: hh.CreateHoldingValidator,
			HandleFunc: func(ctx context.Context, req, res interface{}) error {
				return holdingHandler.CreateHolding(ctx, req.(*presenter.CreateHoldingRequest), res.(*presenter.CreateHoldingResponse))
			},
		},
		Middlewares: []router.Middleware{authMiddleware},
	})

	// update holding
	r.RegisterHttpRoute(&router.HttpRoute{
		Path:   config.PathUpdateHolding,
		Method: http.MethodPost,
		Handler: router.Handler{
			Req:       new(presenter.UpdateHoldingRequest),
			Res:       new(presenter.UpdateHoldingResponse),
			Validator: hh.UpdateHoldingValidator,
			HandleFunc: func(ctx context.Context, req, res interface{}) error {
				return holdingHandler.UpdateHolding(ctx, req.(*presenter.UpdateHoldingRequest), res.(*presenter.UpdateHoldingResponse))
			},
		},
		Middlewares: []router.Middleware{authMiddleware},
	})

	// get holding
	r.RegisterHttpRoute(&router.HttpRoute{
		Path:   config.PathGetHolding,
		Method: http.MethodPost,
		Handler: router.Handler{
			Req:       new(presenter.GetHoldingRequest),
			Res:       new(presenter.GetHoldingResponse),
			Validator: hh.GetHoldingValidator,
			HandleFunc: func(ctx context.Context, req, res interface{}) error {
				return holdingHandler.GetHolding(ctx, req.(*presenter.GetHoldingRequest), res.(*presenter.GetHoldingResponse))
			},
		},
		Middlewares: []router.Middleware{authMiddleware},
	})

	// ========== Lot ========== //

	lotHandler := lh.NewLotHandler(s.lotUseCase)

	// create lot
	r.RegisterHttpRoute(&router.HttpRoute{
		Path:   config.PathCreateLot,
		Method: http.MethodPost,
		Handler: router.Handler{
			Req:       new(presenter.CreateLotRequest),
			Res:       new(presenter.CreateLotResponse),
			Validator: lh.CreateLotValidator,
			HandleFunc: func(ctx context.Context, req, res interface{}) error {
				return lotHandler.CreateLot(ctx, req.(*presenter.CreateLotRequest), res.(*presenter.CreateLotResponse))
			},
		},
		Middlewares: []router.Middleware{authMiddleware},
	})

	// delete lot
	r.RegisterHttpRoute(&router.HttpRoute{
		Path:   config.PathDeleteLot,
		Method: http.MethodPost,
		Handler: router.Handler{
			Req:       new(presenter.DeleteLotRequest),
			Res:       new(presenter.DeleteLotResponse),
			Validator: lh.DeleteLotValidator,
			HandleFunc: func(ctx context.Context, req, res interface{}) error {
				return lotHandler.DeleteLot(ctx, req.(*presenter.DeleteLotRequest), res.(*presenter.DeleteLotResponse))
			},
		},
		Middlewares: []router.Middleware{authMiddleware},
	})

	// update lot
	r.RegisterHttpRoute(&router.HttpRoute{
		Path:   config.PathUpdateLot,
		Method: http.MethodPost,
		Handler: router.Handler{
			Req:       new(presenter.UpdateLotRequest),
			Res:       new(presenter.UpdateLotResponse),
			Validator: lh.UpdateLotValidator,
			HandleFunc: func(ctx context.Context, req, res interface{}) error {
				return lotHandler.UpdateLot(ctx, req.(*presenter.UpdateLotRequest), res.(*presenter.UpdateLotResponse))
			},
		},
		Middlewares: []router.Middleware{authMiddleware},
	})

	// get lot
	r.RegisterHttpRoute(&router.HttpRoute{
		Path:   config.PathGetLot,
		Method: http.MethodPost,
		Handler: router.Handler{
			Req:       new(presenter.GetLotRequest),
			Res:       new(presenter.GetLotResponse),
			Validator: lh.GetLotValidator,
			HandleFunc: func(ctx context.Context, req, res interface{}) error {
				return lotHandler.GetLot(ctx, req.(*presenter.GetLotRequest), res.(*presenter.GetLotResponse))
			},
		},
		Middlewares: []router.Middleware{authMiddleware},
	})

	// get lots
	r.RegisterHttpRoute(&router.HttpRoute{
		Path:   config.PathGetLots,
		Method: http.MethodPost,
		Handler: router.Handler{
			Req:       new(presenter.GetLotsRequest),
			Res:       new(presenter.GetLotsResponse),
			Validator: lh.GetLotsValidator,
			HandleFunc: func(ctx context.Context, req, res interface{}) error {
				return lotHandler.GetLots(ctx, req.(*presenter.GetLotsRequest), res.(*presenter.GetLotsResponse))
			},
		},
		Middlewares: []router.Middleware{authMiddleware},
	})

	return r
}
