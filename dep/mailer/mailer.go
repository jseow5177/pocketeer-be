package mailer

import "context"

type Template uint32

const (
	TemplateOTP Template = iota + 1
)

type SendEmailRequest struct {
	To     string
	Params map[string]interface{}
}

type Mailer interface {
	SendEmail(ctx context.Context, template Template, req *SendEmailRequest) error
	Close(ctx context.Context) error
}
