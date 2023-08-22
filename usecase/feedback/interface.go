package feedback

import (
	"context"

	"github.com/jseow5177/pockteer-be/entity"
)

type UseCase interface {
	CreateFeedback(ctx context.Context, req *CreateFeedbackRequest) (*CreateFeedbackResponse, error)
}

type CreateFeedbackRequest struct {
	UserID *string
	Score  *uint32
	Text   *string
}

func (m *CreateFeedbackRequest) GetUserID() string {
	if m != nil && m.UserID != nil {
		return *m.UserID
	}
	return ""
}

func (m *CreateFeedbackRequest) GetScore() uint32 {
	if m != nil && m.Score != nil {
		return *m.Score
	}
	return 0
}

func (m *CreateFeedbackRequest) GetText() string {
	if m != nil && m.Text != nil {
		return *m.Text
	}
	return ""
}

func (m *CreateFeedbackRequest) ToFeedbackEntity() (*entity.Feedback, error) {
	return entity.NewFeedback(
		m.GetUserID(),
		entity.WithFeedbackScore(m.Score),
		entity.WithFeedbackText(m.Text),
	)
}

type CreateFeedbackResponse struct {
	Feedback *entity.Feedback
}

func (m *CreateFeedbackResponse) GetFeedback() *entity.Feedback {
	if m != nil && m.Feedback != nil {
		return m.Feedback
	}
	return nil
}
