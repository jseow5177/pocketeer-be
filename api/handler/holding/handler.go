package holding

import "github.com/jseow5177/pockteer-be/usecase/holding"

type holdingHandler struct {
	holdingUseCase holding.UseCase
}

func NewHoldingHandler(holdingUseCase holding.UseCase) *holdingHandler {
	return &holdingHandler{
		holdingUseCase,
	}
}
