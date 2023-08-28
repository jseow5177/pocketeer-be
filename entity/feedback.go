package entity

import (
	"errors"
	"time"

	"github.com/jseow5177/pockteer-be/pkg/goutil"
)

var (
	ErrScoreOutOfRange = errors.New("score out of range")
	ErrFeedbackTooLong = errors.New("feedback too long")
)

const (
	MaxFeedbackLength = 500
)

type FeedbackScore uint32

const (
	FeedbackScoreZero FeedbackScore = iota
	FeedbackScoreOne
	FeedbackScoreTwo
)

type FeedbackOption = func(f *Feedback)

func WithFeedbackScore(score *uint32) FeedbackOption {
	return func(f *Feedback) {
		f.SetScore(score)
	}
}

func WithFeedbackText(text *string) FeedbackOption {
	return func(f *Feedback) {
		f.SetText(text)
	}
}

type Feedback struct {
	UserID     *string
	Score      *uint32
	Text       *string
	CreateTime *uint64
}

func NewFeedback(userID string, opts ...FeedbackOption) (*Feedback, error) {
	f := &Feedback{
		UserID:     goutil.String(userID),
		Score:      goutil.Uint32(uint32(FeedbackScoreTwo)),
		Text:       goutil.String(""),
		CreateTime: goutil.Uint64(uint64(time.Now().Unix())),
	}
	for _, opt := range opts {
		opt(f)
	}

	if err := f.checkOpts(); err != nil {
		return nil, err
	}

	return f, nil
}

func (f *Feedback) checkOpts() error {
	if f.GetScore() < uint32(FeedbackScoreZero) ||
		f.GetScore() > uint32(FeedbackScoreTwo) {
		return ErrScoreOutOfRange
	}

	if len(f.GetText()) > MaxFeedbackLength {
		return ErrFeedbackTooLong
	}

	return nil
}

func (f *Feedback) GetUserID() string {
	if f != nil && f.UserID != nil {
		return *f.UserID
	}
	return ""
}

func (f *Feedback) SetUserID(userID *string) {
	f.UserID = userID
}

func (f *Feedback) GetScore() uint32 {
	if f != nil && f.Score != nil {
		return *f.Score
	}
	return 0
}

func (f *Feedback) SetScore(score *uint32) {
	f.Score = score
}

func (f *Feedback) GetText() string {
	if f != nil && f.Text != nil {
		return *f.Text
	}
	return ""
}

func (f *Feedback) SetText(text *string) {
	f.Text = text
}

func (f *Feedback) GetCreateTime() uint64 {
	if f != nil && f.CreateTime != nil {
		return *f.CreateTime
	}
	return 0
}

func (f *Feedback) SetCreateTime(createTime *uint64) {
	f.CreateTime = createTime
}
