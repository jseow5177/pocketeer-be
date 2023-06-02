package user

import (
	"context"

	"github.com/jseow5177/pockteer-be/dep/repo"
)

type userUseCase struct {
	userRepo repo.UserRepo
}

func NewUserUseCase(userRepo repo.UserRepo) UseCase {
	return &userUseCase{
		userRepo: userRepo,
	}
}

func (uc *userUseCase) GetUser(ctx context.Context, req *GetUserRequest) (*GetUserResponse, error) {
	return nil, nil
}

func (uc *userUseCase) SignUp(ctx context.Context, req *SignUpRequest) (*SignUpResponse, error) {
	return nil, nil
}

func (uc *userUseCase) LogIn(ctx context.Context, req *LogInRequest) (*LogInResponse, error) {
	return nil, nil
}
