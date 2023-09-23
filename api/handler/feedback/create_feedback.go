package feedback

import (
	"context"

	"github.com/jseow5177/pockteer-be/api/presenter"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"github.com/jseow5177/pockteer-be/pkg/validator"
	"github.com/rs/zerolog/log"
)

var CreateFeedbackValidator = validator.MustForm(map[string]validator.Validator{
	"score": &validator.UInt32{
		Optional: true,
		Min:      goutil.Uint32(uint32(entity.FeedbackScoreZero)),
		Max:      goutil.Uint32(uint32(entity.FeedbackScoreTwo)),
	},
	"text": &validator.String{
		Optional: true,
		MaxLen:   entity.MaxFeedbackLength,
	},
})

func (h *feedbackHandler) CreateFeedback(ctx context.Context, req *presenter.CreateFeedbackRequest, res *presenter.CreateFeedbackResponse) error {
	user := entity.GetUserFromCtx(ctx)

	useCaseRes, err := h.feedbackUseCase.CreateFeedback(ctx, req.ToUseCaseReq(user.GetUserID()))
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to create feedback, err: %v", err)
		return err
	}

	res.Set(useCaseRes)

	return nil
}
