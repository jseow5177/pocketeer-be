package lot

import "github.com/jseow5177/pockteer-be/usecase/lot"

type lotHandler struct {
	lotUseCase lot.UseCase
}

func NewLotHandler(lotUseCase lot.UseCase) *lotHandler {
	return &lotHandler{
		lotUseCase,
	}
}
