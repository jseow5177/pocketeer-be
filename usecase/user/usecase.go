package user

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jseow5177/pockteer-be/dep/mailer"
	"github.com/jseow5177/pockteer-be/dep/repo"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"github.com/jseow5177/pockteer-be/usecase/token"
	"github.com/rs/zerolog/log"
)

var (
	ErrEmailAlreadyExist = errors.New("email already exists")
	ErrUserInvalid       = errors.New("user invalid")
	ErrOTPInit           = errors.New("otp init fail")
	ErrOTPInvalid        = errors.New("otp invalid")
)

type userUseCase struct {
	txMgr        repo.TxMgr
	userRepo     repo.UserRepo
	otpRepo      repo.OTPRepo
	tokenUseCase token.UseCase
	mailer       mailer.Mailer
	categoryRepo repo.CategoryRepo
	budgetRepo   repo.BudgetRepo
	accountRepo  repo.AccountRepo
	securityRepo repo.SecurityRepo
	holdingRepo  repo.HoldingRepo
	lotRepo      repo.LotRepo
}

func NewUserUseCase(
	txMgr repo.TxMgr,
	userRepo repo.UserRepo,
	otpRepo repo.OTPRepo,
	tokenUseCase token.UseCase,
	mailer mailer.Mailer,
	categoryRepo repo.CategoryRepo,
	budgetRepo repo.BudgetRepo,
	accountRepo repo.AccountRepo,
	securityRepo repo.SecurityRepo,
	holdingRepo repo.HoldingRepo,
	lotRepo repo.LotRepo,
) UseCase {
	return &userUseCase{
		txMgr,
		userRepo,
		otpRepo,
		tokenUseCase,
		mailer,
		categoryRepo,
		budgetRepo,
		accountRepo,
		securityRepo,
		holdingRepo,
		lotRepo,
	}
}

func (uc *userUseCase) UpdateUserMeta(ctx context.Context, req *UpdateUserMetaRequest) (*UpdateUserMetaResponse, error) {
	uf := req.ToUserFilter()

	u, err := uc.userRepo.Get(ctx, uf)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to get user from repo, err: %v", err)
		return nil, err
	}

	uu, err := u.Update(req.ToUserUpdate())
	if err != nil {
		return nil, err
	}

	if uu != nil {
		if err := uc.userRepo.Update(ctx, uf, uu); err != nil {
			log.Ctx(ctx).Error().Msgf("fail to save user updates to repo, err: %v", err)
			return nil, err
		}
	}

	return &UpdateUserMetaResponse{
		User: u,
	}, nil
}

func (uc *userUseCase) InitUser(ctx context.Context, req *InitUserRequest) (*InitUserResponse, error) {
	uf := req.ToUserFilter()

	u, err := uc.userRepo.Get(ctx, uf)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to get user from repo, err: %v", err)
		return nil, err
	}

	// no-op
	if !u.IsNew() {
		log.Ctx(ctx).Info().Msgf("user already init, user_id: %v", u.GetUserID())
		return new(InitUserResponse), nil
	}

	if err := uc.txMgr.WithTx(ctx, func(txCtx context.Context) error {
		// add new categories
		if err := uc.initCategoriesAndBudgets(txCtx, req); err != nil {
			log.Ctx(txCtx).Error().Msgf("fail to init new categories, err: %v", err)
			return err
		}

		// add new accounts
		if err := uc.initAccounts(txCtx, req); err != nil {
			log.Ctx(txCtx).Error().Msgf("fail to init new accounts, err: %v", err)
			return err
		}

		// TODO: change user flag
		// uu, err := u.Update(req.ToUserUpdate())
		// if err != nil {
		// 	return err
		// }

		// if uu != nil {
		// 	if err := uc.userRepo.Update(txCtx, uf, uu); err != nil {
		// 		log.Ctx(txCtx).Error().Msgf("fail update user to flag default, err: %v", err)
		// 		return err
		// 	}
		// }

		return nil
	}); err != nil {
		return nil, err
	}

	return new(InitUserResponse), nil
}

