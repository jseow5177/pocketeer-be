package model

import "github.com/jseow5177/pockteer-be/entity"

type OTP struct {
	Code *string
}

func ToOTPModelFromEntity(otp *entity.OTP) *OTP {
	if otp == nil {
		return nil
	}

	return &OTP{
		Code: otp.Code,
	}
}

func ToOTPEntity(otp *OTP) (*entity.OTP, error) {
	if otp == nil {
		return nil, nil
	}

	return entity.NewOTP(
		entity.WithOTPCode(otp.Code),
	)
}

func (otp *OTP) GetCode() string {
	if otp != nil && otp.Code != nil {
		return *otp.Code
	}
	return ""
}
