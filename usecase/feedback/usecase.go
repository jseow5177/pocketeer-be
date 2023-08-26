package feedback

import (
	"context"

	"github.com/jseow5177/pockteer-be/dep/repo"
)

type feedbackUseCase struct {
	feedbackRepo repo.FeedbackRepo
}

func NewFeedbackUseCase(feedbackRepo repo.FeedbackRepo) UseCase {
	return &feedbackUseCase{
		feedbackRepo: feedbackRepo,
	}
}

func (uc *feedbackUseCase) CreateFeedback(ctx context.Context, req *CreateFeedbackRequest) (*CreateFeedbackResponse, error) {
	f, err := req.ToFeedbackEntity()
	if err != nil {
		return nil, err
	}

	if err := uc.feedbackRepo.Create(ctx, f); err != nil {
		return nil, err
	}

	return &CreateFeedbackResponse{
		Feedback: f,
	}, nil
}