func (uc *userUseCase) VerifyEmail(ctx context.Context, req *VerifyEmailRequest) (*VerifyEmailResponse, error) {
	otp, err := uc.otpRepo.Get(ctx, req.ToOTPFilter())
	if err != nil {
		return nil, err
	}

	if !otp.IsMatch(req.GetCode()) {
		return nil, ErrOTPInvalid
	}

	uf := req.ToUserFilter(req.GetEmail())

	// Check if user exists
	u, err := uc.userRepo.Get(ctx, uf)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to get user from repo, err: %v", err)
		return nil, err
	}

	// Update user to status normal
	uu, err := u.Update(req.ToUserUpdate())
	if err != nil {
		return nil, err
	}

	if uu != nil {
		if err = uc.userRepo.Update(ctx, uf, uu); err != nil {
			log.Ctx(ctx).Error().Msgf("fail to save user updates to repo, err: %v", err)
			return nil, err
		}
	}

	// create access token
	tokenRes, err := uc.tokenUseCase.CreateToken(ctx, &token.CreateTokenRequest{
		TokenType: goutil.Uint32(uint32(entity.TokenTypeAccess)),
		CustomClaims: &entity.CustomClaims{
			UserID: u.UserID,
		},
	})
	if err != nil {
		return nil, err
	}

	// async retry send email, no cancel
	async := goutil.NewAsync(time.Second, 5)
	async.Retry(ctx, func(ctx context.Context) error {
		ctx = goutil.WithoutCancel(ctx)
		if err := uc.mailer.SendEmail(ctx, mailer.TemplateWelcome, &mailer.SendEmailRequest{
			To: u.GetEmail(),
		}); err != nil {
			log.Ctx(ctx).Error().Msgf("fail to send welcome email, user_id: %v, err: %v", u.GetUserID(), err)
			return err
		}
		return nil
	})

	return &VerifyEmailResponse{
		AccessToken: tokenRes.Token,
		User:        u,
	}, nil
}

func (uc *userUseCase) IsAuthenticated(ctx context.Context, req *IsAuthenticatedRequest) (*IsAuthenticatedResponse, error) {
	res, err := uc.tokenUseCase.ValidateToken(ctx, req.ToValidateTokenRequest())
	if err != nil {
		return nil, err
	}

	userID := res.CustomClaims.GetUserID()

	// check if user exists
	u, err := uc.userRepo.Get(ctx, req.ToUserFilter(userID))
	if err != nil {
		return nil, err
	}

	return &IsAuthenticatedResponse{
		User: u,
	}, nil
}

func (uc *userUseCase) GetUser(ctx context.Context, req *GetUserRequest) (*GetUserResponse, error) {
	u, err := uc.userRepo.Get(ctx, req.ToUserFilter())
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to get user from repo, err: %v", err)
		return nil, err
	}

	return &GetUserResponse{
		User: u,
	}, nil
}

func (uc *userUseCase) SendOTP(ctx context.Context, req *SendOTPRequest) (*SendOTPResponse, error) {
	u, err := uc.userRepo.Get(ctx, req.ToUserFilter())
	if err != nil {
		return nil, err
	}

	async := goutil.NewAsync(time.Second, 5)
	async.Retry(ctx, func(ctx context.Context) error {
		ctx = goutil.WithoutCancel(ctx)
		if err := uc.sendOTP(ctx, u.GetEmail()); err != nil {
			log.Ctx(ctx).Error().Msgf("fail to send otp, email: %v, err: %v", req.GetEmail(), err)
			return err
		}
		return nil
	})

	return new(SendOTPResponse), nil
}

func (uc *userUseCase) SignUp(ctx context.Context, req *SignUpRequest) (*SignUpResponse, error) {
	u, err := uc.userRepo.Get(ctx, req.ToUserFilter())
	if err != nil && err != repo.ErrUserNotFound {
		return nil, err
	}

	if u != nil && u.IsNormal() {
		return nil, ErrEmailAlreadyExist
	}

	if u == nil {
		// create a new user
		u, err = req.ToUserEntity()
		if err != nil {
			return nil, err
		}

		_, err = uc.userRepo.Create(ctx, u)
		if err != nil {
			log.Ctx(ctx).Error().Msgf("fail to save new user to repo, err: %v", err)
			return nil, err
		}
	} else {
		// check if user signed up with a different password
		uu, err := u.Update(req.ToUserUpdate())
		if err != nil {
			return nil, err
		}

		if uu != nil {
			if err = uc.userRepo.Update(ctx, req.ToUserFilter(), uu); err != nil {
				log.Ctx(ctx).Error().Msgf("fail to save user updates to repo, err: %v", err)
				return nil, err
			}
		}
	}

	async := goutil.NewAsync(time.Second, 5)
	async.Retry(ctx, func(ctx context.Context) error {
		ctx = goutil.WithoutCancel(ctx)
		if err := uc.sendOTP(ctx, req.GetEmail()); err != nil {
			log.Ctx(ctx).Error().Msgf("fail to send otp, email: %v, err: %v", req.GetEmail(), err)
			return err
		}
		return nil
	})

	return &SignUpResponse{
		User: u,
	}, nil
}

