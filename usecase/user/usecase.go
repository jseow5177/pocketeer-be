package user

import (
	"context"
	"errors"
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
}

func NewUserUseCase(
	txMgr repo.TxMgr,
	userRepo repo.UserRepo,
	otpRepo repo.OTPRepo,
	tokenUseCase token.UseCase,
	mailer mailer.Mailer,
) UseCase {
	return &userUseCase{
		txMgr,
		userRepo,
		otpRepo,
		tokenUseCase,
		mailer,
	}
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

	uu, _ := u.Update(entity.NewUserUpdate(
		entity.WithUpdateUserFlag(goutil.Uint32(uint32(entity.UserFlagDefault))),
	))

	if err := uc.userRepo.Update(ctx, uf, uu); err != nil {
		log.Ctx(ctx).Error().Msgf("fail update user to flag default, err: %v", err)
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

	uu, _ := u.Update(req.ToUserUpdate())

	// Update user to status normal
	if err = uc.userRepo.Update(ctx, uf, uu); err != nil {
		log.Ctx(ctx).Error().Msgf("fail to save user updates to repo, err: %v", err)
		return nil, err
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

func (uc *userUseCase) SignUp(ctx context.Context, req *SignUpRequest) (*SignUpResponse, error) {
	u, err := uc.userRepo.Get(ctx, req.ToUserFilter())
	if err != nil && err != repo.ErrUserNotFound {
		return nil, err
	}

	if u != nil && u.IsNormal() {
		return nil, ErrEmailAlreadyExist
	}

	if u == nil {
		u, err = req.ToUserEntity()
		if err != nil {
			return nil, err
		}

		_, err = uc.userRepo.Create(ctx, u)
		if err != nil {
			log.Ctx(ctx).Error().Msgf("fail to save new user to repo, err: %v", err)
			return nil, err
		}
	}

	// create email otp
	otp, err := entity.NewOTP()
	if err != nil {
		return nil, err
	}
	uc.otpRepo.Set(ctx, req.GetEmail(), otp)

	// async retry send email, no cancel
	async := goutil.NewAsync(time.Second, 5)
	async.Retry(ctx, func(ctx context.Context) error {
		ctx = goutil.WithoutCancel(ctx)
		if err := uc.mailer.SendEmail(ctx, mailer.TemplateVerifyEmail, &mailer.SendEmailRequest{
			To: req.GetEmail(),
			Params: map[string]interface{}{
				"otp": otp.GetCode(),
			},
		}); err != nil {
			log.Ctx(ctx).Error().Msgf("fail to send verification email, email: %v, err: %v", req.GetEmail(), err)
			return err
		}
		return nil
	})

	return &SignUpResponse{
		User: u,
	}, nil
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

	isPasswordCorrect, err := u.IsPasswordCorrect(req.GetPassword())
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
