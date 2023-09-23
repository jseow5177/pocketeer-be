package presenter

import (
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"github.com/jseow5177/pockteer-be/usecase/feedback"
)

type Feedback struct {
	Score *uint32 `json:"score,omitempty"`
	Text  *string `json:"text,omitempty"`
}

func (f *Feedback) GetScore() uint32 {
	if f != nil && f.Score != nil {
		return *f.Score
	}
	return 0
}

func (f *Feedback) GetText() string {
	if f != nil && f.Text != nil {
		return *f.Text
	}
	return ""
}

type CreateFeedbackRequest struct {
	Score *uint32 `json:"score,omitempty"`
	Text  *string `json:"text,omitempty"`
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

func (m *CreateFeedbackRequest) ToUseCaseReq(userID string) *feedback.CreateFeedbackRequest {
	return &feedback.CreateFeedbackRequest{
		UserID: goutil.String(userID),
		Score:  m.Score,
		Text:   m.Text,
	}
}

type CreateFeedbackResponse struct {
	Feedback *Feedback `json:"feedback,omitempty"`
}

func (m *CreateFeedbackResponse) GetFeedback() *Feedback {
	if m != nil && m.Feedback != nil {
		return m.Feedback
	}
	return nil
}

func (m *CreateFeedbackResponse) Set(useCaseRes *feedback.CreateFeedbackResponse) {}