func (uc *userUseCase) sendOTP(ctx context.Context, email string) error {
	f := &repo.OTPFilter{
		Email: goutil.String(email),
	}

	// check if there is an existing otp
	otp, err := uc.otpRepo.Get(ctx, f)
	if err != nil && err != repo.ErrOTPNotFound {
		return err
	}

	if otp == nil {
		// create new otp
		otp, err = entity.NewOTP()
		if err != nil {
			return err
		}
		uc.otpRepo.Set(ctx, email, otp)
	}

	return uc.mailer.SendEmail(ctx, mailer.TemplateOTP, &mailer.SendEmailRequest{
		To: email,
		Params: map[string]interface{}{
			"otp": otp.GetCode(),
		},
	})
}

func (uc *userUseCase) LogIn(ctx context.Context, req *LogInRequest) (*LogInResponse, error) {
	u, err := uc.userRepo.Get(ctx, req.ToUserFilter())
	if err != nil {
		if err == repo.ErrUserNotFound {
			// hide not found error
			return nil, ErrUserInvalid
		}
		return nil, err
	}

	isPasswordCorrect, err := u.IsSamePassword(req.GetPassword())
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to check if password is correct, err: %v", err)
		return nil, err
	}

	if !isPasswordCorrect {
		return nil, ErrUserInvalid
	}

	// create access token
	res, err := uc.tokenUseCase.CreateToken(ctx, &token.CreateTokenRequest{
		TokenType: goutil.Uint32(uint32(entity.TokenTypeAccess)),
		CustomClaims: &entity.CustomClaims{
			UserID: u.UserID,
		},
	})
	if err != nil {
		return nil, err
	}

	return &LogInResponse{
		AccessToken: res.Token,
		User:        u,
	}, nil
}

func (uc *userUseCase) initCategoriesAndBudgets(ctx context.Context, req *InitUserRequest) error {
	if len(req.Categories) == 0 {
		return nil
	}

	cs, err := req.ToCategoryEntities()
	if err != nil {
		return err
	}

	if _, err := uc.categoryRepo.CreateMany(ctx, cs); err != nil {
		return err
	}

	bs := make([]*entity.Budget, 0)
	for _, c := range cs {
		if c.Budget == nil {
			continue
		}
		c.Budget.SetCategoryID(c.CategoryID)
		bs = append(bs, c.Budget)
	}

	if len(bs) == 0 {
		return nil
	}

	if _, err := uc.budgetRepo.CreateMany(ctx, bs); err != nil {
		return err
	}

	return nil
}

func (uc *userUseCase) initAccounts(ctx context.Context, req *InitUserRequest) error {
	if len(req.Accounts) == 0 {
		return nil
	}

	acs, err := req.ToAccountEntities()
	if err != nil {
		return err
	}

	if _, err := uc.accountRepo.CreateMany(ctx, acs); err != nil {
		return err
	}

	for i, ac := range acs {
		var (
			hrs = req.Accounts[i].Holdings
			hs  = ac.Holdings
		)

		if len(hs) == 0 {
			continue
		}

		for j, h := range hs {
			if h.IsDefault() {
				if _, err = uc.securityRepo.Get(ctx, hrs[j].ToSecurityFilter()); err != nil {
					return fmt.Errorf("symbol %v, err: %v", h.GetSymbol(), err)
				}
			}
			h.SetAccountID(ac.AccountID)
		}

		_, err = uc.holdingRepo.CreateMany(ctx, hs)
		if err != nil {
			return err
		}

		for _, h := range hs {
			ls := h.Lots

			if len(ls) == 0 {
				continue
			}

			for _, l := range ls {
				l.SetHoldingID(h.HoldingID)
			}

			_, err := uc.lotRepo.CreateMany(ctx, ls)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
