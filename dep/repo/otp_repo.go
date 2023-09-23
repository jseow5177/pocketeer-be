package repo

import (
	"context"
	"errors"

	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/errutil"
)

var (
	ErrOTPNotFound  = errutil.NotFoundError(errors.New("otp not found"))
	ErrInvalidOTP   = errors.New("invalid otp in mem cache")
	ErrDuplicateOTP = errors.New("duplicate otp")
)

type OTPRepo interface {
	Get(ctx context.Context, of *OTPFilter) (*entity.OTP, error)
	Set(ctx context.Context, email string, otp *entity.OTP)
}

type OTPFilter struct {
	Email *string
}

func (f *OTPFilter) GetEmail() string {
	if f != nil && f.Email != nil {
		return *f.Email
	}
	return ""
}
