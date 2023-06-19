package account

import "github.com/jseow5177/pockteer-be/usecase/account"

type accountHandler struct {
	accountUseCase account.UseCase
}

func NewAccountHandler(accountUseCase account.UseCase) *accountHandler {
	return &accountHandler{
		accountUseCase,
	}
}
