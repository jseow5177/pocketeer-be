package gmail

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"html/template"

	"github.com/jseow5177/pockteer-be/config"
	"github.com/jseow5177/pockteer-be/dep/mailer"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"gopkg.in/gomail.v2"
)

var (
	ErrTemplateUndefined = errors.New("email template undefined")
	ErrEmptySubject      = errors.New("empty email subject")
)

//go:embed tmpl
var emailTmpls embed.FS

var emailPaths = map[uint32]string{
	uint32(mailer.TemplateOTP): "tmpl/verify_otp.html",
}

var emailSubjects = map[uint32]string{
	uint32(mailer.TemplateOTP): "One-Time Password",
}

type GmailMgr struct {
	name   string
	sender string
	dialer *gomail.Dialer
	tmpls  map[uint32]*template.Template
}

func NewGmailMgr(cfg *config.Gmail) (*GmailMgr, error) {
	dialer := gomail.NewDialer(cfg.Host, cfg.Port, cfg.Email, cfg.Password)

	// init templates
	tmpls := make(map[uint32]*template.Template)
	for id, emailPath := range emailPaths {
		tmpl, err := template.ParseFS(emailTmpls, emailPath)
		if err != nil {
			return nil, fmt.Errorf("fail to parse email tmpl, path: %v, err: %v", emailPath, err)
		}
		tmpls[id] = tmpl
	}

	return &GmailMgr{
		name:   cfg.Name,
		sender: cfg.Email,
		dialer: dialer,
		tmpls:  tmpls,
	}, nil
}

func (mgr *GmailMgr) SendEmail(ctx context.Context, template mailer.Template, req *mailer.SendEmailRequest) error {
	msg := gomail.NewMessage()
	msg.SetAddressHeader("From", "no-reply@bytewise.com", mgr.name)
	msg.SetHeader("To", req.To)

	tmpl, ok := mgr.tmpls[uint32(template)]
	if !ok {
		return ErrTemplateUndefined
	}

	subject := emailSubjects[uint32(template)]
	if subject == "" {
		return ErrEmptySubject
	}
	msg.SetHeader("Subject", subject)

	body, err := goutil.ParseTemplate(tmpl, req.Params)
	if err != nil {
		return fmt.Errorf("fail to parse template, id: %v, params: %v", template, req.Params)
	}
	msg.SetBody("text/html", body)

	return mgr.dialer.DialAndSend(msg)
}

func (mgr *GmailMgr) Close(ctx context.Context) error {
	return nil
}
