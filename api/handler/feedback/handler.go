package feedback

import "github.com/jseow5177/pockteer-be/usecase/feedback"

type feedbackHandler struct {
	feedbackUseCase feedback.UseCase
}

func NewFeedbackHandler(feedbackUseCase feedback.UseCase) *feedbackHandler {
	return &feedbackHandler{
		feedbackUseCase,
	}
}
