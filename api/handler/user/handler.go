package user

import "github.com/jseow5177/pockteer-be/usecase/user"

type userHandler struct {
	userUseCase user.UseCase
}

func NewUserHandler(userUseCase user.UseCase) *userHandler {
	return &userHandler{
		userUseCase,
	}
}
