package model

import "github.com/jseow5177/pockteer-be/entity"

type Feedback struct {
	UserID *string
	Score  *uint32
	Text   *string
}

func ToFeedbackModelFromEntity(f *entity.Feedback) *Feedback {
	if f == nil {
		return nil
	}

	return &Feedback{
		UserID: f.UserID,
		Score:  f.Score,
		Text:   f.Text,
	}
}

func (f *Feedback) GetUserID() string {
	if f != nil && f.UserID != nil {
		return *f.UserID
	}
	return ""
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
