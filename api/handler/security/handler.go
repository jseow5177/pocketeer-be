package security

import "github.com/jseow5177/pockteer-be/usecase/security"

type securityHandler struct {
	securityUseCase security.UseCase
}

func NewSecurityHandler(securityUseCase security.UseCase) *securityHandler {
	return &securityHandler{
		securityUseCase,
	}
}
