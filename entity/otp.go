package entity

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/jseow5177/pockteer-be/pkg/goutil"
)

var (
	ErrInvalidOTPCode = errors.New("otp code must have 6 digits")
)

type OTP struct {
	Code *string
}

type OTPOption = func(otp *OTP)

func WithOTPCode(code *string) OTPOption {
	return func(otp *OTP) {
		otp.Code = code
	}
}

func NewOTP(opts ...OTPOption) (*OTP, error) {
	rand.Seed(time.Now().UnixNano())

	randNum := rand.Intn(1_000_000)

	otp := &OTP{
		Code: goutil.String(fmt.Sprintf("%06d", randNum)),
	}
	for _, opt := range opts {
		opt(otp)
	}

	if err := otp.checkOpts(); err != nil {
		return nil, err
	}

	return otp, nil
}

func (otp *OTP) checkOpts() error {
	if len(otp.GetCode()) != 6 {
		return ErrInvalidOTPCode
	}
	return nil
}

func (otp *OTP) GetCode() string {
	if otp != nil && otp.Code != nil {
		return *otp.Code
	}
	return ""
}

func (otp *OTP) IsMatch(code string) bool {
	return otp.GetCode() == code
}
