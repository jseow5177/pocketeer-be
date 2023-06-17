package user

import (
	"context"
	"errors"

	"github.com/jseow5177/pockteer-be/dep/repo"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"github.com/jseow5177/pockteer-be/usecase/token"
	"github.com/rs/zerolog/log"
)

var (
	ErrUsernameAlreadyExist = errors.New("username already exists")
	ErrUserInvalid          = errors.New("user invalid")
)

type userUseCase struct {
	userRepo     repo.UserRepo
	tokenUseCase token.UseCase
}

func NewUserUseCase(userRepo repo.UserRepo, tokenUseCase token.UseCase) UseCase {
	return &userUseCase{
		userRepo,
		tokenUseCase,
	}
}

func (uc *userUseCase) IsAuthenticated(ctx context.Context, req *IsAuthenticatedRequest) (*IsAuthenticatedResponse, error) {
	validateTokenRes, err := uc.tokenUseCase.ValidateToken(ctx, req.ToValidateTokenRequest())
	if err != nil {
		return nil, err
	}

	userID := validateTokenRes.CustomClaims.GetUserID()

	// check if user exists
	getUserRes, err := uc.GetUser(ctx, &GetUserRequest{
		UserID:     goutil.String(userID),
		UserStatus: goutil.Uint32(uint32(entity.UserStatusNormal)),
	})
	if err != nil {
		return nil, err
	}

	return &IsAuthenticatedResponse{
		UserID: getUserRes.UserID,
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
	_, err := uc.GetUser(ctx, req.ToGetUserRequest())
	if err != nil && err != repo.ErrUserNotFound {
		return nil, err
	}

	if err == nil {
		return nil, ErrUsernameAlreadyExist
	}

	u := req.ToUserEntity()
	if err = u.SetHash(req.GetPassword()); err != nil {
		log.Ctx(ctx).Error().Msgf("fail to set password hash, err: %v", err)
		return nil, err
	}

	userID, err := uc.userRepo.Create(ctx, u)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to save new user to repo, err: %v", err)
		return nil, err
	}

	u.UserID = goutil.String(userID)

	return &SignUpResponse{
		User: u,
	}, nil
}

func (uc *userUseCase) LogIn(ctx context.Context, req *LogInRequest) (*LogInResponse, error) {
	getUserRes, err := uc.GetUser(ctx, req.ToGetUserRequest())
	if err != nil {
		if err == repo.ErrUserNotFound {
			// hide not found error
			return nil, ErrUserInvalid
		}
		return nil, err
	}

	isPasswordCorrect, err := getUserRes.IsPasswordCorrect(req.GetPassword())
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to check if password is correct, err: %v", err)
		return nil, err
	}

	if !isPasswordCorrect {
		return nil, ErrUserInvalid
	}

	// create access token
	accessTokenRes, err := uc.tokenUseCase.CreateToken(ctx, &token.CreateTokenRequest{
		TokenType: goutil.Uint32(uint32(entity.TokenTypeAccess)),
		CustomClaims: &entity.CustomClaims{
			UserID: getUserRes.UserID,
		},
	})
	if err != nil {
		return nil, err
	}

	return &LogInResponse{
		AccessToken: accessTokenRes.Token,
	}, nil
}
