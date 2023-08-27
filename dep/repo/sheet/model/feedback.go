package model

import (
	"time"

	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
)

const dateTimeFmt = "2006-01-02 15:04:05"

type Feedback struct {
	UserID     *string
	Score      *uint32
	Text       *string
	CreateTime *string
}

func ToFeedbackModelFromEntity(f *entity.Feedback) *Feedback {
	if f == nil {
		return nil
	}

	var (
		fmtTs string
		ts    = time.Unix(int64(f.GetCreateTime()), 0)
	)

	tz, err := time.LoadLocation("Asia/Singapore")
	if err == nil {
		ts = ts.In(tz)
		fmtTs = ts.Format(dateTimeFmt)
	}

	return &Feedback{
		UserID:     f.UserID,
		Score:      f.Score,
		Text:       f.Text,
		CreateTime: goutil.String(fmtTs),
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

func (f *Feedback) GetCreateTime() string {
	if f != nil && f.CreateTime != nil {
		return *f.CreateTime
	}
	return ""
}
